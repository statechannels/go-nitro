package consensus_channel

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

func TestNewLeaderChannel(t *testing.T) {
	o := ledgerOutcome()
	initialVars := Vars{Outcome: o.clone(), TurnNum: 0}
	aliceSig, _ := initialVars.AsState(fp()).Sign(alice.PrivateKey)
	bobsSig, _ := initialVars.AsState(fp()).Sign(bob.PrivateKey)
	sigs := [2]state.Signature{aliceSig, bobsSig}

	channel, err := NewLeaderChannel(fp(), 0, o.clone(), sigs)
	if err != nil {
		t.Fatal("unable to construct channel")
	}

	vars, _ := channel.latestProposedVars()
	if !vars.equals(initialVars) {
		t.Fatalf("constructed with the wrong initial vars")
	}
}

func TestLeaderChannel(t *testing.T) {
	// SignedProposalVars is a test utility that stores the vars signed along with a SignedProposal
	type SignedProposalVars struct {
		SignedProposal
		Vars
	}

	emptyQueue := func() []SignedProposalVars {
		return []SignedProposalVars{}
	}

	// aliceSignedProposal generates a proposal given the vars & proposed change
	// The proposal is signed by the given actor, using a generic fixed part
	createSignedProposal := func(vars Vars, add Add, actor actor) SignedProposalVars {
		proposalVars := Vars{TurnNum: vars.TurnNum, Outcome: vars.Outcome.clone()}
		_ = proposalVars.Add(add)

		state := proposalVars.AsState(fp())
		sig, _ := state.Sign(actor.PrivateKey)

		return SignedProposalVars{SignedProposal{sig, Proposal{toAdd: add}}, proposalVars}
	}

	aliceSignedProposal := func(vars Vars, add Add) SignedProposalVars {
		return createSignedProposal(vars, add, alice)
	}

	bobSignedProposal := func(vars Vars, add Add) SignedProposalVars {
		return createSignedProposal(vars, add, bob)
	}

	cId, _ := fp().ChannelId()
	testChannel := func(lo LedgerOutcome, testProposalQueue []SignedProposalVars) ConsensusChannel {
		vars := Vars{TurnNum: 0, Outcome: lo}
		aliceSig, _ := vars.AsState(fp()).Sign(alice.PrivateKey)
		bobsSig, _ := vars.AsState(fp()).Sign(bob.PrivateKey)
		sigs := [2]state.Signature{aliceSig, bobsSig}

		current := SignedVars{Vars: vars, Signatures: sigs}

		proposalQueue := []SignedProposal{}
		for _, p := range testProposalQueue {
			proposalQueue = append(proposalQueue, p.SignedProposal)
		}

		return ConsensusChannel{
			fp:            fp(),
			Id:            cId,
			myIndex:       Leader,
			proposalQueue: proposalQueue,
			current:       current,
		}
	}

	const amountAdded = uint64(10)

	createAdd := func(turnNum uint64, target types.Destination) Add {
		return Add{
			turnNum,
			guarantee(amountAdded, target, alice, bob),
			big.NewInt(int64(amountAdded)),
		}
	}

	// ******* //
	// Propose //
	// ******* //

	testPropose := func(
		channel ConsensusChannel,
		add Add,
		expectedSp SignedProposal,
		expectedErr error,
	) func(*testing.T) {
		return func(t *testing.T) {
			currentTurnNum := channel.ConsensusTurnNum()
			latest, _ := channel.latestProposedVars()
			latestTurnNum := latest.TurnNum

			sp, err := channel.Propose(add, alice.PrivateKey)

			if !errors.Is(err, expectedErr) {
				t.Fatalf("expected error %v, got %v", expectedErr, err)
			}

			if !reflect.DeepEqual(sp, expectedSp) {
				diff := cmp.Diff(sp, expectedSp, cmp.AllowUnexported(Proposal{}, Add{}, Guarantee{}, big.Int{}))
				t.Fatalf("expected signed proposal %v", diff)
			}

			proposed, _ := channel.IsProposed(add.Guarantee)

			if expectedErr == nil && !proposed {
				t.Fatalf("failed to propose guarantee in happy case")
			}

			if !errors.Is(err, expectedErr) {
				t.Fatalf("unexpected error: got %v, wanted %v", err, expectedErr)
			}

			if channel.ConsensusTurnNum() != currentTurnNum {
				t.Fatalf("guarantee is not correctly proposed")
			}

			var expectedLatest uint64
			if expectedErr == nil {
				expectedLatest = latestTurnNum + 1
			} else {
				expectedLatest = latestTurnNum
			}

			latest, _ = channel.latestProposedVars()
			if latest.TurnNum != expectedLatest {
				t.Fatalf("turn num malformed")
			}
		}
	}

	{
		msg := "ok:adding with an empty queue"
		startingOutcome := makeOutcome(
			allocation(alice, aBal-amountAdded),
			allocation(bob, bBal),
			guarantee(vAmount, channel1Id, alice, bob),
		)
		c := testChannel(startingOutcome, emptyQueue())
		proposalMade := createAdd(1, targetChannel)
		expectedSp := aliceSignedProposal(c.current.Vars, proposalMade).SignedProposal
		t.Run(msg, testPropose(c, proposalMade, expectedSp, nil))
	}

	{
		msg := "ok:provided turn number is ignored"
		startingOutcome := makeOutcome(
			allocation(alice, aBal-amountAdded),
			allocation(bob, bBal),
			guarantee(vAmount, channel1Id, alice, bob),
		)
		c := testChannel(startingOutcome, emptyQueue())
		proposalMade := createAdd(1, targetChannel)
		expectedSp := aliceSignedProposal(c.current.Vars, proposalMade).SignedProposal

		proposalMade.turnNum = 9001
		t.Run(msg, testPropose(c, proposalMade, expectedSp, nil))
	}

	{
		msg := "ok:adding with a non-empty queue"
		startingOutcome := makeOutcome(
			allocation(alice, aBal-amountAdded),
			allocation(bob, bBal),
			guarantee(vAmount, channel1Id, alice, bob),
		)

		p1 := createAdd(1, types.Destination{2})
		sp1 := aliceSignedProposal(Vars{Outcome: startingOutcome}, p1)
		startingQueue := append(emptyQueue(), sp1)

		c := testChannel(startingOutcome, startingQueue)

		newAdd := add(2, amountAdded, types.Destination{3}, alice, bob)

		currentlyProposed, _ := c.latestProposedVars()
		expectedSp := aliceSignedProposal(currentlyProposed, newAdd).SignedProposal

		t.Run(msg, testPropose(c, newAdd, expectedSp, nil))
	}

	{
		msg := "err:adding a duplicate proposal"
		startingOutcome := makeOutcome(
			allocation(alice, aBal-amountAdded),
			allocation(bob, bBal),
			guarantee(vAmount, channel1Id, alice, bob),
		)

		proposedChan := types.Destination{2}
		p1 := createAdd(1, proposedChan)
		sp1 := aliceSignedProposal(Vars{Outcome: startingOutcome}, p1)

		startingQueue := append(emptyQueue(), sp1)

		c := testChannel(startingOutcome, startingQueue)

		duplicateAdd := add(2, amountAdded, proposedChan, alice, bob)

		t.Run(msg, testPropose(c, duplicateAdd, SignedProposal{}, ErrDuplicateGuarantee))
	}

	{
		msg := "err:overspending"
		startingOutcome := makeOutcome(
			allocation(alice, 0),
			allocation(bob, bBal),
		)

		proposedChan := types.Destination{2}

		c := testChannel(startingOutcome, emptyQueue())

		p := add(2, amountAdded, proposedChan, alice, bob)

		t.Run(msg, testPropose(c, p, SignedProposal{}, ErrInsufficientFunds))
	}

	// // *************** //
	// // UpdateConsensus //
	// // *************** //

	startingOutcome := makeOutcome(
		allocation(alice, aBal-amountAdded),
		allocation(bob, bBal),
	)

	const consensusTurnNum = uint64(0)

	populatedQueue := func() []SignedProposalVars {
		vars := Vars{TurnNum: consensusTurnNum, Outcome: startingOutcome}

		p1 := createAdd(vars.TurnNum+1, types.Destination{byte(vars.TurnNum)})
		sp1 := aliceSignedProposal(vars, p1)

		p2 := createAdd(sp1.TurnNum+1, types.Destination{byte(sp1.TurnNum)})
		sp2 := aliceSignedProposal(sp1.Vars, p2)

		p3 := createAdd(sp2.TurnNum+1, types.Destination{byte(sp2.TurnNum)})
		sp3 := aliceSignedProposal(sp2.Vars, p3)

		return []SignedProposalVars{sp1, sp2, sp3}
	}

	testUpdateConsensusOk := func(
		counterProposal SignedProposalVars,
	) func(*testing.T) {
		channel := testChannel(startingOutcome, populatedQueue())

		return func(t *testing.T) {
			latest, _ := channel.latestProposedVars()
			latestTurnNum := latest.TurnNum

			err := channel.UpdateConsensus(counterProposal.SignedProposal)

			if err != nil {
				t.Fatalf("expected error %v", err)
			}
			g := counterProposal.Proposal.toAdd.Guarantee
			if !channel.Includes(g) {
				t.Fatalf("failed to fund guarantee given successful counterproposal")
			}

			if proposed, _ := channel.IsProposed(g); proposed {
				t.Fatalf("guarantee still proposed given successful counterproposal")
			}

			if channel.ConsensusTurnNum() != counterProposal.Proposal.toAdd.turnNum {
				t.Fatalf("consensus not reached")
			}

			latest, _ = channel.latestProposedVars()
			if latest.TurnNum != latestTurnNum {
				t.Fatalf("latest proposed turn number has changed")
			}
		}
	}

	testUpdateConsensusErr := func(
		counterProposal SignedProposalVars,
		expectedErr error,
	) func(*testing.T) {
		channel := testChannel(startingOutcome, populatedQueue())

		return func(t *testing.T) {
			currentTurnNum := channel.ConsensusTurnNum()
			latest, _ := channel.latestProposedVars()
			latestTurnNum := latest.TurnNum

			err := channel.UpdateConsensus(counterProposal.SignedProposal)

			if !errors.Is(err, expectedErr) {
				t.Fatalf("expected error %v, got %v", expectedErr, err)
			}

			if currentTurnNum != channel.ConsensusTurnNum() {
				t.Fatalf("consensus changed in error case")
			}

			latest, _ = channel.latestProposedVars()
			if latest.TurnNum != latestTurnNum {
				t.Fatalf("latest proposed turn number has changed")
			}

			if !errors.Is(err, expectedErr) {
				t.Fatalf("unexpected error: got %v, wanted %v", err, expectedErr)
			}
		}
	}

	for i, signedbyAlice := range populatedQueue() {
		msg := fmt.Sprintf("ok: receiving a valid counter proposal in position %v", i)

		counterP := bobSignedProposal(signedbyAlice.Vars, signedbyAlice.Proposal.toAdd)
		t.Run(msg, testUpdateConsensusOk(counterP))
	}

	{ // Receiving a valid (but stale) proposal
		initialVars := Vars{TurnNum: consensusTurnNum, Outcome: startingOutcome.clone()}
		p0 := createAdd(0, channel1Id)

		counterP := bobSignedProposal(initialVars, p0).SignedProposal
		channel := testChannel(startingOutcome, populatedQueue())
		err := channel.UpdateConsensus(counterP)
		if err != nil {
			t.Fatalf("unable to update consensus: %v", err)
		}

		if channel.ConsensusTurnNum() != 0 {
			t.Fatalf("incorrectly received stale counterproposal")
		}
	}

	{
		msg := "err:wrong signature"
		p := populatedQueue()[0]
		counterP := createSignedProposal(p.Vars, p.Proposal.toAdd, brian)
		t.Run(msg, testUpdateConsensusErr(counterP, ErrWrongSigner))
	}

	{
		msg := "err:unexpected proposal"
		p := populatedQueue()[2]
		p4 := createAdd(p.TurnNum+10, types.Destination{11})
		counterP := bobSignedProposal(p.Vars, p4)
		t.Run(msg, testUpdateConsensusErr(counterP, ErrProposalQueueExhausted))
	}
}

func TestRestrictedFollowerMethods(t *testing.T) {
	initialVars := Vars{Outcome: ledgerOutcome(), TurnNum: 0}
	aliceSig, _ := initialVars.AsState(fp()).Sign(alice.PrivateKey)
	bobsSig, _ := initialVars.AsState(fp()).Sign(bob.PrivateKey)
	sigs := [2]state.Signature{aliceSig, bobsSig}

	channel, _ := NewLeaderChannel(fp(), 0, ledgerOutcome(), sigs)

	if err := channel.SignNextProposal(Proposal{}, alice.PrivateKey); err != ErrNotFollower {
		t.Errorf("Expected error when calling SignNextProposal as a leader, but found none")
	}

	if err := channel.Receive(SignedProposal{}); err != ErrNotFollower {
		t.Errorf("Expected error when calling Receive as a leader, but found none")
	}
}
