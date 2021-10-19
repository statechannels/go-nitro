// Package engine contains the types and imperative code for the business logic of a go-nitro Client.
package engine // import "github.com/statechannels/go-nitro/client/engine"

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// Engine is the imperative part of the core business logic of a go-nitro Client
type Engine struct {
	// inbound go channels
	API   chan APIEvent
	Chain chan ChainEvent
	Inbox chan Message

	store store.Store // A Store foe persisting important data
}

// APIEvent is an internal representation of an API call
type APIEvent struct {
	ObjectiveToSpawn   protocols.Objective
	ObjectiveToReject  protocols.ObjectiveId
	ObjectiveToApprove protocols.ObjectiveId

	Response chan Response
}

// ChainEvent is an internal representation of a blockchain event
type ChainEvent struct {
	ChannelId          types.Bytes32
	Holdings           map[types.Address]big.Int // indexed by asset
	AdjudicationStatus protocols.AdjudicationStatus
}

// Message is an internal representation of a message from another client
type Message struct {
	ObjectiveId protocols.ObjectiveId
	Sigs        map[types.Bytes32]state.Signature // mapping from state hash to signature
}

// Response is the return type that asynchronous API calls "resolve to". Such a call returns a go channel of type Response.
type Response struct{}

// NewEngine is the constructor for an Engine
func New() Engine {
	e := Engine{}

	// create the engine's inbound channels
	e.API = make(chan APIEvent)
	e.Chain = make(chan ChainEvent)
	e.Inbox = make(chan Message)

	return e
}

// Run kicks of an infinite loop that waits for communications on the supplied channels, and handles them accordingly
func (e *Engine) Run() {
	for {
		select {
		case apiEvent := <-e.API:
			e.handleAPIEvent(apiEvent)

		case chainEvent := <-e.Chain:
			e.handleChainEvent(chainEvent)

		case message := <-e.Inbox:
			e.handleMessage(message)

		}
	}
}

// handleMessage handles a Message from a peer go-nitro Wallet.
// It
// reads an objective from the store,
// generates an updated objective and declaration of side effects,
// commits the updated objective to the store,
// executes the side effects and
// evaluates objecive progress.
func (e *Engine) handleMessage(message Message) {
	objective := e.store.GetObjectiveById(message.ObjectiveId)
	event := protocols.ObjectiveEvent{Sigs: message.Sigs}
	updatedProtocol := objective.Update(event)
	_ = e.store.SetObjective(updatedProtocol)             // TODO handle error
	sideEffects, waitingFor, _ := updatedProtocol.Crank() // TODO handle error
	e.executeSideEffects(sideEffects)
	e.store.UpdateProgressLastMadeAt(message.ObjectiveId, waitingFor)
}

// handleChainEvent handles a Chain Event from the blockchain.
// It
// reads an objective from the store,
// generates an updated objective and declaration of side effects,
// commits the updated objective to the store,
// executes the side effects and
// evaluates objecive progress.
func (e *Engine) handleChainEvent(chainEvent ChainEvent) {
	objective := e.store.GetObjectiveByChannelId(chainEvent.ChannelId)
	event := protocols.ObjectiveEvent{Holdings: chainEvent.Holdings, AdjudicationStatus: chainEvent.AdjudicationStatus}
	updatedProtocol := objective.Update(event)
	_ = e.store.SetObjective(updatedProtocol)             // TODO handle error
	sideEffects, waitingFor, _ := updatedProtocol.Crank() // TODO handle error
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
		e.store.SetObjective(updatedProtocol)
	}
	if apiEvent.ObjectiveToApprove != `` {
		objective := e.store.GetObjectiveById(apiEvent.ObjectiveToReject)
		updatedProtocol := objective.Approve()
		e.store.SetObjective(updatedProtocol)
	}
}

// executeSideEffects executes the SideEffects declared by cranking an Objective
func (e *Engine) executeSideEffects(protocols.SideEffects) {
	// TODO
}
