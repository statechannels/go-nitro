package state

import (
	"bytes"
	"errors"
	"fmt"
)

type SignedState struct {
	State State
	Sigs  map[uint]Signature // keyed by participant index
}

// newSignedState initializes a SignedState struct for the given
// The signedState returned will have no signatures.
func NewSignedState(s State) SignedState {
	return SignedState{s, make(map[uint]Signature, len(s.Participants))}
}

// addSignature adds a participant's signature for the
// An error is thrown if the signature is invalid.
func (ss SignedState) AddSignature(sig Signature) error {
	signer, err := ss.State.RecoverSigner(sig)
	if err != nil {
		return fmt.Errorf("addSignature failed to recover signer %w", err)
	}

	for i, p := range ss.State.Participants {
		if p == signer {
			_, found := ss.Sigs[uint(i)]
			if found {
				return errors.New("signature already exists for participant")
			} else {
				ss.Sigs[uint(i)] = sig
				return nil
			}

		}

	}
	return errors.New("signature does not match any participant")

}

// HasSignatureForParticipant returns true if the participant (at participantIndex) has a valid signature.
func (ss SignedState) HasSignatureForParticipant(participantIndex uint) bool {
	_, found := ss.Sigs[uint(participantIndex)]
	return found
}

// HasAllSignatures returns true if every participant has a valid signature.
func (ss SignedState) HasAllSignatures() bool {
	// Since signatures are validated
	if len(ss.Sigs) == len(ss.State.Participants) {
		return true
	} else {
		return false
	}
}

// Merge checks the passed SignedState's state and the reciever's state for equality, andd adds each signature from the former to the latter.
func (ss SignedState) Merge(ss2 SignedState) error {
	if !ss.State.Equal(ss2.State) {
		return errors.New(`cannot merge signed states with distinct state hashes`)
	}
	for _, sig := range ss2.Sigs {
		err := ss.AddSignature(sig)
		if err != nil {
			return err
		}
	}
	return nil
}

// Equal returns true if the passed SignedState is deeply equal in value to the reciever.
func (ss SignedState) Equal(ss2 SignedState) bool {
	if !ss.State.Equal(ss2.State) {
		return false
	}
	if len(ss.Sigs) != len(ss2.Sigs) {
		return false
	}
	for i, sig := range ss.Sigs {
		if !bytes.Equal(sig.S, ss2.Sigs[i].S) || !bytes.Equal(sig.R, ss2.Sigs[i].R) || sig.V != ss2.Sigs[i].V {
			return false
		}
	}
	return true
}
