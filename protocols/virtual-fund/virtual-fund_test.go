package virtualfund

import (
	"testing"

	"github.com/statechannels/go-nitro/channel/state"
)

// TestNew tests the constructor using a TestState fixture
func TestNew(t *testing.T) {
	// Assert that a valid set of constructor args does not result in an error
	if _, err := New(state.TestState, state.TestState.Participants[0], 0); err != nil {
		t.Error(err)
	}
}
