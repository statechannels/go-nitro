package virtualdefund

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel/ledger"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/testactors"
	ta "github.com/statechannels/go-nitro/internal/testactors"
	. "github.com/statechannels/go-nitro/internal/testhelpers"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// generateLedgers generates the left and right ledger channels based on myRole
// The ledger channels will include a guarantee that funds V
func generateLedgers(myRole uint, vId types.Destination) (left, right *ledger.LedgerChannel) {
	switch myRole {
	case 0:
		{
			return nil, prepareLedgerChannel(uint(ledger.Leader), testactors.Alice, testactors.Irene, generateGuarantee(testactors.Alice, testactors.Irene, vId))
		}
	case 1:
		{

			return prepareLedgerChannel(uint(ledger.Follower), testactors.Alice, testactors.Irene, generateGuarantee(testactors.Alice, testactors.Irene, vId)),
				prepareLedgerChannel(uint(ledger.Leader), testactors.Irene, testactors.Bob, generateGuarantee(testactors.Irene, testactors.Bob, vId))

		}
	case 2:
		{

			return prepareLedgerChannel(uint(ledger.Follower), testactors.Irene, testactors.Bob, generateGuarantee(testactors.Irene, testactors.Bob, vId)), nil

		}
	default:
		panic("invalid myRole")
	}
}

// generateGuarantee generates a guarantee for the given participants and vId
func generateGuarantee(left, right testactors.Actor, vId types.Destination) ledger.Guarantee {
	return ledger.NewGuarantee(big.NewInt(10), vId, left.Destination(), right.Destination())

}

// prepareLedgerChannel prepares a ledger channel with a consensus outcome
//  - allocating 0 to left
//  - allocating 0 to right
//  - including the given guarantees
func prepareLedgerChannel(role uint, left, right testactors.Actor, guarantees ...ledger.Guarantee) *ledger.LedgerChannel {
	fp := state.FixedPart{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{left.Address, right.Address},
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
	}

	leftBal := ledger.NewBalance(left.Destination(), big.NewInt(0))
	rightBal := ledger.NewBalance(right.Destination(), big.NewInt(0))

	lo := *ledger.NewLedgerOutcome(types.Address{}, leftBal, rightBal, guarantees)

	signedVars := ledger.SignedVars{Vars: ledger.Vars{Outcome: lo, TurnNum: 1}}
	leftSig, err := signedVars.Vars.AsState(fp).Sign(left.PrivateKey)
	if err != nil {
		panic(err)
	}
	rightSig, err := signedVars.Vars.AsState(fp).Sign(right.PrivateKey)
	if err != nil {
		panic(err)
	}
	sigs := [2]state.Signature{leftSig, rightSig}

	var cc ledger.LedgerChannel

	if role == 0 {
		cc, err = ledger.NewLeaderChannel(fp, 1, lo, sigs)
	} else {
		cc, err = ledger.NewFollowerChannel(fp, 1, lo, sigs)
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
			rightProposal := ledger.SignedProposal{Proposal: generateRemoveProposal(o.ToMyLeft.Id, td)}
			AssertProposalSent(t, se, rightProposal, alice)
		}
	case 2:
		{
			// Bob should accept a proposal from Irene
			rightProposal := ledger.SignedProposal{Proposal: generateRemoveProposal(o.ToMyLeft.Id, td)}
			AssertProposalSent(t, se, rightProposal, irene)
		}

	}
}

// generateProposalsResponses generates the signed proposals that a participant should expect from the other participants
func generateProposalsResponses(myRole uint, vId types.Destination, o *Objective, td testdata) []ledger.SignedProposal {
	switch myRole {
	case 0:
		{
			// Alice expects Irene to accept her proposal
			p := generateRemoveProposal(o.ToMyRight.Id, td)
			sp, err := signProposal(irene, p, o.ToMyRight)
			if err != nil {
				panic(err)
			}

			return []ledger.SignedProposal{sp}
		}

	case 1:
		{
			// Irene expects Alice to send a proposal
			fromAlice := generateRemoveProposal(o.ToMyLeft.Id, td)
			fromAliceSigned, _ := signProposal(alice, fromAlice, o.ToMyLeft)

			// Irene expects Bob to accept her proposal
			fromBob := generateRemoveProposal(o.ToMyRight.Id, td)
			fromBobSigned, _ := signProposal(bob, fromBob, o.ToMyRight)

			return []ledger.SignedProposal{fromAliceSigned, fromBobSigned}
		}
	case 2:
		{
			// Bob expects Irene to send a proposal
			p := generateRemoveProposal(o.ToMyLeft.Id, td)
			sp, err := signProposal(irene, p, o.ToMyLeft)
			if err != nil {
				panic(err)
			}
			return []ledger.SignedProposal{sp}

		}
	default:
		return []ledger.SignedProposal{}
	}
}

// updateProposals updates the ledger channels on the objective with the given proposals by calling Receive
func updateProposals(o *Objective, proposals ...ledger.SignedProposal) {
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

// checkForLeaderProposals checks that the outgoing message contains the correct proposals from the leader of a ledger channel
func checkForLeaderProposals(t *testing.T, se protocols.SideEffects, o *Objective, td testdata) {

	switch o.MyRole {
	case 0:
		{
			// Alice Proposes to Irene on her right
			rightProposal := ledger.SignedProposal{Proposal: generateRemoveProposal(o.ToMyRight.Id, td)}
			AssertProposalSent(t, se, rightProposal, irene)
		}
	case 1:
		{
			// Irene proposes to Bob on her right
			rightProposal := ledger.SignedProposal{Proposal: generateRemoveProposal(o.ToMyRight.Id, td)}
			AssertProposalSent(t, se, rightProposal, bob)
		}

	}
}

// signProposal signs a proposal with the given actor's private key
func signProposal(me testactors.Actor, p ledger.Proposal, c *ledger.LedgerChannel) (ledger.SignedProposal, error) {

	con := c.ConsensusVars()
	vars := con.Clone()
	err := vars.HandleProposal(p)
	if err != nil {
		return ledger.SignedProposal{}, err
	}

	state := vars.AsState(c.FixedPart())
	sig, err := state.Sign(me.PrivateKey)
	if err != nil {
		return ledger.SignedProposal{}, fmt.Errorf("unable to sign state update: %f", err)
	}

	return ledger.SignedProposal{Signature: sig, Proposal: p}, nil
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
	vFinal         state.State
	initialOutcome outcome.SingleAssetExit

	paid uint
	// finalAliceAmount is the amount we expect to be allocated in the ledger to Alice after defunding is complete
	finalAliceAmount uint
	// finalBobAmount is the amount we expect to be allocated in the ledger to Bob after defunding is complete
	finalBobAmount uint
}

// generateRemoveProposal generates a remove proposal for the given channelId and test data
func generateRemoveProposal(cId types.Destination, td testdata) ledger.Proposal {
	vId, _ := td.vFinal.ChannelId()
	return ledger.NewRemoveProposal(cId, FinalTurnNum, vId, big.NewInt(int64(td.finalAliceAmount)), big.NewInt(int64(td.finalBobAmount)))

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

	finalAliceAmount := uint(6)
	finalBobAmount := uint(4)

	paid := uint(1)
	initialOutcome := makeOutcome(finalAliceAmount+paid, finalBobAmount-paid)
	finalOutcome := makeOutcome(finalAliceAmount, finalBobAmount)
	vFinal := state.StateFromFixedAndVariablePart(vFixed, state.VariablePart{IsFinal: true, Outcome: outcome.Exit{finalOutcome}, TurnNum: FinalTurnNum})

	return testdata{vFinal, initialOutcome, paid, finalAliceAmount, finalBobAmount}
}

// signStateByOthers signs the state by every participant except me
func signStateByOthers(me ta.Actor, signedState state.SignedState) state.SignedState {
	if me.Role != 0 {
		_ = signedState.Sign(&alice.PrivateKey)
	}

	if me.Role != 1 {
		_ = signedState.Sign(&irene.PrivateKey)
	}

	if me.Role != 2 {
		_ = signedState.Sign(&bob.PrivateKey)
	}
	return signedState
}
