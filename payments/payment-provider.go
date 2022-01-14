package payments

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
)

const alice = 0

type (
	// A PaymentProvider can be used to make a payment for a given channel.
	PaymentProvider struct {
		fp             state.FixedPart
		largestVoucher Voucher
		pk             []byte

		payments chan Voucher // gochan used to send vouchers to Bob
	}

	// A receipt provider enables the consumer to listen to incoming
	// vouchers, validate them, and notify the consumer of a payment
	// received
	ReceiptProvider struct {
		fp             state.FixedPart
		largestVoucher Voucher

		vouchers chan Voucher // todo: should be unidirectional
		payments chan uint    // todo: API needs to be clarified
	}
)

// implemented just to make use of fp ...
func NewPaymentProvider(fp state.FixedPart, pk []byte) (PaymentProvider, error) {
	chanId, err := fp.ChannelId()
	pp := PaymentProvider{}

	if err != nil {
		return pp, err
	}

	pp.largestVoucher = Voucher{
		channelId: chanId,
	}

	pp.largestVoucher.sign(pk)

	pp.fp = fp

	return pp, nil
}

// PayBob will use add amount to the biggest voucher to create
// a new voucher. This voucher is then signed and sent to Bob.
func (pp *PaymentProvider) PayBob(amount uint64) error {
	channelId := pp.largestVoucher.channelId
	pp.largestVoucher = Voucher{
		channelId: channelId,
		amount:    uint(amount),
	}
	pp.largestVoucher.sign(pp.pk)

	pp.payments <- pp.largestVoucher

	return nil
}

func (rp *ReceiptProvider) ValidateVouchers(listener chan uint64) {
	for voucher := range rp.vouchers {

		signer, err := voucher.recoverSigner()

		if err != nil {
			// todo: figure out what to do with errors within a goroutine
			panic(err)
		}

		if signer != rp.alice() {
			panic("invalid signature")
		}

		if rp.largestVoucher.amount >= voucher.amount {
			panic("received a stale voucher")
		}

		received := voucher.amount - rp.largestVoucher.amount

		// Clarification needed: What information would a receipt provider need to provide
		// to its consumer?
		rp.payments <- received
	}
}

func (rp ReceiptProvider) alice() common.Address {
	return rp.fp.Participants[alice]
}
