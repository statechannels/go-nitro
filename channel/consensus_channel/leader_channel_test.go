package consensus_channel

import (
	"errors"
	"reflect"
	"testing"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

func TestLeaderChannel(t *testing.T) {
	existingChannel := types.Destination{1}
	aBal := uint64(200)
	bBal := uint64(300)
	vAmount := uint64(5)

	initialVars := Vars{Outcome: ledgerOutcome(), TurnNum: 0}
	aliceSig, _ := initialVars.AsState(fp()).Sign(alice.PrivateKey)
	bobsSig, _ := initialVars.AsState(fp()).Sign(bob.PrivateKey)
	sigs := [2]state.Signature{aliceSig, bobsSig}

	channel, err := NewLeaderChannel(fp(), 0, ledgerOutcome(), sigs)
	if err != nil {
		t.Fatal("unable to construct channel")
	}

	amountAdded := uint64(10)

	if initialVars.TurnNum != 0 {
		t.Fatal("initialized with non-zero turn number")
	}

	p := Proposal{toAdd: add(1, amountAdded, targetChannel, alice, bob)}

	sp, err := channel.Propose(p.toAdd, alice.PrivateKey)
	if err != nil {
		t.Fatalf("failed to add proposal: %v", err)
	}

	success, _ := channel.IsProposed(p.toAdd.Guarantee)
	if !success {
		t.Fatal("incorrect latest proposed vars")
	}
	if channel.ConsensusTurnNum() != 0 || channel.Includes(p.toAdd.Guarantee) {
		t.Fatal("consensus incorrectly updated")
	}

	outcomeSigned := makeOutcome(
		allocation(alice, aBal-amountAdded),
		allocation(bob, bBal),
		guarantee(vAmount, existingChannel, alice, bob),
		guarantee(amountAdded, targetChannel, alice, bob),
	)
	stateSigned := Vars{TurnNum: 1, Outcome: outcomeSigned}
	sig, _ := stateSigned.AsState(fp()).Sign(alice.PrivateKey)
	expected := SignedProposal{Proposal: p, Signature: sig}

	if !reflect.DeepEqual(sp, expected) {
		t.Fatalf("propose failed")
	}

	thirdChannel := types.Destination{3}
	p2 := p
	p2.toAdd.target = thirdChannel
	g2 := p2.toAdd.Guarantee
	secondSigned, err := channel.Propose(p2.toAdd, alice.PrivateKey)
	if err != nil {
		t.Fatalf("failed to add another proposal: %v", err)
	}
	if secondSigned.Proposal.toAdd.turnNum != 2 {
		t.Fatalf("incorrect proposal generated")
	}

	success, _ = channel.IsProposed(g2)
	if !success {
		t.Fatal("incorrect latest proposed vars")
	}
	if channel.ConsensusTurnNum() != 0 {
		t.Fatal("consensus incorrectly updated")
	}
	if channel.Includes(g2) {
		t.Fatal("consensus incorrectly updated")
	}

	latest, _ := channel.latestProposedVars()
	counterSig2, _ := latest.AsState(fp()).Sign(bob.PrivateKey)

	p3 := p
	p3.toAdd.target = types.Destination{4}
	g3 := p3.toAdd.Guarantee
	thirdSigned, _ := channel.Propose(p3.toAdd, alice.PrivateKey)

	p2Returned := SignedProposal{
		Proposal:  secondSigned.Proposal,
		Signature: counterSig2,
	}

	// A counter signature is received on a proposal (but not the latest proposal)
	err = channel.UpdateConsensus(p2Returned)
	if err != nil {
		t.Fatalf("Unable to update consensus: %v", err)
	}

	if channel.ConsensusTurnNum() != 2 {
		t.Fatalf("consensus turn num not updated")
	}
	if !channel.Includes(g2) {
		t.Fatal("consensus incorrectly updated")
	}
	if channel.Includes(g3) {
		t.Fatal("consensus incorrectly updated")
	}

	err = channel.UpdateConsensus(p2Returned)
	if err != nil {
		t.Fatalf("Unable to update consensus: %v", err)
	}

	// The incorrect counter signature is received on the latest proposal
	latest, _ = channel.latestProposedVars()
	wrongCounterSig3, _ := latest.AsState(fp()).Sign(brian.PrivateKey)
	wrongP3Returned := SignedProposal{
		Proposal:  thirdSigned.Proposal,
		Signature: wrongCounterSig3,
	}
	err = channel.UpdateConsensus(wrongP3Returned)
	if !errors.Is(err, ErrWrongSigner) {
		t.Fatalf("ungracefully handled wrong signature: %v", err)
	}

	if channel.ConsensusTurnNum() != 2 {
		t.Fatalf("consensus turn num not updated")
	}

	// The correct counter signature is received on the latest proposal
	latest, _ = channel.latestProposedVars()
	counterSig3, _ := latest.AsState(fp()).Sign(bob.PrivateKey)
	p3Returned := SignedProposal{
		Proposal:  thirdSigned.Proposal,
		Signature: counterSig3,
	}
	_ = channel.UpdateConsensus(p3Returned)

	if channel.ConsensusTurnNum() != 3 {
		t.Fatalf("consensus turn num not updated")
	}
	if !channel.Includes(g2) {
		t.Fatal("consensus incorrectly updated")
	}
	if !channel.Includes(g3) {
		t.Fatal("consensus incorrectly updated")
	}

	// A counter signature is received on an old proposal
	err = channel.UpdateConsensus(p2Returned)
	if err != nil {
		t.Fatalf("Unable to receive old proposal")
	}
	if channel.ConsensusTurnNum() != 3 {
		t.Fatalf("consensus turn num not updated")
	}
	if !channel.Includes(g3) {
		t.Fatal("consensus incorrectly updated")
	}

	// A counter signature is received on an unexpected proposal
	p4Returned := SignedProposal{
		Proposal:  Proposal{toAdd: add(4, 10, targetChannel, alice, bob)},
		Signature: counterSig2,
	}
	err = channel.UpdateConsensus(p4Returned)
	if !errors.Is(err, ErrProposalQueueExhausted) {
		t.Fatalf("did not gracefully handle future proposal: %v", err)
	}
	if channel.ConsensusTurnNum() != 3 {
		t.Fatalf("consensus turn num not updated")
	}
	if !channel.Includes(g3) {
		t.Fatal("consensus incorrectly updated")
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
