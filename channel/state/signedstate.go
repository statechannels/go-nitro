package state

import (
	"errors"
	"fmt"
)

type SignedState struct {
	state State
	sigs  map[uint]Signature // keyed by participant index
}

// newSignedState initializes a SignedState struct for the given
// The signedState returned will have no signatures.
func NewSignedState(s State) SignedState {
	return SignedState{s, make(map[uint]Signature, len(s.Participants))}
}

// addSignature adds a participant's signature for the
// An error is thrown if the signature is invalid.
func (ss SignedState) AddSignature(sig Signature) error {
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
func (ss SignedState) State() State {
	return ss.state
}

// HasSignatureForParticipant returns true if the participant (at participantIndex) has a valid signature.
func (ss SignedState) HasSignatureForParticipant(participantIndex uint) bool {
	_, found := ss.sigs[uint(participantIndex)]
	return found
}

// HasAllSignatures returns true if every participant has a valid signature.
func (ss SignedState) HasAllSignatures() bool {
	// Since signatures are validated
	if len(ss.sigs) == len(ss.state.Participants) {
		return true
	} else {
		return false
	}
}

// Merge checks the passed SignedState's state and the reciever's state for equality, andd adds each signature from the former to the latter.
func (ss SignedState) Merge(ss2 SignedState) error {
	if !ss.state.Equal(ss2.state) {
		return errors.New(`cannot merge signed states with distinct state hashes`)
	}
	for _, sig := range ss2.sigs {
		err := ss.AddSignature(sig)
		if err != nil {
			return err
		}
	}
	return nil
}
