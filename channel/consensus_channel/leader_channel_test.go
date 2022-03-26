package consensus_channel

import (
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
	if channel.ConsensusTurnNum() != 0 {
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
	_, err = channel.Propose(p2, testdata.Actors.Alice.PrivateKey)
	if err != nil {
		t.Fatalf("failed to add another proposal: %v", err)
	}

	latest, _ = channel.latestProposedVars()
	if latest.TurnNum != 2 {
		t.Fatal("incorrect latest proposed vars")
	}
	if channel.ConsensusTurnNum() != 0 {
		t.Fatal("consensus incorrectly updated")
	}

}
