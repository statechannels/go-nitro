package protocols

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

var ErrNotApproved = errors.New("objective not approved")

// ChainTransaction defines the interface that every transaction must implement
type ChainTransaction interface {
	ChannelId() types.Destination
}

// ChainTransactionBase is a convenience struct that is embedded in other transaction structs. It is exported only to allow cmp.Diff to compare transactions
type ChainTransactionBase struct {
	channelId types.Destination
}

func (cct ChainTransactionBase) ChannelId() types.Destination {
	return cct.channelId
}

type DepositTransaction struct {
	ChainTransaction
	Deposit types.Funds
}

func NewDepositTransaction(channelId types.Destination, deposit types.Funds) DepositTransaction {
	return DepositTransaction{ChainTransaction: ChainTransactionBase{channelId: channelId}, Deposit: deposit}
}

type WithdrawAllTransaction struct {
	ChainTransaction
	SignedState state.SignedState
}

func NewWithdrawAllTransaction(channelId types.Destination, signedState state.SignedState) WithdrawAllTransaction {
	return WithdrawAllTransaction{SignedState: signedState, ChainTransaction: ChainTransactionBase{channelId: channelId}}
}

type ChallengeTransaction struct {
	ChainTransaction
	Candidate     state.SignedState
	Proof         []state.SignedState
	ChallengerSig crypto.Signature
}

func NewChallengeTransaction(
	channelId types.Destination,
	candidate state.SignedState,
	proof []state.SignedState,
	challengerSig crypto.Signature,
) ChallengeTransaction {
	return ChallengeTransaction{ChainTransaction: ChainTransactionBase{channelId: channelId}, Candidate: candidate, Proof: proof, ChallengerSig: challengerSig}
}

// SideEffects are effects to be executed by an imperative shell
type SideEffects struct {
	MessagesToSend       []Message
	TransactionsToSubmit []ChainTransaction
	ProposalsToProcess   []consensus_channel.Proposal
}

// WaitingFor is an enumerable "pause-point" computed from an Objective. It describes how the objective is blocked on actions by third parties (i.e. co-participants or the blockchain).
type WaitingFor string

// AdjudicationStatus mirrors the on chain adjudication status of a particular channel.
// Everything that is stored on chain, other than holdings.
type AdjudicationStatus struct {
	TurnNumRecord uint
	// TODO This struct is a placeholder for the time being, until we add a chain-service
	// TODO eventually this struct will contain the other fields stored in (or committed to by) the adjudicator
}

// Storable is an object that can be stored by the store.
type Storable interface {
	json.Marshaler
	json.Unmarshaler
}

// Objective is the interface for off-chain protocols.
// The lifecycle of an objective is as follows:
//   - It is initialized by a single client (passing in various parameters). It is implicitly approved by that client. It is communicated to the other clients.
//   - It is stored and then approved or rejected by the other clients
//   - It is updated with external information arriving to the client
//   - After each update, it is cranked. This generates side effects and other metadata
//   - The metadata will eventually indicate that the Objective has stalled OR the Objective has completed successfully
type Objective interface {
	Id() ObjectiveId

	Approve() Objective                                                  // returns an updated Objective (a copy, no mutation allowed), does not declare effects
	Reject() (Objective, SideEffects)                                    // returns an updated Objective (a copy, no mutation allowed), does not declare effects
	Update(payload ObjectivePayload) (Objective, error)                  // returns an updated Objective (a copy, no mutation allowed), does not declare effects
	Crank(secretKey *[]byte) (Objective, SideEffects, WaitingFor, error) // does *not* accept an event, but *does* accept a pointer to a signing key; declare side effects; return an updated Objective

	// Related returns a slice of related objects that need to be stored along with the objective
	Related() []Storable
	Storable

	// OwnsChannel returns the channel the objective exclusively owns.
	OwnsChannel() types.Destination
	// GetStatus returns the status of the objective.
	GetStatus() ObjectiveStatus
}

// ProposalReceiver is an Objective that receives proposals.
type ProposalReceiver interface {
	Objective
	// ReceiveProposal receives a signed proposal and returns an updated VirtualObjective.
	// It is used to update the objective with a proposal received from a peer.
	ReceiveProposal(signedProposal consensus_channel.SignedProposal) (ProposalReceiver, error)
}

// ObjectiveId is a unique identifier for an Objective.
type ObjectiveId string

type ObjectiveStatus int8

const (
	Unapproved ObjectiveStatus = iota
	Approved
	Rejected
	Completed
)

// ObjectiveRequest is a request to create a new objective.
type ObjectiveRequest interface {
	Id(types.Address, *big.Int) ObjectiveId
	WaitForObjectiveToStart()
	SignalObjectiveStarted()
}
