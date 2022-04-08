package virtualdefund

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func setupStates() (vPrefund, vPostFund, vFinal state.State) {
	alice := testdata.Actors.Alice
	bob := testdata.Actors.Bob
	irene := testdata.Actors.Irene
	vPrefund = state.State{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{alice.Address, irene.Address, bob.Address}, // A single hop virtual channel
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
		AppData:           []byte{},
		Outcome: outcome.Exit{outcome.SingleAssetExit{
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: alice.Destination(),
					Amount:      big.NewInt(6),
				},
				outcome.Allocation{
					Destination: bob.Destination(),
					Amount:      big.NewInt(4),
				},
			},
		}},
		TurnNum: 0,
		IsFinal: false,
	}

	vPostFund = vPrefund.Clone()
	vPostFund.TurnNum = 1

	vFinal = vPostFund.Clone()
	vFinal.TurnNum = 2
	vFinal.IsFinal = true

	return
}

func TestSingleHopVirtualDefund(t *testing.T) {
	alice := testdata.Actors.Alice
	bob := testdata.Actors.Bob
	irene := testdata.Actors.Irene
	vPreFund, vPostFund, vFinal := setupStates()
	TestAs := func(my testdata.Actor, t *testing.T) {
		// Determine my role
		var myRole uint
		for i, p := range vPreFund.Participants {
			if p == my.Address {
				myRole = uint(i)

				break
			}
		}

		// Construct  a virtual channel that has it's post fund setup signed
		signedVPostFund := state.NewSignedState(vPostFund)
		for _, signer := range []testdata.Actor{alice, irene, bob} {
			_ = signedVPostFund.Sign(&signer.PrivateKey)
		}
		virtualChannel, err := channel.NewSingleHopVirtualChannel(vPreFund, myRole)
		virtualChannel.AddSignedState(signedVPostFund)

		if err != nil {
			t.Fatal(err)
		}

		virtualDefund, err := newObjective(false, virtualChannel, my.Address, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		testUpdate := func(t *testing.T) {
			signedFinal := state.NewSignedState(vFinal)

			// Sign the final state by some other participant
			if myRole == 0 {
				_ = signedFinal.Sign(&irene.PrivateKey)

			} else {
				_ = signedFinal.Sign(&alice.PrivateKey)
			}

			e := protocols.ObjectiveEvent{ObjectiveId: virtualDefund.Id(), SignedStates: []state.SignedState{signedFinal}}

			updatedObj, err := virtualDefund.Update(e)
			updated := updatedObj.(*Objective)
			if _, ok := updated.V.SignedStateForTurnNum[2]; !ok {
				t.Fatal("Expected to find signed state for turn 2")

			}
			if err != nil {
				t.Fatal(err)
			}

		}

		t.Run(`testUpdate`, testUpdate)

	}

	t.Run(`AsAlice`, func(t *testing.T) { TestAs(alice, t) })
	t.Run(`AsBob`, func(t *testing.T) { TestAs(bob, t) })
	t.Run(`AsIrene`, func(t *testing.T) { TestAs(irene, t) })

}
