package directfund

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

type actor struct {
	address     types.Address
	destination types.Destination
	privateKey  []byte
}

var alicePK = common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`)
var bobPK = common.Hex2Bytes(`62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2`)

var alice = actor{
	address:     common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`),
	destination: types.AddressToDestination(common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`)),
	privateKey:  alicePK,
}

var bob = actor{
	address:     common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`),
	destination: types.AddressToDestination(common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`)),
	privateKey:  bobPK,
}

var testState = state.State{
	ChainId:           big.NewInt(9001),
	Participants:      []types.Address{alice.address, bob.address},
	ChannelNonce:      big.NewInt(37140676580),
	AppDefinition:     common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`),
	ChallengeDuration: big.NewInt(60),
	AppData:           []byte{},
	Outcome: outcome.Exit{
		outcome.SingleAssetExit{
			Asset: types.Address{},
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: bob.destination, // Bob is first so we can easily test WaitingForMyTurnToFund
					Amount:      big.NewInt(5),
				},
				outcome.Allocation{
					Destination: alice.destination,
					Amount:      big.NewInt(5),
				},
			},
		},
	},
	TurnNum: 2,
	IsFinal: false,
}

// NewFromSignedState constructs a new Channel from the signed state.
func newChannelFromSignedState(ss state.SignedState, myIndex uint) (*channel.Channel, error) {
	s := ss.State()
	prefund := s.Clone()
	prefund.TurnNum = 0
	c, err := channel.New(prefund, myIndex)
	if err != nil {
		return c, err
	}

	sss := make([]state.SignedState, 1)
	sss[0] = ss
	allOk := c.AddSignedStates(sss)
	if !allOk {
		return c, errors.New("Unable to add a state to channel")
	}
	return c, nil
}

func newTestObjective(signByBob bool) (Objective, error) {
	ss := state.NewSignedState(testState)
	o := Objective{}
	sigA, err := ss.State().Sign(alicePK)
	if err != nil {
		return o, err
	}
	sigB, err := ss.State().Sign(bobPK)
	if err != nil {
		return o, err
	}
	err = ss.AddSignature(sigA)
	if err != nil {
		return o, err
	}

	if signByBob {
		err = ss.AddSignature(sigB)
		if err != nil {
			return o, err
		}
	}

	testChannel, err := newChannelFromSignedState(ss, 0)
	if err != nil {
		return o, err
	}

	// Assert that valid constructor args do not result in error
	o, err = NewObjective(false, testChannel)
	if err != nil {
		return o, err
	}
	return o, nil
}

// TestNew tests the constructor using a TestState fixture
func TestNew(t *testing.T) {
	if _, err := newTestObjective(false); err == nil {
		t.Error("expected an error constructing the defund objective from a state without signatures from all participant, but got nil")
	}

	if _, err := newTestObjective(true); err != nil {
		if err != nil {
			t.Error(err)
		}
	}
}
