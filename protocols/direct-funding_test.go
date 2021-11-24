package protocols

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

// TestNew tests the constructor using a TestState fixture
func TestNew(t *testing.T) {
	if _, err := NewDirectFundingObjectiveState(state.TestState, state.TestState.Participants[0]); err != nil {
		t.Error(err)
	}
	invalidState := state.TestState.Clone()
	tooLargeChainId, _ := big.NewInt(0).SetString("49d8e91bd182fb4d489bb2d76a6735d494d5bea24e4b51dd95c9d219293312d949d8e91bd182fb4d489bb2d76a6735d494d5bea24e4b51dd95c9d219293312d9", 16) // This is 128 hexits = 64 bytes = 512 bits. A chainId should be 256 bit
	invalidState.ChainId = tooLargeChainId
	if _, err := NewDirectFundingObjectiveState(invalidState, state.TestState.Participants[0]); err == nil {
		t.Error("Expected an error when constructing with an invalid state, but got nil")
	}

}

// Construct various variables for use in TestUpdate
var s, _ = NewDirectFundingObjectiveState(state.TestState, state.TestState.Participants[0])
var dummySignature = state.Signature{
	R: common.Hex2Bytes(`49d8e91bd182fb4d489bb2d76a6735d494d5bea24e4b51dd95c9d219293312d9`),
	S: common.Hex2Bytes(`22274a3cec23c31e0c073b3c071cf6e0c21260b0d292a10e6a04257a2d8e87fa`),
	V: byte(1),
}
var dummyStateHash = common.Hash{}
var stateToSign state.State = s.ExpectedStates[0]
var stateHash, _ = stateToSign.Hash()
var privateKeyOfParticipant0 = common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`)
var correctSignatureByParticipant, _ = stateToSign.Sign(privateKeyOfParticipant0)

func TestUpdate(t *testing.T) {

	// Prepare an event with a mismatched channelId
	e := ObjectiveEvent{
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
	e.ChannelId = s.ChannelId
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
	if updated.(DirectFundingObjectiveState).PreFundSigned[0] != true {
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
	if !updated.(DirectFundingObjectiveState).OnChainHolding.Equal(e.Holdings) {
		t.Error(`Objective data not updated as expected`, updated.(DirectFundingObjectiveState).OnChainHolding, e.Holdings)
	}

}
