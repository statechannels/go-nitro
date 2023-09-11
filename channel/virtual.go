package channel

import (
	"errors"
	"math/big"

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

	c, err := NewChannel(s, myIndex)

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

func (v *VirtualChannel) GetPaidAndRemaining() (*big.Int, *big.Int) {
	remaining := v.OffChain.SignedStateForTurnNum[v.OffChain.LatestSupportedStateTurnNum].State().Outcome[0].Allocations[0].Amount
	paid := v.OffChain.SignedStateForTurnNum[v.OffChain.LatestSupportedStateTurnNum].State().Outcome[0].Allocations[1].Amount

	return paid, remaining
}
