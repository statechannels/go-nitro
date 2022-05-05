package virtualdefund

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/internal/channel/state"
	"github.com/statechannels/go-nitro/internal/protocols"
	ta "github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testhelpers"
	. "github.com/statechannels/go-nitro/internal/testhelpers"
)

var alice = ta.Alice
var bob = ta.Bob
var irene = ta.Irene
var allActors = []ta.Actor{alice, irene, bob}

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

	virtualDefund := newObjective(false, data.vFinal.FixedPart(), data.initialOutcome, big.NewInt(int64(data.paid)), nil, nil, 0)
	invalidFinal := data.vFinal.Clone()
	invalidFinal.ChannelNonce = big.NewInt(5)

	signedFinal := state.NewSignedState(invalidFinal)

	// Sign the final state by some other participant
	signStateByOthers(alice, signedFinal)

	e := protocols.ObjectiveEvent{ObjectiveId: virtualDefund.Id(), SignedState: signedFinal}
	_, err := virtualDefund.Update(e)
	if err.Error() != "event channelId out of scope of objective" {
		t.Errorf("Expected error for channelId being out of scope, got %v", err)
	}

}

func testUpdateAs(my ta.Actor) func(t *testing.T) {
	return func(t *testing.T) {
		data := generateTestData()
		vId, _ := data.vFinal.ChannelId()
		left, right := generateLedgers(my.Role, vId)

		virtualDefund := newObjective(false, data.vFinal.FixedPart(), data.initialOutcome, big.NewInt(int64(data.paid)), left, right, my.Role)
		signedFinal := state.NewSignedState(data.vFinal)
		// Sign the final state by some other participant
		signStateByOthers(my, signedFinal)

		e := protocols.ObjectiveEvent{ObjectiveId: virtualDefund.Id(), SignedState: signedFinal}

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
		vId, _ := data.vFinal.ChannelId()
		left, right := generateLedgers(my.Role, vId)
		virtualDefund := newObjective(true, data.vFinal.FixedPart(), data.initialOutcome, big.NewInt(int64(data.paid)), left, right, my.Role)

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
		AssertStateSentToEveryone(t, se, signedByMe, my, allActors)

		// Update the signatures on the objective so the final state is fully signed
		signedByOthers := signStateByOthers(my, state.NewSignedState(data.vFinal))
		for i, sig := range signedByOthers.Signatures() {
			if uint(i) != my.Role {
				updated.Signatures[i] = sig
			}
		}

		updatedObj, se, waitingFor, err = updated.Crank(&my.PrivateKey)
		updated = updatedObj.(*Objective)
		testhelpers.Ok(t, err)

		testhelpers.Equals(t, WaitingForCompleteLedgerDefunding, waitingFor)

		checkForLeaderProposals(t, se, updated, data)

		proposals := generateProposalsResponses(my.Role, vId, updated, data)
		updateProposals(updated, proposals...)

		updatedObj, se, waitingFor, err = updated.Crank(&my.PrivateKey)
		updated = updatedObj.(*Objective)
		testhelpers.Ok(t, err)

		testhelpers.Equals(t, waitingFor, WaitingForNothing)
		checkForFollowerProposals(t, se, updated, data)

	}

}
