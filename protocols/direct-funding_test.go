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
		ChannelId: types.Bytes32{},
	}
	_, err := s.Update(e)
	if err == nil {
		t.Error(`ChannelId mismatch -- expected an error but did not get one`)
	}
}
