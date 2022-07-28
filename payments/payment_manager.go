package payments

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

type (
	// Balance stores the remaining and paid funds in a channel.
	Balance struct {
		Remaining *big.Int
		Paid      *big.Int
	}

	// PaymentManager can be used to make a payment for a given channel, issuing a new, signed voucher to be sent to the receiver
	PaymentManager struct {
		signer   common.Address
		channels map[types.Destination]*Balance
	}
)

func NewPaymentManager(signer common.Address) *PaymentManager {
	channels := make(map[types.Destination]*Balance)
	return &PaymentManager{signer, channels}
}

// Register registers a channel with a starting balance
func (pm PaymentManager) Register(channelId types.Destination, startingBalance *big.Int) error {
	balance := &Balance{&big.Int{}, &big.Int{}}
	if _, ok := pm.channels[channelId]; ok {
		return fmt.Errorf("channel already registered")
	}

	balance.Remaining.Set(startingBalance)
	pm.channels[channelId] = balance

	return nil
}

// Remove deletes the channel from the manager
func (pm *PaymentManager) Remove(channelId types.Destination) {
	delete(pm.channels, channelId)
}

// Pay will deduct amount from balance and add it to paid, returning a signed voucher for the
// total amount paid.
func (pm *PaymentManager) Pay(channelId types.Destination, amount *big.Int, pk []byte) (Voucher, error) {
	balance, ok := pm.channels[channelId]
	voucher := Voucher{amount: &big.Int{}}
	if !ok {
		return voucher, fmt.Errorf("channel not found")
	}
	if types.Gt(amount, balance.Remaining) {
		return Voucher{}, fmt.Errorf("unable to pay amount: insufficient funds")
	}

	balance.Remaining.Sub(balance.Remaining, amount)
	balance.Paid.Add(balance.Paid, amount)

	voucher.amount.Set(balance.Paid)
	voucher.channelId = channelId

	if err := voucher.sign(pk); err != nil {
		return Voucher{}, err
	}

	// question: is there a more efficient way to validate the signature against the purported signer?
	// (is this validation even necessary? it's more of a failsafe than an important feature)
	signer, err := voucher.recoverSigner()
	if err != nil {
		return Voucher{}, err
	}

	if signer != pm.signer {
		return Voucher{}, fmt.Errorf("only signer may sign vouchers")
	}

	return voucher, nil
}

// Balance returns the balance of the channel
func (pm *PaymentManager) Balance(channelId types.Destination) (Balance, error) {
	stored, ok := pm.channels[channelId]
	if !ok {
		return Balance{}, fmt.Errorf("channel not found")
	}

	balance := Balance{&big.Int{}, &big.Int{}}
	balance.Paid.Set(stored.Paid)
	balance.Remaining.Set(stored.Remaining)

	return balance, nil
}
