package channel

import (
	"errors"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

type SingleHopVirtualChannel struct {
	Channel
}

// NewSingleHopVirtualChannel returns a new SingleHopVirtualChannel based on the supplied state.
func NewSingleHopVirtualChannel(s state.State, myIndex uint) (*SingleHopVirtualChannel, error) {
	if myIndex > 2 {
		return &SingleHopVirtualChannel{}, errors.New("myIndex in a single hop virtual channel must be 0, 1, or 2")
	}
	if len(s.Participants) != 3 {
		return &SingleHopVirtualChannel{}, errors.New("a single hop virtual channel must have exactly three participants")
	}
	for _, assetExit := range s.Outcome {
		if len(assetExit.Allocations) != 2 {
			return &SingleHopVirtualChannel{}, errors.New("a single hop virtual channel's initial state should only have two allocations")
		}
	}
	c, err := New(s, myIndex)

	return &SingleHopVirtualChannel{*c}, err
}

// amountAtIndex gets allocations at the specified index and returns the amount.
func (v SingleHopVirtualChannel) amountAtIndex(index uint) types.Funds {
	supported, err := v.LatestSupportedState()

	// If there is no supported state we just return an empty amount
	if err != nil {
		return types.Funds{}
	}

	amount := types.Funds{}

	for _, assetExit := range supported.Outcome {
		asset := assetExit.Asset
		allocations := assetExit.Allocations

		if index < uint(len(allocations)) {
			amount[asset] = allocations[index].Amount
		}
	}
	return amount
}

// LeftAmount returns the amount of the first allocation, which allocates to the left.
func (v SingleHopVirtualChannel) LeftAmount() types.Funds {
	return v.amountAtIndex(0)
}

// RightAmount returns the amount of the second allocation, which allocates to the right.
func (v SingleHopVirtualChannel) RightAmount() types.Funds {
	return v.amountAtIndex(1)
}

// Equal returns true if the supplied SingleHopVirtualChannel is deeply equal to the receiver, false otherwise.
func (v *SingleHopVirtualChannel) Equal(w *SingleHopVirtualChannel) bool {
	return v.Channel.Equal(w.Channel)
}

// Clone returns a pointer to a new, deep copy of the receiver, or a nil pointer if the receiver is nil.
func (v *SingleHopVirtualChannel) Clone() *SingleHopVirtualChannel {
	if v == nil {
		return nil
	}
	w := SingleHopVirtualChannel{*v.Channel.Clone()}
	return &w
}
