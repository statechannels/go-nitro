package payments

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

// VoucherStore is an interface for storing voucher information that the voucher manager expects.
// To avoid import cycles, this interface is defined in the payments package, but implemented in the store package.
type VoucherStore interface {
	SetVoucherInfo(channelId types.Destination, vs VoucherInfo) error
	GetVoucherInfo(channelId types.Destination) (v *VoucherInfo, ok bool)
	RemoveVoucherInfo(channelId types.Destination) error
}

// VoucherInfo stores the status of payments for a given payment channel.
// VoucherManager receives and generates vouchers. It is responsible for storing vouchers.
type VoucherManager struct {
	store VoucherStore
	me    common.Address
}

// NewVoucherManager creates a new voucher manager
func NewVoucherManager(me types.Address, store VoucherStore) *VoucherManager {

	return &VoucherManager{store, me}
}

// Register registers a channel for use, given the payer, payee and starting balance of the channel
func (vm *VoucherManager) Register(channelId types.Destination, payer common.Address, payee common.Address, startingBalance *big.Int) error {

	balance := Balance{big.NewInt(0).Set(startingBalance), &big.Int{}}
	voucher := Voucher{ChannelId: channelId, Amount: big.NewInt(0)}
	data := VoucherInfo{payer, payee, big.NewInt(0).Set(startingBalance), voucher, balance}

	if v, _ := vm.store.GetVoucherInfo(channelId); v != nil {
		return fmt.Errorf("channel already registered")
	}
	return vm.store.SetVoucherInfo(channelId, data)

}

// Remove deletes the channel's status
func (vm *VoucherManager) Remove(channelId types.Destination) {
	err := vm.store.RemoveVoucherInfo(channelId)
	// TODO: Return error instead of panicking
	if err != nil {
		panic(err)
	}
}

// Pay will deduct amount from balance and add it to paid, returning a signed voucher for the
// total amount paid.
func (vm *VoucherManager) Pay(channelId types.Destination, amount *big.Int, pk []byte) (Voucher, error) {
	pStatus, ok := vm.store.GetVoucherInfo(channelId)

	voucher := Voucher{Amount: &big.Int{}}
	if !ok {
		return Voucher{}, fmt.Errorf("channel not found")
	}

	if types.Gt(amount, pStatus.CurrentBalance.Remaining) {
		return Voucher{}, fmt.Errorf("unable to pay amount: insufficient funds")
	}

	if pStatus.ChannelPayer != vm.me {
		return Voucher{}, fmt.Errorf("can only sign vouchers if we're the payer")
	}

	pStatus.CurrentBalance.Remaining.Sub(pStatus.CurrentBalance.Remaining, amount)
	pStatus.CurrentBalance.Paid.Add(pStatus.CurrentBalance.Paid, amount)
	pStatus.LargestVoucher = voucher

	voucher.Amount.Set(pStatus.CurrentBalance.Paid)
	voucher.ChannelId = channelId

	if err := voucher.Sign(pk); err != nil {
		return voucher, err
	}

	err := vm.store.SetVoucherInfo(channelId, *pStatus)
	if err != nil {
		return Voucher{}, err
	}
	return voucher, nil
}

// Receive validates the incoming voucher, and returns the total amount received so far
func (vm *VoucherManager) Receive(voucher Voucher) (*big.Int, error) {
	status, ok := vm.store.GetVoucherInfo(voucher.ChannelId)
	if !ok {
		return &big.Int{}, fmt.Errorf("channel not registered")
	}

	// We only care about vouchers when we are the recipient of the payment
	if status.ChannelPayee != vm.me {
		return &big.Int{}, nil
	}
	received := &big.Int{}
	received.Set(voucher.Amount)
	if types.Gt(received, status.StartingBalance) {
		return &big.Int{}, fmt.Errorf("channel has insufficient funds")
	}

	receivedSoFar := status.LargestVoucher.Amount
	if !types.Gt(received, receivedSoFar) {
		return receivedSoFar, nil
	}

	signer, err := voucher.RecoverSigner()
	if err != nil {
		return &big.Int{}, err
	}
	if signer != status.ChannelPayer {
		return &big.Int{}, fmt.Errorf("wrong signer: %+v, %+v", signer, status.ChannelPayer)
	}
	status.CurrentBalance.Paid.Set(received)
	remaining := big.NewInt(0).Sub(status.StartingBalance, received)
	status.CurrentBalance.Remaining.Set(remaining)

	status.LargestVoucher = voucher

	err = vm.store.SetVoucherInfo(voucher.ChannelId, *status)
	if err != nil {
		return nil, err
	}
	return received, nil
}

// ChannelRegistered returns  whether a channel has been registered with the voucher manager or not
func (vm *VoucherManager) ChannelRegistered(channelId types.Destination) bool {
	_, ok := vm.store.GetVoucherInfo(channelId)
	return ok

}

// Balance returns the balance of the channel
func (vm *VoucherManager) Balance(channelId types.Destination) (Balance, error) {
	data, ok := vm.store.GetVoucherInfo(channelId)
	if !ok {
		return Balance{}, fmt.Errorf("channel not found")
	}

	return data.CurrentBalance, nil

}
