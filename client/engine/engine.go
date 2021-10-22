// Package engine contains the types and imperative code for the business logic of a go-nitro Client.
package engine // import "github.com/statechannels/go-nitro/client/engine"

import (
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/protocols"
)

// Engine is the imperative part of the core business logic of a go-nitro Client
type Engine struct {

	// inbound go channels
	FromAPI   chan APIEvent // This one is exported so that the Client can send API calls
	fromChain chan chainservice.Event
	fromMsg   chan protocols.Message

	// outbound go channels
	toMsg   chan protocols.Message
	toChain chan protocols.Transaction

	store store.Store // A Store for persisting and restoring important data
}

// APIEvent is an internal representation of an API call
type APIEvent struct {
	ObjectiveToSpawn   protocols.Objective
	ObjectiveToReject  protocols.ObjectiveId
	ObjectiveToApprove protocols.ObjectiveId

	Response chan Response
}

// Response is the return type that asynchronous API calls "resolve to". Such a call returns a go channel of type Response.
type Response struct{}

// NewEngine is the constructor for an Engine
func New(msg messageservice.MessageService, chain chainservice.ChainService, store store.Store) Engine {
	e := Engine{}

	// bind the store
	e.store = store

	// bind the engine's inbound chans
	e.FromAPI = make(chan APIEvent)
	e.fromChain = chain.GetRecieveChan()
	e.fromMsg = msg.GetRecieveChan()

	// bind the engine's outbound chans
	e.toChain = chain.GetSendChan()
	e.toMsg = msg.GetSendChan()

	return e
}

// Run kicks of an infinite loop that waits for communications on the supplied channels, and handles them accordingly
func (e *Engine) Run() {
	for {
		select {
		case apiEvent := <-e.FromAPI:
			go e.handleAPIEvent(apiEvent)

		case chainEvent := <-e.fromChain:
			go e.handleChainEvent(chainEvent)

		case message := <-e.fromMsg:
			go e.handleMessage(message)

		}
	}
}

// handleMessage handles a Message from a peer go-nitro Wallet.
// It
// reads an objective from the store,
// gets a pointer to a channel secret key from the store,
// generates an updated objective and declaration of side effects,
// commits the updated objective to the store,
// executes the side effects and
// evaluates objecive progress.
func (e *Engine) handleMessage(message protocols.Message) {
	objective := e.store.GetObjectiveById(message.ObjectiveId)
	event := protocols.ObjectiveEvent{Sigs: message.Sigs}
	secretKey := e.store.GetChannelSecretKey()
	updatedProtocol, sideEffects, waitingFor, _ := objective.Update(event).Crank(secretKey) // TODO handle error
	_ = e.store.SetObjective(updatedProtocol)                                               // TODO handle error
	e.executeSideEffects(sideEffects)
	e.store.UpdateProgressLastMadeAt(message.ObjectiveId, waitingFor)
}

// handleChainEvent handles a Chain Event from the blockchain.
// It
// reads an objective from the store,
// gets a pointer to a channel secret key from the store,
// generates an updated objective and declaration of side effects,
// commits the updated objective to the store,
// executes the side effects and
// evaluates objecive progress.
func (e *Engine) handleChainEvent(chainEvent chainservice.Event) {
	objective := e.store.GetObjectiveByChannelId(chainEvent.ChannelId)
	event := protocols.ObjectiveEvent{Holdings: chainEvent.Holdings, AdjudicationStatus: chainEvent.AdjudicationStatus}
	secretKey := e.store.GetChannelSecretKey()
	updatedProtocol, sideEffects, waitingFor, _ := objective.Update(event).Crank(secretKey) // TODO handle error
	_ = e.store.SetObjective(updatedProtocol)                                               // TODO handle error
	e.executeSideEffects(sideEffects)
	e.store.UpdateProgressLastMadeAt(objective.Id(), waitingFor)
}

// handleAPIEvent handles an API Event (triggered by an API call)
// It will attempt to perform all of the following:
// Spawn a new, approved objective (if not null)
// Reject an existing objective (if not null)
// Approve an existing objective (if not null)
func (e *Engine) handleAPIEvent(apiEvent APIEvent) {
	if apiEvent.ObjectiveToSpawn != nil {
		_ = e.store.SetObjective(apiEvent.ObjectiveToSpawn) // TODO handle error
	}
	if apiEvent.ObjectiveToReject != `` {
		objective := e.store.GetObjectiveById(apiEvent.ObjectiveToReject)
		updatedProtocol := objective.Reject()
		_ = e.store.SetObjective(updatedProtocol) // TODO handle error
	}
	if apiEvent.ObjectiveToApprove != `` {
		objective := e.store.GetObjectiveById(apiEvent.ObjectiveToReject)
		updatedProtocol := objective.Approve()
		_ = e.store.SetObjective(updatedProtocol) // TODO handle error
	}
}

// executeSideEffects executes the SideEffects declared by cranking an Objective
func (e *Engine) executeSideEffects(sideEffects protocols.SideEffects) {
	for _, message := range sideEffects.MessagesToSend {
		e.toMsg <- message
	}
	for _, tx := range sideEffects.TransactionsToSubmit {
		e.toChain <- tx
	}
}
