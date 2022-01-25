package channel

import (
	"errors"

	"github.com/statechannels/go-nitro/channel/state"
)

type SignedState struct {
	State state.State
	sigs  map[uint]state.Signature // keyed by participant index
}

// NewSignedState creates a new SignedState from the supplied state and signatures.
// An error returned if there is an invalid signature.
func NewSignedState(s state.State, sigs []state.Signature) (SignedState, error) {
	ss := SignedState{s, make(map[uint]state.Signature, len(sigs))}
	err := ss.AddSignatures(sigs)
	return ss, err
}

/// AddSignatures adds multiple participant's signature for the state.
/// An error is thrown if any signature is invalid.
func (ss SignedState) AddSignatures(sigs []state.Signature) (err error) {
	for _, sig := range sigs {
		err = ss.AddSignature(sig)
	}
	return
}

// AddSignature adds a participant's signature for the state.
// An error is thrown if the signature is invalid.
func (ss SignedState) AddSignature(sig state.Signature) (err error) {

	signer, err := ss.State.RecoverSigner(sig)

	for i, p := range ss.State.Participants {
		if p == signer {
			_, found := ss.sigs[uint(i)]
			if found {
				err = errors.New("signature already exists for participant")
			} else {
				ss.sigs[uint(i)] = sig
			}
			return
		}

	}
	return err
}

// HasSignature returns true if the participant has a valid signature.
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
