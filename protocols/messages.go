package protocols

import (
	"bytes"
	"encoding/json"
	"sort"

	"github.com/ethereum/go-ethereum/common"
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
	To                types.Address
	ObjectivePayloads []ObjectivePayload
}

// An objective payload is either a Signed Propopsal or a Signed state for an objective id.
type ObjectivePayload struct {
	ObjectiveId    ObjectiveId
	SignedState    state.SignedState
	SignedProposal consensus_channel.SignedProposal
}

func (p ObjectivePayload) hasState() bool {
	stateHash, err := p.SignedState.State().Hash()
	if err != nil {
		return false
	}

	return bytes.Compare(stateHash.Bytes(), common.Hash{}.Bytes()) != 0
}

// Type returns the type of the payload, either a SignedProposal or SignedState.
func (p ObjectivePayload) Type() PayloadType {
	emptyProposal := consensus_channel.Proposal{}
	if p.SignedProposal.Proposal != emptyProposal {
		return SignedProposalPayload
	}
	return SignedStatePayload
}

// WithObjectiveId is a struct that contains an objectiveId and some value.
type WithObjectiveId[T interface{}] struct {
	Value       T
	ObjectiveId ObjectiveId
}

// SignedStates returns a slice of signed states with their objectiveId that were contained in the message.
// The states are sorted by channel id then turnNum.
func (m Message) SignedStates() []WithObjectiveId[state.SignedState] {
	signedStates := make([]WithObjectiveId[state.SignedState], 0)
	for _, p := range m.ObjectivePayloads {
		if p.Type() == SignedStatePayload {
			{
				entry := WithObjectiveId[state.SignedState]{p.SignedState, p.ObjectiveId}
				signedStates = append(signedStates, entry)
			}

		}

	}

	// Sort states by channelId then turn number
	sort.Slice(signedStates, func(i, j int) bool {
		s1, s2 := signedStates[i].Value, signedStates[j].Value
		s1CId, _ := s1.State().ChannelId()
		s2CId, _ := s2.State().ChannelId()

		cIdCompare := bytes.Compare(s1CId.Bytes(), s2CId.Bytes())

		if sameChannel := cIdCompare == 0; sameChannel {
			return s1.State().TurnNum < s2.State().TurnNum
		} else {
			return cIdCompare < 0
		}
	})

	return signedStates
}

// SignedProposals returns a slice of signed proposals with their objectiveId that were contained in the message.
// The proposals are sorted by ledger id then turnNum.
func (m Message) SignedProposals() []WithObjectiveId[consensus_channel.SignedProposal] {
	signedProposals := make([]WithObjectiveId[consensus_channel.SignedProposal], 0)
	for _, p := range m.ObjectivePayloads {
		if p.Type() == SignedProposalPayload {

			entry := WithObjectiveId[consensus_channel.SignedProposal]{p.SignedProposal, p.ObjectiveId}
			signedProposals = append(signedProposals, entry)
		}

	}

	// Sort states by channelId then turn number
	sort.Slice(signedProposals, func(i, j int) bool {
		s1, s2 := signedProposals[i], signedProposals[j]
		s1CId, s2CId := s1.Value.Proposal.ChannelID, s2.Value.Proposal.ChannelID

		cIdCompare := bytes.Compare(s1CId.Bytes(), s2CId.Bytes())

		if sameChannel := cIdCompare == 0; sameChannel {
			return s1.Value.Proposal.TurnNum() < s2.Value.Proposal.TurnNum()
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
		payload := ObjectivePayload{
			ObjectiveId: id,
			SignedState: ss,
		}

		message := Message{To: participant, ObjectivePayloads: []ObjectivePayload{payload}}
		messages = append(messages, message)
	}
	return messages
}

// Merge accepts a SideEffects struct that is merged into the the existing SideEffects.
func (se *SideEffects) Merge(other SideEffects) {

	se.MessagesToSend = append(se.MessagesToSend, other.MessagesToSend...)
	se.TransactionsToSubmit = append(se.TransactionsToSubmit, other.TransactionsToSubmit...)

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
			LedgerId:    p.Value.Proposal.ChannelID.String(),
			ObjectiveId: string(p.ObjectiveId),
			Target:      p.Value.Proposal.Target().String(),
			TurnNum:     p.Value.Proposal.TurnNum(),
			Type:        string(p.Value.Proposal.Type()),
		}
	}

	states := make([]StateSummary, len(m.SignedStates()))
	for i, s := range m.SignedStates() {
		channelId, _ := s.Value.State().ChannelId()
		states[i] = StateSummary{
			ObjectiveId: string(s.ObjectiveId),
			ChannelId:   channelId.String(),
			TurnNum:     s.Value.State().TurnNum,
		}
	}

	return MessageSummary{To: m.To.String(), Proposals: proposals, States: states}
}
