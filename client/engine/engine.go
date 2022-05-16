// Package engine contains the types and imperative code for the business logic of a go-nitro Client.
package engine // import "github.com/statechannels/go-nitro/client/engine"

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
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

	// outbound go channels
	toMsg   chan<- protocols.Message
	toChain chan<- protocols.ChainTransaction

	results chan engineResult

	toApi chan ObjectiveChangeEvent

	store store.Store // A Store for persisting and restoring important data

	logger *log.Logger
}

// APIEvent is an internal representation of an API call
type APIEvent struct {
	ObjectiveToSpawn   protocols.ObjectiveRequest
	ObjectiveToReject  protocols.ObjectiveId
	ObjectiveToApprove protocols.ObjectiveId
}

// ObjectiveChangeEvent is a struct that contains a list of changes caused by handling a message/chain event/api event
type ObjectiveChangeEvent struct {
	// These are objectives that are now completed
	CompletedObjectives []protocols.Objective
}

type CompletedObjectiveEvent struct {
	Id protocols.ObjectiveId
}

// Response is the return type that asynchronous API calls "resolve to". Such a call returns a go channel of type Response.
type Response struct{}

// NewEngine is the constructor for an Engine
func New(msg messageservice.MessageService, chain chainservice.ChainService, store store.Store, logDestination io.Writer) Engine {
	e := Engine{}

	e.store = store

	// bind to inbound chans
	e.FromAPI = make(chan APIEvent)
	e.fromChain = chain.Out()
	e.fromMsg = msg.Out()

	e.toApi = make(chan ObjectiveChangeEvent, 100)
	// bind to outbound chans
	e.toChain = chain.In()
	e.toMsg = msg.In()

	e.results = make(chan engineResult)

	// initialize a Logger
	logPrefix := e.store.GetAddress().String()[0:8] + ": "
	e.logger = log.New(logDestination, logPrefix, log.Lmicroseconds|log.Lshortfile)

	e.logger.Println("Constructed Engine")

	return e
}

func (e *Engine) ToApi() <-chan ObjectiveChangeEvent {
	return e.toApi
}

type engineResult struct {
	res ObjectiveChangeEvent
	err error
}

// Run kicks of an infinite loop that waits for communications on the supplied channels, and handles them accordingly
func (e *Engine) Run() {
	for {
		select {
		case apiEvent := <-e.FromAPI:
			go func() { e.results <- e.handleAPIEvent(apiEvent) }()

		case chainEvent := <-e.fromChain:
			go func() { e.results <- e.handleChainEvent(chainEvent) }()

		case message := <-e.fromMsg:
			go func() { e.results <- e.handleMessage(message) }()

		case proposal := <-e.fromLedger:
			go func() { e.results <- e.handleProposal(proposal) }()

		case handled := <-e.results:
			if handled.err != nil {
				e.logger.Panic(handled.err)
			}
			if len(handled.res.CompletedObjectives) > 0 {
				for _, obj := range handled.res.CompletedObjectives {
					e.logger.Printf("Objective %s is complete & returned to API", obj.Id())
				}
				e.toApi <- handled.res
			}
		}
	}
}

// handleProposal handles a Proposal returned to the engine from
// a running ledger channel by pulling its corresponding objective
// from the store and attempting progress.
func (e *Engine) handleProposal(proposal consensus_channel.Proposal) engineResult {
	id := getProposalObjectiveId(proposal)
	obj, err := e.store.GetObjectiveById(id)
	if err != nil {
		return engineResult{ObjectiveChangeEvent{}, err}
	}
	return e.attemptProgress(obj)
}

// handleMessage handles a Message from a peer go-nitro Wallet.
// It:
//  - reads an objective from the store,
//  - generates an updated objective,
//  - attempts progress on the target Objective,
//  - attempts progress on related objectives which may have become unblocked.
func (e *Engine) handleMessage(message protocols.Message) engineResult {

	e.logger.Printf("Handling inbound message %+v", protocols.SummarizeMessage(message))
	allCompleted := ObjectiveChangeEvent{}

	for _, entry := range message.SignedStates() {

		objective, err := e.getOrCreateObjective(entry.ObjectiveId, entry.Payload)
		if err != nil {
			return engineResult{ObjectiveChangeEvent{}, err}
		}
		if objective.GetStatus() == protocols.Completed {
			e.logger.Printf("Ignoring payload for complected objective  %s", objective.Id())
			continue
		}
		event := protocols.ObjectiveEvent{
			ObjectiveId:    entry.ObjectiveId,
			SignedProposal: consensus_channel.SignedProposal{},
			SignedState:    entry.Payload,
		}
		updatedObjective, err := objective.Update(event)
		if err != nil {
			return engineResult{ObjectiveChangeEvent{}, err}
		}

		result := e.attemptProgress(updatedObjective)
		if result.err != nil {
			return engineResult{ObjectiveChangeEvent{}, result.err}
		}
		allCompleted.CompletedObjectives = append(allCompleted.CompletedObjectives, result.res.CompletedObjectives...)

		if err != nil {
			return engineResult{ObjectiveChangeEvent{}, err}
		}

	}

	for _, entry := range message.SignedProposals() {
		e.logger.Printf("handling proposal %+v", protocols.SummarizeProposal(entry.ObjectiveId, entry.Payload))
		objective, err := e.store.GetObjectiveById(entry.ObjectiveId)
		if err != nil {
			return engineResult{ObjectiveChangeEvent{}, err}
		}
		if objective.GetStatus() == protocols.Completed {
			e.logger.Printf("Ignoring payload for complected objective  %s", objective.Id())
			continue
		}

		event := protocols.ObjectiveEvent{
			ObjectiveId:    entry.ObjectiveId,
			SignedProposal: entry.Payload,
			SignedState:    state.SignedState{},
		}
		updatedObjective, err := objective.Update(event)
		if err != nil {
			return engineResult{ObjectiveChangeEvent{}, err}
		}

		progressResult := e.attemptProgress(updatedObjective)
		if err != nil {
			return engineResult{ObjectiveChangeEvent{}, progressResult.err}
		}

		allCompleted.CompletedObjectives = append(allCompleted.CompletedObjectives, progressResult.res.CompletedObjectives...)

		if err != nil {
			return engineResult{ObjectiveChangeEvent{}, err}
		}

	}
	return engineResult{allCompleted, nil}

}

// handleChainEvent handles a Chain Event from the blockchain.
// It:
//  - reads an objective from the store,
//  - generates an updated objective, and
//  - attempts progress.
func (e *Engine) handleChainEvent(chainEvent chainservice.Event) engineResult {
	e.logger.Printf("handling chain event %v", chainEvent)
	objective, ok := e.store.GetObjectiveByChannelId(chainEvent.ChannelID())
	if !ok {
		// TODO: Right now the chain service returns chain events for ALL channels even those we aren't involved in
		// for now we can ignore channels we aren't involved in
		// in the future the chain service should allow us to register for specific channels
		return engineResult{ObjectiveChangeEvent{}, nil}
	}

	eventHandler, ok := objective.(chainservice.ChainEventHandler)
	if !ok {
		return engineResult{ObjectiveChangeEvent{}, &ErrUnhandledChainEvent{event: chainEvent, objective: objective, reason: "objective does not handle chain events"}}
	}
	updatedEventHandler, err := eventHandler.UpdateWithChainEvent(chainEvent)
	if err != nil {
		return engineResult{ObjectiveChangeEvent{}, err}
	}
	return e.attemptProgress(updatedEventHandler)
}

// handleAPIEvent handles an API Event (triggered by a client API call).
// It will attempt to perform all of the following:
//  - Spawn a new, approved objective (if not null)
//  - Reject an existing objective (if not null)
//  - Approve an existing objective (if not null)
func (e *Engine) handleAPIEvent(apiEvent APIEvent) engineResult {
	if apiEvent.ObjectiveToSpawn != nil {

		switch request := (apiEvent.ObjectiveToSpawn).(type) {

		case virtualfund.ObjectiveRequest:
			vfo, err := virtualfund.NewObjective(request, true, *e.store.GetAddress(), e.store.GetConsensusChannel)
			if err != nil {
				return engineResult{ObjectiveChangeEvent{}, fmt.Errorf("handleAPIEvent: Could not create objective for %+v: %w", request, err)}
			}
			return e.attemptProgress(&vfo)

		case virtualdefund.ObjectiveRequest:
			vdfo, err := virtualdefund.NewObjective(request, true, *e.store.GetAddress(), e.store.GetChannelById, e.store.GetConsensusChannel)
			if err != nil {
				return engineResult{ObjectiveChangeEvent{}, fmt.Errorf("handleAPIEvent: Could not create objective for %+v: %w", request, err)}
			}
			return e.attemptProgress(&vdfo)

		case directfund.ObjectiveRequest:
			dfo, err := directfund.NewObjective(request, true, *e.store.GetAddress(), e.store.GetChannelsByParticipant, e.store.GetConsensusChannel)
			if err != nil {
				return engineResult{ObjectiveChangeEvent{}, fmt.Errorf("handleAPIEvent: Could not create objective for %+v: %w", request, err)}
			}
			return e.attemptProgress(&dfo)

		case directdefund.ObjectiveRequest:
			ddfo, err := directdefund.NewObjective(request, true, e.store.GetConsensusChannelById)
			if err != nil {
				return engineResult{ObjectiveChangeEvent{}, fmt.Errorf("handleAPIEvent: Could not create objective for %+v: %w", request, err)}
			}
			// If ddfo creation was successful, destroy the consensus channel to prevent it being used (a Channel will now take over governance)
			e.store.DestroyConsensusChannel(request.ChannelId)
			return e.attemptProgress(&ddfo)

		default:
			return engineResult{ObjectiveChangeEvent{}, fmt.Errorf("handleAPIEvent: Unknown objective type %T", request)}
		}

	}

	if apiEvent.ObjectiveToReject != `` {
		objective, err := e.store.GetObjectiveById(apiEvent.ObjectiveToReject)
		if err != nil {
			return engineResult{ObjectiveChangeEvent{}, err}
		}
		updatedProtocol := objective.Reject()
		err = e.store.SetObjective(updatedProtocol)
		return engineResult{ObjectiveChangeEvent{}, err}
	}
	if apiEvent.ObjectiveToApprove != `` {
		objective, err := e.store.GetObjectiveById(apiEvent.ObjectiveToReject)
		if err != nil {
			return engineResult{ObjectiveChangeEvent{}, err}
		}
		updatedObjective := objective.Approve()
		return e.attemptProgress(updatedObjective)
	}
	return engineResult{ObjectiveChangeEvent{}, nil}

}

// executeSideEffects executes the SideEffects declared by cranking an Objective
func (e *Engine) executeSideEffects(sideEffects protocols.SideEffects) {
	for _, message := range sideEffects.MessagesToSend {
		e.logger.Printf("Sending message %+v", protocols.SummarizeMessage(message))
		e.toMsg <- message
	}
	for _, tx := range sideEffects.TransactionsToSubmit {
		e.logger.Printf("Sending chain transaction for channel %s", tx.ChannelId)
		e.toChain <- tx
	}
	for _, proposal := range sideEffects.ProposalsToProcess {
		e.fromLedger <- proposal
	}
}

// attemptProgress takes a "live" objective in memory and performs the following actions:
//
// 	1. It pulls the secret key from the store
// 	2. It cranks the objective with that key
// 	3. It commits the cranked objective to the store
// 	4. It executes any side effects that were declared during cranking
// 	5. It updates progress metadata in the store
func (e *Engine) attemptProgress(objective protocols.Objective) engineResult {

	secretKey := e.store.GetChannelSecretKey()

	crankedObjective, sideEffects, waitingFor, err := objective.Crank(secretKey)

	if err != nil {
		return engineResult{err: err}
	}

	err = e.store.SetObjective(crankedObjective)

	if err != nil {
		return engineResult{err: err}
	}

	e.logger.Printf("Objective %s is %s", objective.Id(), waitingFor)

	// If our protocol is waiting for nothing then we know the objective is complete
	// TODO: If attemptProgress is called on a completed objective CompletedObjectives would include that objective id
	// Probably should have a better check that only adds it to CompletedObjectives if it was completed in this crank
	if waitingFor == "WaitingForNothing" {
		var res ObjectiveChangeEvent
		res.CompletedObjectives = append(res.CompletedObjectives, crankedObjective)
		e.store.ReleaseChannelFromOwnership(crankedObjective.OwnsChannel())
		err = e.spawnConsensusChannelIfDirectFundObjective(crankedObjective) // Here we assume that every directfund.Objective is for a ledger channel.
		if err != nil {
			return engineResult{res, err}
		}
	}
	e.executeSideEffects(sideEffects)

	return engineResult{err: err}
}

// spawnConsensusChannelIfDirectFundObjective will attempt to create and store a ConsensusChannel derived from the supplied Objective if it is a directfund.Objective.
//
// The associated Channel will remain in the store.
func (e Engine) spawnConsensusChannelIfDirectFundObjective(crankedObjective protocols.Objective) error {
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

// getOrCreateObjective creates the objective if the supplied message is a proposal. Otherwise, it attempts to get the objective from the store.
func (e *Engine) getOrCreateObjective(id protocols.ObjectiveId, ss state.SignedState) (protocols.Objective, error) {

	objective, err := e.store.GetObjectiveById(id)

	if err == nil {
		return objective, nil
	} else if errors.Is(err, store.ErrNoSuchObjective) {

		newObj, err := e.constructObjectiveFromMessage(id, ss)
		if err != nil {
			return nil, fmt.Errorf("error constructing objective from message: %w", err)
		}
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

// constructObjectiveFromMessage Constructs a new objective (of the appropriate concrete type) from the supplied message.
func (e *Engine) constructObjectiveFromMessage(id protocols.ObjectiveId, ss state.SignedState) (protocols.Objective, error) {

	switch {
	case directfund.IsDirectFundObjective(id):
		dfo, err := directfund.ConstructFromState(true, ss.State(), *e.store.GetAddress())

		return &dfo, err
	case virtualfund.IsVirtualFundObjective(id):
		vfo, err := virtualfund.ConstructObjectiveFromState(ss.State(), *e.store.GetAddress(), e.store.GetConsensusChannel)
		if err != nil {
			return &virtualfund.Objective{}, fmt.Errorf("could not create virtual fund objective from message: %w", err)
		}
		return &vfo, nil
	case virtualdefund.IsVirtualDefundObjective(id):
		vdfo, err := virtualdefund.ConstructObjectiveFromState(ss.State(), *e.store.GetAddress(), e.store.GetChannelById, e.store.GetConsensusChannel)
		if err != nil {
			return &virtualfund.Objective{}, fmt.Errorf("could not create virtual fund objective from message: %w", err)
		}
		return &vdfo, nil
	case directdefund.IsDirectDefundObjective(id):
		ddfo, err := directdefund.ConstructObjectiveFromState(ss.State(), e.store.GetConsensusChannelById)
		if err != nil {
			return &directdefund.Objective{}, fmt.Errorf("could not create direct defund objective from message: %w", err)
		}
		// If ddfo creation was successful, destroy the consensus channel to prevent it being used (a Channel will now take over governance)
		e.store.DestroyConsensusChannel(ddfo.C.Id)
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
