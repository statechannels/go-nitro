package chainservice

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// This is the mnemonic for the prefunded accounts on wallaby.
// The first 25 accounts will be prefunded.
const WALLABY_MNEMONIC = "army forest resource shop tray cluster teach cause spice judge link oppose"

// This is the HD path to use when deriving accounts from the mnemonic
const WALLABY_HD_PATH = "m/44'/1'/0'/0"

func TestEthChainServiceFEVM(t *testing.T) {
	// Since this is hitting a contract on a test chain we only want to run it selectively
	t.Skip()
	wallet, err := hdwallet.NewFromMnemonic(WALLABY_MNEMONIC)
	if err != nil {
		panic(err)
	}

	// The 0th account is usually used for deployment so we grab the 1st account
	a, err := wallet.Derive(hdwallet.MustParseDerivationPath(fmt.Sprintf("%s/%d", WALLABY_HD_PATH, 1)), false)
	if err != nil {
		panic(err)
	}

	//PK: 0x1688820ffc6a811e09ff17eccec23d8dec4850c3098ffc03ac4aa38dd8f3a994
	// corresponding ETH address is 0x280c53E2C574418D8d6d8d651d4c3323F4b194Be
	// corresponding f4 address (delegated) is t410ffagfhywforay3dlnrvsr2tbtep2ldff6xuxkrjq.
	pk, err := wallet.PrivateKey(a)

	if err != nil {
		panic(err)
	}

	client, err := ethclient.Dial("https://api.hyperspace.node.glif.io/rpc/v1")

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

	cs, err := NewEthChainService(client, na, naAddress, caAddress, vpaAddress, txSubmitter, NoopLogger{})
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
	var concludeOutcome = outcome.Exit{
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

	var concludeState = state.State{
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

			expectedEvent := NewDepositedEvent(cId, 2, receivedEvent.AssetAddress, big.NewInt(0), testDeposit[receivedEvent.AssetAddress])
			// TODO to validate BlockNum and NowHeld values, chain state prior to transaction must be inspected

			if diff := cmp.Diff(expectedEvent, receivedEvent,
				cmp.AllowUnexported(DepositedEvent{}, commonEvent{}, big.Int{}),
				ignoreBlockNum,
				ignoreNowHeld,
				cmpopts.IgnoreFields(DepositedEvent{}, "assetAndAmount")); diff != "" {
				t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
			}

		case ConcludedEvent:
			expectedEvent := ConcludedEvent{commonEvent: commonEvent{channelID: cId, BlockNum: 3}}
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
