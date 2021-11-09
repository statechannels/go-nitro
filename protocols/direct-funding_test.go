package protocols

import (
	"fmt"
	"testing"

	"github.com/statechannels/go-nitro/channel/state"
)

func TestNew(t *testing.T) {
	s, _ := NewDirectFundingObjectiveState(state.TestState, state.TestState.Participants[0])
	fmt.Println(s)
}
