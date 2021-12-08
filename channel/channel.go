package channel

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

type SignedState struct {
	State state.VariablePart
	Sigs  map[uint]state.Signature // keyed by participant index
}

// isSupported currently returns true if there are numParticipants distinct signatures on the state and false otherwise.
func (ss SignedState) isSupported(numParticipants int) bool {
	if len(ss.Sigs) == numParticipants {
		return true
	} else {
		return false
	}
}

type Channel struct {
	Id types.Destination

	OnChainFunding types.Funds

	state.FixedPart
	Support              []state.VariablePart
	LatestSupportedState state.State // TODO this should be a variable part, with a getter providing the whole state

	IsTwoPartyLedger bool
	MyDestination    types.Destination
	TheirDestination types.Destination // must be nonzero if a two party ledger channel

	SignedStateForTurnNum map[uint64]SignedState // this stores up to 1 state per turn number.
}

// New constructs a new channel from the supplied state
func New(s state.State, isTwoPartyLedger bool, myDestination types.Destination, theirDestination types.Destination) (Channel, error) {
	c := Channel{}
	if s.TurnNum.Cmp(big.NewInt(0)) != 0 {
		return c, errors.New(`Objective must be constructed with a turnNum 0 state`)
	}

	c.OnChainFunding = make(types.Funds)

	c.LatestSupportedState = s.Clone()
	c.FixedPart = c.LatestSupportedState.FixedPart()

	c.Support = make([]state.VariablePart, 0)
	c.MyDestination = myDestination
	c.TheirDestination = theirDestination
	c.IsTwoPartyLedger = isTwoPartyLedger

	// if s.TurnNum != 0 return error // TODO

	c.Id, _ = s.ChannelId() // TODO handle error

	// Store prefund
	c.SignedStateForTurnNum = make(map[uint64]SignedState)
	c.SignedStateForTurnNum[0] = SignedState{s.VariablePart(), make(map[uint]state.Signature)}

	// Store postfund
	post := s.Clone()
	post.TurnNum = big.NewInt(1)
	c.SignedStateForTurnNum[1] = SignedState{post.VariablePart(), make(map[uint]state.Signature)}

	return c, nil
}

func (c Channel) PreFundState() state.State {
	state := state.State{
		ChainId:           c.ChainId,
		Participants:      c.Participants,
		ChannelNonce:      c.ChannelNonce, // uint48 in solidity
		AppDefinition:     c.AppDefinition,
		ChallengeDuration: c.ChallengeDuration,
		AppData:           c.SignedStateForTurnNum[0].State.AppData,
		Outcome:           c.SignedStateForTurnNum[0].State.Outcome,
		TurnNum:           c.SignedStateForTurnNum[0].State.TurnNum,
		IsFinal:           c.SignedStateForTurnNum[0].State.IsFinal,
	}
	return state
}

func (c Channel) PostFundState() state.State {
	state := state.State{
		ChainId:           c.ChainId,
		Participants:      c.Participants,
		ChannelNonce:      c.ChannelNonce, // uint48 in solidity
		AppDefinition:     c.AppDefinition,
		ChallengeDuration: c.ChallengeDuration,
		AppData:           c.SignedStateForTurnNum[1].State.AppData,
		Outcome:           c.SignedStateForTurnNum[1].State.Outcome,
		TurnNum:           c.SignedStateForTurnNum[1].State.TurnNum,
		IsFinal:           c.SignedStateForTurnNum[1].State.IsFinal,
	}
	return state
}

func (c Channel) PreFundSignedByMe() bool {
	myIndex := uint(0) // TODO get this from the channel
	if _, ok := c.SignedStateForTurnNum[0]; ok {
		if _, ok := c.SignedStateForTurnNum[0].Sigs[myIndex]; ok {
			return true
		}
	}
	return false
}
func (c Channel) PostFundSignedByMe() bool {
	myIndex := uint(0) // TODO get this from the channel
	if _, ok := c.SignedStateForTurnNum[1]; ok {
		if _, ok := c.SignedStateForTurnNum[1].Sigs[myIndex]; ok {
			return true
		}
	}
	return false
}
func (c Channel) PreFundComplete() bool {

	return c.SignedStateForTurnNum[0].isSupported(len(c.FixedPart.Participants))

}
func (c Channel) PostFundComplete() bool {
	return c.SignedStateForTurnNum[1].isSupported(len(c.FixedPart.Participants))
}

func (c Channel) Total() types.Funds {
	funds := types.Funds{}
	for _, sae := range c.LatestSupportedState.Outcome {
		funds[sae.Asset] = sae.Allocations.Total()
	}
	return funds
}

// Affords returns true if, for each asset keying the input variables, the channel can afford the allocation given the funding.
// The decision is made based on the latest supported state of the channel.
//
// Both arguments are maps keyed by the same asset
func (c Channel) Affords(
	allocationMap map[common.Address]outcome.Allocation,
	fundingMap types.Funds) bool {
	return c.LatestSupportedState.Outcome.Affords(allocationMap, fundingMap)
}

// AddSignedState adds a signed state to the Channel, updating the LatestSupportedState and Support if appropriate.
// Returns false and does not alter the channel if the state is "stale", belongs to a different channel, or is signed by a non participant
func (c *Channel) AddSignedState(s state.State, sig state.Signature) bool {
	signer, _ := s.RecoverSigner(sig) // TODO handle error
	signerIndex, isParticipant := indexOf(signer, c.FixedPart.Participants)
	turnNum := s.TurnNum.Uint64() // https://github.com/statechannels/go-nitro/issues/95

	if cId, err := s.ChannelId(); cId != c.Id || err != nil || turnNum < c.LatestSupportedState.TurnNum.Uint64() || !isParticipant {
		return false
	}

	// Store the signature. If we have no record yet, add one.
	if signedState, ok := c.SignedStateForTurnNum[turnNum]; !ok {
		c.SignedStateForTurnNum[turnNum] = SignedState{s.VariablePart(), make(map[uint]state.Signature)}
		c.SignedStateForTurnNum[turnNum].Sigs[signerIndex] = sig
	} else {
		signedState.Sigs[signerIndex] = sig
	}

	// Update latest supported state
	if c.SignedStateForTurnNum[turnNum].isSupported(len(c.FixedPart.Participants)) {
		c.LatestSupportedState = s
	}

	// TODO update support

	return true
}

// AddSignedStates adds each signed state in the mapping
func (c Channel) AddSignedStates(mapping map[*state.State]state.Signature) {
	for state, sig := range mapping {
		c.AddSignedState(*state, sig)
	}
}

// indexOf returns the index of the given suspect address in the lineup of addresses. A second return value ("ok") is true the suspect was found, false otherwise.
func indexOf(suspect types.Address, lineup []types.Address) (index uint, ok bool) {

	for index, a := range lineup {
		if bytes.Equal(suspect.Bytes(), a.Bytes()) {
			return uint(index), true
		}
	}
	return ^uint(0), false
}
