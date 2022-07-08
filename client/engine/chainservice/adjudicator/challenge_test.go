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
	ConsensusApp "github.com/statechannels/go-nitro/client/engine/chainservice/consensusapp"
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

// TODO https://github.com/statechannels/go-nitro/issues/785
// This test should cover several different turn numbers, including pre and post fund tun numbers.
var turnNum = uint64(2)

func TestChallenge(t *testing.T) {

	// Setup transacting EOA
	key, _ := crypto.GenerateKey()
	auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337)) // 1337 according to godoc on backends.NewSimulatedBackend
	auth.GasPrice = big.NewInt(10000000000)
	address := auth.From

	// Setup a second transacting EOA
	key2, _ := crypto.GenerateKey()
	auth2, _ := bind.NewKeyedTransactorWithChainID(key2, big.NewInt(1337)) // 1337 according to godoc on backends.NewSimulatedBackend
	address2 := auth2.From

	// Setup "blockchain"
	balance := new(big.Int)
	balance.SetString("10000000000000000000", 10) // 10 eth in wei
	gAlloc := map[common.Address]core.GenesisAccount{
		address:  {Balance: balance},
		address2: {Balance: balance},
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

	// Deploy ConsensusApp
	consensusAppAddress, _, _, err := ConsensusApp.DeployConsensusApp(auth2, sim)

	if err != nil {
		t.Fatal(err)
	}

	// Mine a block
	sim.Commit()

	var s = state.State{
		ChainId: big.NewInt(1337),
		Participants: []types.Address{
			Actors.Alice.Address,
			Actors.Bob.Address,
		},
		ChannelNonce:      big.NewInt(37140676580),
		AppDefinition:     consensusAppAddress,
		ChallengeDuration: big.NewInt(60),
		AppData:           []byte{},
		Outcome:           outcome.Exit{},
		TurnNum:           turnNum,
		IsFinal:           false,
	}

	// Generate Signatures
	aSig, _ := s.Sign(Actors.Alice.PrivateKey)
	bSig, _ := s.Sign(Actors.Bob.PrivateKey)
	challengerSig, err := SignChallengeMessage(s, Actors.Alice.PrivateKey)

	if err != nil {
		t.Fatal(err)
	}

	// Mine a block
	sim.Commit()

	// Fire off a Challenge tx
	tx, err := na.Challenge(
		auth,
		INitroTypesFixedPart(s.FixedPart()),
		[]INitroTypesSignedVariablePart{
			{
				convertVariablePart(s.VariablePart()),
				[]INitroTypesSignature{convertSignature(aSig), convertSignature(bSig)},
				big.NewInt(0b11),
			},
		},
		convertSignature(challengerSig),
	)
	if err != nil {
		t.Log(tx)
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
	cId := s.ChannelId()
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

	// Not sure if this is necessary
	sim.Close()
}
