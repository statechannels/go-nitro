package protocols

import (
	"testing"

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
	e.Sigs[types.Bytes32(e.ChannelId)] = state.Signature{} // Dummmy signature on dummy statehash
	_, err = s.Update(e)
	if err != nil {
		t.Error(`Useless signature -- expected a noop but caught an error`)
	}

}
