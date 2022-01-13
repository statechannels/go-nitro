package state

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

type (
	// State holds all of the data describing the state of a channel
	State struct {
		ChainId           *types.Uint256
		Participants      []types.Address
		ChannelNonce      *types.Uint256 // uint48 in solidity
		AppDefinition     types.Address
		ChallengeDuration *types.Uint256
		AppData           types.Bytes
		Outcome           outcome.Exit
		TurnNum           uint64
		IsFinal           bool
	}

	// FixedPart contains the subset of State data which does not change during a state update.
	FixedPart struct {
		ChainId           *types.Uint256
		Participants      []types.Address
		ChannelNonce      *types.Uint256 // uint48 in solidity
		AppDefinition     types.Address
		ChallengeDuration *types.Uint256
	}

	// VariablePart contains the subset of State data which can change with each state update.
	VariablePart struct {
		AppData types.Bytes
		Outcome outcome.Exit
		TurnNum uint64
		IsFinal bool
	}
)

/// START: ABI ENCODING HELPERS
// To encode objects as bytes, we need to construct an encoder, using abi.Arguments.
// An instance of abi.Arguments implements two functions relevant to us:
// - `Pack`, which packs go values for a given struct into bytes.
// - `unPack`, which unpacks bytes into go values
// To construct an abi.Arguments instance, we need to supply an array of "types", which are
// actually go values. The following types are used when encoding a state

// uint256 is the uint256 type for abi encoding
var uint256, _ = abi.NewType("uint256", "uint256", nil)

// bool is the bool type for abi encoding
var boolTy, _ = abi.NewType("bool", "bool", nil)

// destination is the bytes32 type for abi encoding
var destination, _ = abi.NewType("bytes32", "address", nil)

// bytes is the bytes type for abi encoding
var bytesTy, _ = abi.NewType("bytes", "bytes", nil)

// address is the address[] type for abi encoding
var addressArray, _ = abi.NewType("address[]", "address[]", nil)

// address is the address type for abi encoding
var address, _ = abi.NewType("address", "address", nil)

/// END: ABI ENCODING HELPERS

// FixedPart returns the FixedPart of the State
func (s State) FixedPart() FixedPart {
	return FixedPart{s.ChainId, s.Participants, s.ChannelNonce, s.AppDefinition, s.ChallengeDuration}
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
func (s State) ChannelId() (types.Destination, error) {
	return s.FixedPart().ChannelId()
}

func (fp FixedPart) ChannelId() (types.Destination, error) {

	if fp.ChainId == nil {
		return types.Destination{}, errors.New(`cannot compute ChannelId with nil ChainId`)
	}

	if fp.ChannelNonce == nil {
		return types.Destination{}, errors.New(`cannot compute ChannelId with nil ChannelNonce`)
	}

	encodedChannelPart, error := abi.Arguments{
		{Type: uint256},
		{Type: addressArray},
		{Type: uint256},
		{Type: address},
		{Type: uint256},
	}.Pack(fp.ChainId, fp.Participants, fp.ChannelNonce, fp.AppDefinition, fp.ChallengeDuration)

	channelId := types.Destination(crypto.Keccak256Hash(encodedChannelPart))

	if error == nil && channelId.IsExternal() {
		error = errors.New("channelId is an external destination") // This is extremely unlikely
	}
	return channelId, error

}

// encodes the state into a []bytes value
func (s State) encode() (types.Bytes, error) {
	ChannelId, error := s.ChannelId()
	if error != nil {
		return types.Bytes{}, fmt.Errorf("failed to construct channelId: %w", error)
	}

	if error != nil {
		return types.Bytes{}, fmt.Errorf("failed to encode outcome: %w", error)

	}

	return abi.Arguments{
		{Type: destination},    // channel id (includes ChainID, Participants, ChannelNonce)
		{Type: bytesTy},        // app data
		{Type: outcome.ExitTy}, // outcome
		{Type: uint256},        // turnNum
		{Type: boolTy},         // isFinal
	}.Pack(
		ChannelId,
		[]byte(s.AppData), // Note: even though s.AppData is types.bytes, which is an alias for []byte], Pack will not accept types.bytes
		s.Outcome,
		s.TurnNum,
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
	return SignEthereumMessage(hash.Bytes(), secretKey)
}

// RecoverSigner computes the Ethereum address which generated Signature sig on State state
func (s State) RecoverSigner(sig Signature) (types.Address, error) {
	stateHash, error := s.Hash()
	if error != nil {
		return types.Address{}, error
	}
	return RecoverEthereumMessageSigner(stateHash[:], sig)
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
	return s.ChainId.Cmp(r.ChainId) == 0 &&
		equalParticipants(s.Participants, r.Participants) &&
		s.ChannelNonce.Cmp(r.ChannelNonce) == 0 &&
		bytes.Equal(s.AppDefinition.Bytes(), r.AppDefinition.Bytes()) &&
		s.ChallengeDuration.Cmp(r.ChallengeDuration) == 0 &&
		bytes.Equal(s.AppData, r.AppData) &&
		s.Outcome.Equal(r.Outcome) &&
		s.TurnNum == r.TurnNum &&
		s.IsFinal == r.IsFinal
}

// Clone returns a clone of the state
func (s State) Clone() State {
	clone := State{}

	// Fixed part
	clone.ChainId = new(big.Int).Set(s.ChainId)
	clone.Participants = append(clone.Participants, s.Participants...)
	clone.ChannelNonce = new(big.Int).Set(s.ChannelNonce)
	clone.AppDefinition = s.AppDefinition
	clone.ChallengeDuration = new(big.Int).Set(s.ChallengeDuration)

	// Variable part
	clone.AppData = make(types.Bytes, 0, len(s.AppData))
	copy(clone.AppData, s.AppData)
	clone.Outcome = s.Outcome.Clone()
	clone.TurnNum = s.TurnNum
	clone.IsFinal = s.IsFinal

	return clone
}

// StateFromFixedAndVariablePart constructs a State from a FixedPart and a VariablePart
func StateFromFixedAndVariablePart(f FixedPart, v VariablePart) State {

	return State{
		ChainId:           f.ChainId,
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
