package consensus_channel

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/types"
)

func TestNewLeaderChannel(t *testing.T) {
	o := ledgerOutcome()
	initialVars := Vars{Outcome: o.clone(), TurnNum: 0}
	aliceSig, _ := initialVars.AsState(fp()).Sign(alice.PrivateKey)
	bobsSig, _ := initialVars.AsState(fp()).Sign(bob.PrivateKey)
	sigs := [2]state.Signature{aliceSig, bobsSig}

	channel, err := newLeaderChannel(initialVars.AsState(fp()), 0, o.clone(), sigs)
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

	// createSignedProposal generates a proposal given the vars & proposed change
	// The proposal is signed by the given actor, using a generic fixed part
	createSignedProposal := func(vars Vars, p Proposal, actor testactors.Actor, turnNum uint64) SignedProposalVars {
		proposalVars := Vars{TurnNum: vars.TurnNum, Outcome: vars.Outcome.clone()}
		_ = proposalVars.HandleProposal(p)

		state := proposalVars.AsState(fp())
		sig, _ := state.Sign(actor.PrivateKey)

		return SignedProposalVars{SignedProposal{sig, p, turnNum}, proposalVars}
	}

	aliceSignedProposal := func(vars Vars, p Proposal, turnNum uint64) SignedProposalVars {
		return createSignedProposal(vars, p, alice, turnNum)
	}

	bobSignedProposal := func(vars Vars, p Proposal, turnNum uint64) SignedProposalVars {
		return createSignedProposal(vars, p, bob, turnNum)
	}

	cId := fp().ChannelId()
	testChannel := func(lo LedgerOutcome, testProposalQueue []SignedProposalVars) ConsensusChannel {
		vars := Vars{TurnNum: 0, Outcome: lo}
		aliceSig, _ := vars.AsState(fp()).Sign(alice.PrivateKey)
		bobsSig, _ := vars.AsState(fp()).Sign(bob.PrivateKey)
		sigs := [2]state.Signature{aliceSig, bobsSig}

		current := SignedVars{Vars: vars, Signatures: sigs}

		proposalQueue := []SignedProposal{}
		for _, p := range testProposalQueue {
			p.Proposal.LedgerID = cId
			proposalQueue = append(proposalQueue, p.SignedProposal)
		}

		return ConsensusChannel{
			Channel: channel.Channel{
				FixedPart: fp(),
				Id:        cId,
			},
			proposalQueue: proposalQueue,
			current:       current,
		}
	}
	const aAmount = uint64(6)
	const bAmount = uint64(4)
	const amountAdded = aAmount + bAmount

	createAdd := func(chID types.Destination, target types.Destination) Proposal {
		return NewAddProposal(
			chID,
			guarantee(amountAdded, target, alice, bob),
			big.NewInt(int64(amountAdded)),
		)
	}
	createRemove := func(chID types.Destination, target types.Destination) Proposal {
		return NewRemoveProposal(
			chID,
			target,
			big.NewInt(int64(aAmount)),
		)
	}

	// ******* //
	// Propose //
	// ******* //

	testPropose := func(
		channel ConsensusChannel,
		proposal Proposal,
		expectedSp SignedProposal,
		expectedErr error,
	) func(*testing.T) {
		return func(t *testing.T) {
			currentTurnNum := channel.ConsensusTurnNum()
			latest, _ := channel.latestProposedVars()
			latestTurnNum := latest.TurnNum

			sp, err := channel.Propose(proposal, alice.PrivateKey)
			if err != nil {
				if expectedErr == nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !errors.Is(err, expectedErr) {
					t.Fatalf("expected error %v, got %v", expectedErr, err)
				}
				// If we receive an error we don't want to perform the other checks
				return
			}

			if !reflect.DeepEqual(sp, expectedSp) {
				diff := cmp.Diff(sp, expectedSp, cmp.AllowUnexported(Proposal{}, Add{}, Remove{}, Guarantee{}, big.Int{}))
				t.Fatalf("expected signed proposal %v", diff)
			}

			switch proposal.Type() {
			case AddProposal:
				{
					add := proposal.ToAdd
					proposed, _ := channel.IsProposed(add.Guarantee)

					if expectedErr == nil && !proposed {
						t.Fatalf("failed to propose guarantee in happy case")
					}
				}
			case RemoveProposal:
				{
					remove := proposal.ToRemove
					vars, _ := channel.latestProposedVars()

					for target := range vars.Outcome.guarantees {
						if target == remove.Target {
							t.Fatalf("guarantee still present in proposal for target %s", remove.Target)
						}
					}

				}
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
		proposalMade := createAdd(cId, targetChannel)

		expectedSp := aliceSignedProposal(c.current.Vars, proposalMade, 1).SignedProposal
		t.Run(msg, testPropose(c, proposalMade, expectedSp, nil))
	}

	{
		msg := "ok:adding with a non-empty queue"
		startingOutcome := makeOutcome(
			allocation(alice, aBal-amountAdded),
			allocation(bob, bBal),
			guarantee(vAmount, channel1Id, alice, bob),
		)

		p1 := createAdd(cId, types.Destination{2})
		sp1 := aliceSignedProposal(Vars{TurnNum: 1, Outcome: startingOutcome}, p1, 1)
		startingQueue := append(emptyQueue(), sp1)

		c := testChannel(startingOutcome, startingQueue)

		newAdd := Proposal{LedgerID: p1.LedgerID, ToAdd: add(amountAdded, types.Destination{3}, alice, bob)}

		currentlyProposed, _ := c.latestProposedVars()
		expectedSp := aliceSignedProposal(currentlyProposed, newAdd, 2).SignedProposal

		t.Run(msg, testPropose(c, newAdd, expectedSp, nil))
	}
	{
		msg := "ok:adding a remove proposal"
		startingOutcome := makeOutcome(
			allocation(alice, aBal),
			allocation(bob, bBal),
			guarantee(amountAdded, channel1Id, alice, bob),
		)

		c := testChannel(startingOutcome, emptyQueue())

		newRemove := createRemove(cId, channel1Id)

		currentlyProposed, _ := c.latestProposedVars()
		expectedSp := aliceSignedProposal(currentlyProposed, newRemove, 1).SignedProposal

		t.Run(msg, testPropose(c, newRemove, expectedSp, nil))
	}
	{
		msg := "err:adding a remove proposal with invalid target"
		startingOutcome := makeOutcome(
			allocation(alice, aBal),
			allocation(bob, bBal),
		)

		c := testChannel(startingOutcome, emptyQueue())

		newRemove := createRemove(cId, channel1Id)

		t.Run(msg, testPropose(c, newRemove, SignedProposal{}, ErrGuaranteeNotFound))
	}
	{
		msg := "err:adding a remove proposal with too large left/right amounts"
		startingOutcome := makeOutcome(
			allocation(alice, aBal),
			allocation(bob, bBal),
			guarantee(amountAdded, channel1Id, alice, bob),
		)

		c := testChannel(startingOutcome, emptyQueue())

		// LeftAmount > amountAdded
		newRemove := NewRemoveProposal(cId, channel1Id, big.NewInt(int64(amountAdded+1)))

		t.Run(msg, testPropose(c, newRemove, SignedProposal{}, ErrInvalidAmount))
	}

	{
		msg := "err:adding a duplicate proposal"
		startingOutcome := makeOutcome(
			allocation(alice, aBal-amountAdded),
			allocation(bob, bBal),
			guarantee(vAmount, channel1Id, alice, bob),
		)

		proposedChan := types.Destination{2}
		p1 := createAdd(cId, proposedChan)
		sp1 := aliceSignedProposal(Vars{TurnNum: 1, Outcome: startingOutcome}, p1, 1)

		startingQueue := append(emptyQueue(), sp1)

		c := testChannel(startingOutcome, startingQueue)

		duplicateAdd := Proposal{ToAdd: add(amountAdded, proposedChan, alice, bob), LedgerID: c.Id}

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

		p := Proposal{ToAdd: add(amountAdded, proposedChan, alice, bob), LedgerID: c.Id}
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
		t1 := types.Destination{byte(0)}
		t2 := types.Destination{byte(1)}
		t3 := types.Destination{byte(2)}

		p1 := createAdd(cId, t1)
		sp1 := aliceSignedProposal(vars, p1, vars.TurnNum+1)

		p2 := createAdd(cId, t2)
		sp2 := aliceSignedProposal(sp1.Vars, p2, sp1.Vars.TurnNum+1)

		p3 := createAdd(cId, t3)
		sp3 := aliceSignedProposal(sp2.Vars, p3, sp2.Vars.TurnNum+1)

		p4 := createRemove(cId, t3)
		sp4 := aliceSignedProposal(sp3.Vars, p4, sp3.Vars.TurnNum+1)

		return []SignedProposalVars{sp1, sp2, sp3, sp4}
	}

	testUpdateConsensusOk := func(
		counterProposal SignedProposalVars,
	) func(*testing.T) {
		channel := testChannel(startingOutcome, populatedQueue())

		return func(t *testing.T) {
			latest, _ := channel.latestProposedVars()
			latestTurnNum := latest.TurnNum

			err := channel.Receive(counterProposal.SignedProposal)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}

			switch counterProposal.Proposal.Type() {
			case AddProposal:
				{
					g := counterProposal.Proposal.ToAdd.Guarantee
					if !channel.Includes(g) {
						t.Fatalf("failed to fund guarantee given successful counterproposal")
					}

					if proposed, _ := channel.IsProposed(g); proposed {
						t.Fatalf("guarantee still proposed given successful counterproposal")
					}
				}
			case RemoveProposal:
				{
					r := counterProposal.Proposal.ToRemove
					_, foundGuarantee := channel.current.Outcome.guarantees[r.Target]
					if foundGuarantee {
						t.Fatalf("failed to remove guarantee given successful counterproposal")
					}

				}
			}

			if channel.ConsensusTurnNum() != counterProposal.SignedProposal.TurnNum {
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

			err := channel.Receive(counterProposal.SignedProposal)

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

		counterP := bobSignedProposal(signedbyAlice.Vars, signedbyAlice.Proposal, signedbyAlice.SignedProposal.TurnNum)
		t.Run(msg, testUpdateConsensusOk(counterP))
	}

	{ // Receiving a valid (but stale) proposal
		initialVars := Vars{TurnNum: consensusTurnNum, Outcome: startingOutcome.clone()}
		p0 := createAdd(cId, channel1Id)

		counterP := bobSignedProposal(initialVars, p0, 0).SignedProposal
		channel := testChannel(startingOutcome, populatedQueue())

		err := channel.Receive(counterP)
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
		counterP := createSignedProposal(p.Vars, p.Proposal, brian, p.SignedProposal.TurnNum)
		t.Run(msg, testUpdateConsensusErr(counterP, ErrWrongSigner))
	}

	{
		msg := "err:unexpected proposal"
		p := populatedQueue()[2]
		p4 := createAdd(cId, types.Destination{11})
		counterP := bobSignedProposal(p.Vars, p4, p.Vars.TurnNum+10)
		t.Run(msg, testUpdateConsensusErr(counterP, ErrProposalQueueExhausted))
	}

	{
		msg := "err:wrong channel"
		p := populatedQueue()[2]
		p4 := createAdd(types.Destination{}, types.Destination{11}) // blank ChannelID intentionally different than precomputed cId
		counterP := bobSignedProposal(p.Vars, p4, p.Vars.TurnNum+10)
		t.Run(msg, testUpdateConsensusErr(counterP, ErrIncorrectChannelID))
	}
}

func TestRestrictedFollowerMethods(t *testing.T) {
	initialVars := Vars{Outcome: ledgerOutcome(), TurnNum: 0}
	aliceSig, _ := initialVars.AsState(fp()).Sign(alice.PrivateKey)
	bobsSig, _ := initialVars.AsState(fp()).Sign(bob.PrivateKey)
	sigs := [2]state.Signature{aliceSig, bobsSig}

	channel, _ := newLeaderChannel(initialVars.AsState(fp()), 0, ledgerOutcome(), sigs)

	if _, err := channel.SignNextProposal(Proposal{}, alice.PrivateKey); err != ErrNotFollower {
		t.Errorf("Expected error when calling SignNextProposal as a leader, but found none")
	}

	if err := channel.followerReceive(SignedProposal{}); err != ErrNotFollower {
		t.Errorf("Expected error when calling Receive as a leader, but found none")
	}
}
