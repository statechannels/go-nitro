package channel

import (
	"errors"
	"fmt"

	"github.com/statechannels/go-nitro/channel/state"
)

type signedState struct {
	state state.State
	sigs  map[uint]state.Signature // keyed by participant index
}

// newSignedState initializes a SignedState struct for the given state.
// The signedState returned will have no signatures.
func newSignedState(s state.State) signedState {
	return signedState{s, make(map[uint]state.Signature, len(s.Participants))}
}

// addSignature adds a participant's signature for the state.
// An error is thrown if the signature is invalid.
func (ss signedState) addSignature(sig state.Signature) error {
	signer, err := ss.state.RecoverSigner(sig)
	if err != nil {
		return fmt.Errorf("addSignature failed to recover signer %w", err)
	}

	for i, p := range ss.state.Participants {
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
func (ss signedState) State() state.State {
	return ss.state
}

// HasSignature returns true if the participant (at participantIndex) has a valid signature.
func (ss signedState) HasSignature(participantIndex uint) bool {
	_, found := ss.sigs[uint(participantIndex)]
	return found
}

// HasAllSignatures returns true if every participant has a valid signature.
func (ss signedState) HasAllSignatures() bool {
	// Since signatures are validated
	if len(ss.sigs) == len(ss.state.Participants) {
		return true
	} else {
		return false
	}
}
