package protocols

import (
	"encoding/json"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

// Message is an object to be sent across the wire. It can contain a proposal and signed state hashes, and is addressed to a counterparty.
type Message struct {
	To          types.Address
	ObjectiveId ObjectiveId
	Sigs        map[*state.State]state.Signature // mapping from state hash to signature
	Proposal    Objective
}

// Serialize serializes the message into a string
func (m Message) Serialize() string {
	bytes, _ := json.Marshal(m) // TODO handle error
	return string(bytes)
}

// DeserialiseMessage deserializes the passed string into a protocols.Message
func DeserialiseMessage(s string) Message {
	msg := Message{}
	json.Unmarshal([]byte(s), msg)
	return msg
}
