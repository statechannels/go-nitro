package virtualdefund

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	ta "github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testhelpers"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var alice = ta.Actors.Alice
var bob = ta.Actors.Bob
var irene = ta.Actors.Irene
var allActors = []ta.Actor{alice, irene, bob}

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

type testdata struct {
	vFixed         state.FixedPart
	vFinal         state.State
	initialOutcome outcome.SingleAssetExit
	finalOutcome   outcome.SingleAssetExit
	paid           uint
}

func generateTestData() testdata {
	vFixed := state.FixedPart{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{alice.Address, irene.Address, bob.Address}, // A single hop virtual channel
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
	}

	initialOutcome := makeOutcome(7, 2)
	finalOutcome := makeOutcome(6, 3)
	paid := uint(1)

	vFinal := state.StateFromFixedAndVariablePart(vFixed, state.VariablePart{IsFinal: true, Outcome: outcome.Exit{finalOutcome}, TurnNum: FinalTurnNum})

	return testdata{vFixed, vFinal, initialOutcome, finalOutcome, paid}
}

func signByOthers(my ta.Actor, signedState state.SignedState) state.SignedState {
	if my.Role != 0 {
		_ = signedState.Sign(&alice.PrivateKey)
	}

	if my.Role != 1 {
		_ = signedState.Sign(&irene.PrivateKey)
	}

	if my.Role != 2 {
		_ = signedState.Sign(&bob.PrivateKey)
	}
	return signedState
}

func TestUpdate(t *testing.T) {
	for _, my := range allActors {
		msg := fmt.Sprintf("testing update as %s", my.Name)
		t.Run(msg, testUpdateAs(my))
	}
}

func testUpdateAs(my ta.Actor) func(t *testing.T) {
	return func(t *testing.T) {
		data := generateTestData()
		virtualDefund := newObjective(false, data.vFixed, data.initialOutcome, big.NewInt(int64(data.paid)), my.Role)
		signedFinal := state.NewSignedState(data.vFinal)
		// Sign the final state by some other participant
		signByOthers(my, signedFinal)

		e := protocols.ObjectiveEvent{ObjectiveId: virtualDefund.Id(), SignedStates: []state.SignedState{signedFinal}}

		updatedObj, err := virtualDefund.Update(e)
		updated := updatedObj.(*Objective)
		for _, a := range allActors {
			if a.Role != my.Role {
				testhelpers.Assert(t, !isZero(updated.Signatures[a.Role]), "expected signature for participant %s to be non-zero", a.Name)
			} else {
				testhelpers.Assert(t, isZero(updated.Signatures[a.Role]), "expected signature for current participant %s to be zero", a.Name)
			}
		}
		if err != nil {
			t.Fatal(err)
		}
	}
}
