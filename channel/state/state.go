package state

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
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

	// FixedPart contains the subset of State data which does not change.
	// NOTE: it is a strict superset of ChannelPart.
	FixedPart struct {
		ChainId           *types.Uint256
		Participants      []types.Address
		ChannelNonce      *types.Uint256 // uint48 in solidity
		AppDefinition     types.Address  // This could change (infrequently) without affecting the channel id.
		ChallengeDuration *types.Uint256 // This could change (infrequently) without affecting the channel id.
	}

	// VariablePart contains the subset of State data which can change with each state update.
	VariablePart struct {
		AppData        types.Bytes
		EncodedOutcome types.Bytes
	}
)

// FixedPart returns the FixedPart of the State
func (s State) FixedPart() FixedPart {
	return FixedPart{s.ChainId, s.Participants, s.ChannelNonce, s.AppDefinition, s.ChallengeDuration}
}

// VariablePart returns the VariablePart of the State
func (s State) VariablePart() VariablePart {
	encodedOutcome, _ := s.Outcome.Encode() // TODO here we are swallowing the error
	return VariablePart{s.AppData, encodedOutcome}
}

var uint256, _ = abi.NewType("uint256", "uint256", nil)
var bytesTy, _ = abi.NewType("bytes", "bytes", nil)
var addressArray, _ = abi.NewType("address[]", "address[]", nil)
var address, _ = abi.NewType("address", "address", nil)

// ChannelId computes and returns the id corresponding to a ChannelPart,
// and an error if the id is an external destination.
func (s State) ChannelId() (types.Bytes32, error) {

	encodedChannelPart, error := abi.Arguments{
		{Type: uint256},
		{Type: addressArray},
		{Type: uint256},
	}.Pack(s.ChainId, s.Participants, s.ChannelNonce)

	channelId := crypto.Keccak256Hash(encodedChannelPart)

	if error == nil && outcome.IsExternalDestination(channelId) {
		error = errors.New("channelId is an external destination") // This is extremely unlikely
	}
	return channelId, error

}

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
		ChannelId   types.Bytes32
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
func (s State) Sign(secretKey []byte) (types.Bytes, error) {
	hash, error := s.Hash()
	if error != nil {
		return types.Bytes{}, error
	}
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n32%s", string(hash.Bytes()))
	modifiedHash := crypto.Keccak256([]byte(msg))
	signature, error := secp256k1.Sign(modifiedHash, secretKey)

	if error != nil {
		return types.Bytes{}, error
	}
	return signature, nil
}

// SplitSignature takes a 65 bytes signature in the [R||S||V] format and returns the individual components
func SplitSignature(signature []byte) (r []byte, s []byte, v byte) {
	r = signature[:32]
	s = signature[32:64]
	v = signature[64]
	return
}
