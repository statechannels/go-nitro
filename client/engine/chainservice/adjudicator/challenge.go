package NitroAdjudicator

import (
	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/abi"
	"github.com/statechannels/go-nitro/channel/state"
	nc "github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

// SignChallengeMessage generates the special signature required to launch a challenge. This is used to prevent non-participants from launching challenges.
func SignChallengeMessage(s state.State, privateKey []byte) (state.Signature, error) {
	challengeHash, err := hashChallengeMessage(s)
	if err != nil {
		return state.Signature{}, err
	}
	return nc.SignEthereumMessage(challengeHash[:], privateKey)
}

func hashChallengeMessage(s state.State) (types.Bytes32, error) {
	digest, err := s.Hash()
	if err != nil {
		return types.Bytes32{}, err
	}

	encoded, err := ethAbi.Arguments{{Type: abi.Bytes32}, {Type: abi.String}}.Pack(digest, "forceMove")
	if err != nil {
		return types.Bytes32{}, err
	}

	return crypto.Keccak256Hash(encoded), nil
}
