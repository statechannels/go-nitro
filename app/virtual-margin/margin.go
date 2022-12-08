package virtualmargin

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

type MarginApp struct {
	ChannelId      types.Destination
	LeaderAmount   *big.Int
	FollowerAmount *big.Int
	Version        *big.Int
	LeaderSig      state.Signature
	FollowerSig    state.Signature
}

type Balance struct {
	Leader   *big.Int
	Follower *big.Int
}

func (mv *MarginApp) Hash() (types.Bytes32, error) {
	encoded, err := abi.Arguments{
		{Type: nitroAbi.Destination},
		{Type: nitroAbi.Uint256},
		{Type: nitroAbi.Uint256},
		{Type: nitroAbi.Uint256},
	}.Pack(mv.ChannelId, mv.LeaderAmount, mv.FollowerAmount, mv.Version)

	if err != nil {
		return types.Bytes32{}, fmt.Errorf("failed to encode voucher: %w", err)
	}
	return crypto.Keccak256Hash(encoded), nil
}

func (ma *MarginApp) LeaderSign(pk []byte) error {
	hash, err := ma.Hash()
	if err != nil {
		return err
	}

	sig, err := nitroCrypto.SignEthereumMessage(hash.Bytes(), pk)

	if err != nil {
		return err
	}

	ma.LeaderSig = sig

	return nil
}

func (ma *MarginApp) FollowerSign(pk []byte) error {
	hash, err := ma.Hash()
	if err != nil {
		return err
	}

	sig, err := nitroCrypto.SignEthereumMessage(hash.Bytes(), pk)

	if err != nil {
		return err
	}

	ma.FollowerSig = sig

	return nil
}

func (v *MarginApp) RecoverLeaderSigner() (types.Address, error) {
	h, error := v.Hash()
	if error != nil {
		return types.Address{}, error
	}
	return nitroCrypto.RecoverEthereumMessageSigner(h[:], v.LeaderSig)
}

func (v *MarginApp) RecoverFollowerSigner() (types.Address, error) {
	h, error := v.Hash()
	if error != nil {
		return types.Address{}, error
	}
	return nitroCrypto.RecoverEthereumMessageSigner(h[:], v.FollowerSig)
}

func (v *MarginApp) Equal(other *MarginApp) bool {
	return v.ChannelId == other.ChannelId &&
		v.LeaderAmount.Cmp(other.LeaderAmount) == 0 &&
		v.FollowerAmount.Cmp(other.FollowerAmount) == 0 &&
		v.LeaderSig.Equal(other.LeaderSig) &&
		v.FollowerSig.Equal(other.FollowerSig)
}
