package protocols

import (
	"bytes"
	"encoding/json"
	"sort"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

const (
	SignedStatePayload    PayloadType = "SignedStatePayload"
	SignedProposalPayload PayloadType = "SignedProposalPayload"
)

type PayloadType string

// Message is an object to be sent across the wire. It can contain a proposal and signed states, and is addressed to a counterparty.
type Message struct {
	To       types.Address
	Payloads []MessagePayload
}

// MessagePayload is an objective id and EITHER a SignedState or SignedProposal.
type MessagePayload struct {
	ObjectiveId    ObjectiveId
	SignedState    state.SignedState
	SignedProposal consensus_channel.SignedProposal
}

// hasState returns true if the payload contains a signed state.
func (p MessagePayload) hasState() bool {
	return !p.SignedState.State().Equal(state.State{})
}

// hasProposal returns true if the payload contains a signed proposal.
func (p MessagePayload) hasProposal() bool {
	return p.SignedProposal.Proposal != consensus_channel.Proposal{}
}

// Type returns the type of the payload, either a SignedProposal or SignedState.
func (p MessagePayload) Type() PayloadType {
	switch {
	case p.hasProposal() && !p.hasState():
		return SignedProposalPayload
	case !p.hasProposal() && p.hasState():
		return SignedStatePayload
	case p.hasProposal() && p.hasState():
		panic("payload has both state and proposal %v")
	default:
		panic("payload has neither state nor proposal")
	}
}

// ObjectivePayload is a struct that contains an objectiveId and EITHER a Signed State or Signed Proposal.
type ObjectivePayload[T PayloadValue] struct {
	Payload     T
	ObjectiveId ObjectiveId
}

// SignedStates returns a slice of signed states with their objectiveId that were contained in the message.
// The states are sorted by channel id then turnNum.
func (m Message) SignedStates() []ObjectivePayload[state.SignedState] {
	signedStates := make([]ObjectivePayload[state.SignedState], 0)
	for _, p := range m.Payloads {
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
	for _, p := range m.Payloads {
		if p.Type() == SignedProposalPayload {
			entry := ObjectivePayload[consensus_channel.SignedProposal]{p.SignedProposal, p.ObjectiveId}
			signedProposals = append(signedProposals, entry)
		}
	}

	sortPayloads(signedProposals)

	return signedProposals
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

// CreateSignedStateMessages creates a set of messages containing the signed state.
// A message will be generated for each participant except for the participant at myIndex.
func CreateSignedStateMessages(id ObjectiveId, ss state.SignedState, myIndex uint) []Message {
	messages := make([]Message, 0)

	for i, participant := range ss.State().Participants {

		// Do not generate a message for ourselves
		if uint(i) == myIndex {
			continue
		}
		payload := MessagePayload{
			ObjectiveId: id,
			SignedState: ss,
		}

		message := Message{To: participant, Payloads: []MessagePayload{payload}}
		messages = append(messages, message)
	}
	return messages
}

// Merge accepts a SideEffects struct that is merged into the the existing SideEffects.
func (se *SideEffects) Merge(other SideEffects) {

	se.MessagesToSend = append(se.MessagesToSend, other.MessagesToSend...)
	se.TransactionsToSubmit = append(se.TransactionsToSubmit, other.TransactionsToSubmit...)

}

// PayloadValue is a type constraint that specifies a payload is either a SignedProposal or SignedState.
// It includes functions to get basic info to allow sorting.
type PayloadValue interface {
	state.SignedState | consensus_channel.SignedProposal
	ChannelId() types.Destination
	TurnNum() uint64
}

// sortPayloads sorts the objective payloads by channel id then turnNum.
// This is used to ensure that the payloads can be processed in a deterministic order.
func sortPayloads[T PayloadValue](payloads []ObjectivePayload[T]) {
	sort.Slice(payloads, func(i, j int) bool {
		p1, p2 := payloads[i], payloads[j]

		cId1 := p1.Payload.ChannelId()
		cId2 := p2.Payload.ChannelId()
		cIdCompare := bytes.Compare(cId1.Bytes(), cId2.Bytes())

		if sameChannel := cIdCompare == 0; sameChannel {
			return p1.Payload.TurnNum() < p2.Payload.TurnNum()
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
}

// SummarizeMessage returns a MessageSummary for the provided message.
func SummarizeMessage(m Message) MessageSummary {
	proposals := make([]ProposalSummary, len(m.SignedProposals()))
	for i, p := range m.SignedProposals() {

		proposals[i] = ProposalSummary{
			LedgerId:    p.Payload.Proposal.ChannelID.String(),
			ObjectiveId: string(p.ObjectiveId),
			Target:      p.Payload.Proposal.Target().String(),
			TurnNum:     p.Payload.Proposal.TurnNum(),
			Type:        string(p.Payload.Proposal.Type()),
		}
	}

	states := make([]StateSummary, len(m.SignedStates()))
	for i, s := range m.SignedStates() {
		channelId, _ := s.Payload.State().ChannelId()
		states[i] = StateSummary{
			ObjectiveId: string(s.ObjectiveId),
			ChannelId:   channelId.String(),
			TurnNum:     s.Payload.State().TurnNum,
		}
	}

	return MessageSummary{To: m.To.String(), Proposals: proposals, States: states}
}

// CreateSignedProposalMessage returns a signed proposal message addressed to the counterparty in the given ledger
// It contains the provided signed proposals and any proposals in the proposal queue.
func CreateSignedProposalMessage(ledger *consensus_channel.ConsensusChannel, sp ...consensus_channel.SignedProposal) Message {
	recipient := ledger.Leader()
	if ledger.IsLeader() {
		recipient = ledger.Follower()
	}
	// Append the proposals to any existing proposals in the queue
	proposals := append(ledger.ProposalQueue(), sp...)
	payloads := make([]MessagePayload, len(proposals))
	for i, sp := range proposals {
		id := getProposalObjectiveId(sp.Proposal)
		payloads[i] = MessagePayload{
			ObjectiveId:    id,
			SignedProposal: sp,
		}
	}

	return Message{
		To:       recipient,
		Payloads: payloads,
	}
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
