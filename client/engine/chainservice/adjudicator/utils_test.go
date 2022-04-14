package NitroAdjudicator

import (
	"math/big"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/abi"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	nc "github.com/statechannels/go-nitro/crypto"
)

func convertVariablePart(vp state.VariablePart) IForceMoveAppVariablePart {
	return IForceMoveAppVariablePart{
		AppData: vp.AppData,
		TurnNum: big.NewInt(int64(vp.TurnNum)),
		IsFinal: vp.IsFinal,
		Outcome: convertOutcome(vp.Outcome),
	}
}

func convertOutcome(o outcome.Exit) []ExitFormatSingleAssetExit {
	e := make([]ExitFormatSingleAssetExit, len(o))
	for i, sae := range o {
		e[i].Asset = sae.Asset
		e[i].Metadata = sae.Metadata
		e[i].Allocations = convertAllocations(sae.Allocations)
	}
	return e
}

func convertAllocations(as outcome.Allocations) []ExitFormatAllocation {
	b := make([]ExitFormatAllocation, len(as))
	for i, a := range as {
		b[i].Destination = a.Destination
		b[i].Amount = a.Amount
		b[i].AllocationType = uint8(a.AllocationType)
		b[i].Metadata = a.Metadata
	}
	return b
}

func convertSignature(s nc.Signature) IForceMoveSignature {
	sig := IForceMoveSignature{
		V: s.V,
	}
	copy(sig.R[:], s.R)
	copy(sig.S[:], s.S) // TODO we should just use 32 byte types, which would remove the need for this func
	return sig
}

func generateStatus(state state.State, finalizesAt *big.Int) ([]byte, error) {

	turnNumBytes := big.NewInt(int64(state.TurnNum)).FillBytes(make([]byte, 6))
	finalizesAtBytes := finalizesAt.FillBytes(make([]byte, 6))

	stateHash, err := state.Hash()
	if err != nil {
		return []byte{}, err
	}
	outcomeHash, err := state.Outcome.Hash()
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
