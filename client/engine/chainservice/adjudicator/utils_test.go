package NitroAdjudicator

import (
	"math/big"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/abi"
	"github.com/statechannels/go-nitro/channel/state"
)

func generateStatus(state state.State, finalizesAt uint64) ([]byte, error) {
	turnNumBytes := big.NewInt(int64(state.TurnNum)).FillBytes(make([]byte, 6))
	finalizesAtBytes := new(big.Int).SetUint64(finalizesAt).FillBytes(make([]byte, 6))

	stateHash, err := state.Hash()
	if err != nil {
		return []byte{}, err
	}
	// TODO: Disabling this allows the rpc client to be imported
	// without importing go-ethereum/crypto
	// This allows the rpc client to be used in the boost repo
	outcomeHash := common.Hash{}
	if err != nil {
		return []byte{}, err
	}
	handprintPreimage, err := ethAbi.Arguments{{Type: abi.Bytes32}, {Type: abi.Bytes32}}.Pack(stateHash, outcomeHash)
	handprint := crypto.Keccak256(handprintPreimage)
	if err != nil {
		return []byte{}, err
	}
	fingerprint := handprint[12:]

	status := []byte(string(turnNumBytes) + string(finalizesAtBytes) + string(fingerprint))

	return status, nil
}
