package NitroAdjudicator

import (
	"bytes"
	"context"
	"math/big"
	"testing"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/statechannels/go-nitro/abi"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	nc "github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

type actor struct {
	Address    types.Address
	PrivateKey []byte
}

func (a actor) Destination() types.Destination {
	return types.AddressToDestination(a.Address)
}

// actors namespaces the actors exported for test consumption
type actors struct {
	Alice actor
	Bob   actor
}

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

// Actors is the endpoint for tests to consume constructed statechannel
// network participants (public-key secret-key pairs)
var Actors actors = actors{
	Alice: actor{
		common.HexToAddress(`0xAAA6628Ec44A8a742987EF3A114dDFE2D4F7aDCE`),
		common.Hex2Bytes(`2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d`),
	},
	Bob: actor{
		common.HexToAddress(`0xBBB676f9cFF8D242e9eaC39D063848807d3D1D94`),
		common.Hex2Bytes(`0279651921cd800ac560c21ceea27aab0107b67daf436cdd25ce84cad30159b4`),
	},
}

func TestChallenge(t *testing.T) {

	turnNum := uint64(0)
	s := state.State{
		ChainId: big.NewInt(1337),
		Participants: []types.Address{
			Actors.Alice.Address,
			Actors.Bob.Address,
		},
		ChannelNonce:      big.NewInt(37140676580),
		AppDefinition:     common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`),
		ChallengeDuration: big.NewInt(60),
		AppData:           []byte{},
		Outcome:           outcome.Exit{},
		TurnNum:           turnNum,
		IsFinal:           false,
	}

	aSig, _ := s.Sign(Actors.Alice.PrivateKey)
	bSig, _ := s.Sign(Actors.Bob.PrivateKey)
	challengerSig, err := SignChallengeMessage(s, Actors.Alice.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	address := auth.From
	balance := new(big.Int)
	balance.SetString("10000000000000000000", 10) // 10 eth in wei
	gAlloc := map[common.Address]core.GenesisAccount{
		address: {Balance: balance},
	}

	blockGasLimit := uint64(4712388)

	sim := backends.NewSimulatedBackend(gAlloc, blockGasLimit)

	auth.GasPrice = big.NewInt(10000000000)

	naAddress, _, na, err := DeployNitroAdjudicator(auth, sim)

	if err != nil {
		t.Fatal(err)
	}
	sim.Commit()
	t.Log(naAddress)
	t.Log(na)
	tx, err := na.Challenge(
		auth,
		IForceMoveFixedPart(s.FixedPart()),
		big.NewInt(0),
		[]IForceMoveAppVariablePart{convertVariablePart(s.VariablePart())},
		0,
		[]IForceMoveSignature{convertSignature(aSig), convertSignature(bSig)},
		[]uint8{0, 0},
		convertSignature(challengerSig),
	)
	if err != nil {
		t.Fatal(err)
	}
	sim.Commit()

	receipt, err := sim.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		t.Fatal(err)
	}
	header, err := sim.HeaderByNumber(context.Background(), receipt.BlockNumber)
	if err != nil {
		t.Fatal(err)
	}
	challengeTime := big.NewInt(int64(header.Time))

	expectedFinalizesAt := big.NewInt(0).Add(challengeTime, s.ChallengeDuration)

	cId, _ := s.ChannelId()

	expectedOnChainStatus, err := generateStatus(s, expectedFinalizesAt)
	if err != nil {
		t.Fatal(err)
	}

	statusOnChain, err := na.StatusOf(&bind.CallOpts{}, cId)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(statusOnChain[:], expectedOnChainStatus) {
		t.Fatalf("Adjudicator not updated as expected, got %v wanted %v", common.Bytes2Hex(statusOnChain[:]), common.Bytes2Hex(expectedOnChainStatus[:]))
	}

}
