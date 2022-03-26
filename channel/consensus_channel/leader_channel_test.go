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

func TestLeaderChannel(t *testing.T) {
	var alice = types.Destination(common.HexToHash("0x0a"))
	var bob = types.Destination(common.HexToHash("0x0b"))

	allocation := func(d types.Destination, a uint64) Balance {
		return Balance{destination: d, amount: *big.NewInt(int64(a))}
	}

	guarantee := func(amount uint64, target, left, right types.Destination) Guarantee {
		return Guarantee{
			target: target,
			amount: *big.NewInt(int64(amount)),
			left:   left,
			right:  right,
		}
	}

	makeOutcome := func(left, right Balance, guarantees ...Guarantee) LedgerOutcome {
		mappedGuarantees := make(map[types.Destination]Guarantee)
		for _, g := range guarantees {
			mappedGuarantees[g.target] = g
		}
		return LedgerOutcome{left: left, right: right, guarantees: mappedGuarantees}
	}

	add := func(turnNum, amount uint64, vId, left, right types.Destination) Add {
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

	existingChannel := types.Destination{1}
	targetChannel := types.Destination{2}
	aBal := uint64(200)
	bBal := uint64(300)
	vAmount := uint64(5)

	outcome := func() LedgerOutcome {
		return makeOutcome(
			allocation(alice, aBal),
			allocation(bob, bBal),
			guarantee(vAmount, existingChannel, alice, bob),
		)

	}

	fp := func() state.FixedPart {
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

	initialVars := Vars{Outcome: outcome(), TurnNum: 0}
	aliceSig, _ := initialVars.asState(fp()).Sign(testdata.Actors.Alice.PrivateKey)
	bobsSig, _ := initialVars.asState(fp()).Sign(testdata.Actors.Bob.PrivateKey)
	sigs := [2]state.Signature{aliceSig, bobsSig}

	channel, err := NewLeaderChannel(fp(), outcome(), sigs)
	if err != nil {
		t.Fatal("unable to construct channel")
	}

	amountAdded := uint64(10)

	// latest, _ := channel.latestProposedVars()
	if initialVars.TurnNum != 0 {
		t.Fatal("initialized with non-zero turn number")
	}

	p := add(1, amountAdded, targetChannel, alice, bob)
	sp, err := channel.Propose(p, testdata.Actors.Alice.PrivateKey)
	if err != nil {
		t.Fatalf("failed to add proposal: %v", err)
	}

	success, _ := channel.IsProposed(p.Guarantee)
	if !success {
		t.Fatal("incorrect latest proposed vars")
	}
	if channel.ConsensusTurnNum() != 0 || channel.Includes(p.Guarantee) {
		t.Fatal("consensus incorrectly updated")
	}

	outcomeSigned := makeOutcome(
		allocation(alice, aBal-amountAdded),
		allocation(bob, bBal),
		guarantee(vAmount, existingChannel, alice, bob),
		guarantee(amountAdded, targetChannel, alice, bob),
	)
	stateSigned := Vars{TurnNum: 1, Outcome: outcomeSigned}
	sig, _ := stateSigned.asState(fp()).Sign(testdata.Actors.Alice.PrivateKey)
	expected := SignedProposal{Proposal: p, Signature: sig}

	if !reflect.DeepEqual(sp, expected) {
		t.Fatalf("propose failed")
	}

	thirdChannel := types.Destination{3}
	p2 := p
	p2.target = thirdChannel
	g2 := p2.Guarantee
	secondSigned, err := channel.Propose(p2, testdata.Actors.Alice.PrivateKey)
	if err != nil {
		t.Fatalf("failed to add another proposal: %v", err)
	}
	if secondSigned.Proposal.(Add).turnNum != 2 {
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
	counterSig2, _ := latest.asState(fp()).Sign(testdata.Actors.Bob.PrivateKey)

	p3 := p
	p3.target = types.Destination{4}
	g3 := p3.Guarantee
	thirdSigned, _ := channel.Propose(p3, testdata.Actors.Alice.PrivateKey)

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
	wrongCounterSig3, _ := latest.asState(fp()).Sign(testdata.Actors.Brian.PrivateKey)
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
	counterSig3, _ := latest.asState(fp()).Sign(testdata.Actors.Bob.PrivateKey)
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
		Proposal:  add(4, 10, targetChannel, alice, bob),
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
