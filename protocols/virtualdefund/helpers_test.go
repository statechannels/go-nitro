package virtualdefund

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/testactors"
	ta "github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testhelpers"
	. "github.com/statechannels/go-nitro/internal/testhelpers"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// assertSideEffectsContainsMessageWith fails the test instantly if the supplied side effects does not contain a message for the supplied actor with the supplied expected signed state.
// TODO: This is copied from https://github.com/statechannels/go-nitro/blob/0722a1127241583944f32efa0638012f64b96bf0/protocols/virtualfund/virtualfund_single_hop_test.go#L409
func assertProposalSent(t *testing.T, ses protocols.SideEffects, sp consensus_channel.SignedProposal, to testactors.Actor) {

	Assert(t, len(ses.MessagesToSend) == 1, "expected one message")

	Assert(t, len(ses.MessagesToSend[0].SignedProposals) == 1, "expected one signed proposal")

	msg := ses.MessagesToSend[0]
	sent := msg.SignedProposals[0]

	Assert(t, len(ses.MessagesToSend[0].SignedProposals) == 1, "exp: %+v\n\n\tgot%+v", sent.Proposal, sp.Proposal)

	Assert(t, bytes.Equal(msg.To[:], to.Address[:]), "exp: %+v\n\n\tgot%+v", msg.To.String(), to.Address.String())

}

// generateLedgers generates the left and right ledger channels based on myRole
func generateLedgers(myRole uint, vId types.Destination) (left, right *consensus_channel.ConsensusChannel) {
	switch myRole {
	case 0:
		{
			return nil, prepareConsensusChannel(uint(consensus_channel.Leader), testactors.Alice, testactors.Irene, generateGuarantee(testactors.Alice, testactors.Irene, vId))
		}
	case 1:
		{

			return prepareConsensusChannel(uint(consensus_channel.Follower), testactors.Alice, testactors.Irene, generateGuarantee(testactors.Alice, testactors.Irene, vId)),
				prepareConsensusChannel(uint(consensus_channel.Leader), testactors.Irene, testactors.Bob, generateGuarantee(testactors.Irene, testactors.Bob, vId))

		}
	case 2:
		{

			return prepareConsensusChannel(uint(consensus_channel.Follower), testactors.Irene, testactors.Bob, generateGuarantee(testactors.Irene, testactors.Bob, vId)), nil

		}
	default:
		panic("invalid myRole")
	}
}

func generateGuarantee(left, right testactors.Actor, vId types.Destination) consensus_channel.Guarantee {
	return consensus_channel.NewGuarantee(big.NewInt(10), vId, left.Destination(), right.Destination())

}

// prepareConsensusChannel prepares a consensus channel with a consensus outcome
//  - allocating 6 to left
//  - allocating 4 to right
//  - including the given guarantees
func prepareConsensusChannel(role uint, left, right testactors.Actor, guarantees ...consensus_channel.Guarantee) *consensus_channel.ConsensusChannel {
	fp := state.FixedPart{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{left.Address, right.Address},
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
	}

	leftBal := consensus_channel.NewBalance(left.Destination(), big.NewInt(0))
	rightBal := consensus_channel.NewBalance(right.Destination(), big.NewInt(0))

	lo := *consensus_channel.NewLedgerOutcome(types.Address{}, leftBal, rightBal, guarantees)

	signedVars := consensus_channel.SignedVars{Vars: consensus_channel.Vars{Outcome: lo, TurnNum: 1}}
	leftSig, err := signedVars.Vars.AsState(fp).Sign(left.PrivateKey)
	if err != nil {
		panic(err)
	}
	rightSig, err := signedVars.Vars.AsState(fp).Sign(right.PrivateKey)
	if err != nil {
		panic(err)
	}
	sigs := [2]state.Signature{leftSig, rightSig}

	var cc consensus_channel.ConsensusChannel

	if role == 0 {
		cc, err = consensus_channel.NewLeaderChannel(fp, 1, lo, sigs)
	} else {
		cc, err = consensus_channel.NewFollowerChannel(fp, 1, lo, sigs)
	}
	if err != nil {
		panic(err)
	}

	return &cc
}

// checkForFollowerProposals checks that the follower have signed and sent the appropriate proposal
func checkForFollowerProposals(t *testing.T, se protocols.SideEffects, o *Objective, td testdata) {

	switch o.MyRole {
	case 1:
		{
			// Irene should accept a proposal from Alice
			rightProposal := consensus_channel.SignedProposal{Proposal: consensus_channel.NewRemoveProposal(o.ToMyLeft.Id, FinalTurnNum, o.VId(), td.leftAmount, td.rightAmount)}
			assertProposalSent(t, se, rightProposal, alice)
		}
	case 2:
		{
			// Bob should accept a proposal from Irene
			rightProposal := consensus_channel.SignedProposal{Proposal: consensus_channel.NewRemoveProposal(o.ToMyLeft.Id, FinalTurnNum, o.VId(), td.leftAmount, td.rightAmount)}
			assertProposalSent(t, se, rightProposal, irene)
		}

	}
}

// generateProposalsResponses generates the signed proposals that a participant should expect from the other participants
func generateProposalsResponses(myRole uint, vId types.Destination, o *Objective, td testdata) []consensus_channel.SignedProposal {
	switch myRole {
	case 0:
		{
			// Alice expects Irene to accept her proposal
			p := consensus_channel.NewRemoveProposal(o.ToMyRight.Id, FinalTurnNum, vId, td.leftAmount, td.rightAmount)
			sp, err := signProposal(irene, p, o.ToMyRight)
			if err != nil {
				panic(err)
			}

			return []consensus_channel.SignedProposal{sp}
		}

	case 1:
		{
			// Irene expects Alice to send a proposal
			p := consensus_channel.NewRemoveProposal(o.ToMyLeft.Id, FinalTurnNum, vId, td.leftAmount, td.rightAmount)
			sp, err := signProposal(alice, p, o.ToMyLeft)
			if err != nil {
				panic(err)
			}

			// Irene expects Bob to accept her proposal
			p2 := consensus_channel.NewRemoveProposal(o.ToMyRight.Id, FinalTurnNum, vId, td.leftAmount, td.rightAmount)
			sp2, err := signProposal(bob, p2, o.ToMyRight)
			if err != nil {
				panic(err)
			}

			return []consensus_channel.SignedProposal{sp, sp2}
		}
	case 2:
		{
			// Bob expects Irene to send a proposal
			p := consensus_channel.NewRemoveProposal(o.ToMyLeft.Id, FinalTurnNum, vId, td.leftAmount, td.rightAmount)
			sp, err := signProposal(irene, p, o.ToMyLeft)
			if err != nil {
				panic(err)
			}
			return []consensus_channel.SignedProposal{sp}

		}
	default:
		return []consensus_channel.SignedProposal{}
	}
}

// updateProposals updates the consensus channels on the objective with the given proposals
// It is used to simulate having received a proposal from the other party
func updateProposals(o *Objective, proposals ...consensus_channel.SignedProposal) {
	for _, p := range proposals {
		var err error
		if o.ToMyLeft != nil && o.ToMyLeft.Id == p.Proposal.ChannelID {
			err = o.ToMyLeft.Receive(p)
		}
		if o.ToMyRight != nil && o.ToMyRight.Id == p.Proposal.ChannelID {
			err = o.ToMyRight.Receive(p)
		}
		if err != nil {
			panic(err)
		}
	}
}

// checkForLeaderProposals checks that the outgoing message contains the correct proposals depending on o.MyRole
func checkForLeaderProposals(t *testing.T, se protocols.SideEffects, o *Objective, td testdata) {

	switch o.MyRole {
	case 0:
		{
			// Alice Proposes to Irene on her right
			rightProposal := consensus_channel.SignedProposal{Proposal: consensus_channel.NewRemoveProposal(o.ToMyRight.Id, FinalTurnNum, o.VId(), td.leftAmount, td.rightAmount)}
			assertProposalSent(t, se, rightProposal, irene)
		}
	case 1:
		{
			// Irene proposes to Bob on her right
			rightProposal := consensus_channel.SignedProposal{Proposal: consensus_channel.NewRemoveProposal(o.ToMyRight.Id, FinalTurnNum, o.VId(), td.leftAmount, td.rightAmount)}
			assertProposalSent(t, se, rightProposal, bob)
		}

	}
}

// signProposal signs a proposal with the given actor's private key
func signProposal(me testactors.Actor, p consensus_channel.Proposal, c *consensus_channel.ConsensusChannel) (consensus_channel.SignedProposal, error) {

	vars := c.ConsensusVars().Clone()
	err := vars.HandleProposal(p)
	if err != nil {
		return consensus_channel.SignedProposal{}, err
	}

	state := vars.AsState(c.FixedPart())
	sig, err := state.Sign(me.PrivateKey)
	if err != nil {
		return consensus_channel.SignedProposal{}, fmt.Errorf("unable to sign state update: %f", err)
	}

	return consensus_channel.SignedProposal{Signature: sig, Proposal: p}, nil
}

// makeOutcome creates an outcome allocating to alice and bob
func makeOutcome(aliceAmount uint, bobAmount uint) outcome.SingleAssetExit {
	return outcome.SingleAssetExit{
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: alice.Destination(),
				Amount:      big.NewInt(int64(aliceAmount)),
			},
			outcome.Allocation{
				Destination: bob.Destination(),
				Amount:      big.NewInt(int64(bobAmount)),
			},
		},
	}
}

type testdata struct {
	vFixed         state.FixedPart
	vFinal         state.State
	initialOutcome outcome.SingleAssetExit
	finalOutcome   outcome.SingleAssetExit
	paid           uint
	leftAmount     *big.Int
	rightAmount    *big.Int
}

// generateTestData generates some test data that can be used in a test
func generateTestData() testdata {
	vFixed := state.FixedPart{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{alice.Address, irene.Address, bob.Address}, // A single hop virtual channel
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
	}

	leftAmount := big.NewInt(6)
	rightAmount := big.NewInt(4)
	initialOutcome := makeOutcome(7, 3)
	finalOutcome := makeOutcome(6, 4)
	paid := uint(1)

	vFinal := state.StateFromFixedAndVariablePart(vFixed, state.VariablePart{IsFinal: true, Outcome: outcome.Exit{finalOutcome}, TurnNum: FinalTurnNum})

	return testdata{vFixed, vFinal, initialOutcome, finalOutcome, paid, leftAmount, rightAmount}
}

// signByOthers signs the state by every participant except my
func signByOthers(my ta.Actor, signedState state.SignedState) state.SignedState {
	if my.Role != 0 {
		_ = signedState.Sign(&alice.PrivateKey)
	}

	if my.Role != 1 {
		_ = signedState.Sign(&irene.PrivateKey)
	}

	if my.Role != 2 {
		_ = signedState.Sign(&bob.PrivateKey)
	}
	return signedState
}

// assertStateSentToEveryone asserts that ses contains a message for every participant but from
func assertStateSentToEveryone(t *testing.T, ses protocols.SideEffects, expected state.SignedState, from testactors.Actor) {
	for _, a := range allActors {
		if a.Role != from.Role {
			assertStateSentTo(t, ses, expected, a)
		}
	}
}

// assertStateSentTo asserts that ses contains a message for the participant
func assertStateSentTo(t *testing.T, ses protocols.SideEffects, expected state.SignedState, to testactors.Actor) {
	for _, msg := range ses.MessagesToSend {
		if bytes.Equal(msg.To[:], to.Address[:]) {
			for _, ss := range msg.SignedStates {
				testhelpers.Equals(t, ss, expected)
			}
		}
	}
}
