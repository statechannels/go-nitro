package virtualfund

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	ta "github.com/statechannels/go-nitro/internal/testactors"
	. "github.com/statechannels/go-nitro/internal/testhelpers"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type actorLedgers struct {
	left  *consensus_channel.ConsensusChannel
	right *consensus_channel.ConsensusChannel
}
type (
	ledgerLookup map[types.Destination]actorLedgers
	testData     struct {
		vPreFund        state.State
		vPostFund       state.State
		leaderLedgers   ledgerLookup
		followerLedgers ledgerLookup
	}
)

var (
	alice, p1, bob ta.Actor   = ta.Alice, ta.Irene, ta.Bob
	allActors      []ta.Actor = []ta.Actor{alice, p1, bob}
)

// newTestData returns new copies of consistent test data each time it is called
func newTestData() testData {
	vPreFund := state.State{
		Participants:      []types.Address{alice.Address(), p1.Address(), bob.Address()},
		ChannelNonce:      0,
		AppDefinition:     types.Address{},
		ChallengeDuration: 45,
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
	vPostFund := vPreFund.Clone()
	vPostFund.TurnNum = 1

	leaderLedgers := make(map[types.Destination]actorLedgers)
	leaderLedgers[alice.Destination()] = actorLedgers{
		right: prepareConsensusChannel(uint(consensus_channel.Leader), alice, p1, alice),
	}
	leaderLedgers[p1.Destination()] = actorLedgers{
		left:  prepareConsensusChannel(uint(consensus_channel.Leader), p1, alice, alice),
		right: prepareConsensusChannel(uint(consensus_channel.Leader), p1, bob, p1),
	}
	leaderLedgers[bob.Destination()] = actorLedgers{
		left: prepareConsensusChannel(uint(consensus_channel.Leader), bob, p1, p1),
	}

	followerLedgers := make(map[types.Destination]actorLedgers)
	followerLedgers[alice.Destination()] = actorLedgers{
		right: prepareConsensusChannel(uint(consensus_channel.Follower), p1, alice, alice),
	}
	followerLedgers[p1.Destination()] = actorLedgers{
		left:  prepareConsensusChannel(uint(consensus_channel.Follower), alice, p1, alice),
		right: prepareConsensusChannel(uint(consensus_channel.Follower), bob, p1, p1),
	}
	followerLedgers[bob.Destination()] = actorLedgers{
		left: prepareConsensusChannel(uint(consensus_channel.Follower), p1, bob, p1),
	}

	return testData{vPreFund, vPostFund, leaderLedgers, followerLedgers}
}

type Tester func(t *testing.T)

func testNew(a ta.Actor) Tester {
	return func(t *testing.T) {
		td := newTestData()
		lookup := td.leaderLedgers
		vPreFund := td.vPreFund

		// Assert that a valid set of constructor args does not result in an error
		o, err := constructFromState(
			false,
			vPreFund,
			a.Address(),
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
func diffFromCorrectConnection(c *Connection, left, right ta.Actor) string {
	td := newTestData()
	vPreFund := td.vPreFund

	Id := vPreFund.FixedPart().ChannelId()

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

func testCloneAs(my ta.Actor) Tester {
	return func(t *testing.T) {
		td := newTestData()
		vPreFund := td.vPreFund
		ledgers := td.leaderLedgers

		o, _ := constructFromState(false, vPreFund, my.Address(), ledgers[my.Destination()].left, ledgers[my.Destination()].right)

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

func cloneAndSignSetupStateByPeers(v channel.VirtualChannel, myRole uint, prefund bool) *channel.VirtualChannel {
	withSigs := v.Clone()

	var state state.State
	if prefund {
		state = withSigs.PreFundState()
	} else {
		state = withSigs.PostFundState()
	}

	if myRole != alice.Role {
		aliceSig, _ := state.Sign(alice.PrivateKey)
		withSigs.AddStateWithSignature(state, aliceSig)
	}
	if myRole != p1.Role {
		p1Sig, _ := state.Sign(p1.PrivateKey)
		withSigs.AddStateWithSignature(state, p1Sig)
	}
	if myRole != bob.Role {
		bobSig, _ := state.Sign(bob.PrivateKey)
		withSigs.AddStateWithSignature(state, bobSig)
	}
	return withSigs
}

func TestMisaddressedUpdate(t *testing.T) {
	var (
		td      = newTestData()
		ledgers = td.leaderLedgers
		vfo, _  = constructFromState(false, td.vPreFund, alice.Address(), ledgers[alice.Destination()].left, ledgers[alice.Destination()].right)
		event   = protocols.ObjectivePayload{
			ObjectiveId: "this-is-not-correct",
		}
	)

	if _, err := vfo.Update(event); err == nil {
		t.Fatal("expected error updating vfo with objective ID mismatch, but found none")
	}
}

// TestCrankAsAlice tests the behaviour from a end-user's point of view when they are a leader in the ledger channel
func TestCrankAsAlice(t *testing.T) {
	var (
		my       = alice
		td       = newTestData()
		vPreFund = td.vPreFund
		ledgers  = td.leaderLedgers
		s, _     = constructFromState(false, vPreFund, my.Address(), ledgers[my.Destination()].left, ledgers[my.Destination()].right)
	)
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
	c := cloneAndSignSetupStateByPeers(*o.V, my.Role, true)
	ss := c.SignedPreFundState()
	e, err := protocols.CreateObjectivePayload(o.Id(), SignedStatePayload, &ss)
	Ok(t, err)
	oObj, err = o.Update(e)
	o = oObj.(*Objective)
	Ok(t, err)

	assertSupportedPrefund(o, t)

	// Cranking should move us to the next waiting point, update the ledger channel, and alter the extended state to reflect that
	// TODO: Check that ledger channel is updated as expected
	oObj, effects, waitingFor, err = o.Crank(&my.PrivateKey)
	o = oObj.(*Objective)

	p := consensus_channel.NewAddProposal(o.ToMyRight.Channel.Id, o.ToMyRight.getExpectedGuarantee(), big.NewInt(6))
	sp := consensus_channel.SignedProposal{Proposal: p, Signature: consensusStateSignatures(alice, p1, o.ToMyRight.getExpectedGuarantee())[0], TurnNum: 2}
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
	sp = consensus_channel.SignedProposal{Proposal: p, Signature: consensusStateSignatures(alice, p1, o.ToMyRight.getExpectedGuarantee())[1], TurnNum: 2}

	oObj, err = o.ReceiveProposal(sp)
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
	var (
		my       = bob
		td       = newTestData()
		vPreFund = td.vPreFund
		ledgers  = td.followerLedgers
		s, _     = constructFromState(false, vPreFund, my.Address(), ledgers[my.Destination()].left, ledgers[my.Destination()].right)
	)
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
	c := cloneAndSignSetupStateByPeers(*o.V, my.Role, true)

	e, err := protocols.CreateObjectivePayload(o.Id(), SignedStatePayload, c.SignedPreFundState())
	Ok(t, err)

	oObj, err = o.Update(e)
	o = oObj.(*Objective)
	Ok(t, err)

	assertSupportedPrefund(o, t)

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
	p := consensus_channel.NewAddProposal(o.ToMyLeft.Channel.Id, o.ToMyLeft.getExpectedGuarantee(), big.NewInt(6))
	sp := consensus_channel.SignedProposal{Proposal: p, Signature: consensusStateSignatures(p1, bob, o.ToMyLeft.getExpectedGuarantee())[0], TurnNum: 2}

	oObj, err = o.ReceiveProposal(sp)
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
	var (
		my       = p1
		td       = newTestData()
		vPreFund = td.vPreFund
		left     = td.leaderLedgers[my.Destination()].left
		right    = td.followerLedgers[my.Destination()].right
		s, _     = constructFromState(false, vPreFund, my.Address(), left, right)
	)
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
	c := cloneAndSignSetupStateByPeers(*o.V, my.Role, true)
	e, err := protocols.CreateObjectivePayload(o.Id(), SignedStatePayload, c.SignedPreFundState())
	Ok(t, err)
	oObj, err = o.Update(e)
	o = oObj.(*Objective)
	Ok(t, err)

	assertSupportedPrefund(o, t)

	// Cranking should move us to the next waiting point, update the ledger channel, and alter the extended state to reflect that
	oObj, effects, waitingFor, err = o.Crank(&my.PrivateKey)
	o = oObj.(*Objective)

	p := consensus_channel.NewAddProposal(o.ToMyLeft.Channel.Id, o.ToMyLeft.getExpectedGuarantee(), big.NewInt(6))
	sp := consensus_channel.SignedProposal{Proposal: p, Signature: consensusStateSignatures(p1, alice, o.ToMyLeft.getExpectedGuarantee())[0], TurnNum: 2}
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
	p = consensus_channel.NewAddProposal(o.ToMyLeft.Channel.Id, o.ToMyLeft.getExpectedGuarantee(), big.NewInt(6))
	sp = consensus_channel.SignedProposal{Proposal: p, Signature: consensusStateSignatures(p1, alice, o.ToMyLeft.getExpectedGuarantee())[1], TurnNum: 2}

	oObj, err = o.ReceiveProposal(sp)
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

// assertSupportedPrefund checks that all three participants have signed the prefund. It
// is used to manually inspect the objective after Update receives counterparty signatures.
func assertSupportedPrefund(o *Objective, t *testing.T) {
	if !o.V.OffChain.SignedStateForTurnNum[0].HasSignatureForParticipant(alice.Role) {
		t.Fatal(`Objective prefund state not signed by alice`)
	}
	if !o.V.OffChain.SignedStateForTurnNum[0].HasSignatureForParticipant(bob.Role) {
		t.Fatal(`Objective prefund state not signed by bob`)
	}
	if !o.V.OffChain.SignedStateForTurnNum[0].HasSignatureForParticipant(p1.Role) {
		t.Fatal(`Objective prefund state not signed by p1`)
	}
}

// assertOneProposalSent fails the test instantly if the supplied side effects does not contain a message for the supplied actor with the supplied expected signed proposal.
func assertOneProposalSent(t *testing.T, ses protocols.SideEffects, sp consensus_channel.SignedProposal, to ta.Actor) {
	numProposals := 0
	for _, msg := range ses.MessagesToSend {
		if len(msg.LedgerProposals) > 0 {

			msg := ses.MessagesToSend[0]
			sent := msg.LedgerProposals[0]
			toAddress := to.Address()

			Assert(t, len(ses.MessagesToSend[0].LedgerProposals) == 1, "exp: %+v\n\n\tgot: %+v", sent.Proposal, sp.Proposal)
			Assert(t, bytes.Equal(msg.To[:], toAddress[:]), "exp: %+v\n\n\tgot: %+v", msg.To.String(), to.Address().String())
			Assert(t, compareSignedProposals(sp, sent), "exp: %+v\n\n\tgot: %+v", sp, sent)
			numProposals++
		}
	}
	Assert(t, numProposals == 1, "expected 1 proposal but got %d", numProposals)
}

// assertMessageSentTo asserts that ses contains a message
func assertStateSentTo(t *testing.T, ses protocols.SideEffects, expected state.SignedState, to ta.Actor) {
	found := false
	for _, msg := range ses.MessagesToSend {
		toAddress := to.Address()
		correctAddress := bytes.Equal(msg.To[:], toAddress[:])
		if correctAddress {
			for _, p := range msg.ObjectivePayloads {

				ss := state.SignedState{}
				err := json.Unmarshal(p.PayloadData, &ss)
				if err != nil {
					panic(err)
				}
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

func compareObjectives(a, b Objective) string {
	return cmp.Diff(&a, &b,
		cmp.AllowUnexported(
			Objective{},
			channel.Channel{},
			big.Int{},
			state.SignedState{},
			consensus_channel.ConsensusChannel{},
			consensus_channel.Vars{},
			consensus_channel.LedgerOutcome{},
			consensus_channel.Balance{},
		),
	)
}

func compareGuarantees(a, b consensus_channel.Guarantee) string {
	return cmp.Diff(&a, &b,
		cmp.AllowUnexported(
			consensus_channel.Guarantee{},
			big.Int{},
		),
	)
}
