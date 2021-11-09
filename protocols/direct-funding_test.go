package protocols

import (
	"testing"

	"github.com/statechannels/go-nitro/channel/state"
)

func TestNew(t *testing.T) {
	_, err := NewDirectFundingObjectiveState(state.TestState, state.TestState.Participants[0])
	if err != nil {
		t.Error(err)
	}
}
