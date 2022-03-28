package consensus_channel

import (
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/types"
)

var alice = types.Destination(common.HexToHash("0x0a"))
var bob = types.Destination(common.HexToHash("0x0b"))
var aBal = uint64(200)
var bBal = uint64(300)
var vAmount = uint64(5)
var existingChannel = types.Destination{1}
var targetChannel = types.Destination{2}

// TODO these helpers and the helpers in leader_channel should be shared.
func fp() state.FixedPart {
	participants := [2]types.Address{
		testdata.Actors.Alice.Address, testdata.Actors.Bob.Address,
	}
	return state.FixedPart{
		Participants:      participants[:],
		ChainId:           big.NewInt(0),
		ChannelNonce:      big.NewInt(9001),
		ChallengeDuration: big.NewInt(100),
	}
}

func allocation(d types.Destination, a uint64) Balance {
	return Balance{destination: d, amount: *big.NewInt(int64(a))}
}

func guarantee(amount uint64, target, left, right types.Destination) Guarantee {
	return Guarantee{
		target: target,
		amount: *big.NewInt(int64(amount)),
		left:   left,
		right:  right,
	}
}

func makeOutcome(left, right Balance, guarantees ...Guarantee) LedgerOutcome {
	mappedGuarantees := make(map[types.Destination]Guarantee)
	for _, g := range guarantees {
		mappedGuarantees[g.target] = g
	}
	return LedgerOutcome{left: left, right: right, guarantees: mappedGuarantees}
}

func ledgerOutcome() LedgerOutcome {
	return makeOutcome(
		allocation(alice, aBal),
		allocation(bob, bBal),
		guarantee(vAmount, existingChannel, alice, bob),
	)

}

func add(turnNum, amount uint64, vId, left, right types.Destination) Add {
	bigAmount := *big.NewInt(int64(amount))
	return Add{
		turnNum: turnNum,
		Guarantee: Guarantee{
			amount: bigAmount,
			target: vId,
			left:   left,
			right:  right,
		},
		LeftDeposit: bigAmount,
	}
}

// createSignedProposal generates a signed proposal given the vars, proposal fixed parts and private key
// The vars passed in are NOT mutated!
func createSignedProposal(vars Vars, proposal Add, fp state.FixedPart, pk []byte) SignedProposal {
	proposalVars := Vars{TurnNum: vars.TurnNum, Outcome: vars.Outcome.clone()}
	_ = proposalVars.Add(proposal)

	state := proposalVars.asState(fp)
	sig, _ := state.Sign(pk)

	signedProposal := SignedProposal{
		Proposal:  proposal,
		Signature: sig,
	}

	return signedProposal

}

func TestReceive(t *testing.T) {
	initialVars := Vars{Outcome: ledgerOutcome(), TurnNum: 0}
	aliceSig, _ := initialVars.asState(fp()).Sign(testdata.Actors.Alice.PrivateKey)
	bobsSig, _ := initialVars.asState(fp()).Sign(testdata.Actors.Bob.PrivateKey)
	sigs := [2]state.Signature{aliceSig, bobsSig}

	channel, err := NewFollowerChannel(fp(), ledgerOutcome(), sigs)
	if err != nil {
		t.Fatal("unable to construct channel")
	}

	proposal := add(1, vAmount, targetChannel, alice, bob)

	// Create a proposal with an incorrect signature
	badSigProposal := SignedProposal{bobsSig, proposal}
	err = channel.Receive(badSigProposal)
	if !errors.Is(ErrInvalidProposalSignature, err) {
		t.Fatalf("expected %v, but got %v", ErrInvalidProposalSignature, err)
	}

	valid := createSignedProposal(initialVars, proposal, fp(), testdata.Actors.Alice.PrivateKey)

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

	// Generate a second proposal
	latestProposed, _ := channel.latestProposedVars()
	secondProposal := add(2, vAmount, types.Destination{3}, alice, bob)
	anotherValid := createSignedProposal(latestProposed, secondProposal, fp(), testdata.Actors.Alice.PrivateKey)
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
	stale := createSignedProposal(Vars{TurnNum: 0, Outcome: ledgerOutcome()}, proposal, fp(), testdata.Actors.Alice.PrivateKey)
	err = channel.Receive(stale)
	if !errors.Is(ErrInvalidTurnNum, err) {
		t.Fatalf("expected %v, but got %v", ErrInvalidTurnNum, err)
	}

	// Check that  receive rejects a proposal too far in the future
	tooFar := createSignedProposal(Vars{TurnNum: 10, Outcome: ledgerOutcome()}, proposal, fp(), testdata.Actors.Alice.PrivateKey)
	err = channel.Receive(tooFar)
	if !errors.Is(ErrInvalidTurnNum, err) {
		t.Fatalf("expected %v, but got %v", ErrInvalidTurnNum, err)
	}

}
func TestFollowerChannel(t *testing.T) {
	initialVars := Vars{Outcome: ledgerOutcome(), TurnNum: 0}
	aliceSig, _ := initialVars.asState(fp()).Sign(testdata.Actors.Alice.PrivateKey)
	bobsSig, _ := initialVars.asState(fp()).Sign(testdata.Actors.Bob.PrivateKey)
	sigs := [2]state.Signature{aliceSig, bobsSig}

	channel, err := NewFollowerChannel(fp(), ledgerOutcome(), sigs)
	if err != nil {
		t.Fatal("unable to construct channel")
	}

	proposal := add(1, vAmount, targetChannel, alice, bob)

	err = channel.SignNextProposal(proposal, testdata.Actors.Bob.PrivateKey)
	if !errors.Is(ErrNoProposals, err) {
		t.Fatalf("expected %v, but got %v", ErrNoProposals, err)
	}

	signedProposal := SignedProposal{
		Proposal: proposal,
		// Note that this signature is never checked in SignNextProposal
		Signature: state.Signature{},
	}
	channel.proposalQueue = []SignedProposal{signedProposal}
	proposal2 := add(1, uint64(6), targetChannel, alice, bob)

	err = channel.SignNextProposal(proposal2, testdata.Actors.Bob.PrivateKey)
	if !errors.Is(ErrNonMatchingProposals, err) {
		t.Fatalf("expected %v, but got %v", ErrNonMatchingProposals, err)
	}

	err = channel.SignNextProposal(proposal, testdata.Actors.Bob.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	if channel.ConsensusTurnNum() != 1 {
		t.Fatalf("incorrect turn number: expected 1, got %d", channel.ConsensusTurnNum())
	}
	if !channel.Includes(proposal.Guarantee) {
		t.Fatal("expected the channel to not include the guarantee")
	}
	if len(channel.proposalQueue) != 0 {
		t.Fatal("expected the proposal queue to be empty")
	}
}
