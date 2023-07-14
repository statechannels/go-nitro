package state

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/abi"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	nc "github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

type Signature = nc.Signature

// CloneSignature creates a deep copy of the provided signature.
func CloneSignature(s Signature) Signature {
	clone := Signature{}

	clone.V = s.V
	clone.R = make([]byte, len(s.R))
	clone.S = make([]byte, len(s.S))

	copy(clone.R, s.R)
	copy(clone.S, s.S)

	return clone
}

type (
	// State holds all of the data describing the state of a channel
	State struct {
		Participants      []types.Address
		ChannelNonce      uint64
		AppDefinition     types.Address
		ChallengeDuration uint32
		AppData           types.Bytes
		Outcome           outcome.Exit
		TurnNum           uint64
		IsFinal           bool
	}

	// FixedPart contains the subset of State data which does not change during a state update.
	FixedPart struct {
		Participants      []types.Address
		ChannelNonce      uint64
		AppDefinition     types.Address
		ChallengeDuration uint32
	}

	// VariablePart contains the subset of State data which can change with each state update.
	VariablePart struct {
		AppData types.Bytes
		Outcome outcome.Exit
		TurnNum uint64
		IsFinal bool
	}
)

// FixedPart returns the FixedPart of the State
func (s State) FixedPart() FixedPart {
	return FixedPart{s.Participants, s.ChannelNonce, s.AppDefinition, s.ChallengeDuration}
}

// VariablePart returns the VariablePart of the State
func (s State) VariablePart() VariablePart {
	return VariablePart{s.AppData, s.Outcome, s.TurnNum, s.IsFinal}
}

// ChannelId computes and returns the channel id corresponding to the State,
// and an error if the id is an external destination.
//
// Up to hash collisions, ChannelId distinguishes channels that have different FixedPart
// values
func (s State) ChannelId() types.Destination {
	return s.FixedPart().ChannelId()
}

func (fp FixedPart) ChannelId() types.Destination {
	encodedChannelPart, err := ethAbi.Arguments{
		{Type: abi.AddressArray},
		{Type: abi.Uint256},
		{Type: abi.Address},
		{Type: abi.Uint256},
	}.Pack(fp.Participants, new(big.Int).SetUint64(fp.ChannelNonce), fp.AppDefinition, new(big.Int).SetUint64(uint64(fp.ChallengeDuration)))
	if err != nil {
		panic(err)
	}

	channelId := types.Destination(crypto.Keccak256Hash(encodedChannelPart))

	return channelId
}

// encodes the state into a []bytes value
func (s State) encode() (types.Bytes, error) {
	ChannelId := s.ChannelId()

	return ethAbi.Arguments{
		{Type: abi.Destination}, // channel id (includes ChainID, Participants, ChannelNonce)
		{Type: abi.Bytes},       // app data
		{Type: outcome.ExitTy},  // outcome
		{Type: abi.Uint256},     // turnNum
		{Type: abi.Bool},        // isFinal
	}.Pack(
		ChannelId,
		[]byte(s.AppData), // Note: even though s.AppData is types.bytes, which is an alias for []byte], Pack will not accept types.bytes
		s.Outcome,
		big.NewInt(int64(s.TurnNum)),
		s.IsFinal,
	)
}

// Hash returns the keccak256 hash of the State
func (s State) Hash() (types.Bytes32, error) {
	encoded, err := s.encode()
	if err != nil {
		return types.Bytes32{}, fmt.Errorf("failed to encode state: %w", err)
	}
	return crypto.Keccak256Hash(encoded), nil
}

// Sign generates an ECDSA signature on the state using the supplied private key
// The state hash is prepended with \x19Ethereum Signed Message:\n32 and then rehashed
// to create a digest to sign
func (s State) Sign(secretKey []byte) (Signature, error) {
	hash, error := s.Hash()
	if error != nil {
		return Signature{}, error
	}
	return nc.SignEthereumMessage(hash.Bytes(), secretKey)
}

// RecoverSigner computes the Ethereum address which generated Signature sig on State state
func (s State) RecoverSigner(sig Signature) (types.Address, error) {
	stateHash, error := s.Hash()
	if error != nil {
		return types.Address{}, error
	}
	return nc.RecoverEthereumMessageSigner(stateHash[:], sig)
}

// equalParticipants returns true if the given arrays contain equal addresses (in the same order).
func equalParticipants(p []types.Address, q []types.Address) bool {
	if len(p) != len(q) {
		return false
	}
	for i, a := range p {
		if !bytes.Equal(a.Bytes(), q[i].Bytes()) {
			return false
		}
	}
	return true
}

// Equal returns true if the given State is deeply equal to the receiever.
func (s State) Equal(r State) bool {
	return equalParticipants(s.Participants, r.Participants) &&
		s.ChannelNonce == r.ChannelNonce &&
		bytes.Equal(s.AppDefinition.Bytes(), r.AppDefinition.Bytes()) &&
		s.ChallengeDuration == r.ChallengeDuration &&
		bytes.Equal(s.AppData, r.AppData) &&
		s.Outcome.Equal(r.Outcome) &&
		s.TurnNum == r.TurnNum &&
		s.IsFinal == r.IsFinal
}

// Clone returns a deep copy of the receiver.
func (f FixedPart) Clone() FixedPart {
	clone := FixedPart{}
	clone.Participants = append(clone.Participants, f.Participants...)
	clone.ChannelNonce = f.ChannelNonce
	clone.AppDefinition = f.AppDefinition
	clone.ChallengeDuration = f.ChallengeDuration
	return clone
}

// Validate checks whether the receiver is malformed and returns an error if it is.
func (fp FixedPart) Validate() error {
	if fp.ChannelId().IsExternal() {
		return errors.New("channelId is an external destination") // This is extremely unlikely
	}

	return nil
}

// Validate checks whether the state is malformed and returns an error if it is.
func (s State) Validate() error {
	return s.FixedPart().Validate()
}

// Clone returns a deep copy of the receiver.
func (s State) Clone() State {
	clone := State{}
	// Fixed part
	cloneFixedPart := s.FixedPart().Clone()
	clone.Participants = cloneFixedPart.Participants
	clone.ChannelNonce = cloneFixedPart.ChannelNonce
	clone.AppDefinition = cloneFixedPart.AppDefinition
	clone.ChallengeDuration = cloneFixedPart.ChallengeDuration

	// Variable part
	if s.AppData != nil {
		clone.AppData = make(types.Bytes, len(s.AppData))
		copy(clone.AppData, s.AppData)
	}
	clone.Outcome = s.Outcome.Clone()
	clone.TurnNum = s.TurnNum
	clone.IsFinal = s.IsFinal

	return clone
}

// StateFromFixedAndVariablePart constructs a State from a FixedPart and a VariablePart
func StateFromFixedAndVariablePart(f FixedPart, v VariablePart) State {
	return State{
		Participants:      f.Participants,
		ChannelNonce:      f.ChannelNonce,
		AppDefinition:     f.AppDefinition,
		ChallengeDuration: f.ChallengeDuration,
		AppData:           v.AppData,
		Outcome:           v.Outcome,
		TurnNum:           v.TurnNum,
		IsFinal:           v.IsFinal,
	}
}
