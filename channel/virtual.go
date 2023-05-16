package channel

import (
	"errors"

	"github.com/statechannels/go-nitro/channel/state"
)

type VirtualChannel struct {
	Channel
}

// NewVirtualChannel returns a new VirtualChannel based on the supplied state.
//
// Virtual channel protocol currently presumes exactly two "active" participants,
// Alice and Bob (p[0] and p[last]). They should be the only destinations allocated
// to in the supplied state's Outcome.
func NewVirtualChannel(s state.State, myIndex uint) (*VirtualChannel, error) {
	if int(myIndex) >= len(s.Participants) {
		return &VirtualChannel{}, errors.New("myIndex not in range of the supplied participants")
	}

	for _, assetExit := range s.Outcome {
		if len(assetExit.Allocations) != 2 {
			return &VirtualChannel{}, errors.New("a virtual channel's initial state should only have two allocations")
		}
	}

	c, err := New(s, myIndex)

	return &VirtualChannel{*c}, err
}

// Clone returns a pointer to a new, deep copy of the receiver, or a nil pointer if the receiver is nil.
func (v *VirtualChannel) Clone() *VirtualChannel {
	if v == nil {
		return nil
	}
	w := VirtualChannel{*v.Channel.Clone()}
	return &w
}

func (v *VirtualChannel) Status() Status {
	s := v.Channel.Status()

	// ADR 0009 allows for intermediaries to exit the protocol before receiving all signed post funds
	// So for intermediaries we return Open once they have signed their post fund state
	amIntermediary := v.MyIndex != 0 && v.MyIndex != uint(len(v.Participants)-1)
	if amIntermediary && v.PostFundSignedByMe() {
		s = Open
	}
	return s
}
