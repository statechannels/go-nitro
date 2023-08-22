package channel

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/node/engine/chainservice"
	"github.com/statechannels/go-nitro/types"
)

// Channel contains states and metadata and exposes convenience methods.
type Channel struct {
	Id      types.Destination
	MyIndex uint

	OnChainFunding    types.Funds
	latestBlockNumber uint64 // the latest block number we've seen

	state.FixedPart

	SignedStateForTurnNum map[uint64]state.SignedState
	// Longer term, we should have a more efficient and smart mechanism to store states https://github.com/statechannels/go-nitro/issues/106

	latestSupportedStateTurnNum uint64 // largest uint64 value reserved for "no supported state"
}

// New constructs a new Channel from the supplied state.
func New(s state.State, myIndex uint) (*Channel, error) {
	c := Channel{}
	var err error = s.Validate()

	if err != nil {
		return &c, err
	}

	c.Id = s.ChannelId()

	if err != nil {
		return &c, err
	}
	c.MyIndex = myIndex
	c.OnChainFunding = make(types.Funds)
	c.FixedPart = s.FixedPart().Clone()
	c.latestSupportedStateTurnNum = MaxTurnNum // largest uint64 value reserved for "no supported state"

	// Store prefund
	c.SignedStateForTurnNum = make(map[uint64]state.SignedState)
	c.SignedStateForTurnNum[PreFundTurnNum] = state.NewSignedState(s)

	// Store postfund
	post := s.Clone()
	post.TurnNum = PostFundTurnNum
	c.SignedStateForTurnNum[PostFundTurnNum] = state.NewSignedState(post)

	// Set on chain holdings to zero for each asset
	for asset := range s.Outcome.TotalAllocated() {
		c.OnChainFunding[asset] = big.NewInt(0)
	}

	return &c, nil
}

// jsonChannel replaces Channel's private fields with public ones,
// making it suitable for serialization
type jsonChannel struct {
	Id             types.Destination
	MyIndex        uint
	OnChainFunding types.Funds
	state.FixedPart
	SignedStateForTurnNum map[uint64]state.SignedState

	LatestSupportedStateTurnNum uint64
}

// MarshalJSON returns a JSON representation of the Channel
func (c Channel) MarshalJSON() ([]byte, error) {
	jsonCh := jsonChannel{
		Id:                    c.Id,
		MyIndex:               c.MyIndex,
		OnChainFunding:        c.OnChainFunding,
		FixedPart:             c.FixedPart,
		SignedStateForTurnNum: c.SignedStateForTurnNum,

		LatestSupportedStateTurnNum: c.latestSupportedStateTurnNum,
	}
	return json.Marshal(jsonCh)
}

// UnmarshalJSON populates the calling Channel with the
// json-encoded data
func (c *Channel) UnmarshalJSON(data []byte) error {
	var jsonCh jsonChannel
	err := json.Unmarshal(data, &jsonCh)
	if err != nil {
		return fmt.Errorf("error unmarshaling channel data")
	}

	c.Id = jsonCh.Id
	c.MyIndex = jsonCh.MyIndex
	c.OnChainFunding = jsonCh.OnChainFunding
	c.latestSupportedStateTurnNum = jsonCh.LatestSupportedStateTurnNum
	c.SignedStateForTurnNum = jsonCh.SignedStateForTurnNum

	c.FixedPart = jsonCh.FixedPart

	return nil
}

// MyDestination returns the client's destination
func (c Channel) MyDestination() types.Destination {
	return types.AddressToDestination(c.Participants[c.MyIndex])
}

// Clone returns a pointer to a new, deep copy of the receiver, or a nil pointer if the receiver is nil.
func (c *Channel) Clone() *Channel {
	if c == nil {
		return nil
	}
	d, _ := New(c.PreFundState().Clone(), c.MyIndex)
	d.latestSupportedStateTurnNum = c.latestSupportedStateTurnNum
	for i, ss := range c.SignedStateForTurnNum {
		d.SignedStateForTurnNum[i] = ss.Clone()
	}
	d.OnChainFunding = c.OnChainFunding.Clone()
	d.latestBlockNumber = c.latestBlockNumber
	d.FixedPart = c.FixedPart.Clone()
	return d
}

// PreFundState() returns the pre fund setup state for the channel.
func (c Channel) PreFundState() state.State {
	return c.SignedStateForTurnNum[PreFundTurnNum].State()
}

// SignedPreFundState returns the signed pre fund setup state for the channel.
func (c Channel) SignedPreFundState() state.SignedState {
	return c.SignedStateForTurnNum[PreFundTurnNum]
}

// PostFundState() returns the post fund setup state for the channel.
func (c Channel) PostFundState() state.State {
	return c.SignedStateForTurnNum[PostFundTurnNum].State()
}

// SignedPostFundState() returns the SIGNED post fund setup state for the channel.
func (c Channel) SignedPostFundState() state.SignedState {
	return c.SignedStateForTurnNum[PostFundTurnNum]
}

// PreFundSignedByMe returns true if the calling client has signed the pre fund setup state, false otherwise.
func (c Channel) PreFundSignedByMe() bool {
	if _, ok := c.SignedStateForTurnNum[PreFundTurnNum]; ok {
		if c.SignedStateForTurnNum[PreFundTurnNum].HasSignatureForParticipant(c.MyIndex) {
			return true
		}
	}
	return false
}

// PostFundSignedByMe returns true if the calling client has signed the post fund setup state, false otherwise.
func (c Channel) PostFundSignedByMe() bool {
	if _, ok := c.SignedStateForTurnNum[PostFundTurnNum]; ok {
		if c.SignedStateForTurnNum[PostFundTurnNum].HasSignatureForParticipant(c.MyIndex) {
			return true
		}
	}
	return false
}

// PreFundComplete() returns true if I have a complete set of signatures on  the pre fund setup state, false otherwise.
func (c Channel) PreFundComplete() bool {
	return c.SignedStateForTurnNum[PreFundTurnNum].HasAllSignatures()
}

// PostFundComplete() returns true if I have a complete set of signatures on  the pre fund setup state, false otherwise.
func (c Channel) PostFundComplete() bool {
	return c.SignedStateForTurnNum[PostFundTurnNum].HasAllSignatures()
}

// FinalSignedByMe returns true if the calling client has signed a final state, false otherwise.
func (c Channel) FinalSignedByMe() bool {
	for _, ss := range c.SignedStateForTurnNum {
		if ss.HasSignatureForParticipant(c.MyIndex) && ss.State().IsFinal {
			return true
		}
	}
	return false
}

// FinalCompleted returns true if I have a complete set of signatures on a final state, false otherwise.
func (c Channel) FinalCompleted() bool {
	if c.latestSupportedStateTurnNum == MaxTurnNum {
		return false
	}

	return c.SignedStateForTurnNum[c.latestSupportedStateTurnNum].State().IsFinal
}

// HasSupportedState returns true if the channel has a supported state, false otherwise.
func (c Channel) HasSupportedState() bool {
	return c.latestSupportedStateTurnNum != MaxTurnNum
}

// LatestSupportedState returns the latest supported state. A state is supported if it is signed
// by all participants.
func (c Channel) LatestSupportedState() (state.State, error) {
	if c.latestSupportedStateTurnNum == MaxTurnNum {
		return state.State{}, errors.New(`no state is yet supported`)
	}
	return c.SignedStateForTurnNum[c.latestSupportedStateTurnNum].State(), nil
}

// LatestSignedState fetches the state with the largest turn number signed by at least one participant.
func (c Channel) LatestSignedState() (state.SignedState, error) {
	if len(c.SignedStateForTurnNum) == 0 {
		return state.SignedState{}, errors.New("no states are signed")
	}
	latestTurn := uint64(0)
	for k := range c.SignedStateForTurnNum {
		if k > latestTurn {
			latestTurn = k
		}
	}
	return c.SignedStateForTurnNum[latestTurn], nil
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
	fundingMap types.Funds,
) bool {
	lss, err := c.LatestSupportedState()
	if err != nil {
		return false
	}
	return lss.Outcome.Affords(allocationMap, fundingMap)
}

// AddStateWithSignature constructs a SignedState from the passed state and signature, and calls s.AddSignedState with it.
func (c *Channel) AddStateWithSignature(s state.State, sig state.Signature) bool {
	ss := state.NewSignedState(s)
	if err := ss.AddSignature(sig); err != nil {
		return false
	} else {
		return c.AddSignedState(ss)
	}
}

// AddSignedState adds a signed state to the Channel, updating the LatestSupportedState and Support if appropriate.
// Returns false and does not alter the channel if the state is "stale", belongs to a different channel, or is signed by a non participant.
func (c *Channel) AddSignedState(ss state.SignedState) bool {
	s := ss.State()

	if cId := s.ChannelId(); cId != c.Id {
		// Channel mismatch
		return false
	}

	if c.latestSupportedStateTurnNum != MaxTurnNum && s.TurnNum < c.latestSupportedStateTurnNum {
		// Stale state
		return false
	}

	// Store the signatures. If we have no record yet, add one.
	if signedState, ok := c.SignedStateForTurnNum[s.TurnNum]; !ok {
		c.SignedStateForTurnNum[s.TurnNum] = ss
	} else {
		err := signedState.Merge(ss)
		if err != nil {
			return false
		}
	}

	// Update latest supported state
	if c.SignedStateForTurnNum[s.TurnNum].HasAllSignatures() {
		c.latestSupportedStateTurnNum = s.TurnNum
	}

	return true
}

// SignAndAddPrefund signs and adds the prefund state for the channel, returning a state.SignedState suitable for sending to peers.
func (c *Channel) SignAndAddPrefund(sk *[]byte) (state.SignedState, error) {
	return c.SignAndAddState(c.PreFundState(), sk)
}

// SignAndAddPrefund signs and adds the postfund state for the channel, returning a state.SignedState suitable for sending to peers.
func (c *Channel) SignAndAddPostfund(sk *[]byte) (state.SignedState, error) {
	return c.SignAndAddState(c.PostFundState(), sk)
}

// SignAndAddState signs and adds the state to the channel, returning a state.SignedState suitable for sending to peers.
func (c *Channel) SignAndAddState(s state.State, sk *[]byte) (state.SignedState, error) {
	sig, err := s.Sign(*sk)
	if err != nil {
		return state.SignedState{}, fmt.Errorf("could not sign prefund %w", err)
	}
	ss := state.NewSignedState(s)
	err = ss.AddSignature(sig)
	if err != nil {
		return state.SignedState{}, fmt.Errorf("could not add own signature %w", err)
	}
	ok := c.AddSignedState(ss)
	if !ok {
		return state.SignedState{}, fmt.Errorf("could not add signed state to channel %w", err)
	}
	return ss, nil
}

// UpdateWithChainEvent mutates the receiver if provided with a "new" chain event (with a greater block number than previously seen)
func (c *Channel) UpdateWithChainEvent(event chainservice.Event) (*Channel, error) {
	if event.BlockNum() < c.latestBlockNumber {
		return c, nil // ignore stale information TODO: is this reorg safe?
	}
	c.latestBlockNumber = event.BlockNum()
	switch e := event.(type) {
	case chainservice.AllocationUpdatedEvent:
		c.OnChainFunding[e.AssetAddress] = e.AssetAmount
	case chainservice.DepositedEvent:
		c.OnChainFunding[e.Asset] = e.NowHeld
	case chainservice.ConcludedEvent:
		break
	default:
		return &Channel{}, fmt.Errorf("channel %+v cannot handle event %+v", c, event)
	}
	return c, nil
}
