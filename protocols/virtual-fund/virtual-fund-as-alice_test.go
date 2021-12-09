package virtualfund

import (
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// By running tests as subtests, we can define local variables that only apply to Alice
func TestAsAlice(t *testing.T) {

	/////////////////////
	// BEGIN test data //
	/////////////////////

	myRole := uint(0) // In this test, we play Alice
	my := Alice

	// Alice plays role 0 so has no ledger channel on her left
	var ledgerChannelToMyLeft channel.Channel

	// She has a single ledger channel L_0 connecting her to P_1
	var ledgerChannelToMyRight, _ = channel.New(
		L_0state,
		true,
		0,
		my.destination,
		P_1.destination,
	)

	// Ensure this channel is fully funded on chain
	ledgerChannelToMyRight.OnChainFunding = ledgerChannelToMyRight.PreFundState().Outcome.TotalAllocated()

	// Objective
	var n = uint(2) // number of ledger channels (num_hops + 1)
	var s, _ = New(VPreFund, my.address, n, myRole, ledgerChannelToMyLeft, ledgerChannelToMyRight)
	var expectedGuaranteeMetadata = outcome.GuaranteeMetadata{Left: ledgerChannelToMyRight.MyDestination, Right: ledgerChannelToMyRight.TheirDestination}
	var expectedEncodedGuaranteeMetadata, _ = expectedGuaranteeMetadata.Encode()
	var expectedGuarantee outcome.Allocation = outcome.Allocation{
		Destination:    s.V.Id,
		Amount:         big.NewInt(0).Set(VPreFund.VariablePart().Outcome[0].TotalAllocated()),
		AllocationType: outcome.GuaranteeAllocationType,
		Metadata:       expectedEncodedGuaranteeMetadata,
	}
	var expectedLedgerRequests = []protocols.LedgerRequest{{
		LedgerId:    ledgerChannelToMyRight.Id,
		Destination: s.V.Id,
		Amount:      types.Funds{types.Address{}: s.V.PreFundState().VariablePart().Outcome[0].Allocations.Total()},
		Left:        ledgerChannelToMyRight.MyDestination, Right: ledgerChannelToMyRight.TheirDestination,
	}}
	// TODO Putting garbage in can result in panics -- we should handle these appropriately by doing input validation
	// var dummySignature = state.Signature{
	// 	R: common.Hex2Bytes(`49d8e91bd182fb4d489bb2d76a6735d494d5bea24e4b51dd95c9d219293312d9`),
	// 	S: common.Hex2Bytes(`22274a3cec23c31e0c073b3c071cf6e0c21260b0d292a10e6a04257a2d8e87fa`),
	// 	V: byte(1),
	// }
	// var dummyState = state.State{}

	var correctSignatureByAliceOnVPreFund, _ = s.V.PreFundState().Sign(Alice.privateKey)
	var correctSignatureByP_1OnVPreFund, _ = s.V.PreFundState().Sign(P_1.privateKey)
	var correctSignatureByBobOnVPreFund, _ = s.V.PreFundState().Sign(Bob.privateKey)

	var correctSignatureByAliceOnVPostFund, _ = s.V.PostFundState().Sign(Alice.privateKey)
	var correctSignatureByP_1OnVPostFund, _ = s.V.PostFundState().Sign(P_1.privateKey)
	var correctSignatureByBobOnVPostFund, _ = s.V.PostFundState().Sign(Bob.privateKey)

	var correctSignatureByAliceOnL_0updatedsate, _ = L_0updatedstate.Sign(Alice.privateKey)
	var correctSignatureByP_1OnL_0updatedsate, _ = L_0updatedstate.Sign(P_1.privateKey)

	///////////////////
	// END test data //
	///////////////////

	testNew := func(t *testing.T) {
		// Assert that a valid set of constructor args does not result in an error
		o, err := New(VPreFund, my.address, 2, myRole, ledgerChannelToMyLeft, ledgerChannelToMyRight)
		if err != nil {
			t.Error(err)
		}

		got := o.ToMyRight.ExpectedGuarantees[types.Address{}] // VState only has one (native) asset represented by the zero address
		want := expectedGuarantee

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("TestNew: expectedGuarantee mismatch (-want +got):\n%s", diff)
		}
	}

	testCrank := func(t *testing.T) {

		// Assert that cranking an unapproved objective returns an error
		if _, _, _, err := s.Crank(&my.privateKey); err == nil {
			t.Error(`Expected error when cranking unapproved objective, but got nil`)
		}

		// Approve the objective, so that the rest of the test cases can run.
		o := s.Approve()

		// To test the finite state progression, we are going to progressively mutate o
		// And then crank it to see which "pause point" (WaitingFor) we end up at.

		// Initial Crank
		_, _, waitingFor, err := o.Crank(&my.privateKey)
		if err != nil {
			t.Error(err)
		}
		if waitingFor != WaitingForCompletePrefund {
			t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompletePrefund, waitingFor)
		}

		// Manually progress the extended state by collecting prefund signatures
		o.(VirtualFundObjective).V.AddSignedState(VPreFund, correctSignatureByAliceOnVPreFund)
		o.(VirtualFundObjective).V.AddSignedState(VPreFund, correctSignatureByBobOnVPreFund)
		o.(VirtualFundObjective).V.AddSignedState(VPreFund, correctSignatureByP_1OnVPreFund)

		// Cranking should move us to the next waiting point, generate ledger requests as a side effect, and alter the extended state to reflect that
		o, sideEffects, waitingFor, err := o.Crank(&my.privateKey)
		if err != nil {
			t.Error(err)
		}
		if waitingFor != WaitingForCompleteFunding {
			t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompleteFunding, waitingFor)
		}
		if o.(VirtualFundObjective).requestedLedgerUpdates != true {
			t.Error(`Expected ledger update idempotency flag to be raised, but it wasn't`)
		}

		got, want := sideEffects.LedgerRequests, expectedLedgerRequests

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("TestCrank: side effects mismatch (-want +got):\n%s", diff)
		}

		// Manually progress the extended state by "completing funding" from this wallet's point of view
		o.(VirtualFundObjective).ToMyRight.Channel.AddSignedState(L_0updatedstate, correctSignatureByAliceOnL_0updatedsate)
		o.(VirtualFundObjective).ToMyRight.Channel.AddSignedState(L_0updatedstate, correctSignatureByP_1OnL_0updatedsate)
		o.(VirtualFundObjective).ToMyRight.Channel.OnChainFunding[types.Address{}] = L_0state.Outcome[0].Allocations.Total() // Make this channel fully funded
		// Cranking now should not generate side effects, because we already did that
		o, _, waitingFor, err = o.Crank(&my.privateKey)
		if err != nil {
			t.Error(err)
		}
		if waitingFor != WaitingForCompletePostFund {
			t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompletePostFund, waitingFor)
		}

		// Manually progress the extended state by collecting postfund signatures
		o.(VirtualFundObjective).V.AddSignedState(o.(VirtualFundObjective).V.PostFundState(), correctSignatureByAliceOnVPostFund)
		o.(VirtualFundObjective).V.AddSignedState(o.(VirtualFundObjective).V.PostFundState(), correctSignatureByBobOnVPostFund)
		o.(VirtualFundObjective).V.AddSignedState(o.(VirtualFundObjective).V.PostFundState(), correctSignatureByP_1OnVPostFund)

		// This should be the final crank...
		_, _, waitingFor, err = o.Crank(&my.privateKey)
		if err != nil {
			t.Error(err)
		}
		if waitingFor != WaitingForNothing {
			t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForNothing, waitingFor)
		}

	}

	testUpdate := func(t *testing.T) {

		// Prepare an event with a mismatched channelId
		e := protocols.ObjectiveEvent{
			ChannelId: types.Destination{},
		}
		// Assert that Updating the objective with such an event returns an error
		// TODO is this the behaviour we want? Below with the signatures, we prefer a log + NOOP (no error)
		if _, err := s.Update(e); err == nil {
			t.Error(`ChannelId mismatch -- expected an error but did not get one`)
		}

		// Now modify the event to give it the "correct" channelId (matching the objective),
		// and make a new Sigs map.
		// This prepares us for the rest of the test. We will reuse the same event multiple times
		e.ChannelId = s.V.Id
		e.Sigs = make(map[*state.State]state.Signature)

		// Next, attempt to update the objective with a dummy signature, keyed with a dummy statehash
		// Assert that this results in a NOOP
		// e.Sigs[&dummyState] = dummySignature // Dummmy signature on dummy statehash
		// if _, err := s.Update(e); err != nil {
		// 	t.Error(`dummy signature -- expected a noop but caught an error:`, err)
		// }

		// Next, attempt to update the objective with an invalid signature, keyed with a dummy statehash
		// Assert that this results in a NOOP
		// e.Sigs[&dummyState] = state.Signature{}
		// if _, err := s.Update(e); err != nil {
		// 	t.Error(`faulty signature -- expected a noop but caught an error:`, err)
		// }

		// Next, attempt to update the objective with correct signature by a participant on a relevant state
		// Assert that this results in an appropriate change in the extended state of the objective
		// Part 1: a signature on a state in channel V
		prefundstate := s.V.PreFundState()
		e.Sigs[&prefundstate] = correctSignatureByAliceOnVPreFund
		updated, err := s.Update(e)
		if err != nil {
			t.Error(err)
		}
		if updated.(VirtualFundObjective).V.PreFundSignedByMe() != true {
			t.Error(`Objective data not updated as expected`)
		}

		// Part 2: a signature on Alice's ledger channel (on her right)
		f := protocols.ObjectiveEvent{
			ChannelId: s.ToMyRight.Channel.Id,
		}
		f.Sigs = make(map[*state.State]state.Signature)
		f.Sigs[&L_0updatedstate] = correctSignatureByAliceOnL_0updatedsate
		updated, err = s.Update(f)
		if err != nil {
			t.Error(err)
		}
		if !updated.(VirtualFundObjective).ToMyRight.ledgerChannelAffordsExpectedGuarantees() != true {
			t.Error(`Objective data not updated as expected`)
		}

	}

	t.Run(`New`, testNew)
	t.Run(`Update`, testUpdate)
	t.Run(`Crank`, testCrank)

}
