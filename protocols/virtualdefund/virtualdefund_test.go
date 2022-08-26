package virtualdefund

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	ta "github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testhelpers"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
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

func testUpdateAs(my ta.Actor) func(t *testing.T) {
	return func(t *testing.T) {
		data := generateTestData()
		vId := data.vFinal.ChannelId()
		request := ObjectiveRequest{
			ChannelId: vId,
		}
		var updated *Objective

		getChannel, getConsensusChannel := generateStoreGetters(my.Role, vId, data.vInitial)
		aliceVoucher, _ := payments.NewSignedVoucher(vId, big.NewInt(int64(data.paid)), alice.PrivateKey)
		oldVoucher, _ := payments.NewSignedVoucher(vId, big.NewInt(int64(data.paid-1)), alice.PrivateKey)

		virtualDefund, err := NewObjective(request, false, my.Address(), aliceVoucher, getChannel, getConsensusChannel)

		testhelpers.Ok(t, err)

		var payload protocols.PayloadValue
		if my.Address() != alice.Address() {
			payload = CreateStartPayload(&virtualDefund, aliceVoucher)

		} else {
			payload = CreateStartPayload(&virtualDefund, oldVoucher)
		}
		updatedObj, err := virtualDefund.UpdateWithPayload(payload)
		testhelpers.Ok(t, err)
		updated = updatedObj.(*Objective)
		signedFinal := state.NewSignedState(data.vFinal)
		// Sign the final state by some other participant
		signStateByOthers(my, signedFinal)

		payload = CreateUpdateSigPayload(*updated, signedFinal.Signatures())
		
		updatedObj, err = virtualDefund.UpdateWithPayload(payload)
		updated = updatedObj.(*Objective)
		testhelpers.Ok(t, err)

		for _, a := range allActors {
			if a.Role != my.Role {
				testhelpers.Assert(t, !isZero(updated.Signatures[a.Role]), "expected signature for participant %s to be non-zero", a.Name)
			} else {
				testhelpers.Assert(t, isZero(updated.Signatures[a.Role]), "expected signature for current participant %s to be zero", a.Name)
			}
		}

	}
}

func testCrankAs(my ta.Actor) func(t *testing.T) {
	return func(t *testing.T) {
		data := generateTestData()
		vId := data.vFinal.ChannelId()
		request := ObjectiveRequest{}

		aliceVoucher, _ := payments.NewSignedVoucher(vId, big.NewInt(int64(data.paid)), alice.PrivateKey)
		oldVoucher, _ := payments.NewSignedVoucher(vId, big.NewInt(int64(data.paid-1)), alice.PrivateKey)

		getChannel, getConsensusChannel := generateStoreGetters(my.Role, vId, data.vInitial)
		var submitVoucher *payments.Voucher
		if my.Address() == alice.Address() {
			submitVoucher = aliceVoucher.Clone()
		} else {
			submitVoucher = oldVoucher
		}
		virtualDefund, err := NewObjective(request, true, my.Address(), submitVoucher, getChannel, getConsensusChannel)
		if err != nil {
			t.Fatal(err)

		}

		updatedObj, se, waitingFor, err := virtualDefund.Crank(&my.PrivateKey)
		testhelpers.Ok(t, err)
		updated := updatedObj.(*Objective)

		testhelpers.Equals(t, WaitingForLatestVoucher, waitingFor)
		testhelpers.AssertVoucherSentToEveryone(t, se, submitVoucher, my, allActors)

		updated.Vouchers[0] = aliceVoucher
		updated.Vouchers[1] = oldVoucher
		updated.Vouchers[2] = oldVoucher

		updatedObj, se, waitingFor, err = updated.Crank(&my.PrivateKey)
		testhelpers.Ok(t, err)
		updated = updatedObj.(*Objective)

		for _, a := range allActors {
			if a.Role == my.Role {
				testhelpers.Assert(t, !isZero(updated.Signatures[a.Role]), "expected signature for participant %s to be non-zero", a.Name)
			} else {
				testhelpers.Assert(t, isZero(updated.Signatures[a.Role]), "expected signature for current participant %s to be zero", a.Name)
			}
		}

		testhelpers.Equals(t, waitingFor, WaitingForCompleteFinal)
		signedByMe := state.NewSignedState(data.vFinal)
		testhelpers.SignState(&signedByMe, &my.PrivateKey)
		testhelpers.AssertStateSentToEveryone(t, se, signedByMe, my, allActors)

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
		for _, p := range proposals {
			e := protocols.ObjectiveEvent{ObjectiveId: updated.Id(), SignedProposal: p}
			updatedObj, err = updated.Update(e)
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
	voucher := *payments.NewVoucher(vId, big.NewInt(int64(data.paid)))
	oldVoucher := *payments.NewVoucher(vId, big.NewInt(int64(data.paid)))

	payload := StartMessage{
		VirtualDefundMessage: protocols.VirtualDefundMessage{
			ChannelId: vId},
		Sender:        alice.Address(),
		LatestVoucher: voucher,
	}
	got, err := ConstructObjectiveFromMessagePayload(true, payload, alice.Address(), &oldVoucher, getChannel, getConsensusChannel)
	if err != nil {
		t.Fatal(err)
	}
	left, right := generateLedgers(alice.Role, vId)
	vouchers := [3]*payments.Voucher{&voucher, nil, nil}
	want := Objective{
		Status:         protocols.Approved,
		InitialOutcome: data.vInitial.Outcome[0],
		Vouchers:       vouchers,
		VFixed:         data.vFinal.FixedPart(),
		Signatures:     [3]state.Signature{},
		ToMyLeft:       left,
		ToMyRight:      right,
	}
	if diff := cmp.Diff(want, *got, cmp.AllowUnexported(big.Int{}, consensus_channel.ConsensusChannel{}, consensus_channel.LedgerOutcome{}, consensus_channel.Guarantee{}, payments.Voucher{}, Objective{})); diff != "" {
		t.Errorf("objective mismatch (-want +got):\n%s", diff)
	}
}

func TestApproveReject(t *testing.T) {
	data := generateTestData()
	vId := data.vFinal.ChannelId()
	request := ObjectiveRequest{
		ChannelId: vId,
	}

	getChannel, getConsensusChannel := generateStoreGetters(0, vId, data.vInitial)
	aliceVoucher, _ := payments.NewSignedVoucher(vId, big.NewInt(int64(data.paid)), alice.PrivateKey)

	virtualDefund, err := NewObjective(request, false, alice.Address(), aliceVoucher, getChannel, getConsensusChannel)
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
