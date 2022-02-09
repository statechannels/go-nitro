package virtualfund

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

func TestMarshalJSON(t *testing.T) {
	// left, _ :=
	ts := state.TestState
	ts.TurnNum = channel.PreFundTurnNum

	right, _ := channel.NewTwoPartyLedger(ts, 0)
	vfo, err := New(
		false,
		ts,
		types.Address{},
		2,
		0,
		channel.TwoPartyLedger{},
		right,
	)

	if err != nil {
		err = fmt.Errorf("the test VirtualFundObjective was not initialized: %w", err)
		t.Fatalf("%s", err)
	}

	encodedVfo, err := json.Marshal(vfo)

	if err != nil {
		t.Errorf("error encoding direct-fund objective %v", vfo)
	}

	got := VirtualFundObjective{}
	got.UnmarshalJSON(encodedVfo)

	if !(got.Status == vfo.Status) {
		t.Errorf("expected Status %v but got %v", vfo.Status, got.Status)
	}
	if got.V.Id != vfo.V.Id {
		t.Errorf("expected channel Id %s but got %s", vfo.V.Id, got.V.Id)
	}

	if vfo.ToMyLeft != nil {
		if !reflect.DeepEqual(vfo.ToMyLeft.ExpectedGuarantees, got.ToMyLeft.ExpectedGuarantees) {
			t.Errorf("expected left-channel guarantees %v, but found %v", vfo.ToMyLeft, got.ToMyLeft)
		}

		if got.ToMyLeft.Channel.Id != vfo.ToMyLeft.Channel.Id {
			t.Errorf("expected left channel Id %s but got %s",
				vfo.ToMyLeft.Channel.Id, got.ToMyLeft.Channel.Id)
		}
	} else if (got.ToMyLeft.Channel.Id != types.Destination{}) {
		t.Errorf("recieved a non-blank channel Id where the connection was null")
	}

	if vfo.ToMyRight != nil {
		if !reflect.DeepEqual(vfo.ToMyRight.ExpectedGuarantees, got.ToMyRight.ExpectedGuarantees) {
			t.Errorf("expected right-channel %v, but found %v", vfo.ToMyRight, got.ToMyRight)
		}

		if got.ToMyRight.Channel.Id != vfo.ToMyRight.Channel.Id {
			t.Errorf("expected left channel Id %s but got %s",
				vfo.ToMyRight.Channel.Id, got.ToMyRight.Channel.Id)
		}
	} else if (got.ToMyRight.Channel.Id != types.Destination{}) {
		t.Errorf("recieved a non-blank channel Id where the connection was null")
	}

	if got.n != vfo.n {
		t.Errorf("expected %d channel participants but found %d", vfo.n, got.n)
	}
	if got.MyRole != vfo.MyRole {
		t.Errorf("expected MyRole %d but found %d", vfo.MyRole, got.MyRole)
	}
	if !got.a0.Equal(vfo.a0) {
		t.Errorf("expected alice initial balance of %v but found %v", vfo.a0, got.a0)
	}
	if !got.b0.Equal(vfo.b0) {
		t.Errorf("expected bob initial balance of %v but found %v", vfo.b0, got.a0)
	}
	if got.requestedLedgerUpdates != vfo.requestedLedgerUpdates {
		t.Errorf("expected requestedLedgerUpdates == %t, but found %t",
			vfo.requestedLedgerUpdates, got.requestedLedgerUpdates)
	}
}
