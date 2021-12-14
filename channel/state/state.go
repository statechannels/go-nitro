package state

import (
	"bytes"
	"errors"
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
		TurnNum           *types.Uint256
		IsFinal           bool
	}

	// FixedPart contains the subset of State data which does not change during a state update.
	// NOTE: It is a strict superset of the fields which determine the channel id.
	// It is therefore possible to change some of the fields while preserving said id.
	FixedPart struct {
		ChainId           *types.Uint256
		Participants      []types.Address
		ChannelNonce      *types.Uint256 // uint48 in solidity
		AppDefinition     types.Address  // This could change (infrequently) without affecting the channel id.
		ChallengeDuration *types.Uint256 // This could change (infrequently) without affecting the channel id.
	}

	// VariablePart contains the subset of State data which can change with each state update.
	VariablePart struct {
		AppData types.Bytes
		Outcome outcome.Exit
		TurnNum *types.Uint256
		IsFinal bool
	}
)

// FixedPart returns the FixedPart of the State
func (s State) FixedPart() FixedPart {
	return FixedPart{s.ChainId, s.Participants, s.ChannelNonce, s.AppDefinition, s.ChallengeDuration}
}

// VariablePart returns the VariablePart of the State
func (s State) VariablePart() VariablePart {
	return VariablePart{s.AppData, s.Outcome, s.TurnNum, s.IsFinal}
}

// uint256 is the uint256 type for abi encoding
var uint256, _ = abi.NewType("uint256", "uint256", nil)

// bytesTy is the bytes type for abi encoding
var bytesTy, _ = abi.NewType("bytes", "bytes", nil)

// addressArray is the address[] type for abi encoding
var addressArray, _ = abi.NewType("address[]", "address[]", nil)

// address is the address type for abi encoding
var address, _ = abi.NewType("address", "address", nil)

// ChannelId computes and returns the channel id corresponding to the State,
// and an error if the id is an external destination.
func (s State) ChannelId() (types.Destination, error) {

	encodedChannelPart, error := abi.Arguments{
		{Type: uint256},
		{Type: addressArray},
		{Type: uint256},
	}.Pack(s.ChainId, s.Participants, s.ChannelNonce)

	channelId := types.Destination(crypto.Keccak256Hash(encodedChannelPart))

	if error == nil && channelId.IsExternal() {
		error = errors.New("channelId is an external destination") // This is extremely unlikely
	}
	return channelId, error

}

// appPartHash computes the appPartHash of the State
func (s State) appPartHash() (types.Bytes32, error) {

	encodedAppPart, error := abi.Arguments{
		{Type: uint256},
		{Type: address},
		{Type: bytesTy},
	}.Pack(s.ChallengeDuration, s.AppDefinition, []byte(s.AppData))

	return crypto.Keccak256Hash(encodedAppPart), error

}

// Hash returns the keccak256 hash of the State
func (s State) Hash() (types.Bytes32, error) {

	ChannelId, error := s.ChannelId()
	if error != nil {
		return types.Bytes32{}, error
	}
	OutcomeHash, error := s.Outcome.Hash()
	if error != nil {
		return types.Bytes32{}, error
	}
	AppPartHash, error := s.appPartHash()
	if error != nil {
		return types.Bytes32{}, error
	}

	stateStruct := struct {
		TurnNum     *types.Uint256
		IsFinal     bool
		ChannelId   types.Destination
		AppPartHash types.Bytes32
		OutcomeHash types.Bytes32
	}{
		TurnNum:     s.TurnNum,
		IsFinal:     s.IsFinal,
		ChannelId:   ChannelId,
		AppPartHash: AppPartHash,
		OutcomeHash: OutcomeHash,
	}

	var stateTy, _ = abi.NewType(
		"tuple",
		"struct",
		[]abi.ArgumentMarshaling{
			{Name: "TurnNum", Type: "uint256"},
			{Name: "IsFinal", Type: "bool"},
			{Name: "ChannelId", Type: "bytes32"},
			{Name: "AppPartHash", Type: "bytes32"},
			{Name: "OutcomeHash", Type: "bytes32"},
		},
	)

	encodedState, error := abi.Arguments{{Type: stateTy}}.Pack(stateStruct)
	if error != nil {
		return types.Bytes32{}, error
	}

	return crypto.Keccak256Hash(encodedState), nil
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
		s.TurnNum.Cmp(r.TurnNum) == 0 &&
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
	clone.TurnNum = new(big.Int).Set(s.TurnNum)
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
