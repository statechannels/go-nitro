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
	Payloads: []protocols.MessagePayload{{
		ObjectiveId: testId,
	}},
}

func TestConnect(t *testing.T) {
	bobOut := bobMS.Out()

	aliceMS.Send(aToB)

	got := <-bobOut

	if got.Payloads[0].ObjectiveId != testId {
		t.Fatalf("expected bob to receive ObjectiveId %v, but received %v",
			testId, got.Payloads[0].ObjectiveId)
	}
}
