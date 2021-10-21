package protocols

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

// Message is an object to be sent across the wire. It can contain a proposal and signed state hashes, and is addressed to a counterparty.
type Message struct {
	To          []byte
	ObjectiveId ObjectiveId
	Sigs        map[types.Bytes32]state.Signature // mapping from state hash to signature
	Proposal    Objective
}

// Transaction is an object to be sent to a blockchain provider.
type Transaction struct {
	To   types.Address
	Data []byte
}

// SideEffects are effects to be executed by an imperative shell
type SideEffects struct {
	MessagesToSend       []Message
	TransactionsToSubmit []Transaction
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
	ChannelId          string                            // Must be defined
	Sigs               map[types.Bytes32]state.Signature // mapping from state hash to signature
	Holdings           map[types.Address]big.Int         // mapping from asset identifier to amount
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
	Initialize(initialState state.State) Objective // returns the initial Objective, does not declare effects

	Approve() Objective                    // returns an updated Objective (a copy, no mutation allowed), does not declare effects
	Reject() Objective                     // returns an updated Objective (a copy, no mutation allowed), does not declare effects
	Update(event ObjectiveEvent) Objective // returns an updated Objective (a copy, no mutation allowed), does not declare effects

	Crank(secretKey *[]byte) (Objective, SideEffects, WaitingFor, error) // does *not* accept an event, but *does* accept a pointer to a signing key; declare side effects; return an updated Objective
}

// ObjectiveId is a unique identifier for an Objective.
type ObjectiveId string
