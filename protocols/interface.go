package protocols

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

// TODO these are placeholders for now
type SideEffects []interface{}
type WaitingFor string

type Status struct {
	TurnNumRecord uint // TODO add other fields
}
type ProtocolEvent struct {
	ChannelId          string
	Sigs               map[*state.State]state.Signature // mapping from state to signature TODO consider using a hash of the state
	Holdings           big.Int                          // TODO allow for multiple assets
	AdjudicationStatus Status
}

// Protocol is the interface for off-chain protocols
type Protocol interface {
	Id() ObjectiveId
	Initialize(initialState state.State) Protocol // returns the initial Protocol object, does not declare effects

	Approve()                            // returns an updated Protocol (a copy, no mutation allowed), does not declare effects
	Reject()                             // returns an updated Protocol (a copy, no mutation allowed), does not declare effects
	Update(event ProtocolEvent) Protocol // returns an updated Protocol (a copy, no mutation allowed), does not declare effects

	Crank() (SideEffects, WaitingFor, error) // does *not* accept an event, but *does* declare side effects, does *not* return an updated Protocol

}

// TODO these are placeholders for now (they are the fundamental events the wallet reacts to)
type ObjectiveId string
type APIEvent struct {
	ObjectiveToSpawn   Protocol    // try this first
	ObjectiveToReject  ObjectiveId // then this
	ObjectiveToApprove ObjectiveId // then this

}
type ChainEvent struct {
	ChannelId          types.Bytes32
	Holdings           big.Int
	AdjudicationStatus Status
}
type Message struct {
	ObjectiveId ObjectiveId
	Sigs        map[*state.State]state.Signature // mapping from state to signature TODO consider using a hash of the state
}

type Store interface {
	GetObjectiveById(ObjectiveId) Protocol
	GetObjectiveByChannelId(types.Bytes32) Protocol
	SetObjective(Protocol) error
	ApproveObjective(ObjectiveId)
	RejectObjective(ObjectiveId)

	EvaluateProgress(ObjectiveId, string) // sets waitingFor, checks to see if objective has stalled
	GetWaitingFor(ObjectiveId) string
	SetWaitingFor(ObjectiveId, string)
}

type Engine struct {
	api   chan APIEvent
	chain chan ChainEvent
	inbox chan Message
	Store Store
}

func NewEngine() Engine {
	e := Engine{}
	e.api = make(chan APIEvent)
	e.chain = make(chan ChainEvent)
	e.inbox = make(chan Message)
	return e
}

// Run kicks of an infinite loop that waits for communications on the supplied channels, and handles them accordingly
func (e *Engine) Run() {
	for {
		select {
		case apiEvent := <-e.api:
			e.handleAPIEvent(apiEvent)

		case chainEvent := <-e.chain:
			e.handleChainEvent(chainEvent)

		case message := <-e.inbox:
			e.handleMessage(message)

		}
	}
}

func (e *Engine) handleMessage(message Message) {
	protocol := e.Store.GetObjectiveById(message.ObjectiveId)
	event := ProtocolEvent{Sigs: message.Sigs}
	updatedProtocol := protocol.Update(event)
	e.Store.SetObjective(updatedProtocol)
	sideEffects, waitingFor, _ := updatedProtocol.Crank() // TODO handle error
	e.executeSideEffects(sideEffects)
	e.Store.EvaluateProgress(message.ObjectiveId, waitingFor)
}

func (e *Engine) handleChainEvent(chainEvent ChainEvent) {
	protocol := e.Store.GetObjectiveByChannelId(chainEvent.ChannelId)
	event := ProtocolEvent{Holdings: chainEvent.Holdings, AdjudicationStatus: chainEvent.AdjudicationStatus}
	updatedProtocol := protocol.Update(event)
	e.Store.SetObjective(updatedProtocol)
	sideEffects, waitingFor, _ := updatedProtocol.Crank() // TODO handle error
	e.executeSideEffects(sideEffects)
	e.Store.EvaluateProgress(protocol.Id(), waitingFor)

}
func (e *Engine) handleAPIEvent(apiEvent APIEvent) {
	switch {
	case apiEvent.ObjectiveToSpawn != nil:
		e.Store.SetObjective(apiEvent.ObjectiveToSpawn)
	case apiEvent.ObjectiveToReject != ``:
		e.Store.ApproveObjective(apiEvent.ObjectiveToReject)
	case apiEvent.ObjectiveToApprove != ``:
		e.Store.ApproveObjective(apiEvent.ObjectiveToApprove)
	}
}

func (e *Engine) executeSideEffects(SideEffects) {
	// TODO
}

type Client struct {
	engine Engine
	api    chan APIEvent
}

func NewClient() Client {
	c := Client{}

	c.engine.api = make(chan APIEvent)
	c.engine.chain = make(chan ChainEvent)
	c.engine.inbox = make(chan Message)

	go c.engine.Run()

	return c
}

// CreateChannel creates a channel
func (c *Client) CreateChannel() {
	apiEvent := APIEvent{}
	c.engine.api <- apiEvent // The API call is "converted" into an internal event sent to the engine
}
