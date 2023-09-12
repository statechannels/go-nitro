package virtualdefund

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	ta "github.com/statechannels/go-nitro/internal/testactors"
	. "github.com/statechannels/go-nitro/internal/testhelpers"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const TEST_CHAIN_ID = 1337

// generateLedgers generates the left and right ledger channels based on myRole
// The ledger channels will include a guarantee that funds V
func generateLedgers(myRole uint, vId types.Destination) (left, right *consensus_channel.ConsensusChannel) {
	switch myRole {
	case 0:
		{
			return nil, prepareConsensusChannel(uint(consensus_channel.Leader), ta.Alice, ta.Irene, generateGuarantee(ta.Alice, ta.Irene, vId))
		}
	case 1:
		{
			return prepareConsensusChannel(uint(consensus_channel.Follower), ta.Alice, ta.Irene, generateGuarantee(ta.Alice, ta.Irene, vId)),
				prepareConsensusChannel(uint(consensus_channel.Leader), ta.Irene, ta.Bob, generateGuarantee(ta.Irene, ta.Bob, vId))
		}
	case 2:
		{
			return prepareConsensusChannel(uint(consensus_channel.Follower), ta.Irene, ta.Bob, generateGuarantee(ta.Irene, ta.Bob, vId)), nil
		}
	default:
		panic("invalid myRole")
	}
}

// generateStoreGetters generates mocks for some store methods
func generateStoreGetters(myRole uint, vId types.Destination, vFinal state.State) (GetChannelByIdFunction, GetTwoPartyConsensusLedgerFunction) {
	left, right := generateLedgers(myRole, vId)
	fun1 := func(id types.Destination) (*channel.Channel, bool) {
		c, err := channel.New(vFinal, myRole)
		if err != nil {
			return &channel.Channel{}, false
		}
		return c, true
	}
	fun2 := func(address types.Address) (*consensus_channel.ConsensusChannel, bool) {
		if left != nil && (left.Participants()[0] == address || left.Participants()[1] == address) {
			return left, true
		}
		if right != nil && (right.Participants()[0] == address || right.Participants()[1] == address) {
			return right, true
		}
		return &consensus_channel.ConsensusChannel{}, false
	}
	return fun1, fun2
}

// generateGuarantee generates a guarantee for the given participants and vId
func generateGuarantee(left, right ta.Actor, vId types.Destination) consensus_channel.Guarantee {
	return consensus_channel.NewGuarantee(big.NewInt(10), vId, left.Destination(), right.Destination())
}

// prepareConsensusChannel prepares a consensus channel with a consensus outcome
//   - allocating 0 to left
//   - allocating 0 to right
//   - including the given guarantees
func prepareConsensusChannel(role uint, left, right ta.Actor, guarantees ...consensus_channel.Guarantee) *consensus_channel.ConsensusChannel {
	fp := state.FixedPart{
		Participants:      []types.Address{left.Address(), right.Address()},
		ChannelNonce:      0,
		AppDefinition:     types.Address{},
		ChallengeDuration: 45,
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
			rightProposal := consensus_channel.SignedProposal{Proposal: generateRemoveProposal(o.ToMyLeft.Id, td), TurnNum: 2}
			AssertProposalSent(t, se, rightProposal, alice)
		}
	case 2:
		{
			// Bob should accept a proposal from Irene
			rightProposal := consensus_channel.SignedProposal{Proposal: generateRemoveProposal(o.ToMyLeft.Id, td), TurnNum: 2}
			AssertProposalSent(t, se, rightProposal, irene)
		}

	}
}

// generateProposalsResponses generates the signed proposals that a participant should expect from the other participants
func generateProposalsResponses(myRole uint, vId types.Destination, o *Objective, td testdata) []consensus_channel.SignedProposal {
	switch myRole {
	case 0:
		{
			// Alice expects Irene to accept her proposal
			p := generateRemoveProposal(o.ToMyRight.Id, td)
			sp, err := signProposal(irene, p, o.ToMyRight, 2)
			if err != nil {
				panic(err)
			}

			return []consensus_channel.SignedProposal{sp}
		}

	case 1:
		{
			// Irene expects Alice to send a proposal
			fromAlice := generateRemoveProposal(o.ToMyLeft.Id, td)
			fromAliceSigned, _ := signProposal(alice, fromAlice, o.ToMyLeft, 2)

			// Irene expects Bob to accept her proposal
			fromBob := generateRemoveProposal(o.ToMyRight.Id, td)
			fromBobSigned, _ := signProposal(bob, fromBob, o.ToMyRight, 2)

			return []consensus_channel.SignedProposal{fromAliceSigned, fromBobSigned}
		}
	case 2:
		{
			// Bob expects Irene to send a proposal
			p := generateRemoveProposal(o.ToMyLeft.Id, td)
			sp, err := signProposal(irene, p, o.ToMyLeft, 2)
			if err != nil {
				panic(err)
			}
			return []consensus_channel.SignedProposal{sp}

		}
	default:
		return []consensus_channel.SignedProposal{}
	}
}

// checkForLeaderProposals checks that the outgoing message contains the correct proposals from the leader of a consensus channel
func checkForLeaderProposals(t *testing.T, se protocols.SideEffects, o *Objective, td testdata) {
	switch o.MyRole {
	case 0:
		{
			// Alice Proposes to Irene on her right
			rightProposal := consensus_channel.SignedProposal{Proposal: generateRemoveProposal(o.ToMyRight.Id, td), TurnNum: 2}
			AssertProposalSent(t, se, rightProposal, irene)
		}
	case 1:
		{
			// Irene proposes to Bob on her right
			rightProposal := consensus_channel.SignedProposal{Proposal: generateRemoveProposal(o.ToMyRight.Id, td), TurnNum: 2}
			AssertProposalSent(t, se, rightProposal, bob)
		}

	}
}

// signProposal signs a proposal with the given actor's private key
func signProposal(me ta.Actor, p consensus_channel.Proposal, c *consensus_channel.ConsensusChannel, turnNum uint64) (consensus_channel.SignedProposal, error) {
	con := c.ConsensusVars()
	vars := con.Clone()
	err := vars.HandleProposal(p)
	if err != nil {
		return consensus_channel.SignedProposal{}, err
	}

	state := vars.AsState(c.FixedPart())
	sig, err := state.Sign(me.PrivateKey)
	if err != nil {
		return consensus_channel.SignedProposal{}, fmt.Errorf("unable to sign state update: %f", err)
	}

	return consensus_channel.SignedProposal{Signature: sig, Proposal: p, TurnNum: turnNum}, nil
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
	vInitial state.State
	vFinal   state.State

	paid uint
	// finalAliceAmount is the amount we expect to be allocated in the ledger to Alice after defunding is complete
	finalAliceAmount uint
	// finalBobAmount is the amount we expect to be allocated in the ledger to Bob after defunding is complete
	finalBobAmount uint
}

// generateRemoveProposal generates a remove proposal for the given channelId and test data
func generateRemoveProposal(cId types.Destination, td testdata) consensus_channel.Proposal {
	vId := td.vFinal.ChannelId()
	return consensus_channel.NewRemoveProposal(cId, vId, big.NewInt(int64(td.finalAliceAmount)))
}

// generateTestData generates some test data that can be used in a test
func generateTestData() testdata {
	vFixed := state.FixedPart{
		Participants:      []types.Address{alice.Address(), irene.Address(), bob.Address()}, // A single hop virtual channel
		ChannelNonce:      0,
		AppDefinition:     types.Address{},
		ChallengeDuration: 45,
	}

	finalAliceAmount := uint(6)
	finalBobAmount := uint(4)

	paid := uint(1)
	initialOutcome := makeOutcome(finalAliceAmount+paid, finalBobAmount-paid)
	finalOutcome := makeOutcome(finalAliceAmount, finalBobAmount)

	vInitial := state.StateFromFixedAndVariablePart(vFixed, state.VariablePart{IsFinal: true, Outcome: outcome.Exit{initialOutcome}, TurnNum: 1})
	vFinal := state.StateFromFixedAndVariablePart(vFixed, state.VariablePart{IsFinal: true, Outcome: outcome.Exit{finalOutcome}, TurnNum: FinalTurnNum})

	return testdata{vInitial, vFinal, paid, finalAliceAmount, finalBobAmount}
}

// signStateByOthers signs the state by every participant except me
func signStateByOthers(me ta.Actor, signedState state.SignedState) state.SignedState {
	if me.Role != 0 {
		SignState(&signedState, &alice.PrivateKey)
	}

	if me.Role != 1 {
		SignState(&signedState, &irene.PrivateKey)
	}

	if me.Role != 2 {
		SignState(&signedState, &bob.PrivateKey)
	}
	return signedState
}
