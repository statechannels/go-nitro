package ledger

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

var alice = types.Destination(common.HexToHash("0x0a"))
var bob = types.Destination(common.HexToHash("0x0b"))

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

func add(turnNum, amount uint64, vId, left, right types.Destination) Add {
	return Add{
		turnNum: turnNum,
		Guarantee: Guarantee{
			amount: *big.NewInt(int64(amount)),
			target: vId,
			left:   left,
			right:  right,
		},
	}
}

func TestProposals(t *testing.T) {
	existingChannel := types.Destination{1}
	targetChannel := types.Destination{2}
	aBal := uint64(200)
	bBal := uint64(300)
	vAmount := uint64(5)

	proposal := add(10, vAmount, targetChannel, alice, bob)

	outcome := makeOutcome(
		allocation(alice, aBal),
		allocation(bob, bBal),
		guarantee(vAmount, existingChannel, alice, bob),
	)
	// guarantee(targetChannel, vAmount, alice, bob)) // TODO: this should fail

	before := Vars{TurnNum: 9, Outcome: outcome}

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

	if !after.Outcome.Equal(expected) {
		t.Log(after.Outcome)
		t.Log(expected)
		t.Error("incorrect outcome", err)
	}

	proposal.turnNum += 1
	_, err = after.Add(proposal)

	if err == nil {
		t.Error("expected error when adding duplicate guarantee")
	}
}
