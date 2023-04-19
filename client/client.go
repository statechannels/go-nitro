// Package client contains imperative library code for running a go-nitro client inside another application.
package client // import "github.com/statechannels/go-nitro/client"

import (
	"fmt"
	"io"
	"math/big"
	"runtime/debug"

	eb "github.com/asaskevich/EventBus"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/client/query"
	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/rand"
	"github.com/statechannels/go-nitro/types"
)

// Client provides the interface for the consuming application
type Client struct {
	engine                    engine.Engine // The core business logic of the client
	Address                   *types.Address
	ledgerUpdates             *safesync.Map[chan query.LedgerChannelInfo]
	ledgerUpdatesForRPC       chan query.LedgerChannelInfo
	completedObjectivesForRPC chan protocols.ObjectiveId // This is only used by the RPC server
	completedObjectives       *safesync.Map[chan struct{}]
	failedObjectives          chan protocols.ObjectiveId
	receivedVouchers          chan payments.Voucher
	chainId                   *big.Int
	store                     store.Store
	vm                        *payments.VoucherManager
	updateBus                 eb.Bus
}

// New is the constructor for a Client. It accepts a messaging service, a chain service, and a store as injected dependencies.
func New(messageService messageservice.MessageService, chainservice chainservice.ChainService, store store.Store, logDestination io.Writer, policymaker engine.PolicyMaker, metricsApi engine.MetricsApi) Client {
	c := Client{}
	c.Address = store.GetAddress()
	// If a metrics API is not provided we used the no-op version which does nothing.
	if metricsApi == nil {
		metricsApi = &engine.NoOpMetrics{}
	}
	chainId, err := chainservice.GetChainId()
	if err != nil {
		panic(err)
	}
	c.chainId = chainId
	c.store = store
	c.vm = payments.NewVoucherManager(*store.GetAddress(), store)
	c.engine = engine.New(c.vm, messageService, chainservice, store, logDestination, policymaker, metricsApi)
	c.completedObjectives = &safesync.Map[chan struct{}]{}
	c.completedObjectivesForRPC = make(chan protocols.ObjectiveId, 100)

	c.ledgerUpdates = &safesync.Map[chan query.LedgerChannelInfo]{}
	c.ledgerUpdatesForRPC = make(chan query.LedgerChannelInfo, 100)

	c.completedObjectivesForRPC = make(chan protocols.ObjectiveId, 100)

	c.failedObjectives = make(chan protocols.ObjectiveId, 100)
	// Using a larger buffer since payments can be sent frequently.
	c.receivedVouchers = make(chan payments.Voucher, 1000)

	c.updateBus = eb.New()
	// Start the engine in a go routine
	go c.engine.Run()

	// Start the event handler in a go routine
	// It will listen for events from the engine and dispatch events to client channels
	go c.handleEngineEvents()

	return c
}

// handleEngineEvents is responsible for monitoring the ToApi channel on the engine.
// It parses events from the ToApi chan and then dispatches events to the necessary client chan.
func (c *Client) handleEngineEvents() {
	for update := range c.engine.ToApi() {

		for _, completed := range update.CompletedObjectives {
			d, _ := c.completedObjectives.LoadOrStore(string(completed.Id()), make(chan struct{}))
			close(d)

			// use a nonblocking send to the RPC Client in case no one is listening
			select {
			case c.completedObjectivesForRPC <- completed.Id():
			default:
			}
		}

		for _, erred := range update.FailedObjectives {
			c.failedObjectives <- erred
		}

		for _, payment := range update.ReceivedVouchers {
			c.receivedVouchers <- payment
		}

		for _, ledgerId := range update.UpdatedChannels {

			ledger, err := query.GetLedgerChannelInfo(ledgerId, c.store)
			// TODO: What's the best way of handling this error?
			if err != nil {
				panic(err)
			}

			c.updateBus.Publish(ledgerId.String(), ledger)

			// use a nonblocking send to the RPC Client in case no one is listening
			select {
			case c.ledgerUpdatesForRPC <- ledger:
			default:
			}
		}
	}

	// At this point, the engine ToApi channel has been closed.
	// If there are blocking consumers (for or select channel statements) on any channel for which the client is a producer,
	// those channels need to be closed.
	close(c.completedObjectivesForRPC)
}

// Begin API

// Version returns the go-nitro version
func (c *Client) Version() string {
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
func (c *Client) CompletedObjectives() <-chan protocols.ObjectiveId {
	return c.completedObjectivesForRPC
}

func (c *Client) LedgerUpdates() <-chan query.LedgerChannelInfo {
	return c.ledgerUpdatesForRPC
}

// LedgerUpdatedChan returns a chan that receives an empty struct when the objective with given id is completed
func (c *Client) ObjectiveCompleteChan(id protocols.ObjectiveId) <-chan struct{} {
	d, _ := c.completedObjectives.LoadOrStore(string(id), make(chan struct{}))
	return d
}

func (c *Client) LedgerUpdatedChan(ledgerId types.Destination) <-chan query.LedgerChannelInfo {
	ledgerChan := make(chan query.LedgerChannelInfo, 100)
	err := c.updateBus.Subscribe(ledgerId.String(), func(ledger query.LedgerChannelInfo) {
		fmt.Printf("LedgerUpdatedChan: %+v\n", ledger)
		ledgerChan <- ledger
	})
	if err != nil {
		panic(err)
	}
	return ledgerChan
}

// FailedObjectives returns a chan that receives an objective id whenever that objective has failed
func (c *Client) FailedObjectives() <-chan protocols.ObjectiveId {
	return c.failedObjectives
}

// ReceivedVouchers returns a chan that receives a voucher every time we receive a payment voucher
func (c *Client) ReceivedVouchers() <-chan payments.Voucher {
	return c.receivedVouchers
}

// CreateVirtualChannel creates a virtual channel with the counterParty using ledger channels
// with the supplied intermediaries.
func (c *Client) CreateVirtualPaymentChannel(Intermediaries []types.Address, CounterParty types.Address, ChallengeDuration uint32, Outcome outcome.Exit) virtualfund.ObjectiveResponse {
	objectiveRequest := virtualfund.NewObjectiveRequest(
		Intermediaries,
		CounterParty,
		ChallengeDuration,
		Outcome,
		rand.Uint64(),
		c.engine.GetVirtualPaymentAppAddress(),
	)

	// Send the event to the engine
	c.engine.ObjectiveRequestsFromAPI <- objectiveRequest

	objectiveRequest.WaitForObjectiveToStart()
	return objectiveRequest.Response(*c.Address)
}

// CloseVirtualChannel attempts to close and defund the given virtually funded channel.
func (c *Client) CloseVirtualChannel(channelId types.Destination) protocols.ObjectiveId {
	objectiveRequest := virtualdefund.NewObjectiveRequest(channelId)

	// Send the event to the engine
	c.engine.ObjectiveRequestsFromAPI <- objectiveRequest
	objectiveRequest.WaitForObjectiveToStart()
	return objectiveRequest.Id(*c.Address, c.chainId)
}

// CreateLedgerChannel creates a directly funded ledger channel with the given counterparty.
// The channel will run under full consensus rules (it is not possible to provide a custom AppDefinition or AppData).
func (c *Client) CreateLedgerChannel(Counterparty types.Address, ChallengeDuration uint32, outcome outcome.Exit) directfund.ObjectiveResponse {
	objectiveRequest := directfund.NewObjectiveRequest(
		Counterparty,
		ChallengeDuration,
		outcome,
		rand.Uint64(),
		c.engine.GetConsensusAppAddress(),
		// Appdata implicitly zero
	)

	// Send the event to the engine
	c.engine.ObjectiveRequestsFromAPI <- objectiveRequest
	objectiveRequest.WaitForObjectiveToStart()
	return objectiveRequest.Response(*c.Address, c.chainId)
}

// CloseLedgerChannel attempts to close and defund the given directly funded channel.
func (c *Client) CloseLedgerChannel(channelId types.Destination) protocols.ObjectiveId {
	objectiveRequest := directdefund.NewObjectiveRequest(channelId)

	// Send the event to the engine
	c.engine.ObjectiveRequestsFromAPI <- objectiveRequest
	objectiveRequest.WaitForObjectiveToStart()
	return objectiveRequest.Id(*c.Address, c.chainId)
}

// Pay will send a signed voucher to the payee that they can redeem for the given amount.
func (c *Client) Pay(channelId types.Destination, amount *big.Int) {
	// Send the event to the engine
	c.engine.PaymentRequestsFromAPI <- engine.PaymentRequest{ChannelId: channelId, Amount: amount}
}

// GetPaymentChannel returns the payment channel with the given id.
// If no ledger channel exists with the given id an error is returned.
func (c *Client) GetPaymentChannel(id types.Destination) (query.PaymentChannelInfo, error) {
	return query.GetPaymentChannelInfo(id, c.store, c.vm)
}

// GetLedgerChannel returns the ledger channel with the given id.
// If no ledger channel exists with the given id an error is returned.
func (c *Client) GetLedgerChannel(id types.Destination) (query.LedgerChannelInfo, error) {
	return query.GetLedgerChannelInfo(id, c.store)
}

// Close stops the client from responding to any input.
func (c *Client) Close() error {
	if err := c.engine.Close(); err != nil {
		return err
	}
	return c.store.Close()
}
