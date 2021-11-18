package protocols

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

func TestNew(t *testing.T) {
	_, err := NewDirectFundingObjectiveState(state.TestState, state.TestState.Participants[0])
	if err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	s, _ := NewDirectFundingObjectiveState(state.TestState, state.TestState.Participants[0])
	e := ObjectiveEvent{
		ChannelId: types.Destination{},
	}
	_, err := s.Update(e)
	if err == nil {
		t.Error(`ChannelId mismatch -- expected an error but did not get one`)
	}

	e.ChannelId = s.ChannelId // Fix to correct channelId
	e.Sigs = make(map[types.Bytes32]state.Signature)
	var dummySignature = state.Signature{
		R: common.Hex2Bytes(`49d8e91bd182fb4d489bb2d76a6735d494d5bea24e4b51dd95c9d219293312d9`),
		S: common.Hex2Bytes(`22274a3cec23c31e0c073b3c071cf6e0c21260b0d292a10e6a04257a2d8e87fa`),
		V: byte(1),
	}
	var dummyStateHash = common.Hash{}
	e.Sigs[dummyStateHash] = dummySignature // Dummmy signature on dummy statehash
	_, err = s.Update(e)
	if err != nil {
		t.Error(`dummy signature -- expected a noop but caught an error:`, err)
	}
	e.Sigs[dummyStateHash] = state.Signature{} // Faulty Signature on dummy statehash
	_, err = s.Update(e)
	if err != nil {
		t.Error(`faulty signature -- expected a noop but caught an error:`, err)
	}

	// TODO
	// correct signature by nonparticipant on a dummy state : we want this to be ignored
	// correct signature by participant on a dummy state : we want this to be ignored

	// correct signature by a participant on a relevant state : we ideally want this to be accepted regardless of the enumerable state of the objective !
	stateToSign := s.ExpectedStates[0]
	stateHash, err := stateToSign.Hash()
	if err != nil {
		t.Error(err)
	}
	correctSignatureByParticipant, err := stateToSign.Sign(common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`))
	fmt.Println(correctSignatureByParticipant)
	if err != nil {
		t.Error(err)
	}
	e.Sigs[stateHash] = correctSignatureByParticipant
	updated, err := s.Update(e)
	if err != nil {
		t.Error(err)
	}
	if updated.(DirectFundingObjectiveState).PreFundSigned[0] != true {
		fmt.Printf("%+v\n", updated)
		t.Error(`Objective data not updated as expected`)

	}
}
