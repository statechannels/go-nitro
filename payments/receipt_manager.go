package payments

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

type (
	// paymentStatus stores the status of payments for a given payment channel.
	paymentStatus struct {
		channelSender   common.Address
		startingBalance *big.Int
		largestVoucher  Voucher
	}

	// receiptManager receives vouchers, validates them, and stores the most valuable voucher
	receiptManager struct {
		channels map[types.Destination]*paymentStatus
	}
)

func NewReceiptManager() ReceiptManager {
	channels := make(map[types.Destination]*paymentStatus)
	return &receiptManager{channels}
}

// Register registers a channel for use, given the sender and starting balance of the channel
func (pm receiptManager) Register(channelId types.Destination, sender common.Address, startingBalance *big.Int) error {
	balance := &big.Int{}
	balance.Set(startingBalance)
	voucher := Voucher{channelId: channelId, amount: big.NewInt(0)}
	data := &paymentStatus{sender, balance, voucher}
	if _, ok := pm.channels[channelId]; ok {
		return fmt.Errorf("channel already registered")
	}

	pm.channels[channelId] = data

	return nil
}

// Remove deletes the channel's status
func (pm *receiptManager) Remove(channelId types.Destination) {
	delete(pm.channels, channelId)
}

// Receive validates the incoming voucher, and returns the total amount received so far
func (rm *receiptManager) Receive(voucher Voucher) (*big.Int, error) {
	status, ok := rm.channels[voucher.channelId]
	if !ok {
		return &big.Int{}, fmt.Errorf("channel not registered")
	}

	received := &big.Int{}
	received.Set(voucher.amount)
	if types.Gt(received, status.startingBalance) {
		return &big.Int{}, fmt.Errorf("channel has insufficient funds")
	}

	receivedSoFar := status.largestVoucher.amount
	if !types.Gt(received, receivedSoFar) {
		return receivedSoFar, nil
	}

	signer, err := voucher.recoverSigner()
	if err != nil {
		return &big.Int{}, err
	}
	if signer != status.channelSender {
		return &big.Int{}, fmt.Errorf("wrong signer: %+v, %+v", signer, status.channelSender)
	}

	status.largestVoucher = voucher
	return received, nil
}

// Balance returns the balance of the channel
func (rm *receiptManager) Balance(channelId types.Destination) (Balance, error) {
	data, ok := rm.channels[channelId]
	if !ok {
		return Balance{}, fmt.Errorf("channel not found")
	}

	balance := Balance{&big.Int{}, &big.Int{}}
	balance.Paid.Set(data.largestVoucher.amount)
	balance.Remaining.Sub(data.startingBalance, data.largestVoucher.amount)
	return balance, nil
}
