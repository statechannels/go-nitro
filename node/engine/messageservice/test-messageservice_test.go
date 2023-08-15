package messageservice

import (
	"testing"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/internal/testhelpers"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var (
	broker  = NewBroker()
	aliceMS = NewTestMessageService(types.Address{'a'}, broker, 0)
	bobMS   = NewTestMessageService(types.Address{'b'}, broker, 0)
)

var testId protocols.ObjectiveId = "VirtualDefund-0x0000000000000000000000000000000000000000000000000000000000000000"

var aToB protocols.Message = protocols.CreateSignedProposalMessage(
	bobMS.address,
	consensus_channel.SignedProposal{
		Proposal: consensus_channel.Proposal{LedgerID: types.Destination{1}},
		TurnNum:  1,
	},
)

func TestConnect(t *testing.T) {
	bobOut := bobMS.Out()

	err := aliceMS.Send(aToB)
	if err != nil {
		t.Fatal(err)
	}

	got := <-bobOut

	prop := got.LedgerProposals[0]

	objId, err := protocols.GetProposalObjectiveId(prop.Proposal)

	testhelpers.Ok(t, err)

	if objId != testId {
		t.Fatalf("expected bob to receive ObjectiveId %v, but received %v",
			testId, objId)
	}
}
