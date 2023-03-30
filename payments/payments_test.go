package payments

import (
	"fmt"
	"math/big"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/internal/testactors"

	"github.com/statechannels/go-nitro/types"
)

type balance struct {
	remaining, paid *big.Int
}

// manager lets us implement a getBalancer helper to make test assertions a little neater
type manager interface {
	Paid(chanId types.Destination) (*big.Int, error)
	Remaining(chanId types.Destination) (*big.Int, error)
}

// Since the store package already imports the payments package if we tried to use the mem or persist store
// we get import cycles. So we create a simple store that implements the VoucherStore interface for testing.
func newSimpleVoucherStore() VoucherStore {
	return &simpleVoucherStore{vouchers: safesync.Map[*VoucherInfo]{}}
}

type simpleVoucherStore struct {
	vouchers safesync.Map[*VoucherInfo]
}

func (svs *simpleVoucherStore) SetVoucherInfo(channelId types.Destination, v VoucherInfo) error {
	svs.vouchers.Store(channelId.String(), &v)
	return nil
}

func (svs *simpleVoucherStore) GetVoucherInfo(channelId types.Destination) (v *VoucherInfo, ok bool) {
	return svs.vouchers.Load(channelId.String())
}

func (svs *simpleVoucherStore) RemoveVoucherInfo(channelId types.Destination) error {
	svs.vouchers.Delete(channelId.String())
	return nil
}

func TestPaymentManager(t *testing.T) {
	testVoucher := func(cId types.Destination, amount *big.Int, actor testactors.Actor) Voucher {
		payment := &big.Int{}
		payment.Set(amount)
		voucher := Voucher{ChannelId: cId, Amount: payment}
		_ = voucher.Sign(actor.PrivateKey)
		return voucher
	}

	var (
		channelId        = types.Destination{1}
		wrongChannelId   = types.Destination{2}
		anotherChannelId = types.Destination{3}

		deposit       = big.NewInt(1000)
		payment       = big.NewInt(20)
		doublePayment = big.NewInt(40)
		triplePayment = big.NewInt(60)
		overPayment   = big.NewInt(2000)

		startingBalance = balance{big.NewInt(1000), big.NewInt(0)}
		onePaymentMade  = balance{big.NewInt(980), big.NewInt(20)}
		twoPaymentsMade = balance{big.NewInt(960), big.NewInt(40)}
	)

	getBalance := func(m manager) balance {
		paid, _ := m.Paid(channelId)
		remaining, _ := m.Remaining(channelId)
		return balance{remaining, paid}
	}

	// Happy path: Payment manager can register channels and make payments
	paymentMgr := NewVoucherManager(testactors.Alice.Address(), newSimpleVoucherStore())

	_, err := paymentMgr.Pay(channelId, payment, testactors.Alice.PrivateKey)
	Assert(t, err != nil, "channel must be registered to make payments")

	Ok(t, paymentMgr.Register(channelId, testactors.Alice.Address(), testactors.Bob.Address(), deposit))
	Equals(t, startingBalance, getBalance(paymentMgr))

	firstVoucher, err := paymentMgr.Pay(channelId, payment, testactors.Alice.PrivateKey)
	Ok(t, err)
	Equals(t, testVoucher(channelId, payment, testactors.Alice), firstVoucher)
	Equals(t, onePaymentMade, getBalance(paymentMgr))

	signer, err := firstVoucher.RecoverSigner()
	Ok(t, err)
	Equals(t, testactors.Alice.Address(), signer)

	// Happy path: receipt manager can receive vouchers
	receiptMgr := NewVoucherManager(testactors.Bob.Address(), newSimpleVoucherStore())

	_, err = receiptMgr.Receive(firstVoucher)
	Assert(t, err != nil, "channel must be registered to receive vouchers")

	_ = receiptMgr.Register(channelId, testactors.Alice.Address(), testactors.Bob.Address(), deposit)
	Equals(t, startingBalance, getBalance(receiptMgr))

	received, err := receiptMgr.Receive(firstVoucher)
	Ok(t, err)
	Equals(t, received, payment)
	Equals(t, onePaymentMade, getBalance(receiptMgr))
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
	err = paymentMgr.Register(channelId, testactors.Alice.Address(), testactors.Bob.Address(), deposit)
	Assert(t, err != nil, "expected register to fail")
	Equals(t, twoPaymentsMade, getBalance(paymentMgr))

	err = receiptMgr.Register(channelId, testactors.Alice.Address(), testactors.Bob.Address(), deposit)
	Assert(t, err != nil, "expected register to fail")
	Equals(t, twoPaymentsMade, getBalance(receiptMgr))

	// Receiving old vouchers is ok
	received, err = receiptMgr.Receive(firstVoucher)
	Ok(t, err)
	Equals(t, doublePayment, received)
	Equals(t, twoPaymentsMade, getBalance(receiptMgr))

	// Only the payer can sign vouchers
	err = receiptMgr.Register(anotherChannelId, testactors.Bob.Address(), testactors.Alice.Address(), deposit)
	Ok(t, err)
	_, err = paymentMgr.Pay(anotherChannelId, triplePayment, testactors.Bob.PrivateKey)
	Assert(t, err != nil, "only payer can sign vouchers")

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
	voucher.Amount = triplePayment
	_, err = receiptMgr.Receive(voucher)
	Assert(t, err != nil, "expected an error")
	Equals(t, twoPaymentsMade, getBalance(receiptMgr))
}

// TODO: This is a copy of the test helpers from github.com/statechannels/go-nitro/internal/testactors
// We have a copy of them here to avoid an import cycle.

// makeRed sets the colour to red when printed
const makeRed = "\033[31m"

// makeBlack sets the colour to black when printed.
// as it is intended to be used at the end of a string, it also adds two linebreaks
const makeBlack = "\033[39m\n\n"

// Assert fails the test immediately if the condition is false.
// If the assertion fails the formatted message will be output to the console.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(makeRed+"%s:%d: "+msg+makeBlack, append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// Ok fails the test immediately if an err is not nil.
// If the error is not nil the message containing the error will be outputted to the console
func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(makeRed+"%s:%d: unexpected error: %s"+makeBlack, filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// Equals fails the test if want is not deeply equal to got.
// Equals uses reflect.DeepEqual to compare the two values.
func Equals(tb testing.TB, want, got interface{}) {
	if !reflect.DeepEqual(want, got) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(makeRed+"%s:%d:\n\n\texp: %#v\n\n\tgot: %#v"+makeBlack, filepath.Base(file), line, want, got)
		tb.FailNow()
	}
}
