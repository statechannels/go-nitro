package protocols

import (
	"bytes"
	"encoding/json"
	"sort"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/types"
)

// ObjectivePayload is a message containing a payload of []byte that an objective is responsible for decoding.
type ObjectivePayload struct {
	// PayloadData is the serialized json payload
	PayloadData []byte
	// ObjectiveId is the id of the objective that is responsible for decoding and handling the payload
	ObjectiveId ObjectiveId
	// Type is the type of the payload the message contains.
	// This is useful when a protocol wants to handle different types of payloads.
	Type PayloadType
}

type PayloadType string

// CreateObjectivePayload generates an objective message from the given objective id and payload.
// CreateObjectivePayload handles serializing `p` into json.
func CreateObjectivePayload(id ObjectiveId, payloadType PayloadType, p interface{}) ObjectivePayload {
	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return ObjectivePayload{PayloadData: b, ObjectiveId: id, Type: payloadType}
}

// Message is an object to be sent across the wire.
type Message struct {
	To types.Address
	// ObjectivePayloads contains a collection of payloads for various objectives.
	// Protocols are responsible for parsing the payload.
	ObjectivePayloads []ObjectivePayload
	// LedgerProposals contains a collection of signed proposals.
	// Since proposals need to be handled in order they need to be an explicit part of the message format.
	LedgerProposals []consensus_channel.SignedProposal
	// Payments contains a collection of signed vouchers representing payments.
	// Payments are handled outside of any objective.
	Payments []payments.Voucher
	// RejectedObjectives is a collection of objectives that have been rejected.
	RejectedObjectives []ObjectiveId
}

// SortedProposals sorts the proposals by channelId and then by turn number.
func (m Message) SortedProposals() []consensus_channel.SignedProposal {
	signedProposals := make([]consensus_channel.SignedProposal, len(m.LedgerProposals))
	copy(signedProposals, m.LedgerProposals)

	sort.Slice(signedProposals, func(i, j int) bool {
		cId1, turnNum1 := signedProposals[i].ChannelID(), signedProposals[i].TurnNum
		cId2, turnNum2 := signedProposals[j].ChannelID(), signedProposals[j].TurnNum

		cIdCompare := bytes.Compare(cId1.Bytes(), cId2.Bytes())

		if sameChannel := cIdCompare == 0; sameChannel {
			return turnNum1 < turnNum2
		} else {
			return cIdCompare < 0
		}
	})

	return signedProposals
}

// Serialize serializes the message into a string.
func (m Message) Serialize() (string, error) {
	bytes, err := json.Marshal(m)
	return string(bytes), err
}

// Merge accepts a SideEffects struct that is merged into the the existing SideEffects.
func (se *SideEffects) Merge(other SideEffects) {

	se.MessagesToSend = append(se.MessagesToSend, other.MessagesToSend...)
	se.TransactionsToSubmit = append(se.TransactionsToSubmit, other.TransactionsToSubmit...)

}

// GetProposalObjectiveId returns the objectiveId for a proposal.
func GetProposalObjectiveId(p consensus_channel.Proposal) ObjectiveId {
	switch p.Type() {
	case "AddProposal":
		{
			const prefix = "VirtualFund-"
			channelId := p.ToAdd.Guarantee.Target().String()
			return ObjectiveId(prefix + channelId)

		}
	case "RemoveProposal":
		{
			const prefix = "VirtualDefund-"
			channelId := p.ToRemove.Target.String()
			return ObjectiveId(prefix + channelId)

		}
	default:
		{
			panic("invalid proposal type")
		}
	}
}

// CreateObjectivePayloadMessage returns a message for each recipient tht contains an objective payload.
func CreateObjectivePayloadMessage(id ObjectiveId, p interface{}, payloadType PayloadType, recipients ...types.Address) []Message {
	messages := make([]Message, 0)

	for _, participant := range recipients {
		message := Message{To: participant, ObjectivePayloads: []ObjectivePayload{CreateObjectivePayload(id, payloadType, p)}}
		messages = append(messages, message)
	}
	return messages
}

// CreateSignedProposalMessage returns a signed proposal message addressed to the counterparty in the given ledger
// It contains the provided signed proposals and any proposals in the proposal queue.
func CreateRejectionNoticeMessage(oId ObjectiveId, recipients ...types.Address) []Message {
	messages := make([]Message, 0)
	for _, recipient := range recipients {
		message := Message{To: recipient, RejectedObjectives: []ObjectiveId{oId}}
		messages = append(messages, message)
	}

	return messages
}

// CreateSignedProposalMessage returns a signed proposal message addressed to the counterparty in the given ledger
// It contains the provided signed proposals and any proposals in the proposal queue.
func CreateSignedProposalMessage(recipient types.Address, proposals ...consensus_channel.SignedProposal) Message {
	msg := Message{To: recipient, LedgerProposals: proposals}
	return msg

}

// CreateVoucherMessage returns a signed voucher message for each of the recipients provided.
func CreateVoucherMessage(voucher payments.Voucher, recipients ...types.Address) []Message {
	messages := make([]Message, len(recipients))
	for i, recipient := range recipients {

		messages[i] = Message{To: recipient, Payments: []payments.Voucher{voucher}}

	}

	return messages
}

// DeserializeMessage deserializes the passed string into a protocols.Message.
func DeserializeMessage(s string) (Message, error) {
	msg := Message{}
	err := json.Unmarshal([]byte(s), &msg)

	return msg, err
}
