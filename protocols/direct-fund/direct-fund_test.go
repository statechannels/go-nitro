package directfund

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type actor struct {
	address     types.Address
	destination types.Destination
	privateKey  []byte
}

var alice = actor{
	address:     common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`),
	destination: types.AddressToDestination(common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`)),
	privateKey:  common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`),
}

var bob = actor{
	address:     common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`),
	destination: types.AddressToDestination(common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`)),
	privateKey:  common.Hex2Bytes(`62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2`),
}

var testState = state.State{
	ChainId:           big.NewInt(9001),
	Participants:      []types.Address{alice.address, bob.address},
	ChannelNonce:      big.NewInt(37140676580),
	AppDefinition:     common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`),
	ChallengeDuration: big.NewInt(60),
	AppData:           []byte{},
	Outcome: outcome.Exit{
		outcome.SingleAssetExit{
			Asset: types.Address{},
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: alice.destination,
					Amount:      big.NewInt(5),
				},
				outcome.Allocation{
					Destination: bob.destination,
					Amount:      big.NewInt(5),
				},
			},
		},
	},
	TurnNum: 0,
	IsFinal: false,
}

// TestNew tests the constructor using a TestState fixture
func TestNew(t *testing.T) {
	// Assert that valid constructor args do not result in error
	if _, err := New(false, testState, testState.Participants[0]); err != nil {
		t.Error(err)
	}

	// Construct a final state
	finalState := testState.Clone()
	finalState.IsFinal = true

	if _, err := New(false, finalState, testState.Participants[0]); err == nil {
		t.Error("expected an error when constructing with an intial state marked final, but got nil")
	}

	nonParticipant := common.HexToAddress("0x5b53f71453aeCb03D837bfe170570d40aE736CB4")
	if _, err := New(false, testState, nonParticipant); err == nil {
		t.Error("expected an error when constructing with a participant not in the channel, but got nil")
	}
}

// Construct various variables for use in TestUpdate
var s, _ = New(false, testState, testState.Participants[0])
var dummySignature = state.Signature{
	R: common.Hex2Bytes(`49d8e91bd182fb4d489bb2d76a6735d494d5bea24e4b51dd95c9d219293312d9`),
	S: common.Hex2Bytes(`22274a3cec23c31e0c073b3c071cf6e0c21260b0d292a10e6a04257a2d8e87fa`),
	V: byte(1),
}
var dummyState = state.State{}
var stateToSign state.State = s.C.PreFundState()
var correctSignatureByParticipant, _ = stateToSign.Sign(alice.privateKey)

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
	e.ChannelId = s.C.Id
	e.Sigs = make(map[*state.State]state.Signature)

	// Next, attempt to update the objective with a dummy signature, keyed with a dummy state
	// Assert that this results in a NOOP
	e.Sigs[&dummyState] = dummySignature // Dummmy signature on dummy state
	if _, err := s.Update(e); err != nil {
		t.Error(`dummy signature -- expected a noop but caught an error:`, err)
	}

	// Next, attempt to update the objective with an invalid signature, keyed with a dummy statehash
	// Assert that this results in a NOOP
	e.Sigs[&dummyState] = state.Signature{}
	if _, err := s.Update(e); err != nil {
		t.Error(`faulty signature -- expected a noop but caught an error:`, err)
	}

	// Next, attempt to update the objective with correct signature by a participant on a relevant state
	// Assert that this results in an appropriate change in the extended state of the objective
	e.Sigs[&stateToSign] = correctSignatureByParticipant
	updatedObjective, err := s.Update(e)
	if err != nil {
		t.Error(err)
	}
	updated := updatedObjective.(DirectFundObjective)
	if updated.C.PreFundSignedByMe() != true {
		t.Error(`Objective data not updated as expected`)
	}

	// Finally, add some Holdings information to the event
	// Updating the objective with this event should overwrite the holdings that are stored
	e.Holdings = types.Funds{}
	e.Holdings[common.Address{}] = big.NewInt(3)
	updatedObjective, err = s.Update(e)
	if err != nil {
		t.Error(err)
	}
	updated = updatedObjective.(DirectFundObjective)
	if !updated.C.OnChainFunding.Equal(e.Holdings) {
		t.Error(`Objective data not updated as expected`, updated.C.OnChainFunding, e.Holdings)
	}

}

func TestCrank(t *testing.T) {

	var correctSignatureByAliceOnPreFund, _ = s.C.PreFundState().Sign(alice.privateKey)
	var correctSignatureByBobOnPreFund, _ = s.C.PreFundState().Sign(bob.privateKey)

	var correctSignatureByAliceOnPostFund, _ = s.C.PostFundState().Sign(alice.privateKey)
	var correctSignatureByBobOnPostFund, _ = s.C.PostFundState().Sign(bob.privateKey)

	// Assert that cranking an unapproved objective returns an error
	if _, _, _, err := s.Crank(&alice.privateKey); err == nil {
		t.Error(`Expected error when cranking unapproved objective, but got nil`)
	}

	// Approve the objective, so that the rest of the test cases can run.
	o := s.Approve().(DirectFundObjective)

	// To test the finite state progression, we are going to progressively mutate o
	// And then crank it to see which "pause point" (WaitingFor) we end up at.

	// Initial Crank
	_, _, waitingFor, err := o.Crank(&alice.privateKey)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForCompletePrefund {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompletePrefund, waitingFor)
	}

	// Manually progress the extended state by collecting prefund signatures
	o.C.AddSignedState(o.C.PreFundState(), correctSignatureByAliceOnPreFund)
	o.C.AddSignedState(o.C.PreFundState(), correctSignatureByBobOnPreFund)

	// Cranking should move us to the next waiting point
	_, _, waitingFor, err = o.Crank(&alice.privateKey)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForMyTurnToFund {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForMyTurnToFund, waitingFor)
	}

	// Manually make the first "deposit"
	o.C.OnChainFunding[testState.Outcome[0].Asset] = testState.Outcome[0].Allocations[0].Amount
	_, _, waitingFor, err = o.Crank(&alice.privateKey)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForCompleteFunding {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompleteFunding, waitingFor)
	}

	// Manually make the second "deposit"
	totalAmountAllocated := testState.Outcome[0].TotalAllocated()
	o.C.OnChainFunding[testState.Outcome[0].Asset] = totalAmountAllocated
	_, _, waitingFor, err = o.Crank(&alice.privateKey)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForCompletePostFund {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompletePostFund, waitingFor)
	}

	// Manually progress the extended state by collecting postfund signatures
	o.C.AddSignedState(o.C.PostFundState(), correctSignatureByAliceOnPostFund)
	o.C.AddSignedState(o.C.PostFundState(), correctSignatureByBobOnPostFund)

	// This should be the final crank
	o.C.OnChainFunding[testState.Outcome[0].Asset] = totalAmountAllocated
	_, _, waitingFor, err = o.Crank(&alice.privateKey)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForNothing {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForNothing, waitingFor)
	}

	// TODO Test the returned SideEffects
}
