package consensus_channel

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

func fp() state.FixedPart {
	participants := [2]types.Address{
		alice.Address, bob.Address,
	}
	return state.FixedPart{
		Participants:      participants[:],
		ChainId:           big.NewInt(0),
		ChannelNonce:      big.NewInt(9001),
		ChallengeDuration: big.NewInt(100),
	}
}

func allocation(d actor, a uint64) Balance {
	return Balance{destination: d.Destination(), amount: big.NewInt(int64(a))}
}

func guarantee(amount uint64, target types.Destination, left, right actor) Guarantee {
	return Guarantee{
		target: target,
		amount: big.NewInt(int64(amount)),
		left:   left.Destination(),
		right:  right.Destination(),
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
		allocation(alice, uint64(200)),
		allocation(bob, uint64(300)),
		guarantee(uint64(5), types.Destination{1}, alice, bob),
	)

}

func add(turnNum, amount uint64, vId types.Destination, left, right actor) Add {
	bigAmount := big.NewInt(int64(amount))
	return Add{
		turnNum: turnNum,
		Guarantee: Guarantee{
			amount: bigAmount,
			target: vId,
			left:   left.Destination(),
			right:  right.Destination(),
		},
		LeftDeposit: bigAmount,
	}
}
