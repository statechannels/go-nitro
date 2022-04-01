package protocols

import (
	"bytes"
	"encoding/json"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

// Message is an object to be sent across the wire. It can contain a proposal and signed states, and is addressed to a counterparty.
type Message struct {
	To              types.Address
	ObjectiveId     ObjectiveId
	SignedStates    []state.SignedState
	SignedProposals []consensus_channel.SignedProposal
}

// Serialize serializes the message into a string.
func (m Message) Serialize() (string, error) {
	bytes, err := json.Marshal(m)
	return string(bytes), err
}

// DeserializeMessage deserializes the passed string into a protocols.Message.
func DeserializeMessage(s string) (Message, error) {
	msg := Message{}
	err := json.Unmarshal([]byte(s), &msg)
	return msg, err
}

// Equal returns true if the passed Message is deeply equal in value to the receiver, and false otherwise.
func (m Message) Equal(n Message) bool {
	if !bytes.Equal(m.To.Bytes(), n.To.Bytes()) {
		return false
	}
	if m.ObjectiveId != n.ObjectiveId {
		return false
	}
	if len(m.SignedStates) != len(n.SignedStates) {
		return false
	}

	if !equalLists(m.SignedStates, n.SignedStates) {
		return false
	}

	if !equalLists(m.SignedProposals, n.SignedProposals) {
		return false
	}

	return true
}

// CreateSignedStateMessages creates a set of messages containing the signed state.
// A message will be generated for each participant except for the participant at myIndex.
func CreateSignedStateMessages(id ObjectiveId, ss state.SignedState, myIndex uint) []Message {
	return createMessages(id, ss, consensus_channel.SignedProposal{}, myIndex)
}

// CreateSignedProposalMessages creates an list of messages containing the signed proposal
// A message will be generated for each participant except for the participant at myIndex.
func CreateSignedProposalMessages(id ObjectiveId, sp consensus_channel.SignedProposal, myIndex uint) []Message {
	return createMessages(id, state.SignedState{}, sp, myIndex)
}

// createMessages creates a list of messages with the signed state and the signed proposal
func createMessages(id ObjectiveId, ss state.SignedState, sp consensus_channel.SignedProposal, myIndex uint) []Message {
	messages := make([]Message, 0)
	for i, participant := range ss.State().Participants {

		// Do not generate a message for ourselves
		if uint(i) == myIndex {
			continue
		}
		message := Message{To: participant, ObjectiveId: id, SignedStates: []state.SignedState{ss}, SignedProposals: []consensus_channel.SignedProposal{sp}}
		messages = append(messages, message)
	}
	return messages
}

// Merge accepts a SideEffects struct that is merged into the the existing SideEffects.
func (se *SideEffects) Merge(other SideEffects) {

	se.MessagesToSend = append(se.MessagesToSend, other.MessagesToSend...)
	se.TransactionsToSubmit = append(se.TransactionsToSubmit, other.TransactionsToSubmit...)

}

type comparable[T any] interface {
	Equal(a T) bool
}

func equalLists[T comparable[T]](a []T, b []T) bool {
	for i, sp := range a {
		if !sp.Equal(b[i]) {
			return false
		}
	}
	return true
}
