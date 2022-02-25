// Package engine contains the types and imperative code for the business logic of a go-nitro Client.
package engine // import "github.com/statechannels/go-nitro/client/engine"

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/ledger"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
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

	toApi chan ObjectiveChangeEvent

	store store.Store // A Store for persisting and restoring important data

	ledgerManager *ledger.LedgerManager // A LedgerManager for handling ledger requests.
	logger        *log.Logger
}

// APIEvent is an internal representation of an API call
type APIEvent struct {
	ObjectiveToSpawn   protocols.Objective
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
	e.logger = log.New(logDestination, e.store.GetAddress().String()+": ", log.Ldate|log.Ltime|log.Lshortfile)

	e.ledgerManager = ledger.NewLedgerManager()
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
		select {
		case apiEvent := <-e.FromAPI:
			res = e.handleAPIEvent(apiEvent)

		case chainEvent := <-e.fromChain:

			res = e.handleChainEvent(chainEvent)

		case message := <-e.fromMsg:
			res = e.handleMessage(message)

		}
		// Only send out an event if there are changes
		if len(res.CompletedObjectives) > 0 {
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

	e.logger.Printf("Handling inbound message %v", message)
	objective, err := e.getOrCreateObjective(message)
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
	err = e.store.SetObjective(updatedObjective)
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
	e.logger.Printf("handling chain event %v", chainEvent)
	objective, ok := e.store.GetObjectiveByChannelId(chainEvent.ChannelId)
	if !ok {
		e.logger.Printf("handleChainEvent: No objective in store for channel with id %s", chainEvent.ChannelId)
		return ObjectiveChangeEvent{}
	}
	event := protocols.ObjectiveEvent{Holdings: chainEvent.Holdings, AdjudicationStatus: chainEvent.AdjudicationStatus, ObjectiveId: objective.Id()}
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
		return e.attemptProgress(apiEvent.ObjectiveToSpawn)
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
func (e *Engine) attemptProgress(objective protocols.Objective) (outgoing ObjectiveChangeEvent) {
	secretKey := e.store.GetChannelSecretKey()

	crankedObjective, sideEffects, waitingFor, ledgerRequests, _ := objective.Crank(secretKey) // TODO handle error
	_ = e.store.SetObjective(crankedObjective)                                                 // TODO handle error

	ledgerSideEffects, _ := e.handleLedgerRequests(ledgerRequests, crankedObjective) // TODO: handle error
	sideEffects.Merge(ledgerSideEffects)

	e.executeSideEffects(sideEffects)
	e.logger.Printf("Objective %s is %s", objective.Id(), waitingFor)
	e.store.UpdateProgressLastMadeAt(objective.Id(), waitingFor)

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

	objective, ok := e.store.GetObjectiveById(id)
	if !ok {

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
		return objective, nil
	}
}

// constructObjectiveFromMessage Constructs a new objective (of the appropriate concrete type) from the supplied message.
func (e *Engine) constructObjectiveFromMessage(message protocols.Message) (protocols.Objective, error) {

	switch {
	case strings.Contains(string(message.ObjectiveId), `DirectFund`):
		initialState := message.SignedStates[0].State()
		if initialState.TurnNum != 0 {
			return directfund.DirectFundObjective{}, errors.New("cannot construct direct fund objective without prefund state")
		}
		return directfund.New(
			true, // TODO ensure objective in only approved if the application has given permission somehow
			message.SignedStates[0].State(),
			*e.store.GetAddress(),
		)
	case strings.Contains(string(message.ObjectiveId), "Virtual"):
		initialState := message.SignedStates[0].State()

		if initialState.TurnNum != 0 {
			return virtualfund.VirtualFundObjective{}, errors.New("cannot construct virtual fund objective without prefund state")
		}
		participants := message.SignedStates[0].State().Participants
		if len(participants) != 3 {
			return virtualfund.VirtualFundObjective{}, errors.New("a single hop virtual channel must have exactly 3 participants")
		}

		alice := participants[0]
		intermediary := participants[1]
		bob := participants[2]
		myAddress := *e.store.GetAddress()

		var left *channel.TwoPartyLedger
		var right *channel.TwoPartyLedger
		var ok bool

		if alice != myAddress {
			left, ok = e.store.GetTwoPartyLedger(intermediary, bob)
			if !ok {
				return virtualfund.VirtualFundObjective{}, fmt.Errorf("could not find a left ledger channel between %v and %v", alice, intermediary)
			}
		}
		if bob != myAddress {
			right, ok = e.store.GetTwoPartyLedger(alice, intermediary)
			if !ok {
				return virtualfund.VirtualFundObjective{}, fmt.Errorf("could not find a right ledger channel between %v and %v", intermediary, bob)
			}
		}

		myRole, err := getMyRole(myAddress, participants)
		if err != nil {
			return virtualfund.VirtualFundObjective{}, errors.New("could not determine my role: %w")
		}

		return virtualfund.New(
			true, // TODO ensure objective in only approved if the application has given permission somehow
			message.SignedStates[0].State(),
			*e.store.GetAddress(),
			1, // Always a single hop virtual channel
			myRole,
			left,
			right,
		)
	default:
		return virtualfund.VirtualFundObjective{}, errors.New("cannot handle unimplemented objective type")
	}

}

// handleLedgerRequests handles a collection of ledger requests by submitting them to the ledger manager.
func (e *Engine) handleLedgerRequests(ledgerRequests []protocols.LedgerRequest, objective protocols.Objective) (protocols.SideEffects, error) {
	sideEffects := protocols.SideEffects{}
	for _, req := range ledgerRequests {
		e.logger.Printf("Handling ledger request  %+v", req)

		var ledger *channel.TwoPartyLedger
		found := false
		for _, ch := range objective.Channels() {
			if ch.Id == req.LedgerId {
				ledger = &channel.TwoPartyLedger{Channel: *ch}
				found = true
			}
		}
		if !found {
			return protocols.SideEffects{}, fmt.Errorf("Could not find ledger %s", req.LedgerId)

		}

		se, err := e.ledgerManager.HandleRequest(ledger, req, e.store.GetChannelSecretKey())
		if err != nil {
			return se, fmt.Errorf("could not handle ledger request: %w", err)
		}
		err = e.store.SetChannel(&ledger.Channel)

		if err != nil {
			return se, fmt.Errorf("could not set channel: %w", err)
		}

		sideEffects.Merge(se)
	}
	return sideEffects, nil
}

func getMyRole(myAddress types.Address, participants []types.Address) (uint, error) {
	var myRole uint = ^uint(0)
	for i, p := range participants {
		if p == myAddress {
			myRole = uint(i)
		}
	}
	if myRole == ^uint(0) {
		return myRole, fmt.Errorf("could not find address %v in participants %v", myAddress, participants)
	}
	return myRole, nil
}
