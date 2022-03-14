package messageservice

import (
	"testing"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var broker = NewBroker()
var aliceMS = NewTestMessageService(types.Address{'a'}, broker)
var bobMS = NewTestMessageService(types.Address{'b'}, broker)

var testId protocols.ObjectiveId = "testObjectiveID"

var aToB protocols.Message = protocols.Message{
	To:           bobMS.address,
	ObjectiveId:  testId,
	SignedStates: []state.SignedState{},
}

func TestConnect(t *testing.T) {
	bobIn := bobMS.Inbox()

	aliceMS.outbox <- aToB

	got := <-bobIn

	if got.ObjectiveId != testId {
		t.Errorf("expected bob to recieve ObjectiveId %v, but recieved %v",
			testId, got.ObjectiveId)
	}
}
