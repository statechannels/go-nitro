package client

import (
	"github.com/statechannels/go-nitro/protocols"
)

type Engine struct {
	// inbound go channels
	api   chan APIEvent
	chain chan ChainEvent
	inbox chan Message

	// outbound go channels
	client chan Response

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
	event := protocols.ObjectiveEvent{Sigs: message.Sigs}
	updatedProtocol := protocol.Update(event)
	e.Store.SetObjective(updatedProtocol)
	sideEffects, waitingFor, _ := updatedProtocol.Crank() // TODO handle error
	e.executeSideEffects(sideEffects)
	e.Store.EvaluateProgress(message.ObjectiveId, waitingFor)
}

func (e *Engine) handleChainEvent(chainEvent ChainEvent) {
	protocol := e.Store.GetObjectiveByChannelId(chainEvent.ChannelId)
	event := protocols.ObjectiveEvent{Holdings: chainEvent.Holdings, AdjudicationStatus: chainEvent.AdjudicationStatus}
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

func (e *Engine) executeSideEffects(protocols.SideEffects) {
	// TODO
}
