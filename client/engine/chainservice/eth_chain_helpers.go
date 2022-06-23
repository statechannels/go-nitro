package chainservice

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
)

func getChainHoldings(na *NitroAdjudicator.NitroAdjudicator, tx *types.Transaction, event *NitroAdjudicator.NitroAdjudicatorAllocationUpdated) (common.Address, *big.Int) {
	assetAddress := assetAddressForIndex(na, tx, event.AssetIndex)
	amount, err := na.Holdings(&bind.CallOpts{BlockNumber: new(big.Int).SetUint64(event.Raw.BlockNumber)}, assetAddressForIndex(na, tx, event.AssetIndex), event.ChannelId)
	if err != nil {
		panic(err)
	}
	return assetAddress, amount
}

func assetAddressForIndex(na *NitroAdjudicator.NitroAdjudicator, tx *types.Transaction, index *big.Int) common.Address {
	abi, err := NitroAdjudicator.NitroAdjudicatorMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	params, err := decodeTxParams(abi, tx.Data())
	if err != nil {
		panic(err)
	}
	variablePart := params["latestVariablePart"].(struct {
		Outcome []struct {
			Asset       common.Address "json:\"asset\""
			Metadata    []uint8        "json:\"metadata\""
			Allocations []struct {
				Destination    [32]uint8 "json:\"destination\""
				Amount         *big.Int  "json:\"amount\""
				AllocationType uint8     "json:\"allocationType\""
				Metadata       []uint8   "json:\"metadata\""
			} "json:\"allocations\""
		} "json:\"outcome\""
		AppData []uint8  "json:\"appData\""
		TurnNum *big.Int "json:\"turnNum\""
		IsFinal bool     "json:\"isFinal\""
	})
	return variablePart.Outcome[index.Int64()].Asset

}

func decodeTxParams(abi *abi.ABI, data []byte) (map[string]interface{}, error) {
	v := map[string]interface{}{}
	m, err := abi.MethodById(data[:4])
	if err != nil {
		return map[string]interface{}{}, err

	}
	if err := m.Inputs.UnpackIntoMap(v, data[4:]); err != nil {
		return map[string]interface{}{}, err
	}

	return v, nil
}
