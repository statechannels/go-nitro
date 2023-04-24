package consensus_channel

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

func TestReceive(t *testing.T) {
	vAmount := uint64(5)
	initialVars := Vars{Outcome: ledgerOutcome(), TurnNum: 0}
	aliceSig, _ := initialVars.AsState(fp()).Sign(alice.PrivateKey)
	bobsSig, _ := initialVars.AsState(fp()).Sign(bob.PrivateKey)
	sigs := [2]state.Signature{aliceSig, bobsSig}

	channel, err := newFollowerChannel(initialVars.AsState(fp()), 0, ledgerOutcome(), sigs)
	if err != nil {
		t.Fatal("unable to construct channel")
	}

	proposal := Proposal{LedgerID: channel.Id, ToAdd: add(vAmount, targetChannel, alice, bob)}

	// Create a proposal with an incorrect signature
	badSigProposal := SignedProposal{bobsSig, proposal, 1}
	err = channel.Receive(badSigProposal)
	if !errors.Is(ErrInvalidProposalSignature, err) {
		t.Fatalf("expected %v, but got %v", ErrInvalidProposalSignature, err)
	}

	valid := createSignedProposal(initialVars, proposal, fp(), alice.PrivateKey)

	err = channel.Receive(valid)
	if err != nil {
		t.Fatalf("unable to receive proposal: %v", err)
	}
	// Check that the proposal was queued up properly
	if len(channel.proposalQueue) != 1 {
		t.Fatalf("Expected only one proposal in queue")
	}
	queued := channel.proposalQueue[0]
	if !reflect.DeepEqual(queued.Proposal, proposal) {
		t.Fatalf("Expected proposal to be queued")
	}

	// Generate a second proposal(removal)
	latestProposed, _ := channel.latestProposedVars()
	secondProposal := Proposal{LedgerID: channel.Id, ToRemove: remove(channel1Id, 2)}
	anotherValid := createSignedProposal(latestProposed, secondProposal, fp(), alice.PrivateKey)
	err = channel.Receive(anotherValid)
	if err != nil {
		t.Fatalf("unable to receive proposal: %v", err)
	}

	if len(channel.proposalQueue) != 2 {
		t.Fatalf("Expected both proposals in the queue")
	}
	queued = channel.proposalQueue[1]
	if !reflect.DeepEqual(queued.Proposal, secondProposal) {
		t.Fatalf("Expect the latest proposal to be the last in the queue")
	}

	// Check that receive rejects a stale proposal
	stale := createSignedProposal(Vars{TurnNum: 0, Outcome: ledgerOutcome()}, proposal, fp(), alice.PrivateKey)
	err = channel.Receive(stale)
	if !errors.Is(ErrInvalidTurnNum, err) {
		t.Fatalf("expected %v, but got %v", ErrInvalidTurnNum, err)
	}

	// Check that  receive rejects a proposal too far in the future
	tooFar := createSignedProposal(Vars{TurnNum: 10, Outcome: ledgerOutcome()}, proposal, fp(), alice.PrivateKey)
	err = channel.Receive(tooFar)
	if !errors.Is(ErrInvalidTurnNum, err) {
		t.Fatalf("expected %v, but got %v", ErrInvalidTurnNum, err)
	}
}

func TestFollowerChannel(t *testing.T) {
	initialVars := Vars{Outcome: ledgerOutcome(), TurnNum: 0}
	aliceSig, _ := initialVars.AsState(fp()).Sign(alice.PrivateKey)
	bobsSig, _ := initialVars.AsState(fp()).Sign(bob.PrivateKey)
	sigs := [2]state.Signature{aliceSig, bobsSig}

	channel, err := newFollowerChannel(initialVars.AsState(fp()), 0, ledgerOutcome(), sigs)
	if err != nil {
		t.Fatal("unable to construct channel")
	}

	amountAdded := uint64(5)
	proposal := Proposal{LedgerID: channel.Id, ToAdd: add(amountAdded, targetChannel, alice, bob)}

	_, err = channel.SignNextProposal(proposal, bob.PrivateKey)
	if !errors.Is(ErrNoProposals, err) {
		t.Fatalf("expected %v, but got %v", ErrNoProposals, err)
	}

	signedProposal := SignedProposal{
		Proposal: proposal,
		// Note that this signature is never checked in SignNextProposal
		Signature: state.Signature{},
	}
	channel.proposalQueue = []SignedProposal{signedProposal}
	proposal2 := Proposal{LedgerID: channel.Id, ToAdd: add(amountAdded+1, targetChannel, alice, bob)}

	_, err = channel.SignNextProposal(proposal2, bob.PrivateKey)
	if !errors.Is(ErrNonMatchingProposals, err) {
		t.Fatalf("expected %v, but got %v", ErrNonMatchingProposals, err)
	}

	withMySig, err := channel.SignNextProposal(proposal, bob.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	outcomeSigned := makeOutcome(
		allocation(alice, aBal-amountAdded),
		allocation(bob, bBal),
		guarantee(vAmount, channel1Id, alice, bob),
		guarantee(amountAdded, targetChannel, alice, bob),
	)
	varsSigned := Vars{Outcome: outcomeSigned, TurnNum: 1}
	correctSig, _ := varsSigned.AsState(fp()).Sign(bob.PrivateKey)
	equals(t, withMySig.Signature, correctSig)

	if channel.ConsensusTurnNum() != 1 {
		t.Fatalf("incorrect turn number: expected 1, got %d", channel.ConsensusTurnNum())
	}
	if !channel.Includes(proposal.ToAdd.Guarantee) {
		t.Fatal("expected the channel to not include the guarantee")
	}
	if len(channel.proposalQueue) != 0 {
		t.Fatal("expected the proposal queue to be empty")
	}
}

func TestRestrictedLeaderMethods(t *testing.T) {
	initialVars := Vars{Outcome: ledgerOutcome(), TurnNum: 0}
	aliceSig, _ := initialVars.AsState(fp()).Sign(alice.PrivateKey)
	bobsSig, _ := initialVars.AsState(fp()).Sign(bob.PrivateKey)
	sigs := [2]state.Signature{aliceSig, bobsSig}

	channel, _ := newFollowerChannel(initialVars.AsState(fp()), 0, ledgerOutcome(), sigs)

	if _, err := channel.Propose(Proposal{ToAdd: Add{}}, alice.PrivateKey); err != ErrNotLeader {
		t.Errorf("Expected error when calling Propose() as a follower, but found none")
	}

	if err := channel.leaderReceive(SignedProposal{}); err != ErrNotLeader {
		t.Errorf("Expected error when calling leaderReceive() as a follower, but found none")
	}
}

func TestFollowerIncorrectlyAddressedProposals(t *testing.T) {
	initialVars := Vars{Outcome: ledgerOutcome(), TurnNum: 0}
	aliceSig, _ := initialVars.AsState(fp()).Sign(alice.PrivateKey)
	bobsSig, _ := initialVars.AsState(fp()).Sign(bob.PrivateKey)
	sigs := [2]state.Signature{aliceSig, bobsSig}

	leaderCh, _ := newLeaderChannel(initialVars.AsState(fp()), 0, ledgerOutcome(), sigs)
	followerCh, _ := newFollowerChannel(initialVars.AsState(fp()), 0, ledgerOutcome(), sigs)

	someProposal, _ := leaderCh.Propose(Proposal{ToAdd: add(1, types.Destination{}, alice, bob)}, alice.PrivateKey)
	someProposal.Proposal.LedgerID = types.Destination{} // alter the ChannelID so that it doesn't match

	err := followerCh.Receive(someProposal)

	if err != ErrIncorrectChannelID {
		t.Fatalf("expected error receiving proposal with incorrect ChannelID, but found none")
	}

	_, err = followerCh.SignNextProposal(someProposal.Proposal, bob.PrivateKey)

	if err != ErrIncorrectChannelID {
		t.Fatalf("expected error receiving proposal with incorrect ChannelID, but found none")
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
