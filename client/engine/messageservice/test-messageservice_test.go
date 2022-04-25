package messageservice

import (
	"testing"

	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var broker = NewBroker()
var aliceMS = NewTestMessageService(types.Address{'a'}, broker, 0)
var bobMS = NewTestMessageService(types.Address{'b'}, broker, 0)

var testId protocols.ObjectiveId = "testObjectiveID"

var aToB protocols.Message = protocols.Message{
	To: bobMS.address,
	ObjectivePayloads: []protocols.ObjectivePayload{{
		ObjectiveId: testId,
	}},
}

func TestConnect(t *testing.T) {
	bobOut := bobMS.Out()

	aliceMS.in <- aToB

	got := <-bobOut

	if got.ObjectivePayloads[0].ObjectiveId != testId {
		t.Fatalf("expected bob to recieve ObjectiveId %v, but recieved %v",
			testId, got.ObjectivePayloads[0].ObjectiveId)
	}
}
