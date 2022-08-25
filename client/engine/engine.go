// Package engine contains the types and imperative code for the business logic of a go-nitro Client.
package engine // import "github.com/statechannels/go-nitro/client/engine"

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
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

// Engine is the imperative part of the core business logic of a go-nitro Client
type Engine struct {
	// inbound go channels
	FromAPI    chan APIEvent // This one is exported so that the Client can send API calls
	fromChain  <-chan chainservice.Event
	fromMsg    <-chan protocols.Message
	fromLedger chan consensus_channel.Proposal

	toApi chan EngineEvent

	msg   messageservice.MessageService
	chain chainservice.ChainService

	store       store.Store // A Store for persisting and restoring important data
	policymaker PolicyMaker // A PolicyMaker decides whether to approve or reject objectives

	logger *log.Logger

	metrics *MetricsRecorder

	vm *payments.VoucherManager
}

// PaymentRequest represents a request from the API to make a payment using a channel
type PaymentRequest struct {
	ChannelId types.Destination
	Amount    *big.Int
}

// APIEvent is an internal representation of an API call
type APIEvent struct {
	ObjectiveToSpawn protocols.ObjectiveRequest
	PaymentToMake    PaymentRequest
}

// EngineEvent is a struct that contains a list of changes caused by handling a message/chain event/api event
type EngineEvent struct {
	// These are objectives that are now completed
	CompletedObjectives []protocols.Objective
	// These are objectives that have failed
	FailedObjectives []protocols.ObjectiveId
	// ReceivedVouchers are vouchers we've received from other participants
	ReceivedVouchers []payments.Voucher
}

// Merge merges other into e
func (e *EngineEvent) Merge(other *EngineEvent) {
	e.CompletedObjectives = append(e.CompletedObjectives, other.CompletedObjectives...)
	e.FailedObjectives = append(e.FailedObjectives, other.FailedObjectives...)
	e.ReceivedVouchers = append(e.ReceivedVouchers, other.ReceivedVouchers...)
}

type CompletedObjectiveEvent struct {
	Id protocols.ObjectiveId
}

// Response is the return type that asynchronous API calls "resolve to". Such a call returns a go channel of type Response.
type Response struct{}

// NewEngine is the constructor for an Engine
func New(msg messageservice.MessageService, chain chainservice.ChainService, store store.Store, logDestination io.Writer, policymaker PolicyMaker, metricsApi MetricsApi) Engine {
	e := Engine{}

	e.store = store

	// bind to inbound chans
	e.FromAPI = make(chan APIEvent)
	e.fromChain = chain.EventFeed()
	e.fromMsg = msg.Out()

	e.chain = chain
	e.msg = msg

	e.toApi = make(chan EngineEvent, 100)

	// initialize a Logger
	logPrefix := e.store.GetAddress().String()[0:8] + ": "
	e.logger = log.New(logDestination, logPrefix, log.Lmicroseconds|log.Lshortfile)

	e.policymaker = policymaker

	e.vm = payments.NewVoucherManager(*store.GetAddress())

	e.logger.Println("Constructed Engine")

	if metricsApi == nil {
		metricsApi = &NoOpMetrics{}
	}
	e.metrics = NewMetricsRecorder(*e.store.GetAddress(), metricsApi)
	return e
}

func (e *Engine) ToApi() <-chan EngineEvent {
	return e.toApi
}

// Run kicks of an infinite loop that waits for communications on the supplied channels, and handles them accordingly
func (e *Engine) Run() {
	for {
		var res EngineEvent
		var err error

		e.metrics.RecordQueueLength("api_events_queue", len(e.FromAPI))
		e.metrics.RecordQueueLength("chain_events_queue", len(e.fromChain))
		e.metrics.RecordQueueLength("messages_queue", len(e.fromMsg))
		e.metrics.RecordQueueLength("proposal_queue", len(e.fromLedger))

		select {
		case apiEvent := <-e.FromAPI:
			res, err = e.handleAPIEvent(apiEvent)

			if errors.Is(err, directdefund.ErrNotEmpty) {
				// communicate failure to client & swallow error
				e.toApi <- res
				err = nil
			}

		case chainEvent := <-e.fromChain:
			res, err = e.handleChainEvent(chainEvent)
		case message := <-e.fromMsg:
			res, err = e.handleMessage(message)
		case proposal := <-e.fromLedger:
			res, err = e.handleProposal(proposal)
		}

		// Handle errors
		if err != nil {
			e.logger.Panic(fmt.Errorf("%s, error in run loop: %w", e.store.GetAddress(), err))
			// TODO do not panic if in production.
			// TODO report errors back to the consuming application
		}

		// Only send out an event if there are changes
		if len(res.CompletedObjectives) > 0 || len(res.FailedObjectives) > 0 || len(res.ReceivedVouchers) > 0 {
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
	return e.attemptProgress(obj)
}

// updateObjective handles updating and cranking the objective and dispatching any side effects
func (e *Engine) updateObjective(objective protocols.Objective, event protocols.ObjectiveEvent) (EngineEvent, error) {
	engineEvent := EngineEvent{}
	if objective.GetStatus() == protocols.Unapproved {
		e.logger.Printf("Policymaker is %+v", e.policymaker)
		if e.policymaker.ShouldApprove(objective) {
			objective = objective.Approve()

			ddfo, ok := objective.(*directdefund.Objective)
			if ok {
				// If we just approved a direct defund objective, destroy the consensus channel to prevent it being used (a Channel will now take over governance)
				e.store.DestroyConsensusChannel(ddfo.C.Id)
			}
		} else {
			objective, sideEffects := objective.Reject()
			err := e.store.SetObjective(objective)
			if err != nil {
				return EngineEvent{}, err
			}

			engineEvent.CompletedObjectives = append(engineEvent.CompletedObjectives, objective)
			err = e.executeSideEffects(sideEffects)
			// An error would mean we failed to send a message. But the objective is still "completed".
			// So, we should return allCompleted even if there was an error.
			return engineEvent, err
		}
	}

	if objective.GetStatus() == protocols.Completed {
		e.logger.Printf("Ignoring payload for complected objective  %s", objective.Id())
		return engineEvent, nil
	}
	if objective.GetStatus() == protocols.Rejected {
		e.logger.Printf("Ignoring payload for rejected objective  %s", objective.Id())
		return engineEvent, nil
	}

	updatedObjective, err := objective.Update(event)
	if err != nil {
		return EngineEvent{}, err
	}

	progressEvent, err := e.attemptProgress(updatedObjective)
	if err != nil {
		return EngineEvent{}, err
	}
	engineEvent.CompletedObjectives = append(engineEvent.CompletedObjectives, progressEvent.CompletedObjectives...)

	if err != nil {
		return EngineEvent{}, err
	}

	return engineEvent, nil
}

// handleMessage handles a Message from a peer go-nitro Wallet.
// It:
//   - reads an objective from the store,
//   - generates an updated objective,
//   - attempts progress on the target Objective,
//   - attempts progress on related objectives which may have become unblocked.
func (e *Engine) handleMessage(message protocols.Message) (EngineEvent, error) {
	defer e.metrics.RecordFunctionDuration()()

	e.logger.Printf("Handling inbound message %+v", protocols.SummarizeMessage(message))
	engineEvent := EngineEvent{}

	for _, entry := range message.SignedStates() {

		objective, err := e.getOrCreateObjective(entry.ObjectiveId, entry.Payload)
		if err != nil {
			return EngineEvent{}, err
		}

		event := protocols.ObjectiveEvent{
			ObjectiveId:    entry.ObjectiveId,
			SignedProposal: consensus_channel.SignedProposal{},
			SignedState:    entry.Payload,
		}
		stateEvents, err := e.updateObjective(objective, event)
		if err != nil {
			return EngineEvent{}, err
		}
		engineEvent.Merge(&stateEvents)

	}

	for _, entry := range message.SignedProposals() {
		e.logger.Printf("handling proposal %+v", protocols.SummarizeProposal(entry.ObjectiveId, entry.Payload))
		objective, err := e.store.GetObjectiveById(entry.ObjectiveId)
		if err != nil {
			return EngineEvent{}, err
		}
		event := protocols.ObjectiveEvent{
			ObjectiveId:    entry.ObjectiveId,
			SignedProposal: entry.Payload,
			SignedState:    state.SignedState{},
		}
		proposalEvents, err := e.updateObjective(objective, event)

		if err != nil {
			return EngineEvent{}, err
		}
		engineEvent.Merge(&proposalEvents)

	}

	for _, entry := range message.RejectedObjectives() {
		objective, err := e.store.GetObjectiveById(entry.ObjectiveId)

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

		engineEvent.CompletedObjectives = append(engineEvent.CompletedObjectives, objective)
	}

	for _, entry := range message.Vouchers() {
		// An empty objectiveId indicates a voucher sent as a payment.
		if entry.ObjectiveId == "" {

			_, err := e.vm.Receive(entry.Payload)
			if err != nil {
				return EngineEvent{}, fmt.Errorf("error accepting payment voucher: %w", err)
			}
			engineEvent.ReceivedVouchers = append(engineEvent.ReceivedVouchers, entry.Payload)

		} else {

			objective, err := e.getOrCreateObjectiveFromVoucher(entry.ObjectiveId, entry.Payload)
			if err != nil {
				return EngineEvent{}, err
			}

			event := protocols.ObjectiveEvent{
				ObjectiveId: entry.ObjectiveId,
				Voucher:     entry.Payload,
			}
			e, err := e.updateObjective(objective, event)
			if err != nil {
				return EngineEvent{}, fmt.Errorf("could not update objective: %w", err)
			}
			engineEvent.Merge(&e)
		}

	}
	return engineEvent, nil

}

// handleChainEvent handles a Chain Event from the blockchain.
// It:
//   - reads an objective from the store,
//   - generates an updated objective, and
//   - attempts progress.
func (e *Engine) handleChainEvent(chainEvent chainservice.Event) (EngineEvent, error) {
	defer e.metrics.RecordFunctionDuration()()
	e.logger.Printf("handling chain event %v", chainEvent)
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

// handleAPIEvent handles an API Event (triggered by a client API call).
// It will attempt to perform all of the following:
//   - Spawn a new, approved objective (if not null)
//   - Reject an existing objective (if not null)
//   - Approve an existing objective (if not null)
func (e *Engine) handleAPIEvent(apiEvent APIEvent) (EngineEvent, error) {
	defer e.metrics.RecordFunctionDuration()()

	if apiEvent.ObjectiveToSpawn != nil {

		switch request := (apiEvent.ObjectiveToSpawn).(type) {

		case virtualfund.ObjectiveRequest:
			e.metrics.RecordObjectiveStarted(request.Id(*e.store.GetAddress()))
			vfo, err := virtualfund.NewObjective(request, true, *e.store.GetAddress(), e.store.GetConsensusChannel)
			if err != nil {
				return EngineEvent{}, fmt.Errorf("handleAPIEvent: Could not create objective for %+v: %w", request, err)
			}

			err = e.registerPaymentChannel(vfo)
			if err != nil {
				return EngineEvent{}, fmt.Errorf("could not register channel with payment/receipt manager: %w", err)
			}
			return e.attemptProgress(&vfo)

		case virtualdefund.ObjectiveRequest:
			e.metrics.RecordObjectiveStarted(request.Id(*e.store.GetAddress()))

			voucher, err := e.vm.Voucher(request.ChannelId)

			if err != nil {
				return EngineEvent{}, fmt.Errorf("handleAPIEvent: Could not fetch voucher %+v: %w", request, err)
			}
			vdfo, err := virtualdefund.NewObjective(request, true, *e.store.GetAddress(), &voucher, e.store.GetChannelById, e.store.GetConsensusChannel)
			if err != nil {
				return EngineEvent{}, fmt.Errorf("handleAPIEvent: Could not create objective for %+v: %w", request, err)
			}
			return e.attemptProgress(&vdfo)

		case directfund.ObjectiveRequest:
			e.metrics.RecordObjectiveStarted(request.Id(*e.store.GetAddress()))
			dfo, err := directfund.NewObjective(request, true, *e.store.GetAddress(), e.store.GetChannelsByParticipant, e.store.GetConsensusChannel)
			if err != nil {
				return EngineEvent{}, fmt.Errorf("handleAPIEvent: Could not create objective for %+v: %w", request, err)
			}
			return e.attemptProgress(&dfo)

		case directdefund.ObjectiveRequest:
			e.metrics.RecordObjectiveStarted(request.Id(*e.store.GetAddress()))
			ddfo, err := directdefund.NewObjective(request, true, e.store.GetConsensusChannelById)
			if err != nil {
				return EngineEvent{FailedObjectives: []protocols.ObjectiveId{request.Id(*e.store.GetAddress())}}, fmt.Errorf("handleAPIEvent: Could not create objective for %+v: %w", request, err)
			}
			// If ddfo creation was successful, destroy the consensus channel to prevent it being used (a Channel will now take over governance)
			e.store.DestroyConsensusChannel(request.ChannelId)
			return e.attemptProgress(&ddfo)

		default:
			return EngineEvent{}, fmt.Errorf("handleAPIEvent: Unknown objective type %T", request)
		}

	}

	// TODO: Should this live in the payment manager?
	if cId := apiEvent.PaymentToMake.ChannelId; cId != (types.Destination{}) {
		voucher, err := e.vm.Pay(
			cId,
			apiEvent.PaymentToMake.Amount,
			*e.store.GetChannelSecretKey())
		if err != nil {
			return EngineEvent{}, fmt.Errorf("handleAPIEvent: Error making payment: %w", err)
		}
		c, ok := e.store.GetChannelById(cId)
		if !ok {
			return EngineEvent{}, fmt.Errorf("handleAPIEvent: Could not get channel from the store %s", cId)
		}
		payer, payee := payments.GetPayer(c.Participants), payments.GetPayee(c.Participants)
		if payer != *e.store.GetAddress() {
			return EngineEvent{}, fmt.Errorf("handleAPIEvent: Not the sender in channel %s", cId)
		}
		se := protocols.SideEffects{MessagesToSend: protocols.CreateVoucherMessage(voucher, payer, protocols.EmptyId(), payee)}
		err = e.executeSideEffects(se)
		if err != nil {
			return EngineEvent{}, fmt.Errorf("handleAPIEvent: Error sending payment voucher: %w", err)
		}

	}
	return EngineEvent{}, nil

}

// executeSideEffects executes the SideEffects declared by cranking an Objective
func (e *Engine) executeSideEffects(sideEffects protocols.SideEffects) error {
	defer e.metrics.RecordFunctionDuration()()

	for _, message := range sideEffects.MessagesToSend {

		e.logger.Printf("Sending message %+v", protocols.SummarizeMessage(message))
		e.msg.Send(message)
		e.metrics.RecordOutgoingMessage(message)
	}
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

	if err != nil {
		return
	}

	e.logger.Printf("Objective %s is %s", objective.Id(), waitingFor)

	// If our protocol is waiting for nothing then we know the objective is complete
	// TODO: If attemptProgress is called on a completed objective CompletedObjectives would include that objective id
	// Probably should have a better check that only adds it to CompletedObjectives if it was completed in this crank
	if waitingFor == "WaitingForNothing" {
		outgoing.CompletedObjectives = append(outgoing.CompletedObjectives, crankedObjective)
		e.store.ReleaseChannelFromOwnership(crankedObjective.OwnsChannel())
		err = e.spawnConsensusChannelIfDirectFundObjective(crankedObjective) // Here we assume that every directfund.Objective is for a ledger channel.
		if err != nil {
			return
		}
	}
	err = e.executeSideEffects(sideEffects)
	return
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
		e.store.DestroyChannel(c.Id)
	}
	return nil
}

// getOrCreateObjective retrieves the objective from the store. if the objective does not exist, it creates the objective using the supplied signed state, and stores it in the store
func (e *Engine) getOrCreateObjective(id protocols.ObjectiveId, ss state.SignedState) (protocols.Objective, error) {
	defer e.metrics.RecordFunctionDuration()()

	objective, err := e.store.GetObjectiveById(id)

	if err == nil {
		return objective, nil
	} else if errors.Is(err, store.ErrNoSuchObjective) {

		newObj, err := e.constructObjectiveFromMessage(id, ss)

		if err != nil {
			return nil, fmt.Errorf("error constructing objective from message: %w", err)
		}
		e.metrics.RecordObjectiveStarted(newObj.Id())
		err = e.store.SetObjective(newObj)
		if err != nil {
			return nil, fmt.Errorf("error setting objective in store: %w", err)
		}
		e.logger.Printf("Created new objective from  message %s", newObj.Id())
		return newObj, nil

	} else {
		return nil, fmt.Errorf("unexpected error getting/creating objective %s: %w", id, err)
	}
}

// getOrCreateObjectiveFromVoucher attempts to retrieve the objective from the store.
// If the objective is not found it uses the voucher and the channel store to construct a new objective.
func (e *Engine) getOrCreateObjectiveFromVoucher(id protocols.ObjectiveId, v payments.Voucher) (protocols.Objective, error) {
	defer e.metrics.RecordFunctionDuration()()
	if !virtualdefund.IsVirtualDefundObjective(id) {
		return nil, fmt.Errorf("only virtual defund objectives can be constructed from a voucher")
	}
	objective, err := e.store.GetObjectiveById(id)

	if err == nil {
		return objective, nil
	} else if errors.Is(err, store.ErrNoSuchObjective) {

		c, ok := e.store.GetChannelById(v.ChannelId)
		if !ok {
			return nil, fmt.Errorf("could not find channel for voucher channel id %s", v.ChannelId)
		}

		vdfo, err := virtualdefund.ConstructObjectiveFromVoucher(c.FixedPart, v, false, *e.store.GetAddress(), e.store.GetChannelById, e.store.GetConsensusChannel)
		if err != nil {
			return &virtualfund.Objective{}, fmt.Errorf("could not create virtual fund objective from message: %w", err)
		}

		return &vdfo, nil
	} else {

		return nil, fmt.Errorf("unexpected error getting/creating objective %s: %w", id, err)
	}

}

// constructObjectiveFromMessage Constructs a new objective (of the appropriate concrete type) from the supplied message.
func (e *Engine) constructObjectiveFromMessage(id protocols.ObjectiveId, ss state.SignedState) (protocols.Objective, error) {
	defer e.metrics.RecordFunctionDuration()()

	switch {
	case directfund.IsDirectFundObjective(id):
		dfo, err := directfund.ConstructFromState(false, ss.State(), *e.store.GetAddress())

		return &dfo, err
	case virtualfund.IsVirtualFundObjective(id):
		vfo, err := virtualfund.ConstructObjectiveFromState(ss.State(), false, *e.store.GetAddress(), e.store.GetConsensusChannel)
		if err != nil {
			return &virtualfund.Objective{}, fmt.Errorf("could not create virtual fund objective from message: %w", err)
		}
		err = e.registerPaymentChannel(vfo)
		if err != nil {
			return &virtualfund.Objective{}, fmt.Errorf("could not register channel with payment/receipt manager: %w", err)
		}
		return &vfo, nil
	case virtualdefund.IsVirtualDefundObjective(id):
		voucher, err := e.vm.Voucher(ss.ChannelId())
		if err != nil {
			return &virtualfund.Objective{}, fmt.Errorf("could not create virtual fund objective from message: %w", err)
		}

		vdfo, err := virtualdefund.ConstructObjectiveFromVoucher(ss.State().FixedPart(), voucher, false, *e.store.GetAddress(), e.store.GetChannelById, e.store.GetConsensusChannel)
		if err != nil {
			return &virtualfund.Objective{}, fmt.Errorf("could not create virtual fund objective from message: %w", err)
		}

		return &vdfo, nil

	case directdefund.IsDirectDefundObjective(id):
		ddfo, err := directdefund.ConstructObjectiveFromState(ss.State(), false, e.store.GetConsensusChannelById)
		if err != nil {
			return &directdefund.Objective{}, fmt.Errorf("could not create direct defund objective from message: %w", err)
		}
		return &ddfo, nil

	default:
		return &directfund.Objective{}, errors.New("cannot handle unimplemented objective type")
	}

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
