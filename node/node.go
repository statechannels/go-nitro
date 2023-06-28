// Package node contains imperative library code for running a go-nitro node inside another application.
package node // import "github.com/statechannels/go-nitro/node"

import (
	"context"
	"io"
	"math/big"
	"runtime/debug"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/node/engine"
	"github.com/statechannels/go-nitro/node/engine/chainservice"
	"github.com/statechannels/go-nitro/node/engine/messageservice"
	"github.com/statechannels/go-nitro/node/engine/store"
	"github.com/statechannels/go-nitro/node/notifier"
	"github.com/statechannels/go-nitro/node/query"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/rand"
	"github.com/statechannels/go-nitro/types"
)

// Node provides the interface for the consuming application
type Node struct {
	engine          engine.Engine // The core business logic of the node
	Address         *types.Address
	channelNotifier *notifier.ChannelNotifier

	completedObjectivesForRPC chan protocols.ObjectiveId // This is only used by the RPC server
	completedObjectives       *safesync.Map[chan struct{}]
	failedObjectives          chan protocols.ObjectiveId
	receivedVouchers          chan payments.Voucher
	chainId                   *big.Int
	store                     store.Store
	vm                        *payments.VoucherManager
	logger                    zerolog.Logger
	cancelEventHandler        context.CancelFunc

	wg *sync.WaitGroup
}

// New is the constructor for a Node. It accepts a messaging service, a chain service, and a store as injected dependencies.
func New(messageService messageservice.MessageService, chainservice chainservice.ChainService, store store.Store, logDestination io.Writer, policymaker engine.PolicyMaker, metricsApi engine.MetricsApi) Node {
	n := Node{}
	n.Address = store.GetAddress()
	// If a metrics API is not provided we used the no-op version which does nothing.
	if metricsApi == nil {
		metricsApi = &engine.NoOpMetrics{}
	}
	chainId, err := chainservice.GetChainId()
	if err != nil {
		panic(err)
	}
	n.chainId = chainId
	n.store = store
	n.vm = payments.NewVoucherManager(*store.GetAddress(), store)
	n.logger = zerolog.New(logDestination).With().Timestamp().Str("nitro node", n.Address.String()[0:8]).Caller().Logger()

	n.engine = engine.New(n.vm, messageService, chainservice, store, logDestination, policymaker, metricsApi)
	n.completedObjectives = &safesync.Map[chan struct{}]{}
	n.completedObjectivesForRPC = make(chan protocols.ObjectiveId, 100)

	n.failedObjectives = make(chan protocols.ObjectiveId, 100)
	// Using a larger buffer since payments can be sent frequently.
	n.receivedVouchers = make(chan payments.Voucher, 1000)

	n.channelNotifier = notifier.NewChannelNotifier(store, n.vm)

	// Start the engine in a go routine
	ctx, cancel := context.WithCancel(context.Background())

	n.wg = &sync.WaitGroup{}
	n.wg.Add(1)

	n.cancelEventHandler = cancel
	// Start the event handler in a go routine
	// It will listen for events from the engine and dispatch events to node channels
	go n.handleEngineEvents(ctx)

	return n
}

// handleEngineEvents is responsible for monitoring the ToApi channel on the engine.
// It parses events from the ToApi chan and then dispatches events to the necessary node chan.
func (n *Node) handleEngineEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			n.wg.Done()
			return
		case update := <-n.engine.ToApi():
			for _, completed := range update.CompletedObjectives {
				d, _ := n.completedObjectives.LoadOrStore(string(completed.Id()), make(chan struct{}))
				close(d)

				// use a nonblocking send to the RPC Client in case no one is listening
				select {
				case n.completedObjectivesForRPC <- completed.Id():
				default:
				}
			}

			for _, erred := range update.FailedObjectives {
				n.failedObjectives <- erred
			}

			for _, payment := range update.ReceivedVouchers {
				n.receivedVouchers <- payment
			}

			for _, updated := range update.LedgerChannelUpdates {

				err := n.channelNotifier.NotifyLedgerUpdated(updated)
				n.handleError(err)
			}
			for _, updated := range update.PaymentChannelUpdates {

				err := n.channelNotifier.NotifyPaymentUpdated(updated)
				n.handleError(err)
			}
		}
	}
}

// Begin API

// Version returns the go-nitro version
func (n *Node) Version() string {
	info, _ := debug.ReadBuildInfo()

	version := info.Main.Version
	// Depending on how the binary was built we may get back no version info.
	// In this case we default to "(devel)".
	// See https://github.com/golang/go/issues/51831#issuecomment-1074188363 for more details.
	if version == "" {
		version = "(devel)"
	}

	// If the binary was built with the -buildvcs flag we can get the git commit hash and use that as the version.
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			version = s.Value
			break
		}
	}

	return version
}

// CompletedObjectives returns a chan that receives a objective id whenever that objective is completed. Not suitable fo multiple subscribers.
func (n *Node) CompletedObjectives() <-chan protocols.ObjectiveId {
	return n.completedObjectivesForRPC
}

// LedgerUpdates returns a chan that receives ledger channel info whenever that ledger channel is updated. Not suitable for multiple subscribers.
func (n *Node) LedgerUpdates() <-chan query.LedgerChannelInfo {
	return n.channelNotifier.RegisterForAllLedgerUpdates()
}

// PaymentUpdates returns a chan that receives payment channel info whenever that payment channel is updated. Not suitable fo multiple subscribers.
func (n *Node) PaymentUpdates() <-chan query.PaymentChannelInfo {
	return n.channelNotifier.RegisterForAllPaymentUpdates()
}

// ObjectiveCompleteChan returns a chan that is closed when the objective with given id is completed
func (n *Node) ObjectiveCompleteChan(id protocols.ObjectiveId) <-chan struct{} {
	d, _ := n.completedObjectives.LoadOrStore(string(id), make(chan struct{}))
	return d
}

// LedgerUpdatedChan returns a chan that receives a ledger channel info whenever the ledger with given id is updated
func (n *Node) LedgerUpdatedChan(ledgerId types.Destination) <-chan query.LedgerChannelInfo {
	return n.channelNotifier.RegisterForLedgerUpdates(ledgerId)
}

// PaymentChannelUpdatedChan returns a chan that receives a payment channel info whenever the payment channel with given id is updated
func (n *Node) PaymentChannelUpdatedChan(ledgerId types.Destination) <-chan query.PaymentChannelInfo {
	return n.channelNotifier.RegisterForPaymentChannelUpdates(ledgerId)
}

// FailedObjectives returns a chan that receives an objective id whenever that objective has failed
func (n *Node) FailedObjectives() <-chan protocols.ObjectiveId {
	return n.failedObjectives
}

// ReceivedVouchers returns a chan that receives a voucher every time we receive a payment voucher
func (n *Node) ReceivedVouchers() <-chan payments.Voucher {
	return n.receivedVouchers
}

// CreateVoucher creates and returns a voucher for the given channelId which increments the redeemable balance by amount.
// It is the responsibility of the caller to send the voucher to the payee.
func (c *Node) CreateVoucher(channelId types.Destination, amount *big.Int) (payments.Voucher, error) {
	return c.vm.Pay(channelId, amount, *c.store.GetChannelSecretKey())
}

// ReceiveVoucher receives a voucher and returns the amount that was paid.
// It can be used to add a voucher that was sent outside of the go-nitro system.
func (c *Node) ReceiveVoucher(v payments.Voucher) (total *big.Int, fromVoucher *big.Int, err error) {
	return c.vm.Receive(v)
}

// CreatePaymentChannel creates a virtual channel with the counterParty using ledger channels
// with the supplied intermediaries.
func (n *Node) CreatePaymentChannel(Intermediaries []types.Address, CounterParty types.Address, ChallengeDuration uint32, Outcome outcome.Exit) (virtualfund.ObjectiveResponse, error) {
	objectiveRequest := virtualfund.NewObjectiveRequest(
		Intermediaries,
		CounterParty,
		ChallengeDuration,
		Outcome,
		rand.Uint64(),
		n.engine.GetVirtualPaymentAppAddress(),
	)

	// Send the event to the engine
	n.engine.ObjectiveRequestsFromAPI <- objectiveRequest

	objectiveRequest.WaitForObjectiveToStart()
	return objectiveRequest.Response(*n.Address), nil
}

// ClosePaymentChannel attempts to close and defund the given virtually funded channel.
func (n *Node) ClosePaymentChannel(channelId types.Destination) (protocols.ObjectiveId, error) {
	objectiveRequest := virtualdefund.NewObjectiveRequest(channelId)

	// Send the event to the engine
	n.engine.ObjectiveRequestsFromAPI <- objectiveRequest
	objectiveRequest.WaitForObjectiveToStart()
	return objectiveRequest.Id(*n.Address, n.chainId), nil
}

// CreateLedgerChannel creates a directly funded ledger channel with the given counterparty.
// The channel will run under full consensus rules (it is not possible to provide a custom AppDefinition or AppData).
func (n *Node) CreateLedgerChannel(Counterparty types.Address, ChallengeDuration uint32, outcome outcome.Exit) (directfund.ObjectiveResponse, error) {
	objectiveRequest := directfund.NewObjectiveRequest(
		Counterparty,
		ChallengeDuration,
		outcome,
		rand.Uint64(),
		n.engine.GetConsensusAppAddress(),
		// Appdata implicitly zero
	)

	// Send the event to the engine
	n.engine.ObjectiveRequestsFromAPI <- objectiveRequest
	objectiveRequest.WaitForObjectiveToStart()
	return objectiveRequest.Response(*n.Address, n.chainId), nil
}

// CloseLedgerChannel attempts to close and defund the given directly funded channel.
func (n *Node) CloseLedgerChannel(channelId types.Destination) (protocols.ObjectiveId, error) {
	objectiveRequest := directdefund.NewObjectiveRequest(channelId)

	// Send the event to the engine
	n.engine.ObjectiveRequestsFromAPI <- objectiveRequest
	objectiveRequest.WaitForObjectiveToStart()
	return objectiveRequest.Id(*n.Address, n.chainId), nil
}

// Pay will send a signed voucher to the payee that they can redeem for the given amount.
func (n *Node) Pay(channelId types.Destination, amount *big.Int) {
	// Send the event to the engine
	n.engine.PaymentRequestsFromAPI <- engine.PaymentRequest{ChannelId: channelId, Amount: amount}
}

// GetPaymentChannel returns the payment channel with the given id.
// If no ledger channel exists with the given id an error is returned.
func (n *Node) GetPaymentChannel(id types.Destination) (query.PaymentChannelInfo, error) {
	return query.GetPaymentChannelInfo(id, n.store, n.vm)
}

// GetPaymentChannelsByLedger returns all active payment channels that are funded by the given ledger channel.
func (n *Node) GetPaymentChannelsByLedger(ledgerId types.Destination) ([]query.PaymentChannelInfo, error) {
	return query.GetPaymentChannelsByLedger(ledgerId, n.store, n.vm)
}

// GetAllLedgerChannels returns all ledger channels.
func (n *Node) GetAllLedgerChannels() ([]query.LedgerChannelInfo, error) {
	return query.GetAllLedgerChannels(n.store, n.engine.GetConsensusAppAddress())
}

// GetLedgerChannel returns the ledger channel with the given id.
// If no ledger channel exists with the given id an error is returned.
func (n *Node) GetLedgerChannel(id types.Destination) (query.LedgerChannelInfo, error) {
	return query.GetLedgerChannelInfo(id, n.store)
}

// stopEventHandler stops the event handler goroutine and waits for it to quit successfully.
func (n *Node) stopEventHandler() {
	n.cancelEventHandler()
	n.wg.Wait()
}

// Close stops the node from responding to any input.
func (n *Node) Close() error {
	n.stopEventHandler()

	if err := n.channelNotifier.Close(); err != nil {
		return err
	}
	if err := n.engine.Close(); err != nil {
		return err
	}
	// At this point, the engine ToApi channel has been closed.
	// If there are blocking consumers (for or select channel statements) on any channel for which the node is a producer,
	// those channels need to be closed.
	close(n.completedObjectivesForRPC)

	return n.store.Close()
}

// handleError logs the error and panics
// Eventually it should return the error to the caller
func (n *Node) handleError(err error) {
	if err != nil {

		n.logger.Err(err).Msgf("%s, error in nitro node", n.Address)

		<-time.After(1000 * time.Millisecond) // We wait for a bit so the previous log line has time to complete

		// TODO instead of a panic, errors should be returned to the caller.
		panic(err)

	}
}
