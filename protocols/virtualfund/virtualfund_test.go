package virtualfund

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

func TestMarshalJSON(t *testing.T) {

	type actor struct {
		address     types.Address
		destination types.Destination
		privateKey  []byte
		role        uint
	}

	alice := actor{
		address:     common.HexToAddress(`0xD9995BAE12FEe327256FFec1e3184d492bD94C31`),
		destination: types.AddressToDestination(common.HexToAddress(`0xD9995BAE12FEe327256FFec1e3184d492bD94C31`)),
		privateKey:  common.Hex2Bytes(`7ab741b57e8d94dd7e1a29055646bafde7010f38a900f55bbd7647880faa6ee8`),
		role:        0,
	}

	p1 := actor{ // Aliases: The Hub, Irene
		address:     common.HexToAddress(`0xd4Fa489Eacc52BA59438993f37Be9fcC20090E39`),
		destination: types.AddressToDestination(common.HexToAddress(`0xd4Fa489Eacc52BA59438993f37Be9fcC20090E39`)),
		privateKey:  common.Hex2Bytes(`2030b463177db2da82908ef90fa55ddfcef56e8183caf60db464bc398e736e6f`),
		role:        1,
	}

	bob := actor{
		address:     common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`),
		destination: types.AddressToDestination(common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`)),
		privateKey:  common.Hex2Bytes(`62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2`),
		role:        2,
	}

	vPreFund := state.State{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{alice.address, p1.address, bob.address}, // A single hop virtual channel
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
		AppData:           []byte{},
		Outcome: outcome.Exit{outcome.SingleAssetExit{
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: alice.destination,
					Amount:      big.NewInt(5),
				},
				outcome.Allocation{
					Destination: bob.destination,
					Amount:      big.NewInt(5),
				},
			},
		}},
		TurnNum: 0,
		IsFinal: false,
	}

	ts := state.TestState
	ts.TurnNum = channel.PreFundTurnNum

	right, _ := channel.NewTwoPartyLedger(ts, 0)
	vfo, err := constructFromState(
		false,
		vPreFund,
		alice.address,
		&channel.TwoPartyLedger{},
		right,
	)

	if err != nil {
		err = fmt.Errorf("the test VirtualFundObjective was not initialized: %w", err)
		t.Fatalf("%s", err)
	}

	encodedVfo, err := json.Marshal(vfo)

	if err != nil {
		t.Fatalf("error encoding direct-fund objective %v", vfo)
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
		if !reflect.DeepEqual(vfo.ToMyLeft.getExpectedGuarantees(), got.ToMyLeft.getExpectedGuarantees()) {
			t.Fatalf("expected left-channel guarantees %v, but found %v", vfo.ToMyLeft, got.ToMyLeft)
		}

		if got.ToMyLeft.Channel.Id != vfo.ToMyLeft.Channel.Id {
			t.Fatalf("expected left channel Id %s but got %s",
				vfo.ToMyLeft.Channel.Id, got.ToMyLeft.Channel.Id)
		}
	}

	if vfo.ToMyRight != nil {
		if !reflect.DeepEqual(vfo.ToMyRight.getExpectedGuarantees(), got.ToMyRight.getExpectedGuarantees()) {
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
