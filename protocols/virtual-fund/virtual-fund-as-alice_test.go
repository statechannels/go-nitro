package virtualfund

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
)

// TestNew tests the constructor using a TestState fixture
func TestNew(t *testing.T) {
	// Assert that a valid set of constructor args does not result in an error
	if _, err := New(state.TestState, state.TestState.Participants[0], 0); err != nil {
		t.Error(err)
	}
}

var s, _ = New(state.TestState, state.TestState.Participants[0], 0)
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

	// Cranking should move us to the next waiting point
	_, _, waitingFor, err = o.Crank(&privateKeyOfParticipant0)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForCompleteFunding {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompleteFunding, waitingFor)
	}

	// Manually progress the extended state by "completing funding" from this wallet's point of view
	// In this test

	// This should be the final crank...
	// TODO manually progress the state by ...
	// if waitingFor != WaitingForNothing {
	// 	t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForNothing, waitingFor)
	// }

	// TODO Test the returned SideEffects
}
