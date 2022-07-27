package payments

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/internal/testactors"
	. "github.com/statechannels/go-nitro/internal/testhelpers"
	"github.com/statechannels/go-nitro/types"
)

// manager lets us implement a getBalancer helper to make test assertions a little neater
type manager interface {
	Balance(chanId types.Destination) (Balance, error)
}

func TestPaymentManager(t *testing.T) {
	testVoucher := func(cId types.Destination, amount *big.Int, actor testactors.Actor) Voucher {
		payment := &big.Int{}
		payment.Set(amount)
		voucher := Voucher{channelId: cId, amount: payment}
		_ = voucher.sign(actor.PrivateKey)
		return voucher
	}

	var (
		channelId      = types.Destination{1}
		wrongChannelId = types.Destination{2}

		deposit       = big.NewInt(1000)
		payment       = big.NewInt(20)
		doublePayment = big.NewInt(40)
		triplePayment = big.NewInt(60)
		overPayment   = big.NewInt(2000)

		startingBalance = Balance{big.NewInt(1000), big.NewInt(0)}
		onePaymentMade  = Balance{big.NewInt(980), big.NewInt(20)}
		twoPaymentsMade = Balance{big.NewInt(960), big.NewInt(40)}
	)

	getBalance := func(m manager) Balance {
		bal, _ := m.Balance(channelId)
		return bal
	}

	// Happy path: Payment manager can register channels and make payments
	paymentMgr := NewPaymentManager(testactors.Alice.Address())

	_, err := paymentMgr.Pay(channelId, payment, testactors.Alice.PrivateKey)
	Assert(t, err != nil, "channel must be registered to make payments")

	Ok(t, paymentMgr.Register(channelId, deposit))
	Equals(t, startingBalance, getBalance(paymentMgr))

	firstVoucher, err := paymentMgr.Pay(channelId, payment, testactors.Alice.PrivateKey)
	Ok(t, err)
	Equals(t, testVoucher(channelId, payment, testactors.Alice), firstVoucher)
	Equals(t, onePaymentMade, getBalance(paymentMgr))

	signer, err := firstVoucher.recoverSigner()
	Ok(t, err)
	Equals(t, testactors.Alice.Address(), signer)

	// Happy path: receipt manager can receive vouchers
	receiptMgr := NewReceiptManager()

	_, err = receiptMgr.Receive(firstVoucher)
	Assert(t, err != nil, "channel must be registered to receive vouchers")

	_ = receiptMgr.Register(channelId, testactors.Alice.Address(), deposit)
	Equals(t, startingBalance, getBalance(receiptMgr))

	received, err := receiptMgr.Receive(firstVoucher)
	Ok(t, err)
	Equals(t, received, payment)

	// Receiving a voucher is idempotent
	received, err = receiptMgr.Receive(firstVoucher)
	Ok(t, err)
	Equals(t, received, payment)
	Equals(t, onePaymentMade, getBalance(receiptMgr))

	// paying twice returns a larger voucher
	secondVoucher, err := paymentMgr.Pay(channelId, payment, testactors.Alice.PrivateKey)
	Ok(t, err)
	Equals(t, testVoucher(channelId, doublePayment, testactors.Alice), secondVoucher)
	Equals(t, twoPaymentsMade, getBalance(paymentMgr))

	// Receiving a new voucher increases amount received
	received, err = receiptMgr.Receive(secondVoucher)
	Ok(t, err)
	Equals(t, doublePayment, received)
	Equals(t, twoPaymentsMade, getBalance(receiptMgr))

	// re-registering a channel doesn't reset its balance
	err = paymentMgr.Register(channelId, deposit)
	Assert(t, err != nil, "expected register to fail")
	Equals(t, twoPaymentsMade, getBalance(paymentMgr))

	err = receiptMgr.Register(channelId, testactors.Alice.Address(), deposit)
	Assert(t, err != nil, "expected register to fail")
	Equals(t, twoPaymentsMade, getBalance(receiptMgr))

	// Receiving old vouchers is ok
	received, err = receiptMgr.Receive(firstVoucher)
	Ok(t, err)
	Equals(t, doublePayment, received)
	Equals(t, twoPaymentsMade, getBalance(receiptMgr))

	// Only the signer can sign vouchers
	_, err = paymentMgr.Pay(channelId, triplePayment, testactors.Bob.PrivateKey)
	Assert(t, err != nil, "only Alice can sign vouchers")

	// Receiving a voucher for an unknown channel fails
	_, err = receiptMgr.Receive(testVoucher(wrongChannelId, payment, testactors.Alice))
	Assert(t, err != nil, "expected an error")
	Equals(t, twoPaymentsMade, getBalance(receiptMgr))

	// Receiving a voucher that's too large fails
	_, err = receiptMgr.Receive(testVoucher(channelId, overPayment, testactors.Alice))
	Assert(t, err != nil, "expected an error")
	Equals(t, twoPaymentsMade, getBalance(receiptMgr))

	// Receiving a voucher with the wrong signature fails
	voucher := testVoucher(channelId, payment, testactors.Alice)
	voucher.amount = triplePayment
	_, err = receiptMgr.Receive(voucher)
	Assert(t, err != nil, "expected an error")
	Equals(t, twoPaymentsMade, getBalance(receiptMgr))
}
