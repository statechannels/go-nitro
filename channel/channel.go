package channel

import (
	"bytes"
	"errors"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

// Class containing states and metadata, and exposing convenience methods.
type Channel struct {
	Id      types.Destination
	MyIndex uint

	OnChainFunding types.Funds

	state.FixedPart
	// Support []uint64 // TODO: this property will be important, and allow the Channel to store the necessary data to close out the channel on chain
	// It could be an array of turnNums, which can be used to slice into Channel.SignedStateForTurnNum

	latestSupportedStateTurnNum uint64 // largest uint64 value reserved for "no supported state"

	isTwoPartyLedger bool

	SignedStateForTurnNum map[uint64]SignedState // this stores up to 1 state per turn number.
	// Longer term, we should have a more efficient and smart mechanism to store states https://github.com/statechannels/go-nitro/issues/106
}

type TwoPartyLedger struct {
	Channel
}

func NewTwoPartyLedger(s state.State, myIndex uint) (TwoPartyLedger, error) {
	if myIndex > 1 {
		return TwoPartyLedger{}, errors.New("myIndex in a two party ledger channel must be 0 or 1")
	}
	if len(s.Participants) != 2 {
		return TwoPartyLedger{}, errors.New("two party ledger channels must have exactly two participants")
	}

	c, err := New(s, myIndex)
	c.isTwoPartyLedger = true

	return TwoPartyLedger{c}, err
}

func (lc TwoPartyLedger) Clone() TwoPartyLedger {
	return lc // no pointer methods, so this is sufficient
}

// New constructs a new Channel from the supplied state.
func New(s state.State, myIndex uint) (Channel, error) {
	c := Channel{}
	if s.TurnNum.Cmp(big.NewInt(0)) != 0 {
		return c, errors.New(`objective must be constructed with a turnNum 0 state`)
	}

	var err error
	c.Id, err = s.ChannelId()
	if err != nil {
		return c, err
	}
	c.MyIndex = myIndex
	c.OnChainFunding = make(types.Funds)
	c.FixedPart = s.FixedPart()
	c.latestSupportedStateTurnNum = MAGICTURNNUM // largest uint64 value reserved for "no supported state"
	// c.Support =  // TODO
	c.isTwoPartyLedger = false

	// Store prefund
	c.SignedStateForTurnNum = make(map[uint64]SignedState)
	c.SignedStateForTurnNum[0] = SignedState{s.VariablePart(), make(map[uint]state.Signature)}

	// Store postfund
	post := s.Clone()
	post.TurnNum = big.NewInt(1)
	c.SignedStateForTurnNum[1] = SignedState{post.VariablePart(), make(map[uint]state.Signature)}

	return c, nil
}

func (c Channel) IsTwoPartyLedger() bool {
	return c.isTwoPartyLedger
}

// MyDestination returns the client's destination
func (c Channel) MyDestination() types.Destination {
	return types.AddressToDestination(c.Participants[c.MyIndex])
}

// TheirDestination returns the destination of the ledger counterparty
func (lc TwoPartyLedger) TheirDestination() types.Destination {
	return types.AddressToDestination(lc.Participants[(lc.MyIndex+1)%2])
}

// Clone returns a deep copy of the receiver
func (c Channel) Clone() Channel {
	return c // no pointer members, so this is sufficient
}

// Equal returns true if the channel is deeply equal to the reciever, false otherwise
func (c Channel) Equal(d Channel) bool {
	return reflect.DeepEqual(c, d)
}

// PreFundState() returns the pre fund setup state for the channel.
func (c Channel) PreFundState() state.State {
	return state.StateFromFixedAndVariablePart(c.FixedPart, c.SignedStateForTurnNum[0].State)
}

// PostFundState() returns the post fund setup state for the channel.
func (c Channel) PostFundState() state.State {
	return state.StateFromFixedAndVariablePart(c.FixedPart, c.SignedStateForTurnNum[1].State)

}

// PreFundSignedByMe() returns true if I have signed the pre fund setup state, false otherwise.
func (c Channel) PreFundSignedByMe() bool {
	if _, ok := c.SignedStateForTurnNum[0]; ok {
		if _, ok := c.SignedStateForTurnNum[0].Sigs[c.MyIndex]; ok {
			return true
		}
	}
	return false
}

// PostFundSignedByMe() returns true if I have signed the post fund setup state, false otherwise.
func (c Channel) PostFundSignedByMe() bool {
	if _, ok := c.SignedStateForTurnNum[1]; ok {
		if _, ok := c.SignedStateForTurnNum[1].Sigs[c.MyIndex]; ok {
			return true
		}
	}
	return false
}

// PreFundComplete() returns true if I have a complete set of signatures on  the pre fund setup state, false otherwise.
func (c Channel) PreFundComplete() bool {
	return c.SignedStateForTurnNum[0].hasAllSignatures(len(c.FixedPart.Participants))
}

// PostFundComplete() returns true if I have a complete set of signatures on  the pre fund setup state, false otherwise.
func (c Channel) PostFundComplete() bool {
	return c.SignedStateForTurnNum[1].hasAllSignatures(len(c.FixedPart.Participants))
}

// LatestSupportedState returns the latest supported state.
func (c Channel) LatestSupportedState() (state.State, error) {
	if c.latestSupportedStateTurnNum == MAGICTURNNUM {
		return state.State{}, errors.New(`no state is yet supported`)
	}
	return state.StateFromFixedAndVariablePart(c.FixedPart,
		c.SignedStateForTurnNum[c.latestSupportedStateTurnNum].State), nil
}

// Total() returns the total allocated of each asset allocated by the pre fund setup state of the Channel.
func (c Channel) Total() types.Funds {
	return c.PreFundState().Outcome.TotalAllocated()
}

// Affords returns true if, for each asset keying the input variables, the channel can afford the allocation given the funding.
// The decision is made based on the latest supported state of the channel.
//
// Both arguments are maps keyed by the same asset
func (c Channel) Affords(
	allocationMap map[common.Address]outcome.Allocation,
	fundingMap types.Funds) bool {
	lss, err := c.LatestSupportedState()
	if err != nil {
		return false
	}
	return lss.Outcome.Affords(allocationMap, fundingMap)
}

// AddSignedState adds a signed state to the Channel, updating the LatestSupportedState and Support if appropriate.
// Returns false and does not alter the channel if the state is "stale", belongs to a different channel, or is signed by a non participant.
func (c *Channel) AddSignedState(s state.State, sig state.Signature) bool {
	signer, err := s.RecoverSigner(sig)
	if err != nil {
		// Invalid signature
		return false
	}

	signerIndex, isParticipant := indexOf(signer, c.FixedPart.Participants)
	if !isParticipant {
		// Signature by non participant
		return false
	}
	if cId, err := s.ChannelId(); cId != c.Id || err != nil {
		// Channel mismatch
		return false
	}

	turnNum := s.TurnNum.Uint64() // TODO https://github.com/statechannels/go-nitro/issues/95

	if c.latestSupportedStateTurnNum != MAGICTURNNUM && turnNum < c.latestSupportedStateTurnNum {
		// Stale state
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
	if c.SignedStateForTurnNum[turnNum].hasAllSignatures(len(c.FixedPart.Participants)) {
		c.latestSupportedStateTurnNum = turnNum
	}

	// TODO update support

	return true
}

// AddSignedStates adds each signed state in the mapping. It returns true if all signed states were added successfully, false otherwise.
// If one or more signed states fails to be added, this does not prevent other signed states from being added.
func (c Channel) AddSignedStates(mapping map[*state.State]state.Signature) bool {
	allOk := true
	for state, sig := range mapping {
		ok := c.AddSignedState(*state, sig)
		if !ok {
			allOk = false
		}
	}
	return allOk
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
