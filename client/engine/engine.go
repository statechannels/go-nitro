// Package engine contains the types and imperative code for the business logic of a go-nitro Client.
package engine // import "github.com/statechannels/go-nitro/client/engine"

import (
	"context"
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

// We use buffered channels for communication between the Engine and the API
// This prevents the engine from getting blocked writing to toApi
// This also prevents the client from getting blocked writing to FromAPI
const API_BUFFER_SIZE = 10

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

	channelLocker *ChannelLocker

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
func New(msg messageservice.MessageService, chain chainservice.ChainService, store store.Store, logDestination io.Writer, channelLocker *ChannelLocker) Engine {
	e := Engine{}

	e.store = store
	e.channelLocker = channelLocker
	// bind to inbound chans
	e.FromAPI = make(chan APIEvent, API_BUFFER_SIZE)
	e.fromChain = chain.Out()
	e.fromMsg = msg.Out()

	e.toApi = make(chan ObjectiveChangeEvent, API_BUFFER_SIZE)
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
// It accepts a context that can be used to cancel the loop.
func (e *Engine) Run(ctx context.Context) {
	for {
		var res ObjectiveChangeEvent
		select {
		case <-ctx.Done():
			return
		case apiEvent := <-e.FromAPI:
			res = e.handleAPIEvent(apiEvent)

		case chainEvent := <-e.fromChain:

			res = e.handleChainEvent(chainEvent)

		case message := <-e.fromMsg:
			res = e.handleMessage(message)

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
func (e *Engine) handleMessage(message protocols.Message) ObjectiveChangeEvent {

	e.logger.Printf("Handling inbound message %+v", summarizeMessage(message))
	objective, err := e.getOrCreateObjective(message)
	// Acquire a lock for the channels and then refetch the objective with the latest version
	e.channelLocker.Lock(objective)
	objective, _ = e.getOrCreateObjective(message)
	defer e.channelLocker.Unlock(objective)

	if err != nil {
		e.logger.Print(err)
		return ObjectiveChangeEvent{}
	}
	event := protocols.ObjectiveEvent{ObjectiveId: message.ObjectiveId, SignedStates: message.SignedStates}
	updatedObjective, err := objective.Update(event)
	if err != nil {
		e.logger.Print(err)
		return ObjectiveChangeEvent{}
	}

	return e.attemptProgress(updatedObjective)

}

// handleChainEvent handles a Chain Event from the blockchain.
// It
// reads an objective from the store,
// gets a pointer to a channel secret key from the store,
// generates an updated objective and
// attempts progress.
func (e *Engine) handleChainEvent(chainEvent chainservice.Event) ObjectiveChangeEvent {
	e.logger.Printf("preparing to lock for  chain event %v", chainEvent)
	objective, ok := e.store.GetObjectiveByChannelId(chainEvent.ChannelId)
	if !ok {
		e.logger.Printf("handleChainEvent: No objective in store for channel with id %s", chainEvent.ChannelId)
		return ObjectiveChangeEvent{}
	}
	// Acquire a lock for the channels and then refetch the objective to get the latest version
	e.channelLocker.Lock(objective)
	defer e.channelLocker.Unlock(objective)
	objective, _ = e.store.GetObjectiveByChannelId(chainEvent.ChannelId)
	e.logger.Printf("handling chain event %v", chainEvent)
	event := protocols.ObjectiveEvent{
		Holdings:           chainEvent.Holdings,
		BlockNum:           chainEvent.BlockNum,
		AdjudicationStatus: chainEvent.AdjudicationStatus,
		ObjectiveId:        objective.Id()}
	updatedObjective, err := objective.Update(event)
	if err != nil {
		// TODO handle error
		panic(err)
	}
	return e.attemptProgress(updatedObjective)

}

// handleAPIEvent handles an API Event (triggered by an API call)
// It will attempt to perform all of the following:
// Spawn a new, approved objective (if not null)
// Reject an existing objective (if not null)
// Approve an existing objective (if not null)
func (e *Engine) handleAPIEvent(apiEvent APIEvent) ObjectiveChangeEvent {
	if apiEvent.ObjectiveToSpawn != nil {

		switch request := (apiEvent.ObjectiveToSpawn).(type) {

		case virtualfund.ObjectiveRequest:
			vfo, err := virtualfund.NewObjective(request, e.store.GetTwoPartyLedger)
			if err != nil {
				e.logger.Printf("handleAPIEvent: Could not create objective for  %+v", request)
				return ObjectiveChangeEvent{}
			}
			// Acquire a lock for the channels and then generate the objective with the latest data
			e.channelLocker.Lock(&vfo)
			defer e.channelLocker.Unlock(&vfo)
			vfo, _ = virtualfund.NewObjective(request, e.store.GetTwoPartyLedger)

			return e.attemptProgress(&vfo)

		case directfund.ObjectiveRequest:
			dfo, err := directfund.NewObjective(request, true)
			if err != nil {
				e.logger.Printf("handleAPIEvent: Could not create objective for  %+v", request)
				return ObjectiveChangeEvent{}
			}

			// Acquire a lock for the channels and then generate the objective with the latest data
			e.channelLocker.Lock(&dfo)
			defer e.channelLocker.Unlock(&dfo)
			dfo, _ = directfund.NewObjective(request, true)

			return e.attemptProgress(&dfo)

		default:

			e.logger.Printf("handleAPIEvent: Unknown objective type %T", request)
			return ObjectiveChangeEvent{}

		}

	}

	if apiEvent.ObjectiveToReject != `` {
		objective, _ := e.store.GetObjectiveById(apiEvent.ObjectiveToReject)
		updatedProtocol := objective.Reject()
		_ = e.store.SetObjective(updatedProtocol) // TODO handle error
		return ObjectiveChangeEvent{}
	}
	if apiEvent.ObjectiveToApprove != `` {
		objective, _ := e.store.GetObjectiveById(apiEvent.ObjectiveToReject)
		updatedObjective := objective.Approve()
		return e.attemptProgress(updatedObjective)

	}
	return ObjectiveChangeEvent{}

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
func (e *Engine) attemptProgress(objective protocols.Objective) (outgoing ObjectiveChangeEvent) {
	secretKey := e.store.GetChannelSecretKey()

	crankedObjective, sideEffects, waitingFor, _ := objective.Crank(secretKey) // TODO handle error
	_ = e.store.SetObjective(crankedObjective)                                 // TODO handle error

	// TODO: This is hack to get around the fact that currently each objective in the store has it's own set of channels.
	vfo, isVirtual := crankedObjective.(*virtualfund.Objective)
	if isVirtual {
		if vfo.ToMyLeft != nil {
			_ = e.store.SetChannel(&vfo.ToMyLeft.Channel.Channel)
		}
		if vfo.ToMyRight != nil {
			_ = e.store.SetChannel(&vfo.ToMyRight.Channel.Channel)
		}
	}

	e.executeSideEffects(sideEffects)
	e.logger.Printf("Objective %s is %s", objective.Id(), waitingFor)

	// If our protocol is waiting for nothing then we know the objective is complete
	// TODO: If attemptProgress is called on a completed objective CompletedObjectives would include that objective id
	// Probably should have a better check that only adds it to CompletedObjectives if it was completed in this crank
	if waitingFor == "WaitingForNothing" {
		outgoing.CompletedObjectives = append(outgoing.CompletedObjectives, crankedObjective)
	}
	return
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
		vfo, err := virtualfund.ConstructObjectiveFromMessage(message, *e.store.GetAddress(), e.store.GetTwoPartyLedger)
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
