package virtualdefund

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	ta "github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testhelpers"
	"github.com/statechannels/go-nitro/protocols"
)

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
	vId := data.vFinal.ChannelId()
	request := NewObjectiveRequest(vId)

	getChannel, getConsensusChannel := generateStoreGetters(0, vId, data.vFinal)

	virtualDefund, err := NewObjective(request, false, alice.Address(), nil, getChannel, getConsensusChannel)
	if err != nil {
		t.Fatal(err)
	}
	invalidFinal := data.vFinal.Clone()
	invalidFinal.ChannelNonce = 5

	signedFinal := state.NewSignedState(invalidFinal)

	// Sign the final state by some other participant
	signStateByOthers(alice, signedFinal)

	e := protocols.CreateObjectivePayload(virtualDefund.Id(), SignedStatePayload, signedFinal)
	_, err = virtualDefund.Update(e)
	// TODO: the protocol should probably handle this properly with a nice error
	// if err.Error() != "event channelId out of scope of objective" {
	if err == nil {
		t.Errorf("Expected error for channelId being out of scope, got %v", err)
	}
}

func testUpdateAs(my ta.Actor) func(t *testing.T) {
	return func(t *testing.T) {
		data := generateTestData()
		vId := data.vFinal.ChannelId()
		request := NewObjectiveRequest(vId)
		getChannel, getConsensusChannel := generateStoreGetters(my.Role, vId, data.vInitial)

		virtualDefund, err := NewObjective(request, false, my.Address(), nil, getChannel, getConsensusChannel)
		if err != nil {
			t.Fatal(err)
		}
		signedFinal := state.NewSignedState(data.vFinal)
		// Sign the final state by some other participant
		signStateByOthers(my, signedFinal)

		e := protocols.CreateObjectivePayload(virtualDefund.Id(), SignedStatePayload, signedFinal)
		updatedObj, err := virtualDefund.Update(e)
		testhelpers.Ok(t, err)
		updated := updatedObj.(*Objective)
		ss, ok := updated.V.SignedStateForTurnNum[FinalTurnNum]
		if !ok {
			t.Fatal(err)
		}

		for _, a := range allActors {
			if a.Role != my.Role {
				testhelpers.Assert(t, !isZero(ss.Signatures()[a.Role]), "expected signature for participant %s to be non-zero", a.Name)
			} else {
				testhelpers.Assert(t, isZero(ss.Signatures()[a.Role]), "expected signature for current participant %s to be zero", a.Name)
			}
		}
	}
}

func testCrankAs(my ta.Actor) func(t *testing.T) {
	return func(t *testing.T) {
		data := generateTestData()
		vId := data.vFinal.ChannelId()
		request := NewObjectiveRequest(vId)

		// If we're Alice we should have the latest payment amount
		// Otherwise we have an older or no payment amount
		ourPaymentAmount := big.NewInt(0)
		if my.Role == 0 {
			ourPaymentAmount = big.NewInt(int64(data.paid))
		}
		getChannel, getConsensusChannel := generateStoreGetters(my.Role, vId, data.vInitial)
		virtualDefund, err := NewObjective(request, true, my.Address(), ourPaymentAmount, getChannel, getConsensusChannel)
		if err != nil {
			t.Fatal(err)
		}

		updatedObj, se, waitingFor, err := virtualDefund.Crank(&my.PrivateKey)
		testhelpers.Ok(t, err)
		updated := updatedObj.(*Objective)

		ss, err := updated.signedFinalState()
		if err != nil {
			t.Fatal(err)
		}

		if my.Role != 0 {
			testhelpers.Equals(t, se.MessagesToSend[0].ObjectivePayloads[0].Type, RequestFinalStatePayload)
			testhelpers.Equals(t, waitingFor, WaitingForFinalStateFromAlice)

			// mimic Alice sending the final state by setting PaidToBob to the paid value
			updated.V.AddSignedState(ss)
			updatedObj, se, waitingFor, err = updated.Crank(&my.PrivateKey)
			testhelpers.Ok(t, err)
			updated = updatedObj.(*Objective)
		}

		if my.Role == alice.Role {
			testhelpers.Equals(t, WaitingForSignedFinal, waitingFor)
		} else {
			testhelpers.Equals(t, WaitingForFinalStateFromAlice, waitingFor)
		}

		signedByMe := state.NewSignedState(data.vFinal)
		testhelpers.SignState(&signedByMe, &my.PrivateKey)

		testhelpers.AssertStateSentToEveryone(t, se, signedByMe, my, allActors)

		// Update the signatures on the objective so the final state is fully signed
		signedByOthers := signStateByOthers(my, state.NewSignedState(data.vFinal))
		for i, sig := range signedByOthers.Signatures() {
			if uint(i) != my.Role {
				err := ss.AddSignature(sig)
				if err != nil {
					panic(err)
				}
			}
		}
		updated.V.AddSignedState(ss)

		updatedObj, se, waitingFor, err = updated.Crank(&my.PrivateKey)
		updated = updatedObj.(*Objective)
		testhelpers.Ok(t, err)

		// We wait for the ledger on our left first by default, unless we have no such channel:
		if my.Role == 0 {
			testhelpers.Equals(t, WaitingForDefundingOnMyRight, waitingFor)
		} else {
			testhelpers.Equals(t, WaitingForDefundingOnMyLeft, waitingFor)
		}

		checkForLeaderProposals(t, se, updated, data)

		proposals := generateProposalsResponses(my.Role, vId, updated, data)
		for _, p := range proposals {

			updatedObj, err = updated.ReceiveProposal(p)
			testhelpers.Ok(t, err)
			updated = updatedObj.(*Objective)
		}

		updatedObj, se, waitingFor, err = updated.Crank(&my.PrivateKey)
		updated = updatedObj.(*Objective)
		testhelpers.Ok(t, err)

		testhelpers.Equals(t, waitingFor, WaitingForNothing)
		checkForFollowerProposals(t, se, updated, data)
	}
}

func TestConstructObjectiveFromState(t *testing.T) {
	data := generateTestData()
	vId := data.vFinal.ChannelId()

	getChannel, getConsensusChannel := generateStoreGetters(alice.Role, vId, data.vInitial)
	signedFinal := state.NewSignedState(data.vFinal)
	// Sign the final state by some other participant
	signStateByOthers(alice, signedFinal)
	b, _ := json.Marshal(signedFinal)
	payload := protocols.ObjectivePayload{Type: SignedStatePayload, PayloadData: b, ObjectiveId: protocols.ObjectiveId(fmt.Sprintf("%s%s", ObjectivePrefix, vId))}
	got, err := ConstructObjectiveFromPayload(payload, true, alice.Address(), getChannel, getConsensusChannel, big.NewInt(int64(data.paid)))
	if err != nil {
		t.Fatal(err)
	}
	left, right := generateLedgers(alice.Role, vId)

	s := state.StateFromFixedAndVariablePart(data.vFinal.FixedPart(), data.vInitial.VariablePart())
	v, err := channel.New(s, 0)
	if err != nil {
		panic(err)
	}

	want := Objective{
		Status:               protocols.Approved,
		V:                    v,
		ToMyLeft:             left,
		ToMyRight:            right,
		MinimumPaymentAmount: big.NewInt(int64(data.paid)),
	}
	if diff := cmp.Diff(want, got, cmp.AllowUnexported(big.Int{}, consensus_channel.ConsensusChannel{}, consensus_channel.LedgerOutcome{}, consensus_channel.Guarantee{})); diff != "" {
		t.Errorf("objective mismatch (-want +got):\n%s", diff)
	}
}

func TestApproveReject(t *testing.T) {
	data := generateTestData()
	vId := data.vFinal.ChannelId()
	request := NewObjectiveRequest(vId)

	getChannel, getConsensusChannel := generateStoreGetters(0, vId, data.vInitial)

	virtualDefund, err := NewObjective(request, false, alice.Address(), nil, getChannel, getConsensusChannel)
	if err != nil {
		t.Fatal(err)
	}
	approved := virtualDefund.Approve()
	if approved.GetStatus() != protocols.Approved {
		t.Errorf("Expected approved status, got %v", approved.GetStatus())
	}
	rejected, sideEffects := virtualDefund.Reject()
	if rejected.GetStatus() != protocols.Rejected {
		t.Errorf("Expected rejceted status, got %v", approved.GetStatus())
	}
	if len(sideEffects.MessagesToSend) != 2 {
		t.Errorf("Expected to send 2 messages")
	}
}
