package consensus_channel

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/types"
)

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
	return Balance{destination: d, amount: big.NewInt(int64(a))}
}

func guarantee(amount uint64, target, left, right types.Destination) Guarantee {
	return Guarantee{
		target: target,
		amount: big.NewInt(int64(amount)),
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
		allocation(testdata.Actors.Alice.Destination(), uint64(200)),
		allocation(testdata.Actors.Bob.Destination(), uint64(300)),
		guarantee(uint64(5), types.Destination{1}, testdata.Actors.Alice.Destination(), testdata.Actors.Bob.Destination()),
	)

}

func add(turnNum, amount uint64, vId, left, right types.Destination) Add {
	bigAmount := big.NewInt(int64(amount))
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
