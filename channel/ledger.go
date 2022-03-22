package channel

import (
	"errors"

	"github.com/statechannels/go-nitro/channel/state"
)

const proposerIndex = uint(0)

type TwoPartyLedger struct {
	Channel
}

// NewTwoPartyLedger constructs a new two party ledger channel from the supplied state.
func NewTwoPartyLedger(s state.State, myIndex uint) (*TwoPartyLedger, error) {
	if myIndex > 1 {
		return &TwoPartyLedger{}, errors.New("myIndex in a two party ledger channel must be 0 or 1")
	}
	if len(s.Participants) != 2 {
		return &TwoPartyLedger{}, errors.New("two party ledger channels must have exactly two participants")
	}

	c, err := New(s, myIndex)

	return &TwoPartyLedger{*c}, err
}

// Clone returns a pointer to a new, deep copy of the receiver, or a nil pointer if the receiver is nil.
func (lc *TwoPartyLedger) Clone() *TwoPartyLedger {
	if lc == nil {
		return nil
	}
	w := TwoPartyLedger{*lc.Channel.Clone()}
	return &w
}

// Proposed returns the latest unsupported ledger state signed by the proposer.
//
// If the latest state signed by the proposer is supported Proposed returns a nil state and false.
func (lc *TwoPartyLedger) Proposed() (state.State, bool) {

	highestSignedByProposer := uint64(0)

	for turnNum, signedState := range lc.SignedStateForTurnNum {
		// todo: consider performance (https://github.com/statechannels/go-nitro/issues/307)
		if signedByProposer := signedState.HasSignatureForParticipant(proposerIndex); signedByProposer && turnNum > highestSignedByProposer {
			highestSignedByProposer = turnNum
		}
	}

	if highestSignedByProposer == lc.latestSupportedStateTurnNum {
		return state.State{}, false
	} else {
		return lc.SignedStateForTurnNum[highestSignedByProposer].State(), true
	}
}

// IsProposer returns true if we are responsible for proposing ledger updates, false otherwise
func (lc *TwoPartyLedger) IsProposer() bool {
	return lc.MyIndex == proposerIndex
}
