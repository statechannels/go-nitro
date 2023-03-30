package virtualfund

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testhelpers"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const TEST_CHAIN_ID = 1337

func TestMarshalJSON(t *testing.T) {
	alice, p1, bob := testactors.Alice, testactors.Irene, testactors.Bob

	vPreFund := state.State{
		Participants:      []types.Address{alice.Address(), p1.Address(), bob.Address()}, // A single hop virtual channel
		ChannelNonce:      0,
		AppDefinition:     types.Address{},
		ChallengeDuration: 45,
		AppData:           []byte{},
		Outcome: outcome.Exit{outcome.SingleAssetExit{
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: alice.Destination(),
					Amount:      big.NewInt(5),
				},
				outcome.Allocation{
					Destination: bob.Destination(),
					Amount:      big.NewInt(5),
				},
			},
		}},
		TurnNum: 0,
		IsFinal: false,
	}

	ts := state.TestState
	ts.TurnNum = channel.PreFundTurnNum

	right := prepareConsensusChannel(1, alice, bob, alice)
	vfo, err := constructFromState(
		false,
		vPreFund,
		alice.Address(),
		nil,
		right,
	)
	if err != nil {
		err = fmt.Errorf("the test VirtualFundObjective was not initialized: %w", err)
		t.Fatalf("%s", err)
	}

	encodedVfo, err := json.Marshal(vfo)
	if err != nil {
		t.Fatalf("error encoding virtualfund objective %v", vfo)
	}

	got := Objective{}
	if err := got.UnmarshalJSON(encodedVfo); err != nil {
		t.Fatalf("the test VirtualFundObjective did not deserialize correctly: %s", err.Error())
	}

	if !(got.Status == vfo.Status) {
		t.Fatalf("expected Status %v but got %v", vfo.Status, got.Status)
	}

	// only checking channel ID rather than whole channel because
	// marshal / unmarshal loses channel data
	if got.V.Id != vfo.V.Id {
		t.Fatalf("expected channel Id %s but got %s", vfo.V.Id, got.V.Id)
	}

	if vfo.ToMyLeft != nil {
		if !reflect.DeepEqual(vfo.ToMyLeft.getExpectedGuarantee(), got.ToMyLeft.getExpectedGuarantee()) {
			t.Fatalf("expected left-channel guarantees %v, but found %v", vfo.ToMyLeft, got.ToMyLeft)
		}

		if got.ToMyLeft.Channel.Id != vfo.ToMyLeft.Channel.Id {
			t.Fatalf("expected left channel Id %s but got %s",
				vfo.ToMyLeft.Channel.Id, got.ToMyLeft.Channel.Id)
		}
	}

	if vfo.ToMyRight != nil {
		if !reflect.DeepEqual(vfo.ToMyRight.getExpectedGuarantee(), got.ToMyRight.getExpectedGuarantee()) {
			t.Fatalf("expected right-channel %v, but found %v", vfo.ToMyRight, got.ToMyRight)
		}

		if got.ToMyRight.Channel.Id != vfo.ToMyRight.Channel.Id {
			t.Fatalf("expected left channel Id %s but got %s",
				vfo.ToMyRight.Channel.Id, got.ToMyRight.Channel.Id)
		}
	}

	if got.n != vfo.n {
		t.Fatalf("expected %d channel participants but found %d", vfo.n, got.n)
	}
	if got.MyRole != vfo.MyRole {
		t.Fatalf("expected MyRole %d but found %d", vfo.MyRole, got.MyRole)
	}
	if !got.a0.Equal(vfo.a0) {
		t.Fatalf("expected alice initial balance of %v but found %v", vfo.a0, got.a0)
	}
	if !got.b0.Equal(vfo.b0) {
		t.Fatalf("expected bob initial balance of %v but found %v", vfo.b0, got.a0)
	}
}

func TestApproveReject(t *testing.T) {
	td := newTestData()
	lookup := td.leaderLedgers
	vPreFund := td.vPreFund

	// Assert that a valid set of constructor args does not result in an error
	o, err := constructFromState(
		false,
		vPreFund,
		alice.Address(),
		lookup[alice.Destination()].left,
		lookup[alice.Destination()].right,
	)
	testhelpers.Ok(t, err)

	approved := o.Approve()
	if approved.GetStatus() != protocols.Approved {
		t.Errorf("Expected approved status, got %v", approved.GetStatus())
	}
	rejected, sideEffects := o.Reject()
	if rejected.GetStatus() != protocols.Rejected {
		t.Errorf("Expected rejceted status, got %v", approved.GetStatus())
	}
	if len(sideEffects.MessagesToSend) != 2 {
		t.Errorf("Expected to send two messages")
	}
}
