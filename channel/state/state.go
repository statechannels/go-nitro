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

	encodedChannelPart, error := ethAbi.Arguments{
		{Type: abi.Uint256},
		{Type: abi.AddressArray},
		{Type: abi.Uint256},
		{Type: abi.Address},
		{Type: abi.Uint256},
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

// Clone returns a deep copy of the receiver.
func (f FixedPart) Clone() FixedPart {
	clone := FixedPart{}
	clone.ChainId = new(big.Int).Set(f.ChainId)
	clone.Participants = append(clone.Participants, f.Participants...)
	clone.ChannelNonce = new(big.Int).Set(f.ChannelNonce)
	clone.AppDefinition = f.AppDefinition
	clone.ChallengeDuration = new(big.Int).Set(f.ChallengeDuration)
	return clone
}

// Clone returns a deep copy of the receiver.
func (s State) Clone() State {
	clone := State{}
	// Fixed part
	cloneFixedPart := s.FixedPart().Clone()
	clone.ChainId = cloneFixedPart.ChainId
	clone.Participants = cloneFixedPart.Participants
	clone.ChannelNonce = cloneFixedPart.ChannelNonce
	clone.AppDefinition = cloneFixedPart.AppDefinition
	clone.ChallengeDuration = cloneFixedPart.ChallengeDuration

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
