package payments

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

type (
	// A Voucher signed by Alice can be used by Bob to redeem payments in case of
	// a misbehaving Alice.
	//
	// During normal operation, Alice & Bob would terminate the channel with an
	// outcome reflecting the largest amount signed by Alice. For instance,
	// - if the channel started with balances {alice: 100, bob: 0}
	// - and the biggest voucher signed by alice had amount = 20
	// - then Alice and Bob would cooperatively conclude the channel with outcome
	//   {alice: 80, bob: 20}
	Voucher struct {
		ChannelId types.Destination
		Amount    *big.Int
		Signature state.Signature
	}

	// Balance stores the remaining and paid funds in a channel.
	Balance struct {
		Remaining *big.Int
		Paid      *big.Int
	}

	VoucherManager interface {
		// Register registers a channel with a starting balance
		Register(channelId types.Destination, sender common.Address, receiver common.Address, startingBalance *big.Int) error

		// Remove deletes the channel from the manager
		Remove(channelId types.Destination)

		// Pay will deduct amount from balance and add it to paid, returning a signed voucher for the
		// total amount paid.
		Pay(channelId types.Destination, amount *big.Int, pk []byte) (Voucher, error)

		// Receive validates the incoming voucher, and returns the total amount received so far
		Receive(voucher Voucher) (amountReceived *big.Int, err error)

		// Balance returns the balance of the channel
		Balance(channelId types.Destination) (Balance, error)
	}
)
