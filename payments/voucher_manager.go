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
	SetVoucherInfo(channelId types.Destination, v VoucherInfo) error
	GetVoucherInfo(channelId types.Destination) (v *VoucherInfo, err error)
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
	voucher := Voucher{ChannelId: channelId, Amount: big.NewInt(0)}
	data := VoucherInfo{payer, payee, big.NewInt(0).Set(startingBalance), voucher}

	if v, _ := vm.store.GetVoucherInfo(channelId); v != nil {
		return fmt.Errorf("channel already registered")
	}
	return vm.store.SetVoucherInfo(channelId, data)
}

// Remove deletes the channel's status
func (vm *VoucherManager) Remove(channelId types.Destination) error {
	err := vm.store.RemoveVoucherInfo(channelId)
	if err != nil {
		return err
	}
	return nil
}

// Pay will deduct amount from balance and add it to paid, returning a signed voucher for the
// total amount paid.
func (vm *VoucherManager) Pay(channelId types.Destination, amount *big.Int, pk []byte) (Voucher, error) {
	vInfo, err := vm.store.GetVoucherInfo(channelId)
	if err != nil {
		return Voucher{}, fmt.Errorf("channel not registered: %w", err)
	}

	if types.Gt(amount, vInfo.Remaining()) {
		return Voucher{}, fmt.Errorf("unable to pay amount: insufficient funds")
	}

	if vInfo.ChannelPayer != vm.me {
		return Voucher{}, fmt.Errorf("can only sign vouchers if we're the payer")
	}
	newAmount := big.NewInt(0).Add(vInfo.LargestVoucher.Amount, amount)
	voucher := Voucher{Amount: big.NewInt(0).Set(newAmount), ChannelId: channelId}

	vInfo.LargestVoucher = voucher

	if err := voucher.Sign(pk); err != nil {
		return voucher, err
	}

	err = vm.store.SetVoucherInfo(channelId, *vInfo)
	if err != nil {
		return Voucher{}, err
	}
	return voucher, nil
}

// Receive validates the incoming voucher, and returns the total amount received so far as well as the amount received from the voucher
func (vm *VoucherManager) Receive(voucher Voucher) (total *big.Int, delta *big.Int, err error) {
	vInfo, err := vm.store.GetVoucherInfo(voucher.ChannelId)
	if err != nil {
		return &big.Int{}, &big.Int{}, fmt.Errorf("channel not registered: %w", err)
	}

	// We only care about vouchers when we are the recipient of the payment
	if vInfo.ChannelPayee != vm.me {
		return &big.Int{}, &big.Int{}, nil
	}

	if types.Gt(voucher.Amount, vInfo.StartingBalance) {
		return &big.Int{}, &big.Int{}, fmt.Errorf("channel has insufficient funds")
	}

	total = vInfo.LargestVoucher.Amount
	if !types.Gt(voucher.Amount, total) {
		return total, big.NewInt(0), nil
	}

	signer, err := voucher.RecoverSigner()
	if err != nil {
		return &big.Int{}, &big.Int{}, err
	}
	if signer != vInfo.ChannelPayer {
		return &big.Int{}, &big.Int{}, fmt.Errorf("wrong signer: %+v, %+v", signer, vInfo.ChannelPayer)
	}
	// Check the difference between our largest voucher and this new one
	delta = big.NewInt(0).Sub(voucher.Amount, total)

	total = voucher.Amount
	vInfo.LargestVoucher = voucher

	err = vm.store.SetVoucherInfo(voucher.ChannelId, *vInfo)
	if err != nil {
		return nil, nil, err
	}
	return total, delta, nil
}

// ChannelRegistered returns  whether a channel has been registered with the voucher manager or not
func (vm *VoucherManager) ChannelRegistered(channelId types.Destination) bool {
	_, err := vm.store.GetVoucherInfo(channelId)
	return err == nil
}

// Paid returns the total amount paid so far on a channel
func (vm *VoucherManager) Paid(chanId types.Destination) (*big.Int, error) {
	v, err := vm.store.GetVoucherInfo(chanId)
	if err != nil {
		return &big.Int{}, fmt.Errorf("channel not registered: %w", err)
	}
	return v.LargestVoucher.Amount, nil
}

// Remaining returns the remaining amount of funds in the channel
func (vm *VoucherManager) Remaining(chanId types.Destination) (*big.Int, error) {
	v, err := vm.store.GetVoucherInfo(chanId)
	if err != nil {
		return &big.Int{}, fmt.Errorf("channel not registered: %w", err)
	}
	remaining := big.NewInt(0).Sub(v.StartingBalance, v.LargestVoucher.Amount)
	return remaining, nil
}
