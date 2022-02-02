package protocols

import (
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

// ChainTransaction is an object to be sent to a blockchain provider.
type ChainTransaction struct {
	ChannelId types.Destination
	Deposit   types.Funds
	// TODO support other transaction types (deposit, challenge, respond, conclude, withdraw)
}

// LedgerRequest is an object processed by the ledger cranker
type LedgerRequest struct {
	LedgerId    types.Destination
	Destination types.Destination
	Amount      types.Funds
	Left        types.Destination
	Right       types.Destination
}

// Equal checks for equality between the receiver and a second LedgerRequest
func (l LedgerRequest) Equal(m LedgerRequest) bool {
	return l.LedgerId == m.LedgerId && l.Amount.Equal(m.Amount) && l.Left == m.Left && l.Right == m.Right
}

// SideEffects are effects to be executed by an imperative shell
type SideEffects struct {
	MessagesToSend       []Message
	TransactionsToSubmit []ChainTransaction
	LedgerRequests       []LedgerRequest
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

// ObjectiveEvent holds information used to update an Objective. Some fields may be nil.
type ObjectiveEvent struct {
	ObjectiveId        ObjectiveId
	SignedStates       []state.SignedState
	Holdings           types.Funds // mapping from asset identifier to amount
	AdjudicationStatus AdjudicationStatus
}

// Objective is the interface for off-chain protocols.
// The lifecycle of an objective is as follows:
// 	* It is initialized by a single client (passing in various parameters). It is implicitly approved by that client. It is communicated to the other clients.
// 	* It is stored and then approved or rejected by the other clients
// 	* It is updated with external information arriving to the client
// 	* After each update, it is cranked. This generates side effects and other metadata
// 	* The metadata will eventually indicate that the Objective has stalled OR the Objective has completed successfully
type Objective interface {
	Id() ObjectiveId

	Approve() Objective                             // returns an updated Objective (a copy, no mutation allowed), does not declare effects
	Reject() Objective                              // returns an updated Objective (a copy, no mutation allowed), does not declare effects
	Update(event ObjectiveEvent) (Objective, error) // returns an updated Objective (a copy, no mutation allowed), does not declare effects
	Channels() []types.Destination

	Crank(secretKey *[]byte) (Objective, SideEffects, WaitingFor, error) // does *not* accept an event, but *does* accept a pointer to a signing key; declare side effects; return an updated Objective
}

// ObjectiveId is a unique identifier for an Objective.
type ObjectiveId string

type ObjectiveStatus int8

const (
	Unapproved ObjectiveStatus = iota
	Approved
	Rejected
)
