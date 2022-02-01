package protocols

import (
	"bytes"
	"encoding/json"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

// Message is an object to be sent across the wire. It can contain a proposal and signed state hashes, and is addressed to a counterparty.
type Message struct {
	To           types.Address
	ObjectiveId  ObjectiveId
	SignedStates []state.SignedState
	Proposal     Objective
}

// Serialize serializes the message into a string.
func (m Message) Serialize() (string, error) {
	bytes, err := json.Marshal(m)
	return string(bytes), err
}

// DeserialiseMessage deserializes the passed string into a protocols.Message.
func DeserialiseMessage(s string) (Message, error) {
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
	for i, ss := range m.SignedStates {
		if !ss.Equal(m.SignedStates[i]) {
			return false
		}
	}
	// TODO handle Proposal field :-/
	return true
}
