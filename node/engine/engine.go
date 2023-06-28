// Package engine contains the types and imperative code for the business logic of a go-nitro Node.
package engine // import "github.com/statechannels/go-nitro/node/engine"

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync"

	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/internal/logging"
	"github.com/statechannels/go-nitro/node/engine/chainservice"
	"github.com/statechannels/go-nitro/node/engine/messageservice"
	"github.com/statechannels/go-nitro/node/engine/store"
	"github.com/statechannels/go-nitro/node/query"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

// ErrUnhandledChainEvent is an engine error when the the engine cannot process a chain event
type ErrUnhandledChainEvent struct {
	event     chainservice.Event
	objective protocols.Objective
	reason    string
}

func (uce *ErrUnhandledChainEvent) Error() string {
	return fmt.Sprintf("chain event %#v could not be handled by objective %#v due to: %s", uce.event, uce.objective, uce.reason)
}

type ErrGetObjective struct {
	wrappedError error
	objectiveId  protocols.ObjectiveId
}

func (e *ErrGetObjective) Error() string {
	return fmt.Sprintf("unexpected error getting/creating objective %s: %v", e.objectiveId, e.wrappedError)
}

// nonFatalErrors is a list of errors for which the engine should not panic
var nonFatalErrors = []error{
	&ErrGetObjective{},
	store.ErrLoadVouchers,
	directfund.ErrLedgerChannelExists,
}

// Engine is the imperative part of the core business logic of a go-nitro Node
type Engine struct {
	// inbound go channels

	// From API
	ObjectiveRequestsFromAPI chan protocols.ObjectiveRequest
	PaymentRequestsFromAPI   chan PaymentRequest

	fromChain  <-chan chainservice.Event
	fromMsg    <-chan protocols.Message
	fromLedger chan consensus_channel.Proposal

	toApi chan EngineEvent

	msg   messageservice.MessageService
	chain chainservice.ChainService

	store       store.Store // A Store for persisting and restoring important data
	policymaker PolicyMaker // A PolicyMaker decides whether to approve or reject objectives

	logger zerolog.Logger

	metrics *MetricsRecorder

	vm *payments.VoucherManager

	wg     *sync.WaitGroup
	cancel context.CancelFunc
}

// PaymentRequest represents a request from the API to make a payment using a channel
type PaymentRequest struct {
	ChannelId types.Destination
	Amount    *big.Int
}

// EngineEvent is a struct that contains a list of changes caused by handling a message/chain event/api event
type EngineEvent struct {
	// These are objectives that are now completed
	CompletedObjectives []protocols.Objective
	// These are objectives that have failed
	FailedObjectives []protocols.ObjectiveId
	// ReceivedVouchers are vouchers we've received from other participants
	ReceivedVouchers []payments.Voucher

	// LedgerChannelUpdates contains channel info for ledger channels that have been updated
	LedgerChannelUpdates []query.LedgerChannelInfo
	// PaymentChannelUpdates contains channel info for payment channels that have been updated
	PaymentChannelUpdates []query.PaymentChannelInfo
}

// IsEmpty returns true if the EngineEvent contains no changes
func (ee *EngineEvent) IsEmpty() bool {
	return len(ee.CompletedObjectives) == 0 &&
		len(ee.FailedObjectives) == 0 &&
		len(ee.ReceivedVouchers) == 0 &&
		len(ee.LedgerChannelUpdates) == 0 &&
		len(ee.PaymentChannelUpdates) == 0
}

func (ee *EngineEvent) Merge(other EngineEvent) {
	ee.CompletedObjectives = append(ee.CompletedObjectives, other.CompletedObjectives...)
	ee.FailedObjectives = append(ee.FailedObjectives, other.FailedObjectives...)
	ee.ReceivedVouchers = append(ee.ReceivedVouchers, other.ReceivedVouchers...)
	ee.LedgerChannelUpdates = append(ee.LedgerChannelUpdates, other.LedgerChannelUpdates...)
	ee.PaymentChannelUpdates = append(ee.PaymentChannelUpdates, other.PaymentChannelUpdates...)
}

type CompletedObjectiveEvent struct {
	Id protocols.ObjectiveId
}

// Response is the return type that asynchronous API calls "resolve to". Such a call returns a go channel of type Response.
type Response struct{}

// NewEngine is the constructor for an Engine
func New(vm *payments.VoucherManager, msg messageservice.MessageService, chain chainservice.ChainService, store store.Store, logDestination io.Writer, policymaker PolicyMaker, metricsApi MetricsApi) Engine {
	e := Engine{}

	e.store = store

	e.fromLedger = make(chan consensus_channel.Proposal, 100)
	// bind to inbound chans
	e.ObjectiveRequestsFromAPI = make(chan protocols.ObjectiveRequest)
	e.PaymentRequestsFromAPI = make(chan PaymentRequest)

	e.fromChain = chain.EventFeed()
	e.fromMsg = msg.Out()

	e.chain = chain
	e.msg = msg

	e.toApi = make(chan EngineEvent, 100)

	logging.ConfigureZeroLogger()
	e.logger = zerolog.New(logDestination).With().Timestamp().Str("engine", e.store.GetAddress().String()[0:8]).Caller().Logger()

	e.policymaker = policymaker

	e.vm = vm

	e.logger.Print("Constructed Engine")

	if metricsApi == nil {
		metricsApi = &NoOpMetrics{}
	}
	e.metrics = NewMetricsRecorder(*e.store.GetAddress(), metricsApi)

	e.wg = &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	e.cancel = cancel

	e.wg.Add(1)
	go e.run(ctx)

	return e
}

func (e *Engine) ToApi() <-chan EngineEvent {
	return e.toApi
}

func (e *Engine) Close() error {
	e.cancel()
	e.wg.Wait()
	if err := e.msg.Close(); err != nil {
		return err
	}

	close(e.toApi)
	return e.chain.Close()
}

// run kicks of an infinite loop that waits for communications on the supplied channels, and handles them accordingly
// The loop exits when the context is cancelled.
func (e *Engine) run(ctx context.Context) {
	for {
		var res EngineEvent
		var err error

		e.metrics.RecordQueueLength("api_objective_request_queue", len(e.ObjectiveRequestsFromAPI))
		e.metrics.RecordQueueLength("api_payment_request_queue", len(e.PaymentRequestsFromAPI))
		e.metrics.RecordQueueLength("chain_events_queue", len(e.fromChain))
		e.metrics.RecordQueueLength("messages_queue", len(e.fromMsg))
		e.metrics.RecordQueueLength("proposal_queue", len(e.fromLedger))

		select {

		case or := <-e.ObjectiveRequestsFromAPI:
			res, err = e.handleObjectiveRequest(or)
		case pr := <-e.PaymentRequestsFromAPI:
			res, err = e.handlePaymentRequest(pr)
		case chainEvent := <-e.fromChain:
			res, err = e.handleChainEvent(chainEvent)
		case message := <-e.fromMsg:
			res, err = e.handleMessage(message)
		case proposal := <-e.fromLedger:
			res, err = e.handleProposal(proposal)
		case <-ctx.Done():
			e.wg.Done()
			return
		}

		// Handle errors
		e.checkError(err)

		// Only send out an event if there are changes
		if !res.IsEmpty() {

			for _, obj := range res.CompletedObjectives {
				e.logger.Printf("Objective %s is complete & returned to API", obj.Id())
				e.metrics.RecordObjectiveCompleted(obj.Id())
			}
			e.toApi <- res
		}

	}
}

// handleProposal handles a Proposal returned to the engine from
// a running ledger channel by pulling its corresponding objective
// from the store and attempting progress.
func (e *Engine) handleProposal(proposal consensus_channel.Proposal) (EngineEvent, error) {
	defer e.metrics.RecordFunctionDuration()()

	id := getProposalObjectiveId(proposal)

	obj, err := e.store.GetObjectiveById(id)
	if err != nil {
		return EngineEvent{}, err
	}
	if obj.GetStatus() == protocols.Completed {
		e.logger.Printf("Ignoring proposal for complected objective  %s", obj.Id())
		return EngineEvent{}, nil
	}
	return e.attemptProgress(obj)
}

// handleMessage handles a Message from a peer go-nitro Wallet.
// It:
//   - reads an objective from the store,
//   - generates an updated objective,
//   - attempts progress on the target Objective,
//   - attempts progress on related objectives which may have become unblocked.
func (e *Engine) handleMessage(message protocols.Message) (EngineEvent, error) {
	defer e.metrics.RecordFunctionDuration()()
	e.logMessage(message, Incoming)
	allCompleted := EngineEvent{}

	for _, payload := range message.ObjectivePayloads {

		objective, err := e.getOrCreateObjective(payload)
		if err != nil {
			return EngineEvent{}, err
		}

		if objective.GetStatus() == protocols.Unapproved {
			e.logger.Printf("Policymaker is %+v", e.policymaker)
			if e.policymaker.ShouldApprove(objective) {
				objective = objective.Approve()

				ddfo, ok := objective.(*directdefund.Objective)
				if ok {
					// If we just approved a direct defund objective, destroy the consensus channel to prevent it being used (a Channel will now take over governance)
					err := e.store.DestroyConsensusChannel(ddfo.C.Id)
					if err != nil {
						return EngineEvent{}, err
					}
				}
			} else {
				objective, sideEffects := objective.Reject()
				err = e.store.SetObjective(objective)
				if err != nil {
					return EngineEvent{}, err
				}

				allCompleted.CompletedObjectives = append(allCompleted.CompletedObjectives, objective)

				err = e.executeSideEffects(sideEffects)
				// An error would mean we failed to send a message. But the objective is still "completed".
				// So, we should return allCompleted even if there was an error.
				return allCompleted, err
			}
		}

		if objective.GetStatus() == protocols.Completed {
			e.logger.Printf("Ignoring payload for complected objective  %s", objective.Id())
			continue
		}
		if objective.GetStatus() == protocols.Rejected {
			e.logger.Printf("Ignoring payload for rejected objective  %s", objective.Id())
			continue
		}

		updatedObjective, err := objective.Update(payload)
		if err != nil {
			return EngineEvent{}, err
		}

		progressEvent, err := e.attemptProgress(updatedObjective)
		if err != nil {
			return EngineEvent{}, err
		}

		allCompleted.Merge(progressEvent)

		if err != nil {
			return EngineEvent{}, err
		}

	}

	for _, entry := range message.LedgerProposals { // The ledger protocol requires us to process these proposals in turnNum order.
		// Here we rely on the sender having packed them into the message in that order, and do not apply any checks or sorting of our own.
		id := getProposalObjectiveId(entry.Proposal)

		o, err := e.store.GetObjectiveById(id)
		if err != nil {
			return EngineEvent{}, err
		}
		if o.GetStatus() == protocols.Completed {
			e.logger.Printf("Ignoring payload for complected objective  %s", o.Id())
			continue
		}
		objective, isProposalReceiver := o.(protocols.ProposalReceiver)
		if !isProposalReceiver {
			return EngineEvent{}, fmt.Errorf("received a proposal for an objective which cannot receive proposals %s", objective.Id())
		}

		updatedObjective, err := objective.ReceiveProposal(entry)
		if err != nil {
			return EngineEvent{}, err
		}

		progressEvent, err := e.attemptProgress(updatedObjective)
		if err != nil {
			return EngineEvent{}, err
		}

		allCompleted.Merge(progressEvent)

	}

	for _, entry := range message.RejectedObjectives {
		objective, err := e.store.GetObjectiveById(entry)
		if err != nil {
			return EngineEvent{}, err
		}
		if objective.GetStatus() == protocols.Rejected {
			e.logger.Printf("Ignoring payload for rejected objective  %s", objective.Id())
			continue
		}

		// we are rejecting due to a counterparty message notifying us of their rejection. We
		// do not need to send a message back to that counterparty, and furthermore we assume that
		// counterparty has already notified all other interested parties. We can therefore ignore the side effects
		objective, _ = objective.Reject()
		err = e.store.SetObjective(objective)
		if err != nil {
			return EngineEvent{}, err
		}

		allCompleted.CompletedObjectives = append(allCompleted.CompletedObjectives, objective)
	}

	for _, voucher := range message.Payments {

		// TODO: return the amount we paid?
		_, err := e.vm.Receive(voucher)

		allCompleted.ReceivedVouchers = append(allCompleted.ReceivedVouchers, voucher)
		if err != nil {
			return EngineEvent{}, fmt.Errorf("error accepting payment voucher: %w", err)
		}
		c, ok := e.store.GetChannelById(voucher.ChannelId)
		if !ok {
			return EngineEvent{}, fmt.Errorf("could not fetch channel for voucher %+v", voucher)
		}

		// Vouchers only count as payment channel updates if the channel is open.
		if !c.FinalCompleted() {

			paid, remaining, err := query.GetVoucherBalance(c.Id, e.vm)
			if err != nil {
				return EngineEvent{}, err
			}
			info, err := query.ConstructPaymentInfo(c, paid, remaining)
			if err != nil {
				return EngineEvent{}, err
			}
			allCompleted.PaymentChannelUpdates = append(allCompleted.PaymentChannelUpdates, info)
		}

	}
	return allCompleted, nil
}

// handleChainEvent handles a Chain Event from the blockchain.
// It:
//   - reads an objective from the store,
//   - generates an updated objective, and
//   - attempts progress.
func (e *Engine) handleChainEvent(chainEvent chainservice.Event) (EngineEvent, error) {
	defer e.metrics.RecordFunctionDuration()()
	e.logger.Printf("handling chain event: %v", chainEvent)
	objective, ok := e.store.GetObjectiveByChannelId(chainEvent.ChannelID())
	if !ok {
		// TODO: Right now the chain service returns chain events for ALL channels even those we aren't involved in
		// for now we can ignore channels we aren't involved in
		// in the future the chain service should allow us to register for specific channels
		return EngineEvent{}, nil
	}

	eventHandler, ok := objective.(chainservice.ChainEventHandler)
	if !ok {
		return EngineEvent{}, &ErrUnhandledChainEvent{event: chainEvent, objective: objective, reason: "objective does not handle chain events"}
	}
	updatedEventHandler, err := eventHandler.UpdateWithChainEvent(chainEvent)
	if err != nil {
		return EngineEvent{}, err
	}
	return e.attemptProgress(updatedEventHandler)
}

// handleObjectiveRequest handles an ObjectiveRequest (triggered by a client API call).
// It will attempt to spawn a new, approved objective.
func (e *Engine) handleObjectiveRequest(or protocols.ObjectiveRequest) (EngineEvent, error) {
	defer e.metrics.RecordFunctionDuration()()
	myAddress := *e.store.GetAddress()

	chainId, err := e.chain.GetChainId()
	if err != nil {
		return EngineEvent{}, fmt.Errorf("could not get chain id from chain service: %w", err)
	}

	objectiveId := or.Id(myAddress, chainId)
	failedEngineEvent := EngineEvent{FailedObjectives: []protocols.ObjectiveId{objectiveId}}
	e.logger.Printf("handling new objective request for %s", objectiveId)
	e.metrics.RecordObjectiveStarted(objectiveId)
	defer or.SignalObjectiveStarted()
	switch request := or.(type) {

	case virtualfund.ObjectiveRequest:
		vfo, err := virtualfund.NewObjective(request, true, myAddress, chainId, e.store.GetConsensusChannel)
		if err != nil {
			return failedEngineEvent, fmt.Errorf("handleAPIEvent: Could not create virtualfund objective for %+v: %w", request, err)
		}
		// Only Alice or Bob care about registering the objective and keeping track of vouchers
		lastParticipant := uint(len(vfo.V.Participants) - 1)
		if vfo.MyRole == lastParticipant || vfo.MyRole == payments.PAYER_INDEX {
			err = e.registerPaymentChannel(vfo)
			if err != nil {
				return failedEngineEvent, fmt.Errorf("could not register channel with payment/receipt manager: %w", err)
			}
		}

		if err != nil {
			return failedEngineEvent, fmt.Errorf("could not register channel with payment/receipt manager: %w", err)
		}
		return e.attemptProgress(&vfo)

	case virtualdefund.ObjectiveRequest:
		minAmount := big.NewInt(0)
		if e.vm.ChannelRegistered(request.ChannelId) {
			paid, err := e.vm.Paid(request.ChannelId)
			if err != nil {
				return failedEngineEvent, fmt.Errorf("handleAPIEvent: Could not create virtualdefund objective for %+v: %w", request, err)
			}
			minAmount = paid
		}
		vdfo, err := virtualdefund.NewObjective(request, true, myAddress, minAmount, e.store.GetChannelById, e.store.GetConsensusChannel)
		if err != nil {
			return failedEngineEvent, fmt.Errorf("handleAPIEvent: Could not create virtualdefund objective for %+v: %w", request, err)
		}
		return e.attemptProgress(&vdfo)

	case directfund.ObjectiveRequest:
		dfo, err := directfund.NewObjective(request, true, myAddress, chainId, e.store.GetChannelsByParticipant, e.store.GetConsensusChannel)
		if err != nil {
			return failedEngineEvent, fmt.Errorf("handleAPIEvent: Could not create directfund objective for %+v: %w", request, err)
		}
		return e.attemptProgress(&dfo)

	case directdefund.ObjectiveRequest:
		ddfo, err := directdefund.NewObjective(request, true, e.store.GetConsensusChannelById)
		if err != nil {
			return failedEngineEvent, fmt.Errorf("handleAPIEvent: Could not create directdefund objective for %+v: %w", request, err)
		}
		// If ddfo creation was successful, destroy the consensus channel to prevent it being used (a Channel will now take over governance)
		err = e.store.DestroyConsensusChannel(request.ChannelId)
		if err != nil {
			return failedEngineEvent, fmt.Errorf("handleAPIEvent: Could not destroy consensus channel for %+v: %w", request, err)
		}
		return e.attemptProgress(&ddfo)

	default:
		return failedEngineEvent, fmt.Errorf("handleAPIEvent: Unknown objective type %T", request)
	}
}

// handlePaymentRequest handles an PaymentRequest (triggered by a client API call).
// It prepares and dispatches a payment message to the counterparty.
func (e *Engine) handlePaymentRequest(request PaymentRequest) (EngineEvent, error) {
	ee := EngineEvent{}
	if (request == PaymentRequest{}) {
		return ee, fmt.Errorf("handleAPIEvent: Empty payment request")
	}
	cId := request.ChannelId
	voucher, err := e.vm.Pay(
		cId,
		request.Amount,
		*e.store.GetChannelSecretKey())
	if err != nil {
		return ee, fmt.Errorf("handleAPIEvent: Error making payment: %w", err)
	}
	c, ok := e.store.GetChannelById(cId)
	if !ok {
		return ee, fmt.Errorf("handleAPIEvent: Could not get channel from the store %s", cId)
	}
	payer, payee := payments.GetPayer(c.Participants), payments.GetPayee(c.Participants)
	if payer != *e.store.GetAddress() {
		return ee, fmt.Errorf("handleAPIEvent: Not the sender in channel %s", cId)
	}
	info, err := query.GetPaymentChannelInfo(cId, e.store, e.vm)
	if err != nil {
		return ee, fmt.Errorf("handleAPIEvent: Error querying channel info: %w", err)
	}
	ee.PaymentChannelUpdates = append(ee.PaymentChannelUpdates, info)

	se := protocols.SideEffects{MessagesToSend: protocols.CreateVoucherMessage(voucher, payee)}
	return ee, e.executeSideEffects(se)
}

// sendMessages sends out the messages and records the metrics.
func (e *Engine) sendMessages(msgs []protocols.Message) {
	defer e.metrics.RecordFunctionDuration()()

	for _, message := range msgs {
		message.From = *e.store.GetAddress()
		e.logMessage(message, Outgoing)
		e.recordMessageMetrics(message)
		e.msg.Send(message)
	}
	e.wg.Done()
}

// executeSideEffects executes the SideEffects declared by cranking an Objective or handling a payment request.
func (e *Engine) executeSideEffects(sideEffects protocols.SideEffects) error {
	defer e.metrics.RecordFunctionDuration()()

	e.wg.Add(1)
	// Send messages in a go routine so that we don't block on message delivery
	go e.sendMessages(sideEffects.MessagesToSend)

	for _, tx := range sideEffects.TransactionsToSubmit {
		e.logger.Printf("Sending chain transaction for channel %s", tx.ChannelId())
		err := e.chain.SendTransaction(tx)
		if err != nil {
			return err
		}
	}
	for _, proposal := range sideEffects.ProposalsToProcess {
		e.fromLedger <- proposal
	}
	return nil
}

// attemptProgress takes a "live" objective in memory and performs the following actions:
//
//  1. It pulls the secret key from the store
//  2. It cranks the objective with that key
//  3. It commits the cranked objective to the store
//  4. It executes any side effects that were declared during cranking
//  5. It updates progress metadata in the store
func (e *Engine) attemptProgress(objective protocols.Objective) (outgoing EngineEvent, err error) {
	defer e.metrics.RecordFunctionDuration()()

	secretKey := e.store.GetChannelSecretKey()
	var crankedObjective protocols.Objective
	var sideEffects protocols.SideEffects
	var waitingFor protocols.WaitingFor

	crankedObjective, sideEffects, waitingFor, err = objective.Crank(secretKey)
	if err != nil {
		return
	}

	err = e.store.SetObjective(crankedObjective)

	notifEvents, err := e.generateNotifications(crankedObjective)
	if err != nil {
		return EngineEvent{}, err
	}
	outgoing.Merge(notifEvents)

	e.logger.Printf("Objective %s is %s", objective.Id(), waitingFor)

	// If our protocol is waiting for nothing then we know the objective is complete
	// TODO: If attemptProgress is called on a completed objective CompletedObjectives would include that objective id
	// Probably should have a better check that only adds it to CompletedObjectives if it was completed in this crank
	if waitingFor == "WaitingForNothing" {
		outgoing.CompletedObjectives = append(outgoing.CompletedObjectives, crankedObjective)
		err = e.store.ReleaseChannelFromOwnership(crankedObjective.OwnsChannel())
		if err != nil {
			return
		}
		err = e.spawnConsensusChannelIfDirectFundObjective(crankedObjective) // Here we assume that every directfund.Objective is for a ledger channel.
		if err != nil {
			return
		}
	}
	err = e.executeSideEffects(sideEffects)
	return
}

// generateNotifications takes an objective and constructs notifications for any related channels for that objective.
func (e *Engine) generateNotifications(o protocols.Objective) (EngineEvent, error) {
	outgoing := EngineEvent{}

	for _, rel := range o.Related() {
		switch c := rel.(type) {
		case *channel.VirtualChannel:
			var paid, remaining *big.Int

			if !c.FinalCompleted() {
				// If the channel is open, we inspect vouchers for that channel to get the future resolvable balance
				var err error
				paid, remaining, err = query.GetVoucherBalance(c.Id, e.vm)
				if err != nil {
					return outgoing, err
				}
			} else {
				// If the channel is closed, vouchers have already been resolved.
				// Note that when virtual defunding, this information may in fact be more up to date than
				// the voucher balance due to a race condition https://github.com/statechannels/go-nitro/issues/1323
				paid, remaining = c.GetPaidAndRemaining()
			}
			info, err := query.ConstructPaymentInfo(&c.Channel, paid, remaining)
			if err != nil {
				return outgoing, err
			}
			outgoing.PaymentChannelUpdates = append(outgoing.PaymentChannelUpdates, info)
		case *channel.Channel:
			l, err := query.ConstructLedgerInfoFromChannel(c, *e.store.GetAddress())
			if err != nil {
				return outgoing, err
			}
			outgoing.LedgerChannelUpdates = append(outgoing.LedgerChannelUpdates, l)
		case *consensus_channel.ConsensusChannel:
			l, err := query.ConstructLedgerInfoFromConsensus(c, *e.store.GetAddress())
			if err != nil {
				return outgoing, err
			}
			outgoing.LedgerChannelUpdates = append(outgoing.LedgerChannelUpdates, l)
		default:
			return outgoing, fmt.Errorf("handleNotifications: Unknown related type %T", c)
		}
	}
	return outgoing, nil
}

func (e Engine) registerPaymentChannel(vfo virtualfund.Objective) error {
	postfund := vfo.V.PostFundState()
	startingBalance := big.NewInt(0)
	// TODO: Assumes one asset for now
	startingBalance.Set(postfund.Outcome[0].Allocations[0].Amount)

	return e.vm.Register(vfo.V.Id, payments.GetPayer(postfund.Participants), payments.GetPayee(postfund.Participants), startingBalance)
}

// spawnConsensusChannelIfDirectFundObjective will attempt to create and store a ConsensusChannel derived from the supplied Objective if it is a directfund.Objective.
//
// The associated Channel will remain in the store.
func (e Engine) spawnConsensusChannelIfDirectFundObjective(crankedObjective protocols.Objective) error {
	defer e.metrics.RecordFunctionDuration()()

	if dfo, isDfo := crankedObjective.(*directfund.Objective); isDfo {
		c, err := dfo.CreateConsensusChannel()
		if err != nil {
			return fmt.Errorf("could not create consensus channel for objective %s: %w", crankedObjective.Id(), err)
		}
		err = e.store.SetConsensusChannel(c)
		if err != nil {
			return fmt.Errorf("could not store consensus channel for objective %s: %w", crankedObjective.Id(), err)
		}
		// Destroy the channel since the consensus channel takes over governance:
		err = e.store.DestroyChannel(c.Id)
		if err != nil {
			return fmt.Errorf("could not destroy consensus channel for objective %s: %w", crankedObjective.Id(), err)
		}
	}
	return nil
}

// getOrCreateObjective retrieves the objective from the store.
// If the objective does not exist, it creates the objective using the supplied payload and stores it in the store
func (e *Engine) getOrCreateObjective(p protocols.ObjectivePayload) (protocols.Objective, error) {
	defer e.metrics.RecordFunctionDuration()()
	id := p.ObjectiveId
	objective, err := e.store.GetObjectiveById(id)

	if err == nil {
		return objective, nil
	} else if errors.Is(err, store.ErrNoSuchObjective) {

		newObj, err := e.constructObjectiveFromMessage(id, p)
		if err != nil {
			return nil, fmt.Errorf("error constructing objective from message: %w", err)
		}
		e.metrics.RecordObjectiveStarted(newObj.Id())
		err = e.store.SetObjective(newObj)
		if err != nil {
			return nil, fmt.Errorf("error setting objective in store: %w", err)
		}
		e.logger.Printf("Created new objective from message %s", newObj.Id())
		return newObj, nil

	} else {
		return nil, &ErrGetObjective{err, id}
	}
}

// constructObjectiveFromMessage Constructs a new objective (of the appropriate concrete type) from the supplied payload.
func (e *Engine) constructObjectiveFromMessage(id protocols.ObjectiveId, p protocols.ObjectivePayload) (protocols.Objective, error) {
	e.logger.Printf("Constructing objective %s from message", id)
	defer e.metrics.RecordFunctionDuration()()

	switch {
	case directfund.IsDirectFundObjective(id):

		dfo, err := directfund.ConstructFromPayload(false, p, *e.store.GetAddress())
		return &dfo, err
	case virtualfund.IsVirtualFundObjective(id):
		vfo, err := virtualfund.ConstructObjectiveFromPayload(p, false, *e.store.GetAddress(), e.store.GetConsensusChannel)
		if err != nil {
			return &virtualfund.Objective{}, fromMsgErr(id, err)
		}
		err = e.registerPaymentChannel(vfo)
		if err != nil {
			return &virtualfund.Objective{}, fmt.Errorf("could not register channel with payment/receipt manager.\n\ttarget channel: %s\n\terr: %w", id, err)
		}
		return &vfo, nil
	case virtualdefund.IsVirtualDefundObjective(id):
		vId, err := virtualdefund.GetVirtualChannelFromObjectiveId(id)
		if err != nil {
			return &virtualdefund.Objective{}, fmt.Errorf("could not determine virtual channel id from objective %s: %w", id, err)
		}
		minAmount := big.NewInt(0)
		if e.vm.ChannelRegistered(vId) {
			paid, err := e.vm.Paid(vId)
			if err != nil {
				return &virtualdefund.Objective{}, fmt.Errorf("could not determine virtual channel id from objective %s: %w", id, err)
			}
			minAmount = paid
		}

		vdfo, err := virtualdefund.ConstructObjectiveFromPayload(p, false, *e.store.GetAddress(), e.store.GetChannelById, e.store.GetConsensusChannel, minAmount)
		if err != nil {
			return &virtualfund.Objective{}, fromMsgErr(id, err)
		}
		return &vdfo, nil
	case directdefund.IsDirectDefundObjective(id):
		ddfo, err := directdefund.ConstructObjectiveFromPayload(p, false, e.store.GetConsensusChannelById)
		if err != nil {
			return &directdefund.Objective{}, fromMsgErr(id, err)
		}
		return &ddfo, nil

	default:
		return &directfund.Objective{}, errors.New("cannot handle unimplemented objective type")
	}
}

// fromMsgErr wraps errors from objective construction functions and
// returns an error bundled with the objectiveID
func fromMsgErr(id protocols.ObjectiveId, err error) error {
	return fmt.Errorf("could not create objective from message.\n\ttarget objective: %s\n\terr: %w", id, err)
}

// getProposalObjectiveId returns the objectiveId for a proposal.
func getProposalObjectiveId(p consensus_channel.Proposal) protocols.ObjectiveId {
	switch p.Type() {
	case consensus_channel.AddProposal:
		{
			const prefix = virtualfund.ObjectivePrefix
			channelId := p.ToAdd.Guarantee.Target().String()
			return protocols.ObjectiveId(prefix + channelId)

		}
	case consensus_channel.RemoveProposal:
		{
			const prefix = virtualdefund.ObjectivePrefix
			channelId := p.ToRemove.Target.String()
			return protocols.ObjectiveId(prefix + channelId)

		}
	default:
		{
			panic("invalid proposal type")
		}
	}
}

// GetConsensusAppAddress returns the address of a deployed ConsensusApp (for ledger channels)
func (e *Engine) GetConsensusAppAddress() types.Address {
	return e.chain.GetConsensusAppAddress()
}

// GetVirtualPaymentAppAddress returns the address of a deployed VirtualPaymentApp
func (e *Engine) GetVirtualPaymentAppAddress() types.Address {
	return e.chain.GetVirtualPaymentAppAddress()
}

type messageDirection string

const (
	Incoming messageDirection = "Incoming"
	Outgoing messageDirection = "Outgoing"
)

// logMessage logs a message to the engine's logger
func (e *Engine) logMessage(msg protocols.Message, direction messageDirection) {
	if direction == Incoming {
		e.logger.Trace().EmbedObject(msg.Summarize()).Msg("Received message")
	} else {
		e.logger.Trace().EmbedObject(msg.Summarize()).Msg("Sending message")
	}
}

// recordMessageMetrics records metrics for a message
func (e *Engine) recordMessageMetrics(message protocols.Message) {
	e.metrics.RecordQueueLength(fmt.Sprintf("msg_proposal_count,sender=%s,receiver=%s", e.store.GetAddress(), message.To), len(message.LedgerProposals))
	e.metrics.RecordQueueLength(fmt.Sprintf("msg_payment_count,sender=%s,receiver=%s", e.store.GetAddress(), message.To), len(message.Payments))
	e.metrics.RecordQueueLength(fmt.Sprintf("msg_payload_count,sender=%s,receiver=%s", e.store.GetAddress(), message.To), len(message.ObjectivePayloads))

	totalPayloadsSize := 0
	for _, p := range message.ObjectivePayloads {
		totalPayloadsSize += len(p.PayloadData)
	}
	raw, _ := message.Serialize()
	e.metrics.RecordQueueLength(fmt.Sprintf("msg_payload_size,sender=%s,receiver=%s", e.store.GetAddress(), message.To), totalPayloadsSize)
	e.metrics.RecordQueueLength(fmt.Sprintf("msg_size,sender=%s,receiver=%s", e.store.GetAddress(), message.To), len(raw))
}

func (e *Engine) checkError(err error) {
	if err != nil {
		e.logger.Err(err).Msgf("%s, error in run loop", e.store.GetAddress())

		for _, nonFatalError := range nonFatalErrors {
			if errors.Is(err, nonFatalError) {
				return
			}
		}

		e.logger.Panic().Msg(err.Error())
	}
}
