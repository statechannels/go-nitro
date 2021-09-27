package state

import (
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
		ChallengeDuration *uint
		AppData           types.Bytes
		Outcome           outcome.Exit
		TurnNum           *uint
		IsFinal           bool
	}

	// ChannelPart contains the subset of State data from which the channel id is derived
	ChannelPart struct {
		ChainId      *types.Uint256
		Participants []types.Address
		ChannelNonce *types.Uint256 // uint48 in solidity
	}

	// FixedPart contains the subset of State data which does not change.
	// NOTE: it is a strict superset of ChannelPart.
	FixedPart struct {
		ChainId           *types.Uint256
		Participants      []types.Address
		ChannelNonce      *types.Uint256 // uint48 in solidity
		AppDefinition     types.Address  // This could change (infrequently) without affecting the channel id.
		ChallengeDuration *uint          // This could change (infrequently) without affecting the channel id.
	}

	// VariablePart contains the subset of State data which can change with each state update.
	VariablePart struct {
		AppData        types.Bytes
		EncodedOutcome types.Bytes
	}
)

// ChannelPart returns the ChannelPart of the State
func (s State) ChannelPart() ChannelPart {
	return ChannelPart{s.ChainId, s.Participants, s.ChannelNonce}
}

// FixedPart returns the FixedPart of the State
func (s State) FixedPart() FixedPart {
	return FixedPart{s.ChainId, s.Participants, s.ChannelNonce, s.AppDefinition, s.ChallengeDuration}
}

// VariablePart returns the VariablePart of the State
func (s State) VariablePart() VariablePart {
	encodedOutcome, _ := s.Outcome.Encode() // TODO here we are swallowing the error
	return VariablePart{s.AppData, encodedOutcome}
}

// ChannelId computes and returns the id corresponding to a ChannelPart
func (c ChannelPart) ChannelId() (types.Bytes32, error) {
	uint256, _ := abi.NewType("uint256", "uint256", nil)
	addressArray, _ := abi.NewType("address[]", "address[]", nil)
	encodedChannelPart, error := abi.Arguments{
		{Type: uint256},
		{Type: addressArray},
		{Type: uint256},
	}.Pack(c.ChainId, c.Participants, c.ChannelNonce)

	// TODO return an error if the channelId is an external destination

	return crypto.Keccak256Hash(encodedChannelPart), error

}

// ChannelId computes and returns the id corresponding to a State
func (s State) ChannelId() (types.Bytes32, error) {
	return s.ChannelPart().ChannelId()
}

// TODO hashAppPart
// TODO hashState
