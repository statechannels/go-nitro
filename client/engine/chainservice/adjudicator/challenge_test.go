package NitroAdjudicator

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
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

type backendWithTxReader interface {
	bind.ContractBackend
	ethereum.TransactionReader
}

type preparedChain struct {
	chain               backendWithTxReader
	consensusAppAddress common.Address
	na                  NitroAdjudicator
	txSubmitter         *bind.TransactOpts
}

func TestChallenge(t *testing.T) {
	simulatedBackendPreparedChain := prepareSimulatedBackend(t)

	for _, turnNum := range []uint64{0, 1, 2} {
		t.Run("SimulatedBackend: turnNum = "+fmt.Sprint(turnNum),
			func(t *testing.T) {
				runChallengeWithTurnNum(t, turnNum, simulatedBackendPreparedChain)
			})
	}

	hyperspacePreparedChain := prepareHyperspaceBackend(t)

	for _, turnNum := range []uint64{0, 1, 2} {
		t.Run("HyperSpaceBackend: turnNum = "+fmt.Sprint(turnNum),
			func(t *testing.T) {
				t.Skip() // We run this test manually.
				runChallengeWithTurnNum(t, turnNum, hyperspacePreparedChain)
			})
	}
}

func runChallengeWithTurnNum(t *testing.T, turnNum uint64, pc preparedChain) {
	chain := pc.chain
	consensusAppAddress := pc.consensusAppAddress
	na := pc.na

	// Mine a block if using a simulated backend
	if sim, isSimulatedBackEnd := chain.(*backends.SimulatedBackend); isSimulatedBackEnd {
		sim.Commit()
	}

	var s = state.State{
		Participants: []types.Address{
			Actors.Alice.Address,
			Actors.Bob.Address,
		},
		ChannelNonce:      37140676580,
		AppDefinition:     consensusAppAddress,
		ChallengeDuration: 60,
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

	// Mine a block if using a simulated backend
	if sim, isSimulatedBackEnd := chain.(*backends.SimulatedBackend); isSimulatedBackEnd {
		sim.Commit()
	}

	// Construct support proof
	candidate := INitroTypesSignedVariablePart{
		ConvertVariablePart(s.VariablePart()),
		[]INitroTypesSignature{ConvertSignature(aSig), ConvertSignature(bSig)},
	}
	proof := make([]INitroTypesSignedVariablePart, 0)

	// Fire off a Challenge tx
	tx, err := na.Challenge(
		pc.txSubmitter,
		INitroTypesFixedPart(ConvertFixedPart(s.FixedPart())),
		proof,
		candidate,
		ConvertSignature(challengerSig),
	)
	if err != nil {
		t.Log(tx)
		t.Fatal(err)
	}

	// Mine a block if using a simulated backend
	if sim, isSimulatedBackEnd := chain.(*backends.SimulatedBackend); isSimulatedBackEnd {
		sim.Commit()
	}

	// Compute challenge time
	receipt, err := chain.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		t.Fatal(err)
	}
	header, err := chain.HeaderByNumber(context.Background(), receipt.BlockNumber)
	if err != nil {
		t.Fatal(err)
	}

	// Generate expectation
	expectedFinalizesAt := header.Time + uint64(s.ChallengeDuration)
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

}

func prepareSimulatedBackend(t *testing.T) preparedChain {
	// Setup transacting EOA
	key, _ := crypto.GenerateKey()
	auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337)) // 1337 according to godoc on backends.NewSimulatedBackend
	address := auth.From

	// Setup a second transacting EOA
	key2, _ := crypto.GenerateKey()
	auth2, _ := bind.NewKeyedTransactorWithChainID(key2, big.NewInt(1337)) // 1337 according to godoc on backends.NewSimulatedBackend
	address2 := auth2.From

	// Setup "blockchain" params
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

	sim.Commit()

	// Deploy ConsensusApp
	consensusAppAddress, _, _, err := ConsensusApp.DeployConsensusApp(auth2, sim)

	if err != nil {
		t.Fatal(err)
	}

	return preparedChain{
		sim,
		consensusAppAddress,
		*na,
		auth,
	}

}

func prepareHyperspaceBackend(t *testing.T) preparedChain {
	// This is the mnemonic for the prefunded accounts on wallaby.
	// The first 25 accounts will be prefunded.
	const MNEMONIC = "army forest resource shop tray cluster teach cause spice judge link oppose"

	// This is the HD path to use when deriving accounts from the mnemonic
	const HD_PATH = "m/44'/1'/0'/0"

	hyperspaceChainId := big.NewInt(3141)

	// Setup transacting EOA
	wallet, err := hdwallet.NewFromMnemonic(MNEMONIC)
	if err != nil {
		t.Fatal(err)
	}

	// The 0th account is usually used for deployment so we grab the 1st account
	a, err := wallet.Derive(hdwallet.MustParseDerivationPath(fmt.Sprintf("%s/%d", HD_PATH, 1)), false)
	if err != nil {
		t.Fatal(err)
	}

	//PK: 0x1688820ffc6a811e09ff17eccec23d8dec4850c3098ffc03ac4aa38dd8f3a994
	// corresponding ETH address is 0x280c53E2C574418D8d6d8d651d4c3323F4b194Be
	// corresponding f4 address (delegated) is t410ffagfhywforay3dlnrvsr2tbtep2ldff6xuxkrjq.
	pk, err := wallet.PrivateKey(a)

	if err != nil {
		t.Fatal(err)
	}

	client, err := ethclient.Dial("https://api.hyperspace.node.glif.io/rpc/v1")
	if err != nil {
		t.Fatal(err)
	}

	// When submitting a transaction it's signed against a specific chain id
	// To get the correct signature we need to use the correct chain id that wallaby is expecting
	txSubmitter, err := bind.NewKeyedTransactorWithChainID(pk, hyperspaceChainId)
	if err != nil {
		t.Fatal(err)
	}
	// By setting the GasTipCap we signal this is a type 2 transaction
	// FEVM does NOT support type 1 transactions
	txSubmitter.GasTipCap = big.NewInt(300000)

	// This is the deployed contract on wallaby
	// If wallaby gets reset this will need to be redeployed by running:
	// WALLABY_DEPLOYER_PK="f4d69c36885541f56f4728ddc002a6fa2fcb26c9f608910310a776c83b7fde47" npx hardhat deploy --network hyperspace --deploy-scripts ./hardhat-deploy-fvm --reset
	// The PK corresponds to account 0xE39dce95b1A924E2472E24C20C55eA3559a09251.
	// It should be prefunded after every wallaby reset.
	naAddress := common.HexToAddress("0xd354e517646872dd27d702c6f9f7a6ada3646c47")
	caAddress := common.HexToAddress("0x91315b24Bbf0c2e640d9BD47a2ADe927e72d9173")

	na, err := NewNitroAdjudicator(naAddress, client)
	if err != nil {
		t.Fatal(err)
	}

	balance := new(big.Int)
	balance.SetString("10000000000000000000", 10) // 10 eth in wei

	return preparedChain{
		client,
		caAddress,
		*na,
		txSubmitter,
	}

}
