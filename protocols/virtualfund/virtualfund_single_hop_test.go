package virtualfund

import (
	"bytes"
	"fmt"
	"math/big"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
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

// newTestData returns new copies of consistent test data each time it is called
func newTestData() testData {
	var vPreFund = state.State{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{alice.address, p1.address, bob.address}, // A single hop virtual channel
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
		AppData:           []byte{},
		Outcome: outcome.Exit{outcome.SingleAssetExit{
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: alice.destination,
					Amount:      big.NewInt(6),
				},
				outcome.Allocation{
					Destination: bob.destination,
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
	leaderLedgers[alice.destination] = actorLedgers{
		right: prepareConsensusChannel(uint(consensus_channel.Leader), alice, p1),
	}
	leaderLedgers[p1.destination] = actorLedgers{
		left:  prepareConsensusChannel(uint(consensus_channel.Leader), p1, alice),
		right: prepareConsensusChannel(uint(consensus_channel.Leader), p1, bob),
	}
	leaderLedgers[bob.destination] = actorLedgers{
		left: prepareConsensusChannel(uint(consensus_channel.Leader), bob, p1),
	}

	followerLedgers := make(map[types.Destination]actorLedgers)
	followerLedgers[alice.destination] = actorLedgers{
		right: prepareConsensusChannel(uint(consensus_channel.Follower), alice, p1),
	}
	followerLedgers[p1.destination] = actorLedgers{
		left:  prepareConsensusChannel(uint(consensus_channel.Follower), alice, p1),
		right: prepareConsensusChannel(uint(consensus_channel.Follower), p1, bob),
	}
	followerLedgers[bob.destination] = actorLedgers{
		left: prepareConsensusChannel(uint(consensus_channel.Follower), p1, bob),
	}

	return testData{vPreFund, vPostFund, leaderLedgers, followerLedgers}
}

type Tester func(t *testing.T)

func testNew(a actor) Tester {
	return func(t *testing.T) {
		td := newTestData()
		lookup := td.leaderLedgers
		vPreFund := td.vPreFund

		// Assert that a valid set of constructor args does not result in an error
		o, err := constructFromState(
			false,
			vPreFund,
			a.address,
			lookup[a.destination].left,
			lookup[a.destination].right,
		)
		if err != nil {
			t.Fatal(err)
		}

		switch a.role {
		case alice.role:
			assert(t, o.ToMyLeft == nil, "left connection should be nil")
			assert(t, diffFromCorrectConnection(o.ToMyRight, alice, p1) == "", "incorrect connection")
		case p1.role:
			assert(t, diffFromCorrectConnection(o.ToMyLeft, alice, p1) == "", "incorrect connection")
			assert(t, diffFromCorrectConnection(o.ToMyRight, p1, bob) == "", "incorrect connection")
		case bob.role:
			assert(t, diffFromCorrectConnection(o.ToMyLeft, p1, bob) == "", "incorrect connection")
			assert(t, o.ToMyRight == nil, "right connection should be nil")
		}
	}
}

// diffFromCorrectConnection compares the guarantee stored on a connection with
// the guarantee we expect, given the expected left and right actors
func diffFromCorrectConnection(c *Connection, left, right actor) string {
	td := newTestData()
	vPreFund := td.vPreFund

	Id, _ := vPreFund.FixedPart().ChannelId()

	// HACK: This should really be comparing GuaranteeInfo, but GuaranteeInfo
	// contains types.Funds amounts in their LeftAmount and RightAmount fields.
	// I am not sure how these types are meant to be used, and am
	// comparing the _guarantees_ that we expect to include, instead of the GuaranteeInfo

	expectedAmount := big.NewInt(0).Set(vPreFund.VariablePart().Outcome[0].TotalAllocated())
	want := consensus_channel.NewGuarantee(expectedAmount, Id, left.destination, right.destination)
	got := c.getExpectedGuarantee()

	return compareGuarantees(want, got)
}

func TestNew(t *testing.T) {
	for _, a := range allActors {
		msg := fmt.Sprintf("Testing new as %v", a.name)
		t.Run(msg, testNew(a))
	}
}

func testCloneAs(my actor) Tester {
	return func(t *testing.T) {
		td := newTestData()
		vPreFund := td.vPreFund
		ledgers := td.leaderLedgers

		o, _ := constructFromState(false, vPreFund, my.address, ledgers[my.destination].left, ledgers[my.destination].right)

		clone := o.clone()

		if diff := compareObjectives(o, clone); diff != "" {
			t.Fatalf("Clone: mismatch (-want +got):\n%s", diff)
		}

	}
}

func TestClone(t *testing.T) {
	for _, a := range allActors {
		msg := fmt.Sprintf("Testing clone as %v", a.name)
		t.Run(msg, testCloneAs(a))
	}
}

func collectPeerSignaturesOnSetupState(V *channel.SingleHopVirtualChannel, myRole uint, prefund bool) {
	var state state.State
	if prefund {
		state = V.PreFundState()
	} else {
		state = V.PostFundState()
	}

	if myRole != alice.role {
		aliceSig, _ := state.Sign(alice.privateKey)
		V.AddStateWithSignature(state, aliceSig)
	}
	if myRole != p1.role {
		p1Sig, _ := state.Sign(p1.privateKey)
		V.AddStateWithSignature(state, p1Sig)
	}
	if myRole != bob.role {
		bobSig, _ := state.Sign(bob.privateKey)
		V.AddStateWithSignature(state, bobSig)
	}
}

func TestCrankAsAlice(t *testing.T) {
	my := alice
	td := newTestData()
	vPreFund := td.vPreFund
		ledgers := td.leaderLedgers
	var s, _ = constructFromState(false, vPreFund, my.address, ledgers[my.destination].left, ledgers[my.destination].right)
	// Assert that cranking an unapproved objective returns an error
	_, _, _, err := s.Crank(&my.privateKey)
	assert(t, err != nil, `Expected error when cranking unapproved objective, but got nil`)

	// Approve the objective, so that the rest of the test cases can run.
	o := s.Approve().(*Objective)

	// To test the finite state progression, we are going to progressively mutate o
	// And then crank it to see which "pause point" (WaitingFor) we end up at.
	// NOTE: Because crank returns a protocools.Objective interface, after each crank we
	// need to remember to convert the result back to a virtualfund.Objective struct

	// Initial Crank
	oObj, effects, waitingFor, err := o.Crank(&my.privateKey)
	o = oObj.(*Objective)

	expectedSignedState := state.NewSignedState(o.V.PreFundState())
	mySig, _ := o.V.PreFundState().Sign(my.privateKey)
	_ = expectedSignedState.AddSignature(mySig)

	ok(t, err)
	equals(t, waitingFor, WaitingForCompletePrefund)
	assertStateSentTo(t, effects, expectedSignedState, bob)
	assertStateSentTo(t, effects, expectedSignedState, p1)

	// Manually progress the extended state by collecting prefund signatures
	collectPeerSignaturesOnSetupState(o.V, my.role, true)

	// Cranking should move us to the next waiting point, update the ledger channel, and alter the extended state to reflect that
	// TODO: Check that ledger channel is updated as expected
	oObj, effects, waitingFor, err = o.Crank(&my.privateKey)
	o = oObj.(*Objective)

	p := consensus_channel.NewAddProposal(o.ToMyRight.Channel.Id, 2, o.ToMyRight.getExpectedGuarantee(), big.NewInt(6))
	sp := consensus_channel.SignedProposal{Proposal: p}
	assertProposalSent(t, effects, sp, p1)
	ok(t, err)
	equals(t, waitingFor, WaitingForCompleteFunding)

	// Check idempotency
	emptySideEffects := protocols.SideEffects{}
	oObj, effects, waitingFor, err = o.Crank(&my.privateKey)
	o = oObj.(*Objective)
	ok(t, err)
	equals(t, effects, emptySideEffects)
	equals(t, waitingFor, WaitingForCompleteFunding)

	// If Alice had received a signed counterproposal, she should proceed to postFundSetup
	guaranteeFundingV := consensus_channel.NewGuarantee(big.NewInt(10), o.V.Id, alice.destination, p1.destination)
	o.ToMyRight.Channel = prepareConsensusChannel(my.role, alice, p1, guaranteeFundingV)

	oObj, effects, waitingFor, err = o.Crank(&my.privateKey)
	o = oObj.(*Objective)

	postFS := state.NewSignedState(o.V.PostFundState())
	mySig, _ = postFS.State().Sign(my.privateKey)
	_ = postFS.AddSignature(mySig)

	ok(t, err)
	equals(t, waitingFor, WaitingForCompletePostFund)
	assertStateSentTo(t, effects, postFS, bob)
}

func TestCrankAsBob(t *testing.T) {
	my := bob
	td := newTestData()
	vPreFund := td.vPreFund
	ledgers := td.leaderLedgers
	var s, _ = constructFromState(false, vPreFund, my.address, ledgers[my.destination].left, ledgers[my.destination].right)
	// Assert that cranking an unapproved objective returns an error
	_, _, _, err := s.Crank(&my.privateKey)
	assert(t, err != nil, `Expected error when cranking unapproved objective, but got nil`)

	// Approve the objective, so that the rest of the test cases can run.
	o := s.Approve().(*Objective)

	// To test the finite state progression, we are going to progressively mutate o
	// And then crank it to see which "pause point" (WaitingFor) we end up at.
	// NOTE: Because crank returns a protocools.Objective interface, after each crank we
	// need to remember to convert the result back to a virtualfund.Objective struct

	// Initial Crank
	oObj, effects, waitingFor, err := o.Crank(&my.privateKey)
	o = oObj.(*Objective)

	expectedSignedState := state.NewSignedState(o.V.PreFundState())
	mySig, _ := o.V.PreFundState().Sign(my.privateKey)
	_ = expectedSignedState.AddSignature(mySig)

	ok(t, err)
	equals(t, waitingFor, WaitingForCompletePrefund)
	assertStateSentTo(t, effects, expectedSignedState, alice)
	assertStateSentTo(t, effects, expectedSignedState, p1)

	// Manually progress the extended state by collecting prefund signatures
	collectPeerSignaturesOnSetupState(o.V, my.role, true)

	// Cranking should move us to the next waiting point, update the ledger channel, and alter the extended state to reflect that
	// TODO: Check that ledger channel is updated as expected
	oObj, effects, waitingFor, err = o.Crank(&my.privateKey)
	o = oObj.(*Objective)

	p := consensus_channel.NewAddProposal(o.ToMyLeft.Channel.Id, 2, o.ToMyLeft.getExpectedGuarantee(), big.NewInt(6))
	sp := consensus_channel.SignedProposal{Proposal: p}
	assertProposalSent(t, effects, sp, p1)
	ok(t, err)
	equals(t, waitingFor, WaitingForCompleteFunding)

	// Check idempotency
	emptySideEffects := protocols.SideEffects{}
	oObj, effects, waitingFor, err = o.Crank(&my.privateKey)
	o = oObj.(*Objective)
	ok(t, err)
	equals(t, effects, emptySideEffects)
	equals(t, waitingFor, WaitingForCompleteFunding)

	// If Bob had received a signed counterproposal, she should proceed to postFundSetup
	guaranteeFundingV := consensus_channel.NewGuarantee(big.NewInt(10), o.V.Id, p1.destination, bob.destination)
	o.ToMyLeft.Channel = prepareConsensusChannel(uint(consensus_channel.Leader), bob, p1, guaranteeFundingV)

	oObj, effects, waitingFor, err = o.Crank(&my.privateKey)
	o = oObj.(*Objective)

	postFS := state.NewSignedState(o.V.PostFundState())
	mySig, _ = postFS.State().Sign(my.privateKey)
	_ = postFS.AddSignature(mySig)

	ok(t, err)
	equals(t, waitingFor, WaitingForCompletePostFund)
	assertStateSentTo(t, effects, postFS, p1)
}

// Copied from https://github.com/benbjohnson/testing

// makeRed sets the colour to red when printed
const makeRed = "\033[31m"

// makeBlack sets the colour to black when printed.
// as it is intended to be used at the end of a string, it also adds two linebreaks
const makeBlack = "\033[39m\n\n"

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(makeRed+"%s:%d: "+msg+makeBlack, append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(makeRed+"%s:%d: unexpected error: %s"+makeBlack, filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(makeRed+"%s:%d:\n\n\texp: %#v\n\n\tgot: %#v"+makeBlack, filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

// The following assertions are inspired by the ok, assert and equals above

// assertSideEffectsContainsMessageWith fails the test instantly if the supplied side effects does not contain a message for the supplied actor with the supplied expected signed state.
func assertProposalSent(t *testing.T, ses protocols.SideEffects, sp consensus_channel.SignedProposal, to actor) {
	_, file, line, _ := runtime.Caller(1)
	if len(ses.MessagesToSend) != 1 {
		fmt.Printf(makeRed+"%s:%d:\n\n\texpected one message"+makeBlack, filepath.Base(file), line)
		t.FailNow()
	}
	if len(ses.MessagesToSend[0].SignedProposals) != 1 {
		fmt.Printf(makeRed+"%s:%d:\n\n\texpected one signed proposal"+makeBlack, filepath.Base(file), line)
		t.FailNow()
	}

	msg := ses.MessagesToSend[0]
	sent := msg.SignedProposals[0]

	if !reflect.DeepEqual(sent.Proposal, sp.Proposal) {
		fmt.Printf(makeRed+"%s:%d:\n\n\texp: %+v\n\n\tgot: %+v"+makeBlack, filepath.Base(file), line, sent.Proposal, sp.Proposal)
		t.FailNow()
	}

	if !bytes.Equal(msg.To[:], to.address[:]) {
		fmt.Printf(makeRed+"%s:%d:\n\n\texp: %#v\n\n\tgot: %#v"+makeBlack, filepath.Base(file), line, msg.To.String(), to.address.String())
		t.FailNow()
	}
}

// assertMessageSentTo asserts that ses contains a message
func assertStateSentTo(t *testing.T, ses protocols.SideEffects, expected state.SignedState, to actor) {
	_, file, line, _ := runtime.Caller(1)
	for _, msg := range ses.MessagesToSend {
		for _, ss := range msg.SignedStates {
			correctAddress := bytes.Equal(msg.To[:], to.address[:])

			if correctAddress {
				diff := compareStates(ss, expected)
				if diff == "" {
					return
				}

				fmt.Printf("\033[31m%s:%d:\n\n\tincorrect state\n\ndiff: %v", filepath.Base(file), line, diff)
				t.FailNow()
			}
		}
	}

	fmt.Printf(makeRed+"%s:%d:\n\n\tside effects do not incude signed state"+makeBlack, filepath.Base(file), line)
	t.FailNow()
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
