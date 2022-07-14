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

type (
	// A Voucher signed by Alice can be used by Bob to redeem payments in case of
	// a misbehaving Alice.
	//
	// During normal operation, Alice & Bob would terminate the channel with an
	// outcome reflecting the largest amount signed by Alice. For instance,
	// - if the channel started with balances {alice: 100, bob: 0}
	// - and the biggest voucher signed by alice had amount = 20
	// - then Alice and Bob would cooperatively conclude the channel with outcome
	//   {alice: 80, bob: 20}
	Voucher struct {
		channelId types.Destination
		amount    *big.Int
		signature state.Signature
	}
)

func (v Voucher) hash() (types.Bytes32, error) {
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

func (v Voucher) recoverSigner() (types.Address, error) {
	h, error := v.hash()
	if error != nil {
		return types.Address{}, error
	}
	return nitroCrypto.RecoverEthereumMessageSigner(h[:], v.signature)
}
