package channel

import (
	"errors"

	"github.com/statechannels/go-nitro/channel/state"
)

type SingleHopVirtualChannel struct { // todo904: refactor as VirtualChannel - why not?
	Channel
}

// NewSingleHopVirtualChannel returns a new SingleHopVirtualChannel based on the supplied state.
func NewSingleHopVirtualChannel(s state.State, myIndex uint) (*SingleHopVirtualChannel, error) {
	if int(myIndex) >= len(s.Participants) {
		return &SingleHopVirtualChannel{}, errors.New("myIndex is out of range for the supplied participant list")
	}
	// if myIndex > 2 {
	// 	return &SingleHopVirtualChannel{}, errors.New("myIndex in a single hop virtual channel must be 0, 1, or 2")
	// }
	// if len(s.Participants) != 3 {
	// 	return &SingleHopVirtualChannel{}, errors.New("a single hop virtual channel must have exactly three participants")
	// }
	for _, assetExit := range s.Outcome {
		if len(assetExit.Allocations) != 2 {
			return &SingleHopVirtualChannel{}, errors.New("a virtual channel's initial state should only have two allocations")
		}
	}
	c, err := New(s, myIndex)

	return &SingleHopVirtualChannel{*c}, err
}

// Clone returns a pointer to a new, deep copy of the receiver, or a nil pointer if the receiver is nil.
func (v *SingleHopVirtualChannel) Clone() *SingleHopVirtualChannel {
	if v == nil {
		return nil
	}
	w := SingleHopVirtualChannel{*v.Channel.Clone()}
	return &w
}
