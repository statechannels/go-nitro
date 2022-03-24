package consensus_channel

import (
	"errors"
	"math/big"
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
	testApplyingAddProposalToVars := func(t *testing.T) {
		before := Vars{TurnNum: 9, Outcome: outcome()}

		after, err := before.Add(proposal)

		if err != nil {
			t.Error("unable to compute next state: ", err)
		}

		if after.TurnNum != before.TurnNum+1 {
			t.Error("incorrect state calculation", err)
		}

		expected := makeOutcome(
			allocation(alice, aBal-vAmount),
			allocation(bob, bBal),
			guarantee(vAmount, existingChannel, alice, bob),
			guarantee(vAmount, targetChannel, alice, bob),
		)

		if diff := cmp.Diff(after.Outcome, expected, cmp.AllowUnexported(expected, Balance{}, big.Int{}, Guarantee{})); diff != "" {
			t.Errorf("incorrect outcome: %v", diff)
		}

		largeProposal := proposal
		leftAmount := before.Outcome.left.amount
		largeProposal.amount = *leftAmount.Add(&leftAmount, big.NewInt(1))

		_, err = before.Add(largeProposal)
		if !errors.Is(err, ErrInsufficientFunds) {
			t.Error("expected error when adding too large a guarantee")
		}

		duplicateProposal := proposal
		duplicateProposal.turnNum += 1
		_, err = after.Add(duplicateProposal)

		if !errors.Is(err, ErrDuplicateGuarantee) {
			t.Log(err)
			t.Error("expected error when adding duplicate guarantee")
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
			t.Errorf("channel should check that signer is participant")
		}
	}

	t.Run(`TestConsensusChannelFunctionality`, testConsensusChannelFunctionality)
	t.Run(`TestApplyingAddProposalToVars`, testApplyingAddProposalToVars)
}
