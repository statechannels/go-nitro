package protocols

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/types"
)

type PayloadType string

const (
	SignedStatePayload     PayloadType = "SignedStatePayload"
	SignedProposalPayload  PayloadType = "SignedProposalPayload"
	RejectionNoticePayload PayloadType = "RejectionNoticePayload"
	VoucherPayload         PayloadType = "VoucherPayload"
)

// Message is an object to be sent across the wire. It can contain a proposal and signed states, and is addressed to a counterparty.
type Message struct {
	To       types.Address
	payloads []messagePayload
}

// messagePayload is an objective id and EITHER a SignedState or SignedProposal. This package guarantees that a payload has only one value by:
//  - validating messages that are deserialized from JSON
//  - providing message constructors which create valid messages

type messagePayload struct {
	ObjectiveId    ObjectiveId
	SignedState    state.SignedState
	SignedProposal consensus_channel.SignedProposal
	Voucher        payments.Voucher
	Rejected       bool
}

// hasState returns true if the payload contains a signed state.
func (p messagePayload) hasState() bool {
	return !p.SignedState.State().Equal(state.State{})
}

// hasProposal returns true if the payload contains a signed proposal.
func (p messagePayload) hasProposal() bool {
	return p.SignedProposal.Proposal != consensus_channel.Proposal{}
}

// hasRejection returns true if the payload indicates that the objective is rejected
func (p messagePayload) hasRejection() bool {
	return p.Rejected
}

// hasVoucher returns true if the payload contains a voucher.
func (p messagePayload) hasVoucher() bool {
	return !(&p.Voucher).Equal(&payments.Voucher{})
}

// Type returns the type of the payload, either a SignedProposal or SignedState.
func (p messagePayload) Type() PayloadType {
	if p.hasProposal() {
		return SignedProposalPayload
	} else if p.hasState() {
		return SignedStatePayload
	} else if p.hasVoucher() {
		return VoucherPayload
	} else {
		return RejectionNoticePayload
	}
}

// ObjectivePayload is a struct that contains an objectiveId and EITHER a Signed State or Signed Proposal.
type ObjectivePayload[T PayloadValue] struct {
	Payload     T
	ObjectiveId ObjectiveId
}

// Vouchers returns the collection of vouchers in the message
func (m Message) Vouchers() []payments.Voucher {
	var vouchers []payments.Voucher
	for _, p := range m.payloads {
		if p.hasVoucher() {
			vouchers = append(vouchers, p.Voucher)
		}
	}
	return vouchers
}

// SignedStates returns a slice of signed states with their objectiveId that were contained in the message.
// The states are sorted by channel id then turnNum.
func (m Message) SignedStates() []ObjectivePayload[state.SignedState] {
	signedStates := make([]ObjectivePayload[state.SignedState], 0)
	for _, p := range m.payloads {
		if p.Type() == SignedStatePayload {
			entry := ObjectivePayload[state.SignedState]{p.SignedState, p.ObjectiveId}
			signedStates = append(signedStates, entry)
		}
	}

	sortPayloads(signedStates)

	return signedStates
}

// SignedProposals returns a slice of signed proposals with their objectiveId that were contained in the message.
// The proposals are sorted by ledger id then turnNum.
func (m Message) SignedProposals() []ObjectivePayload[consensus_channel.SignedProposal] {
	signedProposals := make([]ObjectivePayload[consensus_channel.SignedProposal], 0)
	for _, p := range m.payloads {
		if p.Type() == SignedProposalPayload {
			entry := ObjectivePayload[consensus_channel.SignedProposal]{p.SignedProposal, p.ObjectiveId}
			signedProposals = append(signedProposals, entry)
		}
	}

	sortPayloads(signedProposals)

	return signedProposals
}

// rejectedObjective is a placeholder type that holds no value, but implements the PayloadValue interface.
// This allows the RejectedObjectives function to use the generic ObjectivePayload type.
type rejectedObjective struct{}

func (o rejectedObjective) SortInfo() (channelID types.Destination, turnNum uint64) {
	return
}

// RejectedObjectives returns a slice of rejected objectives
func (m Message) RejectedObjectives() []ObjectivePayload[rejectedObjective] {
	rejectedObjectives := make([]ObjectivePayload[rejectedObjective], 0)
	for _, p := range m.payloads {
		if p.Type() == RejectionNoticePayload {
			entry := ObjectivePayload[rejectedObjective]{rejectedObjective{}, p.ObjectiveId}
			rejectedObjectives = append(rejectedObjectives, entry)
		}
	}

	return rejectedObjectives
}

// Serialize serializes the message into a string.
func (m Message) Serialize() (string, error) {
	bytes, err := json.Marshal(jsonMessage{m.To, m.payloads})
	return string(bytes), err
}

// jsonMessage is a private struct with public members, allowing a Message to be easily serialized
type jsonMessage struct {
	To       types.Address
	Payloads []messagePayload
}

// MarshalJSON provides a custom json marshaler that avoids marshaling empty structs
func (p *messagePayload) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["ObjectiveId"] = p.ObjectiveId
	switch p.Type() {
	case SignedStatePayload:
		m["SignedState"] = p.SignedState
	case SignedProposalPayload:
		m["SignedProposal"] = p.SignedProposal
	case RejectionNoticePayload:
		m["Rejected"] = p.Rejected
	case VoucherPayload:
		m["Voucher"] = p.Voucher
	default:
		return []byte{}, fmt.Errorf("unknown payload type")
	}

	return json.Marshal(m)
}

// ErrInvalidPayload is returned when the payload has too many values
var ErrInvalidPayload = fmt.Errorf("payload has too many values")

// DeserializeMessage deserializes the passed string into a protocols.Message.
func DeserializeMessage(s string) (Message, error) {
	msg := jsonMessage{}
	err := json.Unmarshal([]byte(s), &msg)

	for _, p := range msg.Payloads {
		numPresent := 0
		if p.hasProposal() {
			numPresent += 1
		}
		if p.hasState() {
			numPresent += 1
		}
		if p.hasVoucher() {
			numPresent += 1
		}
		if p.hasRejection() {
			numPresent += 1
		}
		if numPresent != 1 {
			return Message{}, ErrInvalidPayload
		}
	}

	return Message{To: msg.To, payloads: msg.Payloads}, err
}

// CreateSignedStateMessages creates a set of messages containing the signed state.
// A message will be generated for each participant except for the participant at myIndex.
func CreateSignedStateMessages(id ObjectiveId, ss state.SignedState, myIndex uint) []Message {
	messages := make([]Message, 0)

	for i, participant := range ss.State().Participants {

		// Do not generate a message for ourselves
		if uint(i) == myIndex {
			continue
		}
		payload := messagePayload{
			ObjectiveId: id,
			SignedState: ss,
		}

		message := Message{To: participant, payloads: []messagePayload{payload}}
		messages = append(messages, message)
	}
	return messages
}

// Merge accepts a SideEffects struct that is merged into the the existing SideEffects.
func (se *SideEffects) Merge(other SideEffects) {

	se.MessagesToSend = append(se.MessagesToSend, other.MessagesToSend...)
	se.TransactionsToSubmit = append(se.TransactionsToSubmit, other.TransactionsToSubmit...)

}

// PayloadValue is a type constraint that specifies a payload is either a SignedProposal, SignedState,
// or RejectedObjective
// It includes functions to get basic info to allow sorting.
type PayloadValue interface {
	state.SignedState | consensus_channel.SignedProposal | rejectedObjective | payments.Voucher
	SortInfo() (channelID types.Destination, turnNum uint64)
}

// sortPayloads sorts the objective payloads by channel id then turnNum.
// This is used to ensure that the payloads can be processed in a deterministic order.
func sortPayloads[T PayloadValue](payloads []ObjectivePayload[T]) {
	sort.Slice(payloads, func(i, j int) bool {
		cId1, turnNum1 := payloads[i].Payload.SortInfo()
		cId2, turnNum2 := payloads[j].Payload.SortInfo()

		cIdCompare := bytes.Compare(cId1.Bytes(), cId2.Bytes())

		if sameChannel := cIdCompare == 0; sameChannel {
			return turnNum1 < turnNum2
		} else {
			return cIdCompare < 0
		}
	})
}

// ProposalSummary contains some basic info about a proposal for logging.
type ProposalSummary struct {
	ObjectiveId string
	LedgerId    string
	Target      string
	Type        string
	TurnNum     uint64
}

// StateSummary contains some basic info about a state for logging.
type StateSummary struct {
	ObjectiveId string
	ChannelId   string
	TurnNum     uint64
}

// MessagarSummary contains some basic info about a message for logging.
type MessageSummary struct {
	To        string
	Proposals []ProposalSummary
	States    []StateSummary
	Rejected  []ObjectiveId
}

// SummarizeMessage returns a MessageSummary for the provided message.
func SummarizeMessage(m Message) MessageSummary {
	proposals := make([]ProposalSummary, len(m.SignedProposals()))
	for i, p := range m.SignedProposals() {

		proposals[i] = SummarizeProposal(p.ObjectiveId, p.Payload)
	}

	states := make([]StateSummary, len(m.SignedStates()))
	for i, s := range m.SignedStates() {
		channelId := s.Payload.State().ChannelId()
		states[i] = StateSummary{
			ObjectiveId: string(s.ObjectiveId),
			ChannelId:   channelId.String(),
			TurnNum:     s.Payload.State().TurnNum,
		}
	}
	rejectedObjectives := []ObjectiveId{}
	for _, p := range m.RejectedObjectives() {
		rejectedObjectives = append(rejectedObjectives, p.ObjectiveId)
	}

	return MessageSummary{To: m.To.String(), Proposals: proposals, States: states, Rejected: rejectedObjectives}
}

// SummarizeProposal returns a ProposalSummary for the provided signed proposal.
func SummarizeProposal(oId ObjectiveId, sp consensus_channel.SignedProposal) ProposalSummary {

	return ProposalSummary{
		LedgerId:    sp.Proposal.LedgerID.String(),
		ObjectiveId: string(oId),
		Target:      sp.Proposal.Target().String(),
		TurnNum:     sp.TurnNum,
		Type:        string(sp.Proposal.Type()),
	}
}

// CreateSignedProposalMessage returns a signed proposal message addressed to the counterparty in the given ledger
// It contains the provided signed proposals and any proposals in the proposal queue.
func CreateSignedProposalMessage(recipient types.Address, proposals ...consensus_channel.SignedProposal) Message {

	payloads := make([]messagePayload, len(proposals))
	for i, sp := range proposals {
		id := getProposalObjectiveId(sp.Proposal)
		payloads[i] = messagePayload{
			ObjectiveId:    id,
			SignedProposal: sp,
		}
	}

	return Message{
		To:       recipient,
		payloads: payloads,
	}
}

// CreateSignedProposalMessage returns a signed proposal message addressed to the counterparty in the given ledger
// It contains the provided signed proposals and any proposals in the proposal queue.
func CreateRejectionNoticeMessage(oId ObjectiveId, recipients ...types.Address) []Message {
	messages := make([]Message, len(recipients))
	for i, recipient := range recipients {
		payload := messagePayload{
			ObjectiveId: oId,
			Rejected:    true,
		}
		payloads := []messagePayload{payload}
		messages[i] = Message{To: recipient, payloads: payloads}

	}

	return messages
}

// CreateVoucherMessage returns a signed voucher message for each of the recipients provided.
func CreateVoucherMessage(voucher payments.Voucher, recipients ...types.Address) []Message {
	messages := make([]Message, len(recipients))
	for i, recipient := range recipients {
		payload := messagePayload{
			Voucher: voucher,
		}
		payloads := []messagePayload{payload}
		messages[i] = Message{To: recipient, payloads: payloads}

	}

	return messages
}

// getProposalObjectiveId returns the objectiveId for a proposal.
func getProposalObjectiveId(p consensus_channel.Proposal) ObjectiveId {
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
