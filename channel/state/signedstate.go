package state

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

type SignedState struct {
	state State
	sigs  map[uint]Signature // keyed by participant index
}

// NewSignedState initializes a SignedState struct for the given
// The signedState returned will have no signatures.
func NewSignedState(s State) SignedState {
	return SignedState{s, make(map[uint]Signature, len(s.Participants))}
}

// AddSignature adds a participant's signature to the SignedState.
//
// An error is returned if
//   - the signer is not a participant, or
//   - OR the signature was already stored
func (ss SignedState) AddSignature(sig Signature) error {
	signer, err := ss.state.RecoverSigner(sig)
	if err != nil {
		return fmt.Errorf("AddSignature failed to recover signer %w", err)
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

// State returns the State part of the SignedState.
func (ss SignedState) State() State {
	return ss.state
}

// Signatures returns a slice of the signatures stored in the SignedState.
// There will be one signature per participant, in order of channel's Participants.
// Returned signatures are expected either to be valid or zero-valued.
func (ss SignedState) Signatures() []Signature {
	sigs := make([]Signature, len(ss.state.Participants))
	for i := 0; i < len(ss.state.Participants); i++ {
		sigs[i] = ss.sigs[uint(i)]
	}
	return sigs
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

// GetParticipantSignature returns the signature for the participant specified by participantIndex
func (ss SignedState) GetParticipantSignature(participantIndex uint) (crypto.Signature, error) {
	sig, found := ss.sigs[uint(participantIndex)]
	if !found {
		return crypto.Signature{}, fmt.Errorf("participant %d does not have a signature", participantIndex)
	} else {
		return sig, nil
	}
}

// Merge checks the passed SignedState's state and the receiver's state for equality, and adds each signature from the former to the latter.
func (ss SignedState) Merge(ss2 SignedState) error {
	if !ss.state.Equal(ss2.state) {
		return errors.New(`cannot merge signed states with distinct state hashes`)
	}
	for i, sig := range ss2.sigs {
		existing, found := ss.sigs[uint(i)]
		if found { // if the signature is already present, check that it is the same
			if !existing.Equal(sig) {
				return errors.New(`cannot merge signed states with conflicting signatures`)
			}
		} else { // otherwise add the signature
			err := ss.AddSignature(sig)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Clone returns a deep copy of the receiver.
func (ss SignedState) Clone() SignedState {
	clonedSigs := make(map[uint]Signature, len(ss.sigs))
	for i, sig := range ss.sigs {
		clonedSigs[i] = sig
	}
	return SignedState{
		state: ss.state.Clone(),
		sigs:  clonedSigs,
	}
}

// MarshalJSON marshals the SignedState into JSON, implementing the Marshaler interface.
func (ss SignedState) MarshalJSON() ([]byte, error) {
	rr := struct {
		State State
		Sigs  map[uint]Signature // keyed by participant index
	}{
		ss.state, ss.sigs,
	}

	return json.Marshal(rr)
}

// UnmarshalJSON unmarshals the passed JSON into a SignedState, implementing the Unmarshaler interface.
func (ss *SignedState) UnmarshalJSON(j []byte) error {
	rr := struct {
		State State
		Sigs  map[uint]Signature // keyed by participant index
	}{}

	err := json.Unmarshal(j, &rr)

	ss.state = rr.State
	ss.sigs = rr.Sigs
	return err
}

// ChannelId returns the channel id of the state.
func (ss SignedState) ChannelId() types.Destination {
	cId := ss.state.ChannelId()
	return cId
}

// SortInfo returns the channel id and turn number of the state, so the state can be easily sorted.
func (ss SignedState) SortInfo() (types.Destination, uint64) {
	cId := ss.State().ChannelId()
	turnNum := ss.State().TurnNum
	return cId, turnNum
}
