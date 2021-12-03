package virtualfund

import (
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel"
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
	var ledgerChannelToMyRight channel.Channel = channel.New(
		L_0state,
		true,
		my.destination,
		P_1.destination,
	)

	// Objective
	var s, _ = New(VState, my.address, myRole, ledgerChannelToMyLeft, ledgerChannelToMyRight)
	var expectedGuaranteeMetadata = outcome.GuaranteeMetadata{Left: ledgerChannelToMyRight.MyDestination, Right: ledgerChannelToMyRight.TheirDestination}
	var expectedEncodedGuaranteeMetadata, _ = expectedGuaranteeMetadata.Encode()
	var expectedGuarantee outcome.Allocation = outcome.Allocation{
		Destination:    s.V.Id,
		Amount:         big.NewInt(0).Set(VState.VariablePart().Outcome[0].TotalAllocated()),
		AllocationType: outcome.GuaranteeAllocationType,
		Metadata:       expectedEncodedGuaranteeMetadata,
	}
	var expectedLedgerRequests = []protocols.LedgerRequest{{
		LedgerId:    ledgerChannelToMyRight.Id,
		Destination: s.V.Id,
		Amount:      types.Funds{types.Address{}: s.V.PreFund.VariablePart().Outcome[0].Allocations.Total()},
		Left:        ledgerChannelToMyRight.MyDestination, Right: ledgerChannelToMyRight.TheirDestination,
	}}

	///////////////////
	// END test data //
	///////////////////

	testNew := func(t *testing.T) {
		// Assert that a valid set of constructor args does not result in an error
		o, err := New(VState, my.address, myRole, ledgerChannelToMyLeft, ledgerChannelToMyRight)
		if err != nil {
			t.Error(err)
		}

		got := o.ExpectedGuarantees[0][types.Address{}] // VState only has one (native) asset represented by the zero address
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
		o.(VirtualFundObjective).preFundSigned[0] = true
		o.(VirtualFundObjective).preFundSigned[1] = true
		o.(VirtualFundObjective).preFundSigned[2] = true

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
		var UpdatedL0Outcome = o.(VirtualFundObjective).L[0].LatestSupportedState.Outcome // TODO clone this?
		UpdatedL0Outcome[0].Allocations, _ = UpdatedL0Outcome[0].Allocations.DivertToGuarantee(my.destination, P_1.destination, s.a0[types.Address{}], s.b0[types.Address{}], s.V.Id)
		var UpdatedL0State = o.(VirtualFundObjective).L[0].LatestSupportedState
		UpdatedL0State.Outcome = UpdatedL0Outcome
		var UpdatedL0Channel = o.(VirtualFundObjective).L[0]
		UpdatedL0Channel.LatestSupportedState = UpdatedL0State
		o.(VirtualFundObjective).L[0] = UpdatedL0Channel
		o.(VirtualFundObjective).L[0].OnChainFunding[types.Address{}] = UpdatedL0Outcome[0].Allocations.Total() // Make this channel fully funded

		// Cranking now should not generate side effects, because we already did that
		o, _, waitingFor, err = o.Crank(&my.privateKey)
		if err != nil {
			t.Error(err)
		}
		if waitingFor != protocols.WaitingForCompletePostFund {
			t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompletePostFund, waitingFor)
		}

		// Manually progress the extended state by collecting postfund signatures
		o.(VirtualFundObjective).postFundSigned[0] = true
		o.(VirtualFundObjective).postFundSigned[1] = true
		o.(VirtualFundObjective).postFundSigned[2] = true

		// This should be the final crank...
		_, _, waitingFor, err = o.Crank(&my.privateKey)
		if err != nil {
			t.Error(err)
		}
		if waitingFor != WaitingForNothing {
			t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForNothing, waitingFor)
		}

	}

	t.Run(`New`, testNew)
	t.Run(`Crank`, testCrank)

}
