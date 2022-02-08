package messageservice

import (
	"testing"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/protocols"
	directfund "github.com/statechannels/go-nitro/protocols/direct-fund"
	"github.com/statechannels/go-nitro/types"
)

var aliceMS = NewTestMessageService(types.Address{'a'})
var bobMS = NewTestMessageService(types.Address{'b'})

var objective, _ = directfund.New(false, state.TestState, aliceMS.address)
var testId protocols.ObjectiveId = "testObjectiveID"

var aToB protocols.Message = protocols.Message{
	To:          bobMS.address,
	ObjectiveId: testId,
	Sigs:        make(map[*state.State]state.Signature),
	Proposal:    &objective,
}

func TestConnect(t *testing.T) {
	bobIn := bobMS.GetReceiveChan()

	aliceMS.Connect(bobMS)
	aliceMS.Send(aToB)

	got := <-bobIn

	if got.ObjectiveId != testId {
		t.Errorf("expected bob to recieve ObjectiveId %v, but recieved %v",
			got.ObjectiveId, testId)
	}
}
