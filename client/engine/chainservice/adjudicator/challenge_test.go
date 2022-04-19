package NitroAdjudicator

import (
	"bytes"
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
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

var turnNum = uint64(0)

var s = state.State{
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

func TestChallenge(t *testing.T) {

	// Generate Signatures
	aSig, _ := s.Sign(Actors.Alice.PrivateKey)
	bSig, _ := s.Sign(Actors.Bob.PrivateKey)
	challengerSig, err := SignChallengeMessage(s, Actors.Alice.PrivateKey)

	if err != nil {
		t.Fatal(err)
	}

	// Setup transacting EOA
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	auth.GasPrice = big.NewInt(10000000000)
	address := auth.From
	balance := new(big.Int)
	balance.SetString("10000000000000000000", 10) // 10 eth in wei

	// Setup "blockchain"
	gAlloc := map[common.Address]core.GenesisAccount{
		address: {Balance: balance},
	}
	blockGasLimit := uint64(4712388)
	sim := backends.NewSimulatedBackend(gAlloc, blockGasLimit)

	// Deploy Adjudicator
	_, _, na, err := DeployNitroAdjudicator(auth, sim)

	if err != nil {
		t.Fatal(err)
	}

	// Mine a block
	sim.Commit()

	// Fire off a Challenge tx
	tx, err := na.Challenge(
		auth,
		IForceMoveFixedPart(s.FixedPart()),
		[]IForceMoveAppVariablePart{convertVariablePart(s.VariablePart())},
		[]IForceMoveSignature{convertSignature(aSig), convertSignature(bSig)},
		[]uint8{0, 0},
		convertSignature(challengerSig),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Mine a block
	sim.Commit()

	// Compute challenge time
	receipt, err := sim.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		t.Fatal(err)
	}
	header, err := sim.HeaderByNumber(context.Background(), receipt.BlockNumber)
	if err != nil {
		t.Fatal(err)
	}
	challengeTime := big.NewInt(int64(header.Time))

	// Generate expectation
	expectedFinalizesAt := big.NewInt(0).Add(challengeTime, s.ChallengeDuration)
	cId, _ := s.ChannelId()
	expectedOnChainStatus, err := generateStatus(s, expectedFinalizesAt)
	if err != nil {
		t.Fatal(err)
	}

	// Inspect state of chain (call StatusOf)
	statusOnChain, err := na.StatusOf(&bind.CallOpts{}, cId)
	if err != nil {
		t.Fatal(err)
	}

	// Make assertion
	if !bytes.Equal(statusOnChain[:], expectedOnChainStatus) {
		t.Fatalf("Adjudicator not updated as expected, got %v wanted %v", common.Bytes2Hex(statusOnChain[:]), common.Bytes2Hex(expectedOnChainStatus[:]))
	}

}
