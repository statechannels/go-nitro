package consensus_channel

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

var targetChannel = types.Destination{2}

type actor struct {
	Address    types.Address
	PrivateKey []byte
}

func (a actor) Destination() types.Destination {
	return types.AddressToDestination(a.Address)
}

var alice = actor{
	common.HexToAddress(`0xAAA6628Ec44A8a742987EF3A114dDFE2D4F7aDCE`),
	common.Hex2Bytes(`2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d`),
}
var bob = actor{
	common.HexToAddress(`0xBBB676f9cFF8D242e9eaC39D063848807d3D1D94`),
	common.Hex2Bytes(`0279651921cd800ac560c21ceea27aab0107b67daf436cdd25ce84cad30159b4`),
}
var brian = actor{
	common.HexToAddress("0xB2B22ec3889d11f2ddb1A1Db11e80D20EF367c01"),
	common.Hex2Bytes("0aca28ba64679f63d71e671ab4dbb32aaa212d4789988e6ca47da47601c18fe2"),
}

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
