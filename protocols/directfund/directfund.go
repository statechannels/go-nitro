// Package directfund implements an off-chain protocol to directly fund a channel.
package directfund // import "github.com/statechannels/go-nitro/directfund"

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const (
	WaitingForCompletePrefund  protocols.WaitingFor = "WaitingForCompletePrefund"
	WaitingForMyTurnToFund     protocols.WaitingFor = "WaitingForMyTurnToFund"
	WaitingForCompleteFunding  protocols.WaitingFor = "WaitingForCompleteFunding"
	WaitingForCompletePostFund protocols.WaitingFor = "WaitingForCompletePostFund"
	WaitingForNothing          protocols.WaitingFor = "WaitingForNothing" // Finished
)

const ObjectivePrefix = "DirectFunding-"

func FundOnChainEffect(cId types.Destination, asset string, amount types.Funds) string {
	return "deposit" + amount.String() + "into" + cId.String()
}

// Objective is a cache of data computed by reading from the store. It stores (potentially) infinite data
type Objective struct {
	Status protocols.ObjectiveStatus
	C      *channel.Channel

	myDepositSafetyThreshold types.Funds // if the on chain holdings are equal to this amount it is safe for me to deposit
	myDepositTarget          types.Funds // I want to get the on chain holdings up to this much
	fullyFundedThreshold     types.Funds // if the on chain holdings are equal
	latestBlockNumber        uint64      // the latest block number we've seen
}

// NewObjective creates a new direct funding objective from a given request.
func NewObjective(request ObjectiveRequest, preApprove bool) (Objective, error) {

	objective, err := ConstructFromState(preApprove,
		state.State{
			ChainId:           big.NewInt(9001), // TODO https://github.com/statechannels/go-nitro/issues/601
			Participants:      []types.Address{request.MyAddress, request.CounterParty},
			ChannelNonce:      big.NewInt(request.Nonce),
			AppDefinition:     request.AppDefinition,
			ChallengeDuration: request.ChallengeDuration,
			AppData:           request.AppData,
			Outcome:           request.Outcome,
			TurnNum:           0,
			IsFinal:           false,
		},
		request.MyAddress,
	)
	if err != nil {
		return Objective{}, fmt.Errorf("could not create new objective: %w", err)
	}
	return objective, nil
}

// ConstructFromState initiates a Objective with data calculated from
// the supplied initialState and client address
func ConstructFromState(
	preApprove bool,
	initialState state.State,
	myAddress types.Address,
) (Objective, error) {
	if initialState.TurnNum != 0 {
		return Objective{}, errors.New("cannot construct direct fund objective without prefund state")
	}
	if initialState.IsFinal {
		return Objective{}, errors.New("attempted to initiate new direct-funding objective with IsFinal == true")
	}

	var init = Objective{}
	var err error

	if preApprove {
		init.Status = protocols.Approved
	} else {
		init.Status = protocols.Unapproved
	}

	var myIndex uint
	foundMyAddress := false
	for i, v := range initialState.Participants {
		if v == myAddress {
			myIndex = uint(i)
			foundMyAddress = true
			break
		}
	}
	if !foundMyAddress {
		return Objective{}, errors.New("my address not found in participants")
	}

	init.C = &channel.Channel{}
	init.C, err = channel.New(initialState, myIndex)

	if err != nil {
		return Objective{}, fmt.Errorf("failed to initialize channel for direct-fund objective: %w", err)
	}

	myAllocatedAmount := initialState.Outcome.TotalAllocatedFor(
		types.AddressToDestination(myAddress),
	)

	init.fullyFundedThreshold = initialState.Outcome.TotalAllocated()
	init.myDepositSafetyThreshold = initialState.Outcome.DepositSafetyThreshold(
		types.AddressToDestination(myAddress),
	)
	init.myDepositTarget = init.myDepositSafetyThreshold.Add(myAllocatedAmount)

	return init, nil
}

// OwnsChannel returns the channel that the objective is funding.
func (dfo Objective) OwnsChannel() types.Destination {
	return dfo.C.Id
}

// GetStatus returns the status of the objective.
func (dfo Objective) GetStatus() protocols.ObjectiveStatus {
	return dfo.Status
}

// CreateConsensusChannel creates a ConsensusChannel from the Objective by extracting signatures and a single asset outcome from the post fund state.
func (dfo *Objective) CreateConsensusChannel() (*consensus_channel.ConsensusChannel, error) {
	ledger := dfo.C

	if !ledger.PostFundComplete() {
		return nil, fmt.Errorf("expected funding for channel %s to be complete", dfo.C.Id)
	}
	signedPostFund := ledger.SignedPostFundState()
	leaderSig, err := signedPostFund.GetParticipantSignature(uint(consensus_channel.Leader))
	if err != nil {
		return nil, fmt.Errorf("could not get leader signature: %w", err)
	}
	followerSig, err := signedPostFund.GetParticipantSignature(uint(consensus_channel.Follower))
	if err != nil {
		return nil, fmt.Errorf("could not get follower signature: %w", err)
	}
	signatures := [2]state.Signature{leaderSig, followerSig}

	if len(signedPostFund.State().Outcome) != 1 {
		return nil, fmt.Errorf("a consensus channel only supports a single asset")
	}
	assetExit := signedPostFund.State().Outcome[0]
	turnNum := signedPostFund.State().TurnNum
	outcome, err := consensus_channel.FromExit(assetExit)

	if err != nil {
		return nil, fmt.Errorf("could not create ledger outcome from channel exit: %w", err)
	}

	if ledger.MyIndex == uint(consensus_channel.Leader) {
		con, err := consensus_channel.NewLeaderChannel(ledger.FixedPart, turnNum, outcome, signatures)
		con.OnChainFunding = ledger.OnChainFunding.Clone() // Copy OnChainFunding so we don't lose this information
		if err != nil {
			return nil, fmt.Errorf("could not create consensus channel as leader: %w", err)
		}
		return &con, nil

	} else {
		con, err := consensus_channel.NewFollowerChannel(ledger.FixedPart, turnNum, outcome, signatures)
		con.OnChainFunding = ledger.OnChainFunding.Clone() // Copy OnChainFunding so we don't lose this information
		if err != nil {
			return nil, fmt.Errorf("could not create consensus channel as follower: %w", err)
		}
		return &con, nil
	}

}

// Public methods on the DirectFundingObjectiveState

func (o Objective) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId(ObjectivePrefix + o.C.Id.String())
}

func (o Objective) Approve() protocols.Objective {
	updated := o.clone()
	// todo: consider case of s.Status == Rejected
	updated.Status = protocols.Approved

	return &updated
}

func (o Objective) Reject() protocols.Objective {
	updated := o.clone()
	updated.Status = protocols.Rejected
	return &updated
}

// Update receives an ObjectiveEvent, applies all applicable event data to the DirectFundingObjectiveState,
// and returns the updated state
func (o Objective) Update(event protocols.ObjectiveEvent) (protocols.Objective, error) {
	if o.Id() != event.ObjectiveId {
		return &o, fmt.Errorf("event and objective Ids do not match: %s and %s respectively", string(event.ObjectiveId), string(o.Id()))
	}

	updated := o.clone()
	updated.C.AddSignedState(event.SignedState)

	return &updated, nil
}

// UpdateWithChainEvent updates the objective with observed on-chain data.
//
// Only Channel Deposit events are currently handled.
func (o Objective) UpdateWithChainEvent(event chainservice.Event) (protocols.Objective, error) {
	updated := o.clone()

	de, ok := event.(chainservice.DepositedEvent)
	if !ok {
		return &updated, fmt.Errorf("objective %+v cannot handle event %+v", updated, event)
	}
	if de.Holdings != nil && de.BlockNum > updated.latestBlockNumber {
		updated.C.OnChainFunding = de.Holdings.Clone()
		updated.latestBlockNumber = de.BlockNum
	}

	return &updated, nil

}

// Crank inspects the extended state and declares a list of Effects to be executed
// It's like a state machine transition function where the finite / enumerable state is returned (computed from the extended state)
// rather than being independent of the extended state; and where there is only one type of event ("the crank") with no data on it at all
func (o Objective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	updated := o.clone()

	sideEffects := protocols.SideEffects{}
	// Input validation
	if updated.Status != protocols.Approved {
		return &updated, protocols.SideEffects{}, WaitingForNothing, protocols.ErrNotApproved
	}

	// Prefunding
	if !updated.C.PreFundSignedByMe() {
		ss, err := updated.C.SignAndAddPrefund(secretKey)
		if err != nil {
			return &updated, protocols.SideEffects{}, WaitingForCompletePrefund, fmt.Errorf("could not sign prefund %w", err)
		}
		messages := protocols.CreateSignedStateMessages(updated.Id(), ss, updated.C.MyIndex)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)
	}

	if !updated.C.PreFundComplete() {
		return &updated, sideEffects, WaitingForCompletePrefund, nil
	}

	// Funding
	fundingComplete := updated.fundingComplete() // note all information stored in state (since there are no real events)
	amountToDeposit := updated.amountToDeposit()
	safeToDeposit := updated.safeToDeposit()

	if !fundingComplete && !safeToDeposit {
		return &updated, sideEffects, WaitingForMyTurnToFund, nil
	}

	if !fundingComplete && safeToDeposit && amountToDeposit.IsNonZero() {
		deposit := protocols.ChainTransaction{Type: protocols.DepositTransactionType, ChannelId: updated.C.Id, Deposit: amountToDeposit}
		sideEffects.TransactionsToSubmit = append(sideEffects.TransactionsToSubmit, deposit)
	}

	if !fundingComplete {
		return &updated, sideEffects, WaitingForCompleteFunding, nil
	}

	// Postfunding
	if !updated.C.PostFundSignedByMe() {

		ss, err := updated.C.SignAndAddPostfund(secretKey)

		if err != nil {
			return &updated, protocols.SideEffects{}, WaitingForCompletePostFund, fmt.Errorf("could not sign postfund %w", err)
		}
		messages := protocols.CreateSignedStateMessages(updated.Id(), ss, updated.C.MyIndex)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)
	}

	if !updated.C.PostFundComplete() {
		return &updated, sideEffects, WaitingForCompletePostFund, nil
	}

	// Completion
	return &updated, sideEffects, WaitingForNothing, nil
}

func (o Objective) Related() []protocols.Storable {
	return []protocols.Storable{o.C}
}

//  Private methods on the DirectFundingObjectiveState

// fundingComplete returns true if the recorded OnChainHoldings are greater than or equal to the threshold for being fully funded.
func (o Objective) fundingComplete() bool {
	for asset, threshold := range o.fullyFundedThreshold {
		chainHolding, ok := o.C.OnChainFunding[asset]

		if !ok {
			return false
		}

		if types.Gt(threshold, chainHolding) {
			return false
		}
	}

	return true
}

// safeToDeposit returns true if the recorded OnChainHoldings are greater than or equal to the threshold for safety.
func (o Objective) safeToDeposit() bool {
	for asset, safetyThreshold := range o.myDepositSafetyThreshold {

		chainHolding, ok := o.C.OnChainFunding[asset]

		if !ok {
			panic("nil chainHolding for asset in myDepositSafetyThreshold")
		}

		if types.Gt(safetyThreshold, chainHolding) {
			return false
		}
	}

	return true
}

// amountToDeposit computes the appropriate amount to deposit given the current recorded OnChainHoldings
func (o Objective) amountToDeposit() types.Funds {
	deposits := make(types.Funds, len(o.C.OnChainFunding))

	for asset, target := range o.myDepositTarget {
		holding, ok := o.C.OnChainFunding[asset]
		if !ok {
			panic("nil chainHolding for asset in myDepositTarget")
		}
		deposits[asset] = big.NewInt(0).Sub(target, holding)
	}

	return deposits
}

// clone returns a deep copy of the receiver.
func (o Objective) clone() Objective {
	clone := Objective{}
	clone.Status = o.Status

	cClone := o.C.Clone()
	clone.C = cClone

	clone.myDepositSafetyThreshold = o.myDepositSafetyThreshold.Clone()
	clone.myDepositTarget = o.myDepositTarget.Clone()
	clone.fullyFundedThreshold = o.fullyFundedThreshold.Clone()
	clone.latestBlockNumber = o.latestBlockNumber
	return clone
}

// IsDirectFundObjective inspects a objective id and returns true if the objective id is for a direct fund objective.
func IsDirectFundObjective(id protocols.ObjectiveId) bool {
	return strings.HasPrefix(string(id), ObjectivePrefix)
}

// ObjectiveRequest represents a request to create a new direct funding objective.
type ObjectiveRequest struct {
	MyAddress         types.Address
	CounterParty      types.Address
	AppDefinition     types.Address
	AppData           types.Bytes
	ChallengeDuration *types.Uint256
	Outcome           outcome.Exit
	Nonce             int64
}

// Id returns the objective id for the request.
func (r ObjectiveRequest) Id() protocols.ObjectiveId {
	fixedPart := state.FixedPart{ChainId: big.NewInt(9001), // TODO add this field to the request and pull it from there. https://github.com/statechannels/go-nitro/issues/601
		Participants:      []types.Address{r.MyAddress, r.CounterParty},
		ChannelNonce:      big.NewInt(r.Nonce),
		ChallengeDuration: r.ChallengeDuration}

	channelId, _ := fixedPart.ChannelId()
	return protocols.ObjectiveId(ObjectivePrefix + channelId.String())
}

// ObjectiveResponse is the type returned across the API in response to the ObjectiveRequest.
type ObjectiveResponse struct {
	Id        protocols.ObjectiveId
	ChannelId types.Destination
}

// Response computes and returns the appropriate response from the request.
func (r ObjectiveRequest) Response() ObjectiveResponse {
	fixedPart := state.FixedPart{ChainId: big.NewInt(9001), // TODO add this field to the request and pull it from there. https://github.com/statechannels/go-nitro/issues/601
		Participants:      []types.Address{r.MyAddress, r.CounterParty},
		ChannelNonce:      big.NewInt(r.Nonce),
		ChallengeDuration: r.ChallengeDuration}

	channelId, _ := fixedPart.ChannelId()

	return ObjectiveResponse{
		Id:        protocols.ObjectiveId(ObjectivePrefix + channelId.String()),
		ChannelId: channelId,
	}
}

// mermaid diagram
// key:
// - effect!
// - waiting...
//
// https://mermaid-js.github.io/mermaid-live-editor/edit/#eyJjb2RlIjoiZ3JhcGggVERcbiAgICBTdGFydCAtLT4gQ3tJbnZhbGlkIElucHV0P31cbiAgICBDIC0tPnxZZXN8IEVbZXJyb3JdXG4gICAgQyAtLT58Tm98IEQwXG4gICAgXG4gICAgRDB7U2hvdWxkU2lnblByZUZ1bmR9XG4gICAgRDAgLS0-fFllc3wgUjFbU2lnblByZWZ1bmQhXVxuICAgIEQwIC0tPnxOb3wgRDFcbiAgICBcbiAgICBEMXtTYWZlVG9EZXBvc2l0ICY8YnI-ICFGdW5kaW5nQ29tcGxldGV9XG4gICAgRDEgLS0-IHxZZXN8IFIyW0Z1bmQgb24gY2hhaW4hXVxuICAgIEQxIC0tPiB8Tm98IEQyXG4gICAgXG4gICAgRDJ7IVNhZmVUb0RlcG9zaXQgJjxicj4gIUZ1bmRpbmdDb21wbGV0ZX1cbiAgICBEMiAtLT4gfFllc3wgUjNbXCJteSB0dXJuLi4uXCJdXG4gICAgRDIgLS0-IHxOb3wgRDNcblxuICAgIEQze1NhZmVUb0RlcG9zaXQgJjxicj4gIUZ1bmRpbmdDb21wbGV0ZX1cbiAgICBEMyAtLT4gfFllc3wgUjRbRGVwb3NpdCFdXG4gICAgRDMgLS0-IHxOb3wgRDRcblxuICAgIEQ0eyFGdW5kaW5nQ29tcGxldGV9XG4gICAgRDQgLS0-IHxZZXN8IFI1W1wiY29tcGxldGUgZnVuZGluZy4uLlwiXVxuICAgIEQ0IC0tPiB8Tm98IEQ1XG5cbiAgICBENXtTaG91bGRTaWduUHJlRnVuZH1cbiAgICBENSAtLT58WWVzfCBSNltTaWduUG9zdGZ1bmQhXVxuICAgIEQ1IC0tPnxOb3wgRDZcblxuICAgIEQ2eyFQb3N0RnVuZENvbXBsZXRlfVxuICAgIEQ2IC0tPnxZZXN8IFI3W1wiY29tcGxldGUgcG9zdGZ1bmQuLi5cIl1cbiAgICBENiAtLT58Tm98IFI4XG5cbiAgICBSOFtcImZpbmlzaFwiXVxuICAgIFxuXG5cbiIsIm1lcm1haWQiOiJ7fSIsInVwZGF0ZUVkaXRvciI6ZmFsc2UsImF1dG9TeW5jIjp0cnVlLCJ1cGRhdGVEaWFncmFtIjp0cnVlfQ
// graph TD
//     Start --> C{Invalid Input?}
//     C -->|Yes| E[error]
//     C -->|No| D0

//     D0{ShouldSignPreFund}
//     D0 -->|Yes| R1[SignPrefund!]
//     D0 -->|No| D1

//     D1{SafeToDeposit &<br> !FundingComplete}
//     D1 --> |Yes| R2[Fund on chain!]
//     D1 --> |No| D2

//     D2{!SafeToDeposit &<br> !FundingComplete}
//     D2 --> |Yes| R3["wait my turn..."]
//     D2 --> |No| D3

//     D3{SafeToDeposit &<br> !FundingComplete}
//     D3 --> |Yes| R4[Deposit!]
//     D3 --> |No| D4

//     D4{!FundingComplete}
//     D4 --> |Yes| R5["wait for complete funding..."]
//     D4 --> |No| D5

//     D5{ShouldSignPostFund}
//     D5 -->|Yes| R6[SignPostfund!]
//     D5 -->|No| D6

//     D6{!PostFundComplete}
//     D6 -->|Yes| R7["wait for complete postfund..."]
//     D6 -->|No| R8

//     R8["finish"]
