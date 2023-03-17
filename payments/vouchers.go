package payments

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	nitroAbi "github.com/statechannels/go-nitro/abi"
	"github.com/statechannels/go-nitro/channel/state"
	nitroCrypto "github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

// A Voucher signed by Alice can be used by Bob to redeem payments in case of
// a misbehaving Alice.
//
// During normal operation, Alice & Bob would terminate the channel with an
// outcome reflecting the largest amount signed by Alice. For instance,
//   - if the channel started with balances {alice: 100, bob: 0}
//   - and the biggest voucher signed by alice had amount = 20
//   - then Alice and Bob would cooperatively conclude the channel with outcome
//     {alice: 80, bob: 20}
type Voucher struct {
	ChannelId types.Destination
	Amount    *big.Int
	Signature state.Signature
}

// VoucherInfo contains the largest voucher we've received on a channel.
// As well as details about the balance and who the payee/payer is.
type VoucherInfo struct {
	ChannelPayer    common.Address
	ChannelPayee    common.Address
	StartingBalance *big.Int
	LargestVoucher  Voucher
	Remaining       *big.Int
	Paid            *big.Int
}

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
