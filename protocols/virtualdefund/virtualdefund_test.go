package virtualdefund

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var alice = testdata.Actors.Alice
var bob = testdata.Actors.Bob
var irene = testdata.Actors.Irene

func makeOutcome(aliceAmount uint, bobAmount uint) outcome.SingleAssetExit {
	return outcome.SingleAssetExit{
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: alice.Destination(),
				Amount:      big.NewInt(int64(aliceAmount)),
			},
			outcome.Allocation{
				Destination: bob.Destination(),
				Amount:      big.NewInt(int64(bobAmount)),
			},
		},
	}
}

func TestSingleHopVirtualDefund(t *testing.T) {

	vFixed := state.FixedPart{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{alice.Address, irene.Address, bob.Address}, // A single hop virtual channel
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
	}

	initialOutcome := makeOutcome(7, 2)
	finalOutcome := makeOutcome(6, 3)
	paid := 1

	vFinal := state.StateFromFixedAndVariablePart(vFixed, state.VariablePart{IsFinal: true, Outcome: outcome.Exit{finalOutcome}, TurnNum: FinalTurnNum})

	TestAs := func(my testdata.Actor, t *testing.T) {
		// Determine my role
		var myRole uint
		for i, p := range vFixed.Participants {
			if p == my.Address {
				myRole = uint(i)

				break
			}
		}

		virtualDefund := newObjective(false, vFixed, initialOutcome, big.NewInt(int64(paid)), myRole)

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
			if myRole == 0 {
				if isZero(updated.Signatures[1]) {
					t.Fatalf("Expected signature for participant irene to be non-zero")
				}

			} else {
				if isZero(updated.Signatures[0]) {
					t.Fatalf("Expected signature for participant alice to be non-zero")
				}
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
