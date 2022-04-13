package virtualdefund

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/testactors"
	ta "github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testhelpers"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var alice = ta.Alice
var bob = ta.Bob
var irene = ta.Irene
var allActors = []ta.Actor{alice, irene, bob}

// makeOutcome creates an outcome allocating to alice and bob
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

// generateTestData generates some test data that can be used in a test
func generateTestData() testdata {
	vFixed := state.FixedPart{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{alice.Address, irene.Address, bob.Address}, // A single hop virtual channel
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
	}

	initialOutcome := makeOutcome(7, 3)
	finalOutcome := makeOutcome(6, 4)
	paid := uint(1)

	vFinal := state.StateFromFixedAndVariablePart(vFixed, state.VariablePart{IsFinal: true, Outcome: outcome.Exit{finalOutcome}, TurnNum: FinalTurnNum})

	return testdata{vFixed, vFinal, initialOutcome, finalOutcome, paid}
}

// signByOthers signs the state by every participant except my
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

// assertStateSentToEveryone asserts that ses contains a message for every participant but from
func assertStateSentToEveryone(t *testing.T, ses protocols.SideEffects, expected state.SignedState, from testactors.Actor) {
	for _, a := range allActors {
		if a.Role != from.Role {
			assertStateSentTo(t, ses, expected, a)
		}
	}
}

// assertStateSentTo asserts that ses contains a message for the participant
func assertStateSentTo(t *testing.T, ses protocols.SideEffects, expected state.SignedState, to testactors.Actor) {
	for _, msg := range ses.MessagesToSend {
		if bytes.Equal(msg.To[:], to.Address[:]) {
			for _, ss := range msg.SignedStates {
				testhelpers.Equals(t, ss, expected)
			}
		}
	}
}

func TestUpdate(t *testing.T) {
	for _, my := range allActors {
		msg := fmt.Sprintf("testing update as %s", my.Name)
		t.Run(msg, testUpdateAs(my))
	}
}

func TestCrank(t *testing.T) {
	for _, my := range allActors {
		msg := fmt.Sprintf("testing crank as %s", my.Name)
		t.Run(msg, testCrankAs(my))
	}
}

func TestInvalidUpdate(t *testing.T) {
	data := generateTestData()

	virtualDefund := newObjective(false, data.vFixed, data.initialOutcome, big.NewInt(int64(data.paid)), nil, nil, 0)
	invalidFinal := data.vFinal.Clone()
	invalidFinal.ChannelNonce = big.NewInt(5)

	signedFinal := state.NewSignedState(invalidFinal)

	// Sign the final state by other participant
	signByOthers(alice, signedFinal)

	e := protocols.ObjectiveEvent{ObjectiveId: virtualDefund.Id(), SignedStates: []state.SignedState{signedFinal}}
	_, err := virtualDefund.Update(e)
	if err.Error() != "event channelId out of scope of objective" {
		t.Errorf("Expected error for channelId being out of scope, got %v", err)
	}

}

func testUpdateAs(my ta.Actor) func(t *testing.T) {
	return func(t *testing.T) {
		data := generateTestData()
		vId, _ := data.vFixed.ChannelId()
		left, right := generateLedgers(my.Role, vId, true)

		virtualDefund := newObjective(false, data.vFixed, data.initialOutcome, big.NewInt(int64(data.paid)), left, right, my.Role)
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
		testhelpers.Ok(t, err)
	}
}

func testCrankAs(my ta.Actor) func(t *testing.T) {
	return func(t *testing.T) {
		data := generateTestData()
		vId, _ := data.vFixed.ChannelId()
		left, right := generateLedgers(my.Role, vId, true)
		virtualDefund := newObjective(true, data.vFixed, data.initialOutcome, big.NewInt(int64(data.paid)), left, right, my.Role)

		updatedObj, se, waitingFor, err := virtualDefund.Crank(&my.PrivateKey)
		testhelpers.Ok(t, err)
		updated := updatedObj.(*Objective)

		for _, a := range allActors {
			if a.Role == my.Role {
				testhelpers.Assert(t, !isZero(updated.Signatures[a.Role]), "expected signature for participant %s to be non-zero", a.Name)
			} else {
				testhelpers.Assert(t, isZero(updated.Signatures[a.Role]), "expected signature for current participant %s to be zero", a.Name)
			}
		}

		testhelpers.Equals(t, waitingFor, WaitingForCompleteFinal)
		signedByMe := state.NewSignedState(data.vFinal)
		_ = signedByMe.Sign(&my.PrivateKey)
		assertStateSentToEveryone(t, se, signedByMe, my)

		// Update the signatures on the objective so the final state is fully signed
		signedByOthers := signByOthers(my, state.NewSignedState(data.vFinal))
		for i, sig := range signedByOthers.Signatures() {
			if uint(i) != my.Role {
				updated.Signatures[i] = sig
			}
		}

		updatedObj, se, waitingFor, err = updated.Crank(&my.PrivateKey)
		updated = updatedObj.(*Objective)
		testhelpers.Ok(t, err)

		testhelpers.Equals(t, WaitingForCompleteLedgerDefunding, waitingFor)

		checkForProposals(t, se, updated)

		// Generate ledger channels that have the guarantee removed
		defundedLeft, defundedRight := generateLedgers(my.Role, vId, false)
		updated.ToMyLeft = defundedLeft
		updated.ToMyRight = defundedRight

		updatedObj, se, waitingFor, err = updated.Crank(&my.PrivateKey)
		testhelpers.Ok(t, err)
		testhelpers.Assert(t, len(se.MessagesToSend) == 0, "expected no messages to send")
		testhelpers.Equals(t, waitingFor, WaitingForNothing)

	}

}

// checkForProposals checks that the outgoing message contains the correct proposals depending on o.MyRole
func checkForProposals(t *testing.T, se protocols.SideEffects, o *Objective) {

	leftAmount := big.NewInt(6)
	rightAmount := big.NewInt(4)

	switch o.MyRole {
	case 0:
		{
			// Alice Proposes to Irene on her right
			rightProposal := consensus_channel.SignedProposal{Proposal: consensus_channel.NewRemoveProposal(o.ToMyRight.Id, FinalTurnNum, o.VId(), leftAmount, rightAmount)}
			assertProposalSent(t, se, rightProposal, irene)
		}
	case 1:
		{
			// Irene proposes to Bob on her right
			rightProposal := consensus_channel.SignedProposal{Proposal: consensus_channel.NewRemoveProposal(o.ToMyRight.Id, FinalTurnNum, o.VId(), leftAmount, rightAmount)}
			assertProposalSent(t, se, rightProposal, bob)
		}

	}
}
