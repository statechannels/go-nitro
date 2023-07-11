// Package directdefund implements an off-chain protocol to defund a directly-funded channel.
package directdefund // import "github.com/statechannels/go-nitro/directfund"

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/node/engine/chainservice"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const (
	WaitingForFinalization protocols.WaitingFor = "WaitingForFinalization"
	WaitingForWithdraw     protocols.WaitingFor = "WaitingForWithdraw"
	WaitingForNothing      protocols.WaitingFor = "WaitingForNothing" // Finished
)

const (
	SignedStatePayload protocols.PayloadType = "SignedStatePayload"
)

const ObjectivePrefix = "DirectDefunding-"

const (
	ErrChannelUpdateInProgress = types.ConstError("can only defund a channel when the latest state is supported or when the channel has a final state")
	ErrNoFinalState            = types.ConstError("cannot spawn direct defund objective without a final state")
	ErrNotEmpty                = types.ConstError("ledger channel has running guarantees")
)

// Objective is a cache of data computed by reading from the store. It stores (potentially) infinite data
type Objective struct {
	Status               protocols.ObjectiveStatus
	C                    *channel.Channel
	finalTurnNum         uint64
	transactionSubmitted bool // whether a transition for the objective has been submitted or not
}

// isInConsensusOrFinalState returns true if the channel has a final state or latest state that is supported
func isInConsensusOrFinalState(c *channel.Channel) (bool, error) {
	latestSS, err := c.LatestSignedState()
	// There are no signed states. We consider this as consensus
	if err != nil && err.Error() == "No states are signed" {
		return true, nil
	}
	if latestSS.State().IsFinal {
		return true, nil
	}

	latestSupportedState, err := c.LatestSupportedState()
	if err != nil {
		return false, err
	}

	return cmp.Equal(latestSS.State(), latestSupportedState), nil
}

// GetChannelByIdFunction specifies a function that can be used to retrieve channels from a store.
type GetChannelByIdFunction func(id types.Destination) (channel *channel.Channel, ok bool)

// GetConsensusChannel describes functions which return a ConsensusChannel ledger channel for a channel id.
type GetConsensusChannel func(channelId types.Destination) (ledger *consensus_channel.ConsensusChannel, err error)

// NewObjective initiates an Objective with the supplied channel
func NewObjective(
	request ObjectiveRequest,
	preApprove bool,
	getConsensusChannel GetConsensusChannel,
) (Objective, error) {
	cc, err := getConsensusChannel(request.ChannelId)
	if err != nil {
		return Objective{}, fmt.Errorf("could not find channel %s; %w", request.ChannelId, err)
	}

	if len(cc.FundingTargets()) != 0 {
		return Objective{}, ErrNotEmpty
	}

	c, err := CreateChannelFromConsensusChannel(*cc)
	if err != nil {
		return Objective{}, fmt.Errorf("could not create Channel from ConsensusChannel; %w", err)
	}

	// We choose to disallow creating an objective if the channel has an in-progress update.
	// We allow the creation of of an objective if the channel has some final states.
	// In the future, we can add a restriction that only defund objectives can add final states to the channel.
	canCreateObjective, err := isInConsensusOrFinalState(c)
	if err != nil {
		return Objective{}, err
	}
	if !canCreateObjective {
		return Objective{}, ErrChannelUpdateInProgress
	}

	init := Objective{}

	if preApprove {
		init.Status = protocols.Approved
	} else {
		init.Status = protocols.Unapproved
	}
	init.C = c.Clone()

	latestSS, err := c.LatestSupportedState()
	if err != nil {
		return init, err
	}

	if !latestSS.IsFinal {
		init.finalTurnNum = latestSS.TurnNum + 1
	} else {
		init.finalTurnNum = latestSS.TurnNum
	}

	return init, nil
}

// ConstructObjectiveFromPayload takes in a state and constructs an objective from it.
func ConstructObjectiveFromPayload(
	p protocols.ObjectivePayload,
	preapprove bool,
	getConsensusChannel GetConsensusChannel,
) (Objective, error) {
	ss, err := getSignedStatePayload(p.PayloadData)
	if err != nil {
		return Objective{}, fmt.Errorf("could not get signed state payload: %w", err)
	}
	s := ss.State()

	// Implicit in the wire protocol is that the message signalling
	// closure of a channel includes an isFinal state (in the 0 slot of the message)
	//
	if !s.IsFinal {
		return Objective{}, ErrNoFinalState
	}

	err = s.FixedPart().Validate()
	if err != nil {
		return Objective{}, err
	}

	cId := s.ChannelId()
	request := NewObjectiveRequest(cId)
	return NewObjective(request, preapprove, getConsensusChannel)
}

// Public methods on the DirectDefundingObjective

// Id returns the unique id of the objective
func (o *Objective) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId(ObjectivePrefix + o.C.Id.String())
}

func (o *Objective) Approve() protocols.Objective {
	updated := o.clone()
	// todo: consider case of o.Status == Rejected
	updated.Status = protocols.Approved

	return &updated
}

func (o *Objective) Reject() (protocols.Objective, protocols.SideEffects) {
	updated := o.clone()
	updated.Status = protocols.Rejected
	peer := o.C.Participants[1-o.C.MyIndex]

	sideEffects := protocols.SideEffects{MessagesToSend: protocols.CreateRejectionNoticeMessage(o.Id(), peer)}
	return &updated, sideEffects
}

// OwnsChannel returns the channel that the objective is funding.
func (o Objective) OwnsChannel() types.Destination {
	return o.C.Id
}

// GetStatus returns the status of the objective.
func (o Objective) GetStatus() protocols.ObjectiveStatus {
	return o.Status
}

func (o *Objective) Related() []protocols.Storable {
	return []protocols.Storable{o.C}
}

// Update receives an ObjectiveEvent, applies all applicable event data to the DirectDefundingObjective,
// and returns the updated objective
func (o *Objective) Update(p protocols.ObjectivePayload) (protocols.Objective, error) {
	if o.Id() != p.ObjectiveId {
		return o, fmt.Errorf("event and objective Ids do not match: %s and %s respectively", string(p.ObjectiveId), string(o.Id()))
	}
	ss, err := getSignedStatePayload(p.PayloadData)
	if err != nil {
		return o, fmt.Errorf("could not get signed state payload: %w", err)
	}
	if len(ss.Signatures()) != 0 {

		if !ss.State().IsFinal {
			return o, errors.New("direct defund objective can only be updated with final states")
		}
		if o.finalTurnNum != ss.State().TurnNum {
			return o, fmt.Errorf("expected state with turn number %d, received turn number %d", o.finalTurnNum, ss.State().TurnNum)
		}
	} else {
		return o, fmt.Errorf("event does not contain a signed state")
	}

	updated := o.clone()
	updated.C.AddSignedState(ss)

	return &updated, nil
}

// UpdateWithChainEvent updates the objective with observed on-chain data.
//
// Only Allocation Updated events are currently handled.
func (o *Objective) UpdateWithChainEvent(event chainservice.Event) (protocols.Objective, error) {
	updated := o.clone()
	switch e := event.(type) {
	case chainservice.AllocationUpdatedEvent:
		// todo: check block number
		updated.C.OnChainFunding[e.AssetAddress] = e.AssetAmount
	case chainservice.ConcludedEvent:
		break
	default:
		return &updated, fmt.Errorf("objective %+v cannot handle event %+v", updated, event)
	}
	return &updated, nil
}

// Crank inspects the extended state and declares a list of Effects to be executed
func (o *Objective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	updated := o.clone()

	sideEffects := protocols.SideEffects{}

	if updated.Status != protocols.Approved {
		return &updated, sideEffects, WaitingForNothing, protocols.ErrNotApproved
	}

	latestSignedState, err := updated.C.LatestSignedState()
	if err != nil {
		return &updated, sideEffects, WaitingForNothing, errors.New("the channel must contain at least one signed state to crank the defund objective")
	}

	// Finalize and sign a state if no supported, finalized state exists
	if !latestSignedState.State().IsFinal || !latestSignedState.HasSignatureForParticipant(updated.C.MyIndex) {
		stateToSign := latestSignedState.State().Clone()
		if !stateToSign.IsFinal {
			stateToSign.TurnNum += 1
			stateToSign.IsFinal = true
		}
		ss, err := updated.C.SignAndAddState(stateToSign, secretKey)
		if err != nil {
			return &updated, protocols.SideEffects{}, WaitingForFinalization, fmt.Errorf("could not sign final state %w", err)
		}
		messages, err := protocols.CreateObjectivePayloadMessage(updated.Id(), ss, SignedStatePayload, o.otherParticipants()...)
		if err != nil {
			return &updated, protocols.SideEffects{}, WaitingForFinalization, fmt.Errorf("could not create payload message %w", err)
		}
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)
	}

	latestSupportedState, err := updated.C.LatestSupportedState()
	if err != nil {
		return &updated, sideEffects, WaitingForFinalization, fmt.Errorf("error finding a supported state: %w", err)
	}
	if !latestSupportedState.IsFinal {
		return &updated, sideEffects, WaitingForFinalization, nil
	}

	// Withdrawal of funds
	if !updated.fullyWithdrawn() {
		// The first participant in the channel submits the withdrawAll transaction
		if updated.C.MyIndex == 0 && !updated.transactionSubmitted {
			withdrawAll := protocols.NewWithdrawAllTransaction(updated.C.Id, latestSignedState)
			sideEffects.TransactionsToSubmit = append(sideEffects.TransactionsToSubmit, withdrawAll)
			updated.transactionSubmitted = true
		}
		// Every participant waits for all channel funds to be distributed, even if the participant has no funds in the channel
		return &updated, sideEffects, WaitingForWithdraw, nil
	}

	updated.Status = protocols.Completed
	return &updated, sideEffects, WaitingForNothing, nil
}

// IsDirectDefundObjective inspects a objective id and returns true if the objective id is for a direct defund objective.
func IsDirectDefundObjective(id protocols.ObjectiveId) bool {
	return strings.HasPrefix(string(id), ObjectivePrefix)
}

//  Private methods on the DirectDefundingObjective

// CreateChannelFromConsensusChannel creates a Channel with (an appropriate latest supported state) from the supplied ConsensusChannel.
func CreateChannelFromConsensusChannel(cc consensus_channel.ConsensusChannel) (*channel.Channel, error) {
	c, err := channel.New(cc.ConsensusVars().AsState(cc.SupportedSignedState().State().FixedPart()), uint(cc.MyIndex))
	if err != nil {
		return &channel.Channel{}, err
	}
	c.OnChainFunding = cc.OnChainFunding.Clone()
	c.AddSignedState(cc.SupportedSignedState())

	return c, nil
}

// fullyWithdrawn returns true if the channel contains no assets on chain
func (o *Objective) fullyWithdrawn() bool {
	return !o.C.OnChainFunding.IsNonZero()
}

// clone returns a deep copy of the receiver.
func (o *Objective) clone() Objective {
	clone := Objective{}
	clone.Status = o.Status

	cClone := o.C.Clone()
	clone.C = cClone
	clone.finalTurnNum = o.finalTurnNum
	clone.transactionSubmitted = o.transactionSubmitted

	return clone
}

// ObjectiveRequest represents a request to create a new direct defund objective.
type ObjectiveRequest struct {
	ChannelId        types.Destination
	objectiveStarted chan struct{}
}

// NewObjectiveRequest creates a new ObjectiveRequest.
func NewObjectiveRequest(channelId types.Destination) ObjectiveRequest {
	return ObjectiveRequest{
		ChannelId:        channelId,
		objectiveStarted: make(chan struct{}),
	}
}

// SignalObjectiveStarted is used by the engine to signal the objective has been started.
func (r ObjectiveRequest) SignalObjectiveStarted() {
	close(r.objectiveStarted)
}

// WaitForObjectiveToStart blocks until the objective starts
func (r ObjectiveRequest) WaitForObjectiveToStart() {
	<-r.objectiveStarted
}

// Id returns the objective id for the request.
func (r ObjectiveRequest) Id(myAddress types.Address, chainId *big.Int) protocols.ObjectiveId {
	return protocols.ObjectiveId(ObjectivePrefix + r.ChannelId.String())
}

// getSignedStatePayload takes in a serialized signed state payload and returns the deserialized SignedState.
func getSignedStatePayload(b []byte) (state.SignedState, error) {
	ss := state.SignedState{}
	err := json.Unmarshal(b, &ss)
	if err != nil {
		return ss, fmt.Errorf("could not unmarshal signed state: %w", err)
	}
	return ss, nil
}

// otherParticipants returns the participants in the channel that are not the current participant.
func (o *Objective) otherParticipants() []types.Address {
	others := make([]types.Address, 0)
	for i, p := range o.C.Participants {
		if i != int(o.C.MyIndex) {
			others = append(others, p)
		}
	}
	return others
}
