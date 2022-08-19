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
	"github.com/statechannels/go-nitro/types"
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
	vId := data.vFinal.ChannelId()
	request := ObjectiveRequest{
		ChannelId: vId,
	}

	getChannel, getConsensusChannel := generateStoreGetters(0, vId, data.vFinal)

	voucherFetch := func(types.Destination) (payments.Voucher, error) {
		return payments.Voucher{}, nil
	}

	virtualDefund, err := NewObjective(request, false, alice.Address(), getChannel, getConsensusChannel, voucherFetch)
	if err != nil {
		t.Fatal(err)
	}
	invalidFinal := data.vFinal.Clone()
	invalidFinal.ChannelNonce = big.NewInt(5)

	signedFinal := state.NewSignedState(invalidFinal)

	// Sign the final state by some other participant
	signStateByOthers(alice, signedFinal)

	e := protocols.ObjectiveEvent{ObjectiveId: virtualDefund.Id(), SignedState: signedFinal}
	_, err = virtualDefund.Update(e)
	if err.Error() != "event channelId out of scope of objective" {
		t.Errorf("Expected error for channelId being out of scope, got %v", err)
	}

}

func testUpdateAs(my ta.Actor) func(t *testing.T) {
	return func(t *testing.T) {
		data := generateTestData()
		vId := data.vFinal.ChannelId()
		request := ObjectiveRequest{
			ChannelId: vId,
		}

		getChannel, getConsensusChannel := generateStoreGetters(my.Role, vId, data.vInitial)

		voucherFetch := func(types.Destination) (payments.Voucher, error) {
			return payments.Voucher{}, nil
		}
		virtualDefund, err := NewObjective(request, false, my.Address(), getChannel, getConsensusChannel, voucherFetch)
		testhelpers.Ok(t, err)
		from := alice.Address()
		if my.Address() == alice.Address() {
			from = bob.Address()
		}
		e := protocols.ObjectiveEvent{ObjectiveId: virtualDefund.Id(), From: from, Voucher: data.voucher}

		updatedObj, err := virtualDefund.Update(e)
		testhelpers.Ok(t, err)
		updated := updatedObj.(*Objective)
		signedFinal := state.NewSignedState(data.vFinal)
		// Sign the final state by some other participant
		signStateByOthers(my, signedFinal)

		e = protocols.ObjectiveEvent{ObjectiveId: virtualDefund.Id(), SignedState: signedFinal}

		updatedObj, err = updated.Update(e)
		testhelpers.Ok(t, err)
		updated = updatedObj.(*Objective)
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

		voucherFetch := func(types.Destination) (payments.Voucher, error) {
			return payments.Voucher{}, nil
		}
		if my.Role == 0 {
			voucherFetch = func(types.Destination) (payments.Voucher, error) {
				return *aliceVoucher, nil
			}
		}

		getChannel, getConsensusChannel := generateStoreGetters(my.Role, vId, data.vInitial)

		virtualDefund, err := NewObjective(request, true, my.Address(), getChannel, getConsensusChannel, voucherFetch)
		if err != nil {
			t.Fatal(err)

		}
		//
		updatedObj, se, waitingFor, err := virtualDefund.Crank(&my.PrivateKey)
		testhelpers.Ok(t, err)
		updated := updatedObj.(*Objective)

		testhelpers.Equals(t, WaitingForLatestVoucher, waitingFor)
		testhelpers.AssertVoucherSentToEveryone(t, se, updated.Vouchers[updated.MyRole], my, allActors)

		// Set all the vouchers. This mimics all the parties exchanging the latest voucher they have.
		for i := range updated.Vouchers {
			updated.Vouchers[i] = aliceVoucher.Clone()
		}

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
	voucherFetch := func(types.Destination) (payments.Voucher, error) {
		return voucher, nil
	}

	// TODO: Move voucher to data
	got, err := ConstructObjectiveFromVoucher(data.vFinal.FixedPart(), voucher, true, alice.Address(), getChannel, getConsensusChannel, voucherFetch)
	if err != nil {
		t.Fatal(err)
	}
	left, right := generateLedgers(alice.Role, vId)
	want := Objective{
		Status:         protocols.Approved,
		InitialOutcome: data.vInitial.Outcome[0],
		Vouchers:       [3]*payments.Voucher{voucher.Clone()}, // TODO: We expect the largest voucher we have to be there
		VFixed:         data.vFinal.FixedPart(),
		Signatures:     [3]state.Signature{},
		ToMyLeft:       left,
		ToMyRight:      right,
	}
	if diff := cmp.Diff(want, got, cmp.AllowUnexported(big.Int{}, consensus_channel.ConsensusChannel{}, consensus_channel.LedgerOutcome{}, consensus_channel.Guarantee{}, payments.Voucher{}, Objective{})); diff != "" {
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
	voucherFetch := func(types.Destination) (payments.Voucher, error) {
		return payments.Voucher{}, nil
	}
	virtualDefund, err := NewObjective(request, false, alice.Address(), getChannel, getConsensusChannel, voucherFetch)
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
