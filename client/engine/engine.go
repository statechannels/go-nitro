// Package engine contains the types and imperative code for the business logic of a go-nitro Client.
package engine // import "github.com/statechannels/go-nitro/client/engine"

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
)

// ErrUnhandledChainEvent is an engine error when the the engine cannot process a chain event
type ErrUnhandledChainEvent struct {
	event     chainservice.Event
	objective protocols.Objective
}

func (uce *ErrUnhandledChainEvent) Error() string {
	return fmt.Sprintf("chain event %#v could not be handled by objective %#v", uce.event, uce.objective)
}

// Engine is the imperative part of the core business logic of a go-nitro Client
type Engine struct {
	// inbound go channels
	FromAPI   chan APIEvent // This one is exported so that the Client can send API calls
	fromChain <-chan chainservice.Event
	fromMsg   <-chan protocols.Message

	// outbound go channels
	toMsg   chan<- protocols.Message
	toChain chan<- protocols.ChainTransaction

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

	// initialize a Logger
	logPrefix := e.store.GetAddress().String()[0:8] + ": "
	e.logger = log.New(logDestination, logPrefix, log.Lmicroseconds|log.Lshortfile)

	e.logger.Println("Constructed Engine")

	return e
}

func (e *Engine) ToApi() <-chan ObjectiveChangeEvent {
	return e.toApi
}

// Run kicks of an infinite loop that waits for communications on the supplied channels, and handles them accordingly
func (e *Engine) Run() {
	for {
		var res ObjectiveChangeEvent
		var err error
		select {
		case apiEvent := <-e.FromAPI:
			res, err = e.handleAPIEvent(apiEvent)

		case chainEvent := <-e.fromChain:

			res, err = e.handleChainEvent(chainEvent)

		case message := <-e.fromMsg:
			res, err = e.handleMessage(message)
		}

		// Handle errors
		if err != nil {
			e.logger.Panic(err)
			// TODO do not panic if in production.
			// TODO report errors back to the consuming application
		}

		// Only send out an event if there are changes
		if len(res.CompletedObjectives) > 0 {
			for _, obj := range res.CompletedObjectives {
				e.logger.Printf("Objective %s is complete & returned to API", obj.Id())
			}
			e.toApi <- res
		}

	}
}

// handleMessage handles a Message from a peer go-nitro Wallet.
// It
// reads an objective from the store,
// gets a pointer to a channel secret key from the store,
// generates an updated objective and
// attempts progress.
func (e *Engine) handleMessage(message protocols.Message) (ObjectiveChangeEvent, error) {

	e.logger.Printf("Handling inbound message %+v", summarizeMessage(message))
	objective, err := e.getOrCreateObjective(message)
	if err != nil {
		return ObjectiveChangeEvent{}, err
	}
	event := protocols.ObjectiveEvent{
		ObjectiveId:     message.ObjectiveId,
		SignedStates:    message.SignedStates,
		SignedProposals: message.SignedProposals,
	}
	updatedObjective, err := objective.Update(event)
	if err != nil {
		return ObjectiveChangeEvent{}, err
	}

	return e.attemptProgress(updatedObjective)

}

// handleChainEvent handles a Chain Event from the blockchain.
// It
// reads an objective from the store,
// gets a pointer to a channel secret key from the store,
// generates an updated objective and
// attempts progress.
func (e *Engine) handleChainEvent(chainEvent chainservice.Event) (ObjectiveChangeEvent, error) {
	e.logger.Printf("handling chain event %v", chainEvent)
	objective, ok := e.store.GetObjectiveByChannelId(chainEvent.GetChannelId())
	if !ok {
		return ObjectiveChangeEvent{}, &ErrUnhandledChainEvent{event: chainEvent, objective: objective}
	}

	eventHandler, ok := objective.(chainservice.ChainEventHandler)
	if !ok {
		return ObjectiveChangeEvent{}, &ErrUnhandledChainEvent{event: chainEvent, objective: objective}
	}
	updatedEventHandler, err := eventHandler.UpdateWithChainEvent(chainEvent)
	if err != nil {
		return ObjectiveChangeEvent{}, err
	}
	return e.attemptProgress(updatedEventHandler)
}

// handleAPIEvent handles an API Event (triggered by an API call)
// It will attempt to perform all of the following:
// Spawn a new, approved objective (if not null)
// Reject an existing objective (if not null)
// Approve an existing objective (if not null)
func (e *Engine) handleAPIEvent(apiEvent APIEvent) (ObjectiveChangeEvent, error) {
	if apiEvent.ObjectiveToSpawn != nil {

		switch request := (apiEvent.ObjectiveToSpawn).(type) {

		case virtualfund.ObjectiveRequest:
			vfo, err := virtualfund.NewObjective(request, e.store.GetTwoPartyLedger, e.store.GetConsensusChannel)
			if err != nil {
				return ObjectiveChangeEvent{}, fmt.Errorf("handleAPIEvent: Could not create objective for %+v: %w", request, err)
			}
			return e.attemptProgress(&vfo)

		case directfund.ObjectiveRequest:
			dfo, err := directfund.NewObjective(request, true)
			if err != nil {
				return ObjectiveChangeEvent{}, fmt.Errorf("handleAPIEvent: Could not create objective for %+v: %w", request, err)
			}
			return e.attemptProgress(&dfo)

		default:
			return ObjectiveChangeEvent{}, fmt.Errorf("handleAPIEvent: Unknown objective type %T", request)
		}

	}

	if apiEvent.ObjectiveToReject != `` {
		objective, err := e.store.GetObjectiveById(apiEvent.ObjectiveToReject)
		if err != nil {
			return ObjectiveChangeEvent{}, err
		}
		updatedProtocol := objective.Reject()
		err = e.store.SetObjective(updatedProtocol)
		return ObjectiveChangeEvent{}, err
	}
	if apiEvent.ObjectiveToApprove != `` {
		objective, err := e.store.GetObjectiveById(apiEvent.ObjectiveToReject)
		if err != nil {
			return ObjectiveChangeEvent{}, err
		}
		updatedObjective := objective.Approve()
		return e.attemptProgress(updatedObjective)
	}
	return ObjectiveChangeEvent{}, nil

}

// executeSideEffects executes the SideEffects declared by cranking an Objective
func (e *Engine) executeSideEffects(sideEffects protocols.SideEffects) {
	for _, message := range sideEffects.MessagesToSend {
		e.logger.Printf("Sending message %+v", summarizeMessage(message))
		e.toMsg <- message
	}
	for _, tx := range sideEffects.TransactionsToSubmit {
		e.logger.Printf("Sending chain transaction for channel %s", tx.ChannelId)
		e.toChain <- tx
	}
}

// attemptProgress takes a "live" objective in memory and performs the following actions:
//
// 	1. It pulls the secret key from the store
// 	2. It cranks the objective with that key
// 	3. It commits the cranked objective to the store
// 	4. It executes any side effects that were declared during cranking
// 	5. It updates progress metadata in the store
func (e *Engine) attemptProgress(objective protocols.Objective) (outgoing ObjectiveChangeEvent, err error) {

	secretKey := e.store.GetChannelSecretKey()

	crankedObjective, sideEffects, waitingFor, err := objective.Crank(secretKey)

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
		err = e.SpawnConsensusChannelIfDirectFundObjective(crankedObjective) // Here we assume that every directfund.Objective is for a ledger channel.
		if err != nil {
			return
		}
	}
	e.executeSideEffects(sideEffects)
	return
}

// SpawnConsensusChannelIfDirectFundObjective will attempt to create and store a ConsensusChannel derived from the supplied Objective iff it is a directfund.Objective.
func (e Engine) SpawnConsensusChannelIfDirectFundObjective(crankedObjective protocols.Objective) error {
	if dfo, isDfo := crankedObjective.(*directfund.Objective); isDfo {
		c, err := dfo.CreateConsensusChannel()
		if err != nil {
			return fmt.Errorf("could not create consensus channel for objective %s: %w", crankedObjective.Id(), err)
		}
		err = e.store.SetConsensusChannel(c)
		if err != nil {
			return fmt.Errorf("could not store consensus channel for objective %s: %w", crankedObjective.Id(), err)
		}
	}
	return nil
}

// getOrCreateObjective creates the objective if the supplied message is a proposal. Otherwise, it attempts to get the objective from the store.
func (e *Engine) getOrCreateObjective(message protocols.Message) (protocols.Objective, error) {
	id := message.ObjectiveId

	objective, err := e.store.GetObjectiveById(id)

	if err == nil {
		return objective, nil
	} else if errors.Is(err, store.ErrNoSuchObjective) {

		newObj, err := e.constructObjectiveFromMessage(message)
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
		return nil, fmt.Errorf("unexpected error getting/creating objective %s: %w", message.ObjectiveId, err)
	}
}

// constructObjectiveFromMessage Constructs a new objective (of the appropriate concrete type) from the supplied message.
func (e *Engine) constructObjectiveFromMessage(message protocols.Message) (protocols.Objective, error) {

	switch {
	case directfund.IsDirectFundObjective(message.ObjectiveId):
		dfo, err := directfund.ConstructObjectiveFromMessage(
			message, *e.store.GetAddress())

		return &dfo, err
	case virtualfund.IsVirtualFundObjective(message.ObjectiveId):
		vfo, err := virtualfund.ConstructObjectiveFromMessage(message, *e.store.GetAddress(), e.store.GetTwoPartyLedger, e.store.GetConsensusChannel)
		if err != nil {
			return &virtualfund.Objective{}, fmt.Errorf("could not create virtual fund objective from message: %w", err)
		}
		return &vfo, nil

	default:
		return &directfund.Objective{}, errors.New("cannot handle unimplemented objective type")
	}

}

type messageSummary struct {
	to           string
	objectiveId  protocols.ObjectiveId
	signedStates []signedStateSummary
}
type signedStateSummary struct {
	turnNum   uint64
	channelId string
}

// summarizeMessage returns a basic summary of a message suitable for logging
func summarizeMessage(message protocols.Message) messageSummary {
	summary := messageSummary{to: message.To.String(), objectiveId: message.ObjectiveId, signedStates: []signedStateSummary{}}
	for _, signedState := range message.SignedStates {
		channelId, err := signedState.State().ChannelId()
		if err != nil {
			panic(err)
		}
		summary.signedStates = append(summary.signedStates, signedStateSummary{turnNum: signedState.State().TurnNum, channelId: channelId.String()})
	}
	return summary
}
