package virtualfund

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"

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
	vPreFund  state.State
	vPostFund state.State
	ledgers   ledgerLookup
}

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

	ledgers := make(map[types.Destination]actorLedgers)
	ledgers[alice.destination] = actorLedgers{
		right: prepareConsensusChannel(uint(consensus_channel.Leader), alice, p1),
	}
	ledgers[p1.destination] = actorLedgers{
		left:  prepareConsensusChannel(uint(consensus_channel.Follower), alice, p1),
		right: prepareConsensusChannel(uint(consensus_channel.Leader), p1, bob),
	}
	ledgers[bob.destination] = actorLedgers{
		left: prepareConsensusChannel(uint(consensus_channel.Follower), p1, bob),
	}

	return testData{vPreFund, vPostFund, ledgers}
}

func assertNilConnection(t *testing.T, c *Connection) {
	if c != nil {
		t.Fatalf("TestNew: unexpected connection")
	}
}

func assertCorrectConnection(t *testing.T, c *Connection, left, right actor) {
	td := newTestData()
	vPreFund := td.vPreFund

	Id, _ := vPreFund.FixedPart().ChannelId()

	expectedAmount := big.NewInt(0).Set(vPreFund.VariablePart().Outcome[0].TotalAllocated())
	want := consensus_channel.NewGuarantee(expectedAmount, Id, left.destination, right.destination)
	got := c.getExpectedGuarantee()
	if diff := compareGuarantees(want, got); diff != "" {
		t.Fatalf("TestNew: expectedGuarantee mismatch (-want +got):\n%s", diff)
	}
}

func testNew(t *testing.T, a actor) {
	td := newTestData()
	lookup := td.ledgers
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
		assertNilConnection(t, o.ToMyLeft)
		assertCorrectConnection(t, o.ToMyRight, alice, p1)
	case p1.role:
		assertCorrectConnection(t, o.ToMyLeft, alice, p1)
		assertCorrectConnection(t, o.ToMyRight, p1, bob)
	case bob.role:
		assertCorrectConnection(t, o.ToMyLeft, p1, bob)
		assertNilConnection(t, o.ToMyRight)
	}

}

func TestNew(t *testing.T) {
	for _, a := range allActors {
		testNew(t, a)
	}
}

func testClone(t *testing.T, my actor) {
	td := newTestData()
	vPreFund := td.vPreFund
	ledgers := td.ledgers

	o, _ := constructFromState(false, vPreFund, my.address, ledgers[my.destination].left, ledgers[my.destination].right)

	clone := o.clone()

	if diff := compareObjectives(o, clone); diff != "" {
		t.Fatalf("Clone: mismatch (-want +got):\n%s", diff)
	}
}

func TestClone(t *testing.T) {
	for _, a := range allActors {
		testClone(t, a)
	}
}

// assertSideEffectsContainsMessageWith fails the test instantly if the supplied side effects does not contain a message for the supplied actor with the supplied expected signed state.
func assertSideEffectsContainsMessageWith (ses protocols.SideEffects, expectedSignedState state.SignedState, to actor, t *testing.T) {
	for _, msg := range ses.MessagesToSend {
		for _, ss := range msg.SignedStates {
			if reflect.DeepEqual(ss, expectedSignedState) && bytes.Equal(msg.To[:], to.address[:]) {
				return
			}
		}
	}
	t.Fatalf("side effects %v do not contain signed state %v for %v", ses, expectedSignedState, to)
}

// assertSideEffectsContainsMessageWith calls assertSideEffectsContainsMessageWith for all peers of the actor with role myRole.
func assertSideEffectsContainsMessagesForPeersWith(ses protocols.SideEffects, expectedSignedState state.SignedState, myRole uint, t *testing.T) {
	if myRole != alice.role {
		assertSideEffectsContainsMessageWith(ses, expectedSignedState, alice, t)
	}
	if myRole != p1.role {
		assertSideEffectsContainsMessageWith(ses, expectedSignedState, p1, t)
	}
	if myRole != bob.role {
		assertSideEffectsContainsMessageWith(ses, expectedSignedState, bob, t)
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

func testCrank(t *testing.T, my actor) {
	td := newTestData()
	vPreFund := td.vPreFund
	ledgers := td.ledgers
	var s, _ = constructFromState(false, vPreFund, my.address, ledgers[my.destination].left, ledgers[my.destination].right) // todo: #420 deprecate TwoPartyLedgers
	// Assert that cranking an unapproved objective returns an error
	if _, _, _, err := s.Crank(&my.privateKey); err == nil {
		t.Fatal(`Expected error when cranking unapproved objective, but got nil`)
	}

	// Approve the objective, so that the rest of the test cases can run.
	o := s.Approve().(*Objective)
	// To test the finite state progression, we are going to progressively mutate o
	// And then crank it to see which "pause point" (WaitingFor) we end up at.

	// Initial Crank
	oObj, got, waitingFor, err := o.Crank(&my.privateKey)
	o = oObj.(*Objective)
	if err != nil {
		t.Fatal(err)
	}
	if waitingFor != WaitingForCompletePrefund {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForCompletePrefund, waitingFor)
	}

	expectedSignedState := state.NewSignedState(o.V.PreFundState())
	mySig, _ := o.V.PreFundState().Sign(my.privateKey)
	_ = expectedSignedState.AddSignature(mySig)
	assertSideEffectsContainsMessagesForPeersWith(got, expectedSignedState, my.role, t)

	// Manually progress the extended state by collecting prefund signatures
	collectPeerSignaturesOnSetupState(o.V, my.role, true)

	// Cranking should move us to the next waiting point, update the ledger channel, and alter the extended state to reflect that
	// TODO: Check that ledger channel is updated as expected
	oObj, got, waitingFor, _ = o.Crank(&my.privateKey)

	if waitingFor != WaitingForCompleteFunding {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForCompleteFunding, waitingFor)
	}

	o = oObj.(*Objective)

	// ...

	// Cranking now should not generate side effects, because we already did that
	oObj, got, waitingFor, err = o.Crank(&my.privateKey)
	o = oObj.(*Objective)
	if err != nil {
		t.Fatal(err)
	}
	if waitingFor != WaitingForCompletePostFund {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForCompletePostFund, waitingFor)
	}

	// Check that the messsages contain the expected ledger acceptances
	// We only expect an acceptance in the left ledger channel as we will be the follower in that ledger channel
	switch my.role {
	case 1:
		{
			// supported, _ := o.ToMyLeft.Channel.LatestSupportedState()
			// expectedSignedState := state.NewSignedState(supported)
			// _ = expectedSignedState.Sign(&my.privateKey)

			assertSideEffectsContainsMessageWith(got, expectedSignedState, alice, t)

		}
	case 2:
		{
			// supported, _ := o.ToMyLeft.Channel.LatestSupportedState()
			// expectedSignedState := state.NewSignedState(supported)
			// _ = expectedSignedState.Sign(&my.privateKey)

			assertSideEffectsContainsMessageWith(got, expectedSignedState, p1, t)
		}
	}

	expectedSignedState = state.NewSignedState(o.V.PostFundState())
	mySig, _ = o.V.PostFundState().Sign(my.privateKey)
	_ = expectedSignedState.AddSignature(mySig)
	assertSideEffectsContainsMessagesForPeersWith(got, expectedSignedState, my.role, t)

	// Manually progress the extended state by collecting postfund signatures
	collectPeerSignaturesOnSetupState(o.V, my.role, false)

	// This should be the final crank...
	_, _, waitingFor, err = o.Crank(&my.privateKey)
	if err != nil {
		t.Fatal(err)
	}
	if waitingFor != WaitingForNothing {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForNothing, waitingFor)
	}

}
