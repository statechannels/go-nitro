package payments

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
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

// Balance stores the remaining and paid funds in a channel.
type Balance struct {
	Remaining *big.Int
	Paid      *big.Int
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

// IsZero returns true if the voucher is the default zero value
func (v *Voucher) IsZero() bool {
	return v.Equal(&Voucher{})
}

func (v *Voucher) Clone() *Voucher {
	var amount *big.Int
	if v.Amount != nil {
		amount = new(big.Int).Set(v.Amount)
	}
	signature := state.CloneSignature(v.Signature)
	return &Voucher{
		ChannelId: v.ChannelId,
		Amount:    amount,
		Signature: signature,
	}
}

func (v Voucher) SortInfo() (types.Destination, uint64) {
	return v.ChannelId, 0
}

// NewVoucher constructs a voucher with the given channel id and amount
func NewVoucher(channelId types.Destination, amount *big.Int) *Voucher {
	v := Voucher{
		ChannelId: channelId,
		Amount:    amount,
	}
	return &v
}

// NewSignedVoucher constructs a voucher with the given channel id and amount signed by the provided private key
func NewSignedVoucher(channelId types.Destination, amount *big.Int, pk []byte) (*Voucher, error) {
	v := NewVoucher(channelId, amount)
	err := v.Sign(pk)
	return v, err
}
