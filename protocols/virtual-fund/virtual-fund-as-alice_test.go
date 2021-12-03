package virtualfund

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// Define constants for test
var myRole = uint(0) // In this test, we play Alice
var VirtualChannelState = state.TestState.Clone()
var VId, _ = VirtualChannelState.ChannelId()

var P0P1LedgerChannelState = state.TestState.Clone()

var ledgerChannelToMyLeft channel.Channel // this is null for Alice
var ledgerChannelToMyRight channel.Channel = channel.New(
	P0P1LedgerChannelState, // This should not have Bob in it
	true,
	P0P1LedgerChannelState.Outcome[0].Allocations[0].Destination,
	P0P1LedgerChannelState.Outcome[0].Allocations[1].Destination,
) // this connects Alice to the first intermediary P_1
var expectedEncodedGuaranteeMetadata, _ = outcome.GuaranteeMetadata{Left: ledgerChannelToMyRight.MyDestination, Right: ledgerChannelToMyRight.TheirDestination}.Encode()
var expectedGuarantee outcome.Allocation = outcome.Allocation{
	Destination:    VId,
	Amount:         VirtualChannelState.VariablePart().Outcome[0].TotalAllocated(),
	AllocationType: outcome.GuaranteeAllocationType,
	Metadata:       expectedEncodedGuaranteeMetadata,
}

// TestNew tests the constructor using a TestState fixture
func TestNew(t *testing.T) {
	// Assert that a valid set of constructor args does not result in an error
	o, err := New(state.TestState, state.TestState.Participants[0], myRole, ledgerChannelToMyLeft, ledgerChannelToMyRight)
	if err != nil {
		t.Error(err)
	}

	got := o.ExpectedGuarantees[0][types.Address{}] // The TestOutcome only has one (native) asset represented by the zero address
	want := expectedGuarantee

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestNew: expectedGuarantee mismatch (-want +got):\n%s", diff)
	}
}

var s, _ = New(state.TestState, state.TestState.Participants[0], myRole, ledgerChannelToMyLeft, ledgerChannelToMyRight)
var privateKeyOfParticipant0 = common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`)

func TestCrank(t *testing.T) {
	// Assert that cranking an unapproved objective returns an error
	if _, _, _, err := s.Crank(&privateKeyOfParticipant0); err == nil {
		t.Error(`Expected error when cranking unapproved objective, but got nil`)
	}

	// Approve the objective, so that the rest of the test cases can run.
	o := s.Approve()

	// To test the finite state progression, we are going to progressively mutate o
	// And then crank it to see which "pause point" (WaitingFor) we end up at.

	// Initial Crank
	_, _, waitingFor, err := o.Crank(&privateKeyOfParticipant0)
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
	o, _, waitingFor, err = o.Crank(&privateKeyOfParticipant0)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForCompleteFunding {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompleteFunding, waitingFor)
	}
	if o.(VirtualFundObjective).requestedLedgerUpdates != true { // TODO && sideeffects are as expected
		t.Error(`Expected ledger updates to be requested, but they weren't`)
	}

	// Manually progress the extended state by "completing funding" from this wallet's point of view
	expectedGuarantee := o.(VirtualFundObjective).ExpectedGuarantees[0][types.Address{}] // The TestOutcome only has one (native) asset represented by the zero address
	var UpdatedL0Outcome = outcome.Exit{
		outcome.SingleAssetExit{ // TODO this is not realistic as it does not contain allocations for either Alice (P_0) nor P_1
			Asset: types.Address{},
			Allocations: outcome.Allocations{
				expectedGuarantee,
			},
		},
	}
	var UpdatedL0State = o.(VirtualFundObjective).L[0].LatestSupportedState
	UpdatedL0State.Outcome = UpdatedL0Outcome
	var UpdatedL0Channel = o.(VirtualFundObjective).L[0]
	UpdatedL0Channel.LatestSupportedState = UpdatedL0State
	o.(VirtualFundObjective).L[0] = UpdatedL0Channel
	o.(VirtualFundObjective).L[0].OnChainFunding[types.Address{}] = UpdatedL0Outcome[0].Allocations.Total() // Make this channel fully funded

	// Cranking now should not generate side effects, because we already did that
	o, _, waitingFor, err = o.Crank(&privateKeyOfParticipant0)
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
	_, _, waitingFor, err = o.Crank(&privateKeyOfParticipant0)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForNothing {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForNothing, waitingFor)
	}

}
