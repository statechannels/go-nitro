package payments

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

type (
	// A PaymentManager can be used to make a payment for a given channel.
	PaymentManager struct {
		channelId      types.Destination
		largestVoucher Voucher
		pk             []byte

		payments chan Voucher // gochan used to send vouchers to Bob
	}

	// A receipt provider enables the consumer to listen to incoming
	// vouchers, validate them, and notify the consumer of a payment
	// received
	ReceiptManager struct {
		channelId      types.Destination
		alice          common.Address
		largestVoucher Voucher

		vouchers chan Voucher // todo: should be unidirectional
		payments chan uint    // todo: API needs to be clarified
	}
)

// implemented just to make use of fp ...
func NewPaymentManager(channelId types.Destination, pk []byte) (PaymentManager, error) {
	pp := PaymentManager{}

	pp.largestVoucher = Voucher{
		channelId: channelId,
	}

	if err := pp.largestVoucher.sign(pk); err != nil {
		panic(err)
	}

	pp.channelId = channelId

	return pp, nil
}

// PayBob will add `amount` to the biggest voucher to create
// a new voucher. This voucher is then signed and sent to Bob.
func (pp *PaymentManager) PayBob(amount uint) error {
	channelId := pp.largestVoucher.channelId
	pp.largestVoucher = Voucher{
		channelId: channelId,
		amount:    amount + pp.largestVoucher.amount,
	}

	if err := pp.largestVoucher.sign(pp.pk); err != nil {
		panic(err)
	}

	pp.payments <- pp.largestVoucher

	return nil
}

func (rp *ReceiptManager) ValidateVouchers(listener chan uint64) {
	for voucher := range rp.vouchers {
		if voucher.channelId != rp.channelId {
			panic("wrong channel!")
		}

		signer, err := voucher.recoverSigner()

		if err != nil {
			// todo: figure out what to do with errors within a goroutine
			panic(err)
		}

		if signer != rp.alice {
			panic("invalid signature")
		}

		if rp.largestVoucher.amount >= voucher.amount {
			panic("received a stale voucher")
		}

		received := voucher.amount - rp.largestVoucher.amount

		// Clarification needed: What information would a receipt manager need to provide
		// to its consumer?
		rp.payments <- received
	}
}
