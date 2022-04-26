package virtualfund

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/testactors"
	actors "github.com/statechannels/go-nitro/internal/testactors"
	. "github.com/statechannels/go-nitro/internal/testhelpers"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type actorLedgers struct {
	left  *consensus_channel.ConsensusChannel
	right *consensus_channel.ConsensusChannel
}
type ledgerLookup map[types.Destination]actorLedgers
type testData struct {
	vPreFund        state.State
	vPostFund       state.State
	leaderLedgers   ledgerLookup
	followerLedgers ledgerLookup
}

var alice, p1, bob actors.Actor = actors.Alice, actors.Irene, actors.Bob
var allActors []actors.Actor = []actors.Actor{alice, p1, bob}

// newTestData returns new copies of consistent test data each time it is called
func newTestData() testData {
	var vPreFund = state.State{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{alice.Address, p1.Address, bob.Address},
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
		AppData:           []byte{},
		Outcome: outcome.Exit{outcome.SingleAssetExit{
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: alice.Destination(),
					Amount:      big.NewInt(6),
				},
				outcome.Allocation{
					Destination: bob.Destination(),
					Amount:      big.NewInt(4),
				},
			},
		}},
		TurnNum: 0,
		IsFinal: false,
	}
	var vPostFund = vPreFund.Clone()
	vPostFund.TurnNum = 1

	leaderLedgers := make(map[types.Destination]actorLedgers)
	leaderLedgers[alice.Destination()] = actorLedgers{
		right: prepareConsensusChannel(uint(consensus_channel.Leader), alice, p1),
	}
	leaderLedgers[p1.Destination()] = actorLedgers{
		left:  prepareConsensusChannel(uint(consensus_channel.Leader), p1, alice),
		right: prepareConsensusChannel(uint(consensus_channel.Leader), p1, bob),
	}
	leaderLedgers[bob.Destination()] = actorLedgers{
		left: prepareConsensusChannel(uint(consensus_channel.Leader), bob, p1),
	}

	followerLedgers := make(map[types.Destination]actorLedgers)
	followerLedgers[alice.Destination()] = actorLedgers{
		right: prepareConsensusChannel(uint(consensus_channel.Follower), alice, p1),
	}
	followerLedgers[p1.Destination()] = actorLedgers{
		left:  prepareConsensusChannel(uint(consensus_channel.Follower), alice, p1),
		right: prepareConsensusChannel(uint(consensus_channel.Follower), p1, bob),
	}
	followerLedgers[bob.Destination()] = actorLedgers{
		left: prepareConsensusChannel(uint(consensus_channel.Follower), p1, bob),
	}

	return testData{vPreFund, vPostFund, leaderLedgers, followerLedgers}
}

type Tester func(t *testing.T)

func testNew(a actors.Actor) Tester {
	return func(t *testing.T) {
		td := newTestData()
		lookup := td.leaderLedgers
		vPreFund := td.vPreFund

		// Assert that a valid set of constructor args does not result in an error
		o, err := constructFromState(
			false,
			vPreFund,
			a.Address,
			lookup[a.Destination()].left,
			lookup[a.Destination()].right,
		)
		if err != nil {
			t.Fatal(err)
		}

		switch a.Role {
		case alice.Role:
			Assert(t, o.ToMyLeft == nil, "left connection should be nil")
			Assert(t, diffFromCorrectConnection(o.ToMyRight, alice, p1) == "", "incorrect connection")
		case p1.Role:
			Assert(t, diffFromCorrectConnection(o.ToMyLeft, alice, p1) == "", "incorrect connection")
			Assert(t, diffFromCorrectConnection(o.ToMyRight, p1, bob) == "", "incorrect connection")
		case bob.Role:
			Assert(t, diffFromCorrectConnection(o.ToMyLeft, p1, bob) == "", "incorrect connection")
			Assert(t, o.ToMyRight == nil, "right connection should be nil")
		}
	}
}

// diffFromCorrectConnection compares the guarantee stored on a connection with
// the guarantee we expect, given the expected left and right actors
func diffFromCorrectConnection(c *Connection, left, right actors.Actor) string {
	td := newTestData()
	vPreFund := td.vPreFund

	Id, _ := vPreFund.FixedPart().ChannelId()

	// HACK: This should really be comparing GuaranteeInfo, but GuaranteeInfo
	// contains types.Funds amounts in their LeftAmount and RightAmount fields.
	// I am not sure how these types are meant to be used, and am
	// comparing the _guarantees_ that we expect to include, instead of the GuaranteeInfo

	expectedAmount := big.NewInt(0).Set(vPreFund.VariablePart().Outcome[0].TotalAllocated())
	want := consensus_channel.NewGuarantee(expectedAmount, Id, left.Destination(), right.Destination())
	got := c.getExpectedGuarantee()

	return compareGuarantees(want, got)
}

func TestNew(t *testing.T) {
	for _, a := range allActors {
		msg := fmt.Sprintf("Testing new as %v", a.Name)
		t.Run(msg, testNew(a))
	}
}

func testCloneAs(my actors.Actor) Tester {
	return func(t *testing.T) {
		td := newTestData()
		vPreFund := td.vPreFund
		ledgers := td.leaderLedgers

		o, _ := constructFromState(false, vPreFund, my.Address, ledgers[my.Destination()].left, ledgers[my.Destination()].right)

		clone := o.clone()

		if diff := compareObjectives(o, clone); diff != "" {
			t.Fatalf("Clone: mismatch (-want +got):\n%s", diff)
		}

	}
}

func TestClone(t *testing.T) {
	for _, a := range allActors {
		msg := fmt.Sprintf("Testing clone as %v", a.Name)
		t.Run(msg, testCloneAs(a))
	}
}

func cloneAndSignSetupStateByPeers(v *channel.SingleHopVirtualChannel, myRole uint, prefund bool) *channel.SingleHopVirtualChannel {
	var state state.State
	if prefund {
		state = v.PreFundState()
	} else {
		state = v.PostFundState()
	}

	updated := v.Clone()

	if myRole != alice.Role {
		aliceSig, _ := state.Sign(alice.PrivateKey)
		updated.AddStateWithSignature(state, aliceSig)
	}
	if myRole != p1.Role {
		p1Sig, _ := state.Sign(p1.PrivateKey)
		updated.AddStateWithSignature(state, p1Sig)
	}
	if myRole != bob.Role {
		bobSig, _ := state.Sign(bob.PrivateKey)
		updated.AddStateWithSignature(state, bobSig)
	}
	return updated
}

// TestCrankAsAlice tests the behaviour from a end-user's point of view when they are a leader in the ledger channel
func TestCrankAsAlice(t *testing.T) {
	my := alice
	td := newTestData()
	vPreFund := td.vPreFund
	ledgers := td.leaderLedgers
	var s, _ = constructFromState(false, vPreFund, my.Address, ledgers[my.Destination()].left, ledgers[my.Destination()].right)
	// Assert that cranking an unapproved objective returns an error
	_, _, _, err := s.Crank(&my.PrivateKey)
	Assert(t, err != nil, `Expected error when cranking unapproved objective, but got nil`)

	// Approve the objective, so that the rest of the test cases can run.
	o := s.Approve().(*Objective)

	// To test the finite state progression, we are going to progressively mutate o
	// And then crank it to see which "pause point" (WaitingFor) we end up at.
	// NOTE: Because crank returns a protocools.Objective interface, after each crank we
	// need to remember to convert the result back to a virtualfund.Objective struct

	// Initial Crank
	oObj, effects, waitingFor, err := o.Crank(&my.PrivateKey)
	o = oObj.(*Objective)

	expectedSignedState := state.NewSignedState(o.V.PreFundState())
	mySig, _ := o.V.PreFundState().Sign(my.PrivateKey)
	_ = expectedSignedState.AddSignature(mySig)

	Ok(t, err)
	Equals(t, waitingFor, WaitingForCompletePrefund)
	assertStateSentTo(t, effects, expectedSignedState, bob)
	assertStateSentTo(t, effects, expectedSignedState, p1)

	// Update the objective with prefund signatures
	c := cloneAndSignSetupStateByPeers(o.V, my.Role, true)
	e := protocols.ObjectiveEvent{ObjectiveId: o.Id(), SignedStates: []state.SignedState{c.SignedPreFundState()}}
	oObj, err = o.Update(e)
	o = oObj.(*Objective)
	Ok(t, err)

	// Cranking should move us to the next waiting point, update the ledger channel, and alter the extended state to reflect that
	// TODO: Check that ledger channel is updated as expected
	oObj, effects, waitingFor, err = o.Crank(&my.PrivateKey)
	o = oObj.(*Objective)

	p := consensus_channel.NewAddProposal(o.ToMyRight.Channel.Id, 2, o.ToMyRight.getExpectedGuarantee(), big.NewInt(6))
	sp := consensus_channel.SignedProposal{Proposal: p, Signature: consensusStateSignatures(alice, p1, o.ToMyRight.getExpectedGuarantee())[0]}
	Ok(t, err)
	assertOneProposalSent(t, effects, sp, p1)
	Equals(t, waitingFor, WaitingForCompleteFunding)

	// Check idempotency
	emptySideEffects := protocols.SideEffects{}
	oObj, effects, waitingFor, err = o.Crank(&my.PrivateKey)
	o = oObj.(*Objective)
	Ok(t, err)
	Equals(t, effects, emptySideEffects)
	Equals(t, waitingFor, WaitingForCompleteFunding)

	// If Alice had received a signed counterproposal, she should proceed to postFundSetup
	sp = consensus_channel.SignedProposal{Proposal: p, Signature: consensusStateSignatures(alice, p1, o.ToMyRight.getExpectedGuarantee())[1]}
	e = protocols.ObjectiveEvent{ObjectiveId: o.Id(), SignedProposals: []consensus_channel.SignedProposal{sp}}
	oObj, err = o.Update(e)
	o = oObj.(*Objective)
	Ok(t, err)

	oObj, effects, waitingFor, err = o.Crank(&my.PrivateKey)
	o = oObj.(*Objective)

	postFS := state.NewSignedState(o.V.PostFundState())
	mySig, _ = postFS.State().Sign(my.PrivateKey)
	_ = postFS.AddSignature(mySig)

	Ok(t, err)
	Equals(t, waitingFor, WaitingForCompletePostFund)
	assertStateSentTo(t, effects, postFS, bob)
}

// TestCrankAsBob tests the behaviour from a end-user's point of view when they are a follower in the ledger channel
func TestCrankAsBob(t *testing.T) {
	my := bob
	td := newTestData()
	vPreFund := td.vPreFund
	ledgers := td.followerLedgers
	var s, _ = constructFromState(false, vPreFund, my.Address, ledgers[my.Destination()].left, ledgers[my.Destination()].right)
	// Assert that cranking an unapproved objective returns an error
	_, _, _, err := s.Crank(&my.PrivateKey)
	Assert(t, err != nil, `Expected error when cranking unapproved objective, but got nil`)

	// Approve the objective, so that the rest of the test cases can run.
	o := s.Approve().(*Objective)

	// To test the finite state progression, we are going to progressively mutate o
	// And then crank it to see which "pause point" (WaitingFor) we end up at.
	// NOTE: Because crank returns a protocools.Objective interface, after each crank we
	// need to remember to convert the result back to a virtualfund.Objective struct

	// Initial Crank
	oObj, effects, waitingFor, err := o.Crank(&my.PrivateKey)
	o = oObj.(*Objective)

	expectedSignedState := state.NewSignedState(o.V.PreFundState())
	mySig, _ := o.V.PreFundState().Sign(my.PrivateKey)
	_ = expectedSignedState.AddSignature(mySig)

	Ok(t, err)
	Equals(t, waitingFor, WaitingForCompletePrefund)
	assertStateSentTo(t, effects, expectedSignedState, alice)
	assertStateSentTo(t, effects, expectedSignedState, p1)

	// Update the objective with prefund signatures
	c := cloneAndSignSetupStateByPeers(o.V, my.Role, true)
	e := protocols.ObjectiveEvent{ObjectiveId: o.Id(), SignedStates: []state.SignedState{c.SignedPreFundState()}}
	oObj, err = o.Update(e)
	o = oObj.(*Objective)
	Ok(t, err)

	// Cranking should move us to the next waiting point, update the ledger channel, and alter the extended state to reflect that
	// TODO: Check that ledger channel is updated as expected
	oObj, effects, waitingFor, err = o.Crank(&my.PrivateKey)
	o = oObj.(*Objective)

	emptySideEffects := protocols.SideEffects{}
	Ok(t, err)
	Equals(t, effects, emptySideEffects)
	Equals(t, waitingFor, WaitingForCompleteFunding)

	// Check idempotency
	oObj, effects, waitingFor, err = o.Crank(&my.PrivateKey)
	o = oObj.(*Objective)
	Ok(t, err)
	Equals(t, effects, emptySideEffects)
	Equals(t, waitingFor, WaitingForCompleteFunding)

	// If Bob had received a signed counterproposal, he should proceed to postFundSetup
	p := consensus_channel.NewAddProposal(o.ToMyLeft.Channel.Id, 2, o.ToMyLeft.getExpectedGuarantee(), big.NewInt(6))
	sp := consensus_channel.SignedProposal{Proposal: p, Signature: consensusStateSignatures(p1, bob, o.ToMyLeft.getExpectedGuarantee())[0]}
	e = protocols.ObjectiveEvent{ObjectiveId: o.Id(), SignedProposals: []consensus_channel.SignedProposal{sp}}
	oObj, err = o.Update(e)
	o = oObj.(*Objective)
	Ok(t, err)

	oObj, effects, waitingFor, err = o.Crank(&my.PrivateKey)
	o = oObj.(*Objective)

	postFS := state.NewSignedState(o.V.PostFundState())
	mySig, _ = postFS.State().Sign(my.PrivateKey)
	_ = postFS.AddSignature(mySig)

	Ok(t, err)
	Equals(t, waitingFor, WaitingForCompletePostFund)
	assertStateSentTo(t, effects, postFS, p1)
	sp.Signature = consensusStateSignatures(p1, bob, o.ToMyLeft.getExpectedGuarantee())[1]
	assertOneProposalSent(t, effects, sp, p1)
}

// TestCrankAsP1 tests the behaviour from an intermediary's point of view when they are a leader in one ledger channel and a follower in the other
func TestCrankAsP1(t *testing.T) {
	my := p1
	td := newTestData()
	vPreFund := td.vPreFund
	left := td.leaderLedgers[my.Destination()].left
	right := td.followerLedgers[my.Destination()].right
	var s, _ = constructFromState(false, vPreFund, my.Address, left, right)
	// Assert that cranking an unapproved objective returns an error
	_, _, _, err := s.Crank(&my.PrivateKey)
	Assert(t, err != nil, `Expected error when cranking unapproved objective, but got nil`)

	// Approve the objective, so that the rest of the test cases can run.
	o := s.Approve().(*Objective)

	// To test the finite state progression, we are going to progressively mutate o
	// And then crank it to see which "pause point" (WaitingFor) we end up at.
	// NOTE: Because crank returns a protocools.Objective interface, after each crank we
	// need to remember to convert the result back to a virtualfund.Objective struct

	// Initial Crank
	oObj, effects, waitingFor, err := o.Crank(&my.PrivateKey)
	o = oObj.(*Objective)

	expectedSignedState := state.NewSignedState(o.V.PreFundState())
	mySig, _ := o.V.PreFundState().Sign(my.PrivateKey)
	_ = expectedSignedState.AddSignature(mySig)

	Ok(t, err)
	Equals(t, waitingFor, WaitingForCompletePrefund)
	assertStateSentTo(t, effects, expectedSignedState, alice)
	assertStateSentTo(t, effects, expectedSignedState, bob)

	// Update the objective with prefund signatures
	c := cloneAndSignSetupStateByPeers(o.V, my.Role, true)
	e := protocols.ObjectiveEvent{ObjectiveId: o.Id(), SignedStates: []state.SignedState{c.SignedPreFundState()}}
	oObj, err = o.Update(e)
	o = oObj.(*Objective)
	Ok(t, err)

	// Cranking should move us to the next waiting point, update the ledger channel, and alter the extended state to reflect that
	oObj, effects, waitingFor, err = o.Crank(&my.PrivateKey)
	o = oObj.(*Objective)

	p := consensus_channel.NewAddProposal(o.ToMyLeft.Channel.Id, 2, o.ToMyLeft.getExpectedGuarantee(), big.NewInt(6))
	sp := consensus_channel.SignedProposal{Proposal: p, Signature: consensusStateSignatures(p1, alice, o.ToMyLeft.getExpectedGuarantee())[0]}
	Ok(t, err)
	assertOneProposalSent(t, effects, sp, alice)
	Equals(t, waitingFor, WaitingForCompleteFunding)

	// Check idempotency
	emptySideEffects := protocols.SideEffects{}
	oObj, effects, waitingFor, err = o.Crank(&my.PrivateKey)
	o = oObj.(*Objective)
	Ok(t, err)
	Equals(t, effects, emptySideEffects)
	Equals(t, waitingFor, WaitingForCompleteFunding)

	// If P1 had received a signed counterproposal, she should proceed to postFundSetup
	p = consensus_channel.NewAddProposal(o.ToMyLeft.Channel.Id, 2, o.ToMyLeft.getExpectedGuarantee(), big.NewInt(6))
	sp = consensus_channel.SignedProposal{Proposal: p, Signature: consensusStateSignatures(p1, alice, o.ToMyLeft.getExpectedGuarantee())[1]}
	e = protocols.ObjectiveEvent{ObjectiveId: o.Id(), SignedProposals: []consensus_channel.SignedProposal{sp}}
	oObj, err = o.Update(e)
	o = oObj.(*Objective)
	Ok(t, err)

	oObj, effects, waitingFor, err = o.Crank(&my.PrivateKey)
	o = oObj.(*Objective)

	postFS := state.NewSignedState(o.V.PostFundState())
	mySig, _ = postFS.State().Sign(my.PrivateKey)
	_ = postFS.AddSignature(mySig)

	Ok(t, err)

	// We need to receive a proposal from Bob before funding is completed!
	Equals(t, waitingFor, WaitingForCompleteFunding)
	Equals(t, effects, emptySideEffects)
}

// assertOneProposalSent fails the test instantly if the supplied side effects does not contain a message for the supplied actor with the supplied expected signed proposal.
func assertOneProposalSent(t *testing.T, ses protocols.SideEffects, sp consensus_channel.SignedProposal, to actors.Actor) {
	numProposals := 0
	for _, msg := range ses.MessagesToSend {
		if len(msg.SignedProposals) > 0 {

			msg := ses.MessagesToSend[0]
			sent := msg.SignedProposals[0]

			Assert(t, len(ses.MessagesToSend[0].SignedProposals) == 1, "exp: %+v\n\n\tgot%+v", sent.Proposal, sp.Proposal)
			Assert(t, bytes.Equal(msg.To[:], to.Address[:]), "exp: %+v\n\n\tgot%+v", msg.To.String(), to.Address.String())
			Assert(t, compareSignedProposals(sp, sent), "exp: %+v\n\n\tgot%+v", sp, sent)
			numProposals++
		}
	}
	Assert(t, numProposals == 1, "expected 1 proposal but got %d", numProposals)
}

// assertMessageSentTo asserts that ses contains a message
func assertStateSentTo(t *testing.T, ses protocols.SideEffects, expected state.SignedState, to testactors.Actor) {
	found := false
	for _, msg := range ses.MessagesToSend {
		for _, ss := range msg.SignedStates {
			correctAddress := bytes.Equal(msg.To[:], to.Address[:])

			if correctAddress {
				diff := compareStates(ss, expected)
				Assert(t, diff == "", "incorrect state\n\ndiff: %v", diff)
				found = true
				break
			}
		}
	}
	Assert(t, found, "side effects do not include signed state")

}

func compareStates(a, b state.SignedState) string {
	return cmp.Diff(&a, &b,
		cmp.AllowUnexported(
			big.Int{},
			state.SignedState{},
			state.Signature{},
			state.State{},
		),
	)
}

func compareSignedProposals(a, b consensus_channel.SignedProposal) bool {
	return cmp.Equal(&a, &b,
		cmp.AllowUnexported(
			consensus_channel.Add{},
			consensus_channel.Remove{},
			consensus_channel.Guarantee{},
			big.Int{},
		),
	)
}
