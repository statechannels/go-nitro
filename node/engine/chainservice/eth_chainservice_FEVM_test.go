package chainservice

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/testactors"
	NitroAdjudicator "github.com/statechannels/go-nitro/node/engine/chainservice/adjudicator"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// This is the mnemonic for the prefunded accounts on hyperspace.
// The first 25 accounts will be prefunded.
const HYPERSPACE_MNEMONIC = "army forest resource shop tray cluster teach cause spice judge link oppose"

// This is the HD path to use when deriving accounts from the mnemonic
const HYPERSPACE_HD_PATH = "m/44'/1'/0'/0"

func TestEthChainServiceFEVM(t *testing.T) {
	// We only run the test if the RUN_FEVM_TESTS env var is set to true
	if key, found := os.LookupEnv("RUN_FEVM_TESTS"); !found || strings.ToLower(key) != "true" {
		t.Skip()
	}

	testAgainstEndpoint(t, "https://api.hyperspace.node.glif.io/rpc/v1", "test_fevm_https_endpoint.log", getPK(1))
	testAgainstEndpoint(t, "wss://wss.hyperspace.node.glif.io/apigw/lotus/rpc/v0", "test_fevm_wss_endpoint.log", getPK(2))
}

// getPK returns the private key for the account at the given index using the hyperspace mnemonic and path
func getPK(index uint) *ecdsa.PrivateKey {
	wallet, err := hdwallet.NewFromMnemonic(HYPERSPACE_MNEMONIC)
	if err != nil {
		panic(err)
	}

	a, err := wallet.Derive(hdwallet.MustParseDerivationPath(fmt.Sprintf("%s/%d", HYPERSPACE_HD_PATH, 2)), false)
	if err != nil {
		panic(err)
	}
	pk, err := wallet.PrivateKey(a)
	if err != nil {
		panic(err)
	}
	return pk
}

// testAgainstEndpoint runs a simple chain test against the provided endpoint
func testAgainstEndpoint(t *testing.T, endpoint string, logFile string, pk *ecdsa.PrivateKey) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	hyperspaceChainId := big.NewInt(3141)
	// When submitting a transaction it's signed against a specific chain id
	// To get the correct signature we need to use the correct chain id that wallaby is expecting
	txSubmitter, err := bind.NewKeyedTransactorWithChainID(pk, hyperspaceChainId)
	if err != nil {
		log.Fatal(err)
	}
	// By setting the GasTipCap we signal this is a type 2 transaction
	// FEVM does NOT support type 1 transactions
	txSubmitter.GasTipCap = big.NewInt(300000)

	// This is the deployed contract on wallaby
	// If wallaby gets reset this will need to be redeployed by running:
	// WALLABY_DEPLOYER_PK="f4d69c36885541f56f4728ddc002a6fa2fcb26c9f608910310a776c83b7fde47" npx hardhat deploy --network hyperspace --deploy-scripts ./hardhat-deploy-fvm --reset
	// The PK corresponds to account 0xE39dce95b1A924E2472E24C20C55eA3559a09251.
	// It should be prefunded after every wallaby reset.
	naAddress := common.HexToAddress("0xb095A67b76179dAFB5a56628378b919052f978c9")
	caAddress := common.HexToAddress("0x64D444a7B99f07d3a1d69F82798Eaa0a98E04543")
	vpaAddress := common.HexToAddress("0x76E957873156526D8BA482c4b1CE26a60dF639Ba")
	na, err := NitroAdjudicator.NewNitroAdjudicator(naAddress, client)
	if err != nil {
		t.Fatal(err)
	}

	cs, err := newEthChainService(client, na, naAddress, caAddress, vpaAddress, txSubmitter)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare test data to trigger EthChainService
	testDeposit := types.Funds{
		common.HexToAddress("0x00"): big.NewInt(10),
	}
	var (
		Alice = testactors.Alice
		Bob   = testactors.Bob
	)
	concludeOutcome := outcome.Exit{
		outcome.SingleAssetExit{
			Asset: types.Address{},
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: types.AddressToDestination(Alice.Address()),
					Amount:      big.NewInt(1),
				},
				outcome.Allocation{
					Destination: types.AddressToDestination(Bob.Address()),
					Amount:      big.NewInt(1),
				},
			},
		},
	}

	concludeState := state.State{
		Participants: []types.Address{
			Alice.Address(),
			Bob.Address(),
		},
		ChannelNonce:      37140676580,
		AppDefinition:     types.Address{},
		ChallengeDuration: 0,
		AppData:           []byte{},
		Outcome:           concludeOutcome,
		TurnNum:           uint64(2),
		IsFinal:           true,
	}
	cId := concludeState.ChannelId()
	testTx := protocols.NewDepositTransaction(cId, testDeposit)

	out := cs.EventFeed()

	// Submit transactiom
	err = cs.SendTransaction(testTx)
	if err != nil {
		t.Fatal(err)
	}

	// Generate Signatures
	aSig, _ := concludeState.Sign(Alice.PrivateKey)
	bSig, _ := concludeState.Sign(Bob.PrivateKey)

	signedConcludeState := state.NewSignedState(concludeState)
	err = signedConcludeState.AddSignature(aSig)
	if err != nil {
		t.Fatal(err)
	}
	err = signedConcludeState.AddSignature(bSig)
	if err != nil {
		t.Fatal(err)
	}
	concludeTx := protocols.NewWithdrawAllTransaction(cId, signedConcludeState)
	err = cs.SendTransaction(concludeTx)
	if err != nil {
		t.Fatal(err)
	}

	// Inspect state of chain (call StatusOf)
	statusOnChain, err := na.StatusOf(&bind.CallOpts{}, cId)
	if err != nil {
		t.Fatal(err)
	}

	emptyBytes := [32]byte{}
	// Make assertion
	if !bytes.Equal(statusOnChain[:], emptyBytes[:]) {
		t.Fatalf("Adjudicator not updated as expected, got %v wanted %v", common.Bytes2Hex(statusOnChain[:]), common.Bytes2Hex(emptyBytes[:]))
	}
	for i := 0; i < 3; i++ {
		ignoreBlockNum := cmpopts.IgnoreFields(commonEvent{}, "BlockNum")
		ignoreNowHeld := cmpopts.IgnoreFields(DepositedEvent{}, "NowHeld")

		receivedEvent := <-out
		fmt.Printf("Checking event %+v\n", receivedEvent)
		switch receivedEvent := receivedEvent.(type) {
		case DepositedEvent:

			expectedEvent := NewDepositedEvent(cId, 2, receivedEvent.Asset, testDeposit[receivedEvent.Asset])
			// TODO to validate BlockNum and NowHeld values, chain state prior to transaction must be inspected

			if diff := cmp.Diff(expectedEvent, receivedEvent,
				cmp.AllowUnexported(DepositedEvent{}, commonEvent{}, big.Int{}),
				ignoreBlockNum,
				ignoreNowHeld,
				cmpopts.IgnoreFields(DepositedEvent{}, "assetAndAmount")); diff != "" {
				t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
			}

		case ConcludedEvent:
			expectedEvent := ConcludedEvent{commonEvent: commonEvent{channelID: cId, blockNum: 3}}
			if diff := cmp.Diff(expectedEvent, receivedEvent,
				cmp.AllowUnexported(ConcludedEvent{}, commonEvent{}),
				ignoreBlockNum, ignoreNowHeld); diff != "" {
				t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
			}

		case AllocationUpdatedEvent:
			expectedEvent2 := NewAllocationUpdatedEvent(cId, 3, common.Address{}, new(big.Int).SetInt64(1))

			if diff := cmp.Diff(expectedEvent2, receivedEvent,
				cmp.AllowUnexported(AllocationUpdatedEvent{}, commonEvent{}, big.Int{}),
				ignoreBlockNum,
				ignoreNowHeld,
				cmpopts.IgnoreFields(AllocationUpdatedEvent{}, "assetAndAmount")); diff != "" {
				t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
			}

		default:
			fmt.Printf("Ignoring event  %+v\n", receivedEvent)

		}

	}
}
