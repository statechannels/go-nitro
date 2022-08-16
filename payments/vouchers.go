package payments

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	nitroAbi "github.com/statechannels/go-nitro/abi"
	nitroCrypto "github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

// Returns the channel id of the voucher
func (v Voucher) ChannelId() types.Destination {
	return v.channelId
}

func (v *Voucher) hash() (types.Bytes32, error) {
	encoded, err := abi.Arguments{
		{Type: nitroAbi.Destination},
		{Type: nitroAbi.Uint256},
	}.Pack(v.channelId, v.amount)

	if err != nil {
		return types.Bytes32{}, fmt.Errorf("failed to encode voucher: %w", err)
	}
	return crypto.Keccak256Hash(encoded), nil
}

func (v *Voucher) sign(pk []byte) error {
	hash, err := v.hash()
	if err != nil {
		return err
	}

	sig, err := nitroCrypto.SignEthereumMessage(hash.Bytes(), pk)

	if err != nil {
		return err
	}

	v.signature = sig

	return nil
}

func (v *Voucher) recoverSigner() (types.Address, error) {
	h, error := v.hash()
	if error != nil {
		return types.Address{}, error
	}
	return nitroCrypto.RecoverEthereumMessageSigner(h[:], v.signature)
}

// Equal returns true if the two vouchers have the same channel id, amount and signatures
func (v *Voucher) Equal(other *Voucher) bool {
	return v.channelId == other.channelId && v.amount.Cmp(other.amount) == 0 && v.signature.Equal(other.signature)
}

// NewVoucher constructs a voucher with the given channel id and amount
func NewVoucher(channelId types.Destination, amount *big.Int) *Voucher {
	v := Voucher{
		channelId: channelId,
		amount:    amount,
	}
	return &v
}
