package payments

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	nitroAbi "github.com/statechannels/go-nitro/abi"
	nitroCrypto "github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

func (v *Voucher) Hash() (types.Bytes32, error) {
	encoded, err := abi.Arguments{
		{Type: nitroAbi.Destination},
		{Type: nitroAbi.Uint256},
	}.Pack(v.ChannelId, v.Amount)

	if err != nil {
		return types.Bytes32{}, fmt.Errorf("failed to encode voucher: %w", err)
	}
	return crypto.Keccak256Hash(encoded), nil
}

func (v *Voucher) Sign(pk []byte) error {
	hash, err := v.Hash()
	if err != nil {
		return err
	}

	sig, err := nitroCrypto.SignEthereumMessage(hash.Bytes(), pk)

	if err != nil {
		return err
	}

	v.Signature = sig

	return nil
}

func (v *Voucher) RecoverSigner() (types.Address, error) {
	h, error := v.Hash()
	if error != nil {
		return types.Address{}, error
	}
	return nitroCrypto.RecoverEthereumMessageSigner(h[:], v.Signature)
}

// Equal returns true if the two vouchers have the same channel id, amount and signatures
func (v *Voucher) Equal(other *Voucher) bool {
	return v.ChannelId == other.ChannelId && v.Amount.Cmp(other.Amount) == 0 && v.Signature.Equal(other.Signature)
}
