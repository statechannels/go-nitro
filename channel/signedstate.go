package channel

import (
	"errors"

	"github.com/statechannels/go-nitro/channel/state"
)

type SignedState struct {
	State state.State
	sigs  map[uint]state.Signature // keyed by participant index
}

// NewSignedState creates a signed state for the given state with 0 signatures.
func NewSignedState(s state.State) SignedState {
	ss := SignedState{s, make(map[uint]state.Signature, len(s.Participants))}

	return ss
}

// AddSignature adds a participant's signature for the state.
// An error is thrown if the signature is invalid.
func (ss SignedState) AddSignature(sig state.Signature) error {
	signer, err := ss.State.RecoverSigner(sig)
	if err != nil {
		return err
	}

	for i, p := range ss.State.Participants {
		if p == signer {
			_, found := ss.sigs[uint(i)]
			if found {
				return errors.New("signature already exists for participant")
			} else {
				ss.sigs[uint(i)] = sig
				return nil
			}

		}

	}
	return errors.New("signature does not match any participant")

}

// HasSignature returns true if the participant (at participantIndex) has a valid signature.
func (ss SignedState) HasSignature(participantIndex uint) bool {
	_, found := ss.sigs[uint(participantIndex)]
	return found
}

// HasAllSignatures returns true if every participant has a valid signature.
func (ss SignedState) HasAllSignatures() bool {
	if len(ss.sigs) == len(ss.State.Participants) {
		return true
	} else {
		return false
	}
}
