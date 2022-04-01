package consensus_channel

import (
	"errors"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

func TestCloning(t *testing.T) {
	var existingChannel = types.Destination{1}
	outcome := makeOutcome(
		allocation(alice, 10),
		allocation(bob, 20),
		guarantee(30, existingChannel, alice, bob),
	)
	var uniqueNum int64 = 111 // should be different from the numbers above
	vars := Vars{TurnNum: 9, Outcome: outcome}

	f1 := fingerprint(vars)
	clone1 := vars.Outcome.clone()

	if fingerprint(Vars{TurnNum: vars.TurnNum, Outcome: clone1}) != f1 {
		t.Fatal("vars incorrectly cloned")
	}

	clone1.guarantees[existingChannel].amount.SetInt64(uniqueNum)
	if f1 != fingerprint(vars) {
		t.Fatal("vars shares data with clone")
	}

	clone2 := vars.Outcome.clone()
	clone2.left.amount.SetInt64(uniqueNum)
	if f1 != fingerprint(vars) {
		t.Fatal("vars shares data with clone")
	}

	clone3 := vars.Outcome.clone()
	clone3.right.amount.SetInt64(uniqueNum)
	if f1 != fingerprint(vars) {
		t.Fatal("vars shares data with clone")
	}
}

func TestConsensusChannelFunctionality(t *testing.T) {
	var outcome = makeOutcome(allocation(alice, 10), allocation(bob, 20))
	initialVars := Vars{Outcome: outcome, TurnNum: 0}
	aliceSig, _ := initialVars.AsState(fp()).Sign(alice.PrivateKey)
	bobsSig, _ := initialVars.AsState(fp()).Sign(bob.PrivateKey)
	sigs := [2]state.Signature{aliceSig, bobsSig}

	channel, err := newConsensusChannel(fp(), leader, 0, outcome, sigs)

	if err != nil {
		t.Fatalf("unable to construct a new consensus channel: %v", err)
	}

	_, err = channel.sign(initialVars, bob.PrivateKey)
	if err == nil {
		t.Fatalf("channel should check that signer is participant")
	}

	f := fingerprint(channel.current.Vars)

	latest, err := channel.latestProposedVars()
	if err != nil {
		t.Fatalf("latest proposed vars returned err: %v", err)
	}

	latest.Outcome.left.amount.SetInt64(90210)
	if f != fingerprint(channel.current.Vars) {
		t.Fatalf("latestProposedVars did not return a copy")
	}

	briansSig, _ := initialVars.AsState(fp()).Sign(brian.PrivateKey)
	wrongSigs := [2]state.Signature{sigs[1], briansSig}
	_, err = newConsensusChannel(fp(), leader, 0, outcome, wrongSigs)
	if err == nil {
		t.Fatalf("channel should check that signers are participants")
	}
}

func TestApplyingAddProposalToVars(t *testing.T) {
	channel2 := types.Destination{2}
	channel3 := types.Destination{3}
	aBal := uint64(100)
	bBal := uint64(200)
	vAmount := uint64(5)
	proposal := add(10, vAmount, channel3, alice, bob)
	startingTurnNum := uint64(9)

	// Testing the happy path
	{
		outcome := makeOutcome(
			allocation(alice, aBal),
			allocation(bob, bBal),
			guarantee(vAmount, channel2, alice, bob),
		)

		vars := Vars{TurnNum: startingTurnNum, Outcome: outcome}
		// We should be able to compute the next state
		err := vars.Add(proposal)

		if err != nil {
			t.Fatalf("unable to compute next state: %v", err)
		}

		if vars.TurnNum != startingTurnNum+1 {
			t.Fatalf("incorrect state calculation: %v", err)
		}

		expected := makeOutcome(
			allocation(alice, aBal-vAmount),
			allocation(bob, bBal),
			guarantee(vAmount, channel2, alice, bob),
			guarantee(vAmount, channel3, alice, bob),
		)

		if diff := cmp.Diff(vars.Outcome, expected, cmp.AllowUnexported(expected, Balance{}, big.Int{}, Guarantee{})); diff != "" {
			t.Fatalf("incorrect outcome: %v", diff)
		}
	}

	// Trying to add a guarantee targeting an existing target channel should fail
	{
		outcome := makeOutcome(
			allocation(alice, aBal),
			allocation(bob, bBal),
			guarantee(vAmount, channel2, alice, bob),
		)
		vars := Vars{TurnNum: startingTurnNum, Outcome: outcome}
		duplicateProposal := proposal
		duplicateProposal.target = channel2
		err := vars.Add(duplicateProposal)

		if !errors.Is(err, ErrDuplicateGuarantee) {
			t.Fatalf("expected error when adding duplicate guarantee: %v", err)
		}

	}

	// Proposing a change that depletes a balance should fail
	{
		outcome := makeOutcome(allocation(alice, aBal), allocation(bob, bBal))
		vars := Vars{TurnNum: startingTurnNum, Outcome: outcome}
		largeProposal := proposal
		leftAmount := big.NewInt(0).Set(vars.Outcome.left.amount)
		largeProposal.amount = leftAmount.Add(leftAmount, big.NewInt(1))
		err := vars.Add(largeProposal)
		if !errors.Is(err, ErrInsufficientFunds) {
			t.Fatalf("expected error when adding too large a guarantee: %v", err)
		}
	}
}
