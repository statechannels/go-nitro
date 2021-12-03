package directfund

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// TestNew tests the constructor using a TestState fixture
func TestNew(t *testing.T) {
	// Assert that a valid set of constructor args does not result in an error
	if _, err := New(state.TestState, state.TestState.Participants[0]); err != nil {
		t.Error(err)
	}

	// Construct a final state
	finalState := state.TestState.Clone()
	finalState.IsFinal = true

	// Assert that constructing with a final state should return an error
	if _, err := New(finalState, state.TestState.Participants[0]); err == nil {
		t.Error("Expected an error when constructing with an invalid state, but got nil")
	}

}

// Construct various variables for use in TestUpdate
var s, _ = New(state.TestState, state.TestState.Participants[0])
var dummySignature = state.Signature{
	R: common.Hex2Bytes(`49d8e91bd182fb4d489bb2d76a6735d494d5bea24e4b51dd95c9d219293312d9`),
	S: common.Hex2Bytes(`22274a3cec23c31e0c073b3c071cf6e0c21260b0d292a10e6a04257a2d8e87fa`),
	V: byte(1),
}
var dummyStateHash = common.Hash{}
var stateToSign state.State = s.expectedStates[0]
var stateHash, _ = stateToSign.Hash()
var privateKeyOfParticipant0 = common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`)
var correctSignatureByParticipant, _ = stateToSign.Sign(privateKeyOfParticipant0)

func TestUpdate(t *testing.T) {

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
	e.ChannelId = s.channelId
	e.Sigs = make(map[types.Bytes32]state.Signature)

	// Next, attempt to update the objective with a dummy signature, keyed with a dummy statehash
	// Assert that this results in a NOOP
	e.Sigs[dummyStateHash] = dummySignature // Dummmy signature on dummy statehash
	if _, err := s.Update(e); err != nil {
		t.Error(`dummy signature -- expected a noop but caught an error:`, err)
	}

	// Next, attempt to update the objective with an invalid signature, keyed with a dummy statehash
	// Assert that this results in a NOOP
	e.Sigs[dummyStateHash] = state.Signature{}
	if _, err := s.Update(e); err != nil {
		t.Error(`faulty signature -- expected a noop but caught an error:`, err)
	}

	// Next, attempt to update the objective with correct signature by a participant on a relevant state
	// Assert that this results in an appropriate change in the extended state of the objective
	e.Sigs[stateHash] = correctSignatureByParticipant
	updated, err := s.Update(e)
	if err != nil {
		t.Error(err)
	}
	if updated.(DirectFundObjective).preFundSigned[0] != true {
		t.Error(`Objective data not updated as expected`)
	}

	// Finally, add some Holdings information to the event
	// Updating the objective with this event should overwrite the holdings that are stored
	e.Holdings = types.Funds{}
	e.Holdings[common.Address{}] = big.NewInt(3)
	updated, err = s.Update(e)
	if err != nil {
		t.Error(err)
	}
	if !updated.(DirectFundObjective).onChainHolding.Equal(e.Holdings) {
		t.Error(`Objective data not updated as expected`, updated.(DirectFundObjective).onChainHolding, e.Holdings)
	}

}

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
	o.(DirectFundObjective).preFundSigned[0] = true
	o.(DirectFundObjective).preFundSigned[1] = true
	o.(DirectFundObjective).preFundSigned[2] = true

	// Cranking should move us to the next waiting point
	_, _, waitingFor, err = o.Crank(&privateKeyOfParticipant0)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForMyTurnToFund {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForMyTurnToFund, waitingFor)
	}

	// Manually make the first "deposit"
	o.(DirectFundObjective).onChainHolding[state.TestState.Outcome[0].Asset] = state.TestState.Outcome[0].Allocations[0].Amount
	_, _, waitingFor, err = o.Crank(&privateKeyOfParticipant0)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForCompleteFunding {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompleteFunding, waitingFor)
	}

	// Manually make the second "deposit"
	totalAmountAllocated := state.TestState.Outcome[0].TotalAllocated()
	o.(DirectFundObjective).onChainHolding[state.TestState.Outcome[0].Asset] = totalAmountAllocated
	_, _, waitingFor, err = o.Crank(&privateKeyOfParticipant0)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForCompletePostFund {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompletePostFund, waitingFor)
	}

	// Manually progress the extended state by collecting postfund signatures
	o.(DirectFundObjective).postFundSigned[0] = true
	o.(DirectFundObjective).postFundSigned[1] = true
	o.(DirectFundObjective).postFundSigned[2] = true

	// This should be the final crank
	o.(DirectFundObjective).onChainHolding[state.TestState.Outcome[0].Asset] = totalAmountAllocated
	_, _, waitingFor, err = o.Crank(&privateKeyOfParticipant0)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForNothing {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForNothing, waitingFor)
	}

	// TODO Test the returned SideEffects
}
