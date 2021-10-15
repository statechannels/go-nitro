package engine

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

	// outbound go channels
	client chan Response

	store store.Store // A Store foe persisting important data
	// TODO to truly make this private (e.g. to prevent the client accessing the store directly), we need to put engine in its own package
}

// APIEvent is an internal representation of an API call
type APIEvent struct {
	ObjectiveToSpawn   protocols.Objective   // try this first
	ObjectiveToReject  protocols.ObjectiveId // then this
	ObjectiveToApprove protocols.ObjectiveId // then this

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
	e.store.SetObjective(updatedProtocol)
	sideEffects, waitingFor, _ := updatedProtocol.Crank() // TODO handle error
	e.executeSideEffects(sideEffects)
	e.store.EvaluateProgress(message.ObjectiveId, waitingFor)
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
	e.store.SetObjective(updatedProtocol)
	sideEffects, waitingFor, _ := updatedProtocol.Crank() // TODO handle error
	e.executeSideEffects(sideEffects)
	e.store.EvaluateProgress(objective.Id(), waitingFor)

}

// handleAPIEvent handles an API Event (triggered by an API call)
// It will perform one of the following, in priority order:
// Spawn a new, approved objective
// Reject an existing objective
// Approve an existing objective
func (e *Engine) handleAPIEvent(apiEvent APIEvent) {
	switch {
	case apiEvent.ObjectiveToSpawn != nil:
		e.store.SetObjective(apiEvent.ObjectiveToSpawn)
	case apiEvent.ObjectiveToReject != ``:
		objective := e.store.GetObjectiveById(apiEvent.ObjectiveToReject)
		updatedProtocol := objective.Reject()
		e.store.SetObjective(updatedProtocol)
	case apiEvent.ObjectiveToApprove != ``:
		objective := e.store.GetObjectiveById(apiEvent.ObjectiveToReject)
		updatedProtocol := objective.Approve()
		e.store.SetObjective(updatedProtocol)
	}
}

// executeSideEffects executes the SideEffects declared by cranking an Objective
func (e *Engine) executeSideEffects(protocols.SideEffects) {
	// TODO
}
