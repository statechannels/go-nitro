package NitroAdjudicator

import (
	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/abi"
	"github.com/statechannels/go-nitro/channel/state"
	nc "github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

func SignChallengeMessage(s state.State, privateKey []byte) (state.Signature, error) {

	challengeHash, err := hashChallengeMessage(s)

	if err != nil {
		return state.Signature{}, err
	}

	ecdsaKey, err := crypto.ToECDSA(privateKey)
	if err != nil {
		return state.Signature{}, err
	}

	sig, err := crypto.Sign(challengeHash[:], ecdsaKey)
	if err != nil {
		return state.Signature{}, err
	}

	return nc.SplitSignature(sig), nil
}

func hashChallengeMessage(s state.State) (types.Bytes32, error) {

	digest, err := s.Hash()
	if err != nil {
		return types.Bytes32{}, err
	}

	encoded, err := ethAbi.Arguments{{Type: abi.Destination}, {Type: abi.Bytes}}.Pack(digest, "forceMove")
	if err != nil {
		return types.Bytes32{}, err
	}

	return crypto.Keccak256Hash(encoded), nil

}
