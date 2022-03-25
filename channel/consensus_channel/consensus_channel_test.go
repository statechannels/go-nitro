package consensus_channel

import (
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/types"
)

func TestConsensusChannel(t *testing.T) {
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

	proposal := add(10, vAmount, targetChannel, alice, bob)

	outcome := func() LedgerOutcome {
		return makeOutcome(
			allocation(alice, aBal),
			allocation(bob, bBal),
			guarantee(vAmount, existingChannel, alice, bob),
		)

	}

	fingerprint := func(v Vars) string {
		h, err := v.asState(state.TestState.FixedPart()).Hash()

		if err != nil {
			panic(err)
		}
		return h.String()
	}

	vars := Vars{TurnNum: 9, Outcome: outcome()}

	f1 := fingerprint(vars)
	clone1 := vars.clone()

	if fingerprint(clone1) != f1 {
		t.Fatal("vars incorrectly cloned: ", f1, fingerprint(clone1))
	}

	clone1.Outcome.guarantees[targetChannel] = guarantee(999, bob, bob, bob)
	if f1 != fingerprint(vars) {
		t.Fatal("vars shares data with clone")
	}

	clone2 := vars.clone()
	clone2.Outcome.left = allocation(bob, 111)
	if vars.Outcome.left.destination == bob {
		t.Fatal("vars shares data with clone")
	}

	clone3 := vars.clone()
	clone3.Outcome.right = allocation(alice, 111)
	if f1 != fingerprint(vars) {
		t.Fatal("vars shares data with clone")
	}

	testApplyingAddProposalToVars := func(t *testing.T) {
		startingTurnNum := uint64(9)
		vars := Vars{TurnNum: startingTurnNum, Outcome: outcome()}

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
			guarantee(vAmount, existingChannel, alice, bob),
			guarantee(vAmount, targetChannel, alice, bob),
		)

		if diff := cmp.Diff(vars.Outcome, expected, cmp.AllowUnexported(expected, Balance{}, big.Int{}, Guarantee{})); diff != "" {
			t.Fatalf("incorrect outcome: %v", diff)
		}

		// Proposing the same change again should fail
		duplicateProposal := proposal
		duplicateProposal.turnNum += 1
		err = vars.Add(duplicateProposal)

		if !errors.Is(err, ErrDuplicateGuarantee) {
			t.Fatalf("expected error when adding duplicate guarantee: %v", err)
		}

		// Proposing a change that depletes a balance should fail
		vars = Vars{TurnNum: startingTurnNum, Outcome: outcome()}
		largeProposal := proposal
		leftAmount := vars.Outcome.left.amount
		largeProposal.amount = *leftAmount.Add(&leftAmount, big.NewInt(1))
		err = vars.Add(largeProposal)
		if !errors.Is(err, ErrInsufficientFunds) {
			t.Fatalf("expected error when adding too large a guarantee: %v", err)
		}
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

	testConsensusChannelFunctionality := func(t *testing.T) {
		channel, err := NewConsensusChannel(fp(), leader, outcome(), sigs)

		if err != nil {
			t.Fatalf("unable to construct a new consensus channel: %v", err)
		}

		_, err = channel.sign(initialVars, testdata.Actors.Bob.PrivateKey)
		if err == nil {
			t.Fatalf("channel should check that signer is participant")
		}

		f := fingerprint(channel.current.Vars)

		latest, err := channel.latestProposedVars()
		if err != nil {
			t.Fatalf("latest proposed vars returned err: %v", err)
		}

		latest.Outcome.guarantees[targetChannel] = guarantee(10, targetChannel, alice, bob)
		if f != fingerprint(channel.current.Vars) {
			t.Fatalf("latestProposedVars did not return a copy")
		}

		briansSig, _ := initialVars.asState(fp()).Sign(testdata.Actors.Brian.PrivateKey)
		wrongSigs := [2]state.Signature{sigs[1], briansSig}
		_, err = NewConsensusChannel(fp(), leader, outcome(), wrongSigs)
		if err == nil {
			t.Fatalf("channel should check that signers are participants")
		}
	}

	testAsLeader := func(t *testing.T) {
		channel, err := NewConsensusChannel(fp(), leader, outcome(), sigs)
		if err != nil {
			t.Fatal("unable to construct channel")
		}

		amountAdded := uint64(10)

		latest, _ := channel.latestProposedVars()
		if initialVars.TurnNum != 0 {
			t.Fatal("initialized with non-zero turn number")
		}

		p := add(1, amountAdded, targetChannel, alice, bob)
		sp, err := channel.Propose(p, testdata.Actors.Alice.PrivateKey)
		if err != nil {
			t.Fatalf("failed to add proposal: %v", err)
		}

		latest, _ = channel.latestProposedVars()
		if latest.TurnNum != 1 {
			t.Fatal("incorrect latest proposed vars")
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
		_, err = channel.Propose(p2, testdata.Actors.Alice.PrivateKey)
		if err != nil {
			t.Fatalf("failed to add another proposal: %v", err)
		}

		latest, _ = channel.latestProposedVars()
		if latest.TurnNum != 2 {
			t.Fatal("incorrect latest proposed vars")
		}
	}

	t.Run(`TestApplyingAddProposalToVars`, testApplyingAddProposalToVars)
	t.Run(`TestConsensusChannelFunctionality`, testConsensusChannelFunctionality)
	t.Run(`TestAsLeader`, testAsLeader)
}
