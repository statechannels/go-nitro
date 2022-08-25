package payments

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

// paymentStatus stores the status of payments for a given payment channel.
type paymentStatus struct {
	channelPayer    common.Address
	channelPayee    common.Address
	startingBalance *big.Int
	largestVoucher  Voucher
	currentBalance  Balance
}

// VoucherManager receives and generates vouchers. It is responsible for storing vouchers.
type VoucherManager struct {
	channels map[types.Destination]*paymentStatus
	me       common.Address
}

// NewVoucherManager creates a new voucher manager
func NewVoucherManager(me types.Address) *VoucherManager {
	channels := make(map[types.Destination]*paymentStatus)
	return &VoucherManager{channels, me}
}

// Register registers a channel for use, given the payer, payee and starting balance of the channel
func (vm VoucherManager) Register(channelId types.Destination, payer common.Address, payee common.Address, startingBalance *big.Int) error {

	balance := Balance{big.NewInt(0).Set(startingBalance), &big.Int{}}
	voucher := Voucher{ChannelId: channelId, Amount: big.NewInt(0)}
	data := &paymentStatus{payer, payee, big.NewInt(0).Set(startingBalance), voucher, balance}
	if _, ok := vm.channels[channelId]; ok {
		return fmt.Errorf("channel already registered")
	}

	vm.channels[channelId] = data

	return nil
}

// Remove deletes the channel's status
func (vm *VoucherManager) Remove(channelId types.Destination) {
	delete(vm.channels, channelId)
}

// Pay will deduct amount from balance and add it to paid, returning a signed voucher for the
// total amount paid.
func (vm *VoucherManager) Pay(channelId types.Destination, amount *big.Int, pk []byte) (Voucher, error) {
	pStatus, ok := vm.channels[channelId]

	voucher := Voucher{Amount: &big.Int{}}
	if !ok {
		return Voucher{}, fmt.Errorf("channel not found")
	}

	if types.Gt(amount, pStatus.currentBalance.Remaining) {
		return Voucher{}, fmt.Errorf("unable to pay amount: insufficient funds")
	}

	if pStatus.channelPayer != vm.me {
		return Voucher{}, fmt.Errorf("can only sign vouchers if we're the payer")
	}

	pStatus.currentBalance.Remaining.Sub(pStatus.currentBalance.Remaining, amount)
	pStatus.currentBalance.Paid.Add(pStatus.currentBalance.Paid, amount)
	pStatus.largestVoucher = voucher

	voucher.Amount.Set(pStatus.currentBalance.Paid)
	voucher.ChannelId = channelId

	if err := voucher.Sign(pk); err != nil {
		return voucher, err
	}

	return voucher, nil
}

// Receive validates the incoming voucher, and returns the total amount received so far
func (vm *VoucherManager) Receive(voucher Voucher) (*big.Int, error) {
	status, ok := vm.channels[voucher.ChannelId]
	if !ok {
		return &big.Int{}, fmt.Errorf("channel not registered")
	}

	// We only care about vouchers when we are the recipient of the payment
	if status.channelPayee != vm.me {
		return &big.Int{}, nil
	}
	received := &big.Int{}
	received.Set(voucher.Amount)
	if types.Gt(received, status.startingBalance) {
		return &big.Int{}, fmt.Errorf("channel has insufficient funds")
	}

	receivedSoFar := status.largestVoucher.Amount
	if !types.Gt(received, receivedSoFar) {
		return receivedSoFar, nil
	}

	signer, err := voucher.RecoverSigner()
	if err != nil {
		return &big.Int{}, err
	}
	if signer != status.channelPayer {
		return &big.Int{}, fmt.Errorf("wrong signer: %+v, %+v", signer, status.channelPayer)
	}
	status.currentBalance.Paid.Set(received)
	remaining := big.NewInt(0).Sub(status.startingBalance, received)
	status.currentBalance.Remaining.Set(remaining)

	status.largestVoucher = voucher
	return received, nil
}

// Balance returns the balance of the channel
func (vm *VoucherManager) Balance(channelId types.Destination) (Balance, error) {
	data, ok := vm.channels[channelId]
	if !ok {
		return Balance{}, fmt.Errorf("channel not found")
	}

	return data.currentBalance, nil

}

// Voucher returns the latest sent voucher for a channel
func (vm *VoucherManager) Voucher(channelId types.Destination, pk []byte) (Voucher, error) {

	bal, err := vm.Balance(channelId)
	voucher := Voucher{Amount: &big.Int{}}
	if err != nil {
		return voucher, fmt.Errorf("unable to get balance to construct voucher: %w", err)
	}
	voucher.Amount.Set(bal.Paid)
	voucher.ChannelId = channelId

	if err := voucher.Sign(pk); err != nil {
		return Voucher{}, err
	}

	return voucher, nil

}
