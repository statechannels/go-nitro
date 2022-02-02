// Package engine contains the types and imperative code for the business logic of a go-nitro Client.
package engine // import "github.com/statechannels/go-nitro/client/engine"

import (
	"io"
	"log"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/protocols"
)

// Engine is the imperative part of the core business logic of a go-nitro Client
type Engine struct {
	// inbound go channels
	FromAPI   chan APIEvent // This one is exported so that the Client can send API calls
	fromChain <-chan chainservice.Event
	fromMsg   <-chan protocols.Message

	// outbound go channels
	toMsg   chan<- protocols.Message
	toChain chan<- protocols.ChainTransaction

	store store.Store // A Store for persisting and restoring important data

	logger *log.Logger
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
func New(msg messageservice.MessageService, chain chainservice.ChainService, store store.Store, logDestination io.Writer) Engine {
	e := Engine{}

	e.store = store

	// bind to inbound chans
	e.FromAPI = make(chan APIEvent)
	e.fromChain = chain.Out()
	e.fromMsg = msg.Out()

	// bind to outbound chans
	e.toChain = chain.In()
	e.toMsg = msg.In()

	// initialize a Logger
	e.logger = log.New(logDestination, e.store.GetAddress().String()+": ", log.Ldate|log.Ltime|log.Lshortfile)

	e.logger.Println("Constructed Engine")

	return e
}

// Run kicks of an infinite loop that waits for communications on the supplied channels, and handles them accordingly
func (e *Engine) Run() {
	for {
		select {
		case apiEvent := <-e.FromAPI:
			e.handleAPIEvent(apiEvent)

		case chainEvent := <-e.fromChain:
			e.handleChainEvent(chainEvent)

		case message := <-e.fromMsg:
			e.handleMessage(message)

		}
	}
}

// handleMessage handles a Message from a peer go-nitro Wallet.
// It
// reads an objective from the store,
// gets a pointer to a channel secret key from the store,
// generates an updated objective and
// attempts progress.
func (e *Engine) handleMessage(message protocols.Message) {
	var objective protocols.Objective
	ok := true
	if message.Proposal != nil {
		objective = message.Proposal
		// TODO ensure objective in only approved if the application has given permission somehow
	} else {
		objective, ok = e.store.GetObjectiveById(message.ObjectiveId)
	}
	if ok {
		event := protocols.ObjectiveEvent{SignedStates: message.SignedStates}
		updatedObjective, error := objective.Update(event)
		if error == nil {
			e.attemptProgress(updatedObjective)
		}
	}
}

// handleChainEvent handles a Chain Event from the blockchain.
// It
// reads an objective from the store,
// gets a pointer to a channel secret key from the store,
// generates an updated objective and
// attempts progress.
func (e *Engine) handleChainEvent(chainEvent chainservice.Event) {
	event := protocols.ObjectiveEvent{Holdings: chainEvent.Holdings, AdjudicationStatus: chainEvent.AdjudicationStatus}
	objective, _ := e.store.GetObjectiveByChannelId(chainEvent.ChannelId)
	updatedObjective, _ := objective.Update(event) // TODO handle error
	e.attemptProgress(updatedObjective)
}

// handleAPIEvent handles an API Event (triggered by an API call)
// It will attempt to perform all of the following:
// Spawn a new, approved objective (if not null)
// Reject an existing objective (if not null)
// Approve an existing objective (if not null)
func (e *Engine) handleAPIEvent(apiEvent APIEvent) {
	if apiEvent.ObjectiveToSpawn != nil {
		e.attemptProgress(apiEvent.ObjectiveToSpawn)
	}
	if apiEvent.ObjectiveToReject != `` {
		objective, _ := e.store.GetObjectiveById(apiEvent.ObjectiveToReject)
		updatedProtocol := objective.Reject()
		_ = e.store.SetObjective(updatedProtocol) // TODO handle error
	}
	if apiEvent.ObjectiveToApprove != `` {
		objective, _ := e.store.GetObjectiveById(apiEvent.ObjectiveToReject)
		updatedObjective := objective.Approve()
		e.attemptProgress(updatedObjective)
	}
}

// executeSideEffects executes the SideEffects declared by cranking an Objective
func (e *Engine) executeSideEffects(sideEffects protocols.SideEffects) {
	for _, message := range sideEffects.MessagesToSend {
		e.logger.Printf("Sending message to %s", message.To)
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
func (e *Engine) attemptProgress(objective protocols.Objective) {
	secretKey := e.store.GetChannelSecretKey()
	crankedObjective, sideEffects, waitingFor, _ := objective.Crank(secretKey) // TODO handle error
	_ = e.store.SetObjective(crankedObjective)                                 // TODO handle error
	e.executeSideEffects(sideEffects)
	e.logger.Printf("Objective %s is %s", objective.Id(), waitingFor)
	e.store.UpdateProgressLastMadeAt(objective.Id(), waitingFor)
}
