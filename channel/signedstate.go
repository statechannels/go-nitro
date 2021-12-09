package channel

import (
	"github.com/statechannels/go-nitro/channel/state"
)

type SignedState struct {
	State state.VariablePart
	Sigs  map[uint]state.Signature // keyed by participant index
}

// hasAllSignatures returns true if there are numParticipants distinct signatures on the state and false otherwise.
func (ss SignedState) hasAllSignatures(numParticipants int) bool {
	if len(ss.Sigs) == numParticipants {
		return true
	} else {
		return false
	}
}
