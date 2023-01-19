package chainservice

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
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

const DUMMY_TX_DATA = `{"type":"0x2","nonce":"0x5","gasPrice":null,"maxPriorityFeePerGas":"0x1","maxFeePerGas":"0x63d59447","gas":"0x29a4d","value":"0x0","input":"0x180257e6000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000001600000000000000000000000000000000000000000000000000000000000000180000000000000000000000000000000000000000000000000000000000000053900000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000008a5c1bfe4000000000000000000000000ff64107479e1fab0865e131331f79fdde3be877500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce000000000000000000000000bbb676f9cff8d242e9eac39d063848807d3d1d94000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000340000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000002e0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000f5a1bb5607c9d079e46d1b3dc33f257d937b43bd0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ee18ff1575055691009aa246ae608132c57a422c000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000001c25b0ef29e428213994cf93298afd77a7056f41074beeb87457f758b7c1a85ef5722be65e000735f98fbcf913972f32c2210a1fce9d435f0aec85277b0dcb479a000000000000000000000000000000000000000000000000000000000000001c5ce0129385cda8d34d91b57e0ede5838eaa24592667e202abb21f24c64bc2d4a3dcacae1708d41daa65f3e9a90d9f2946c5432e020a4ec00afcacef9e4f2d700","v":"0x1","r":"0xe1faa10220886233a77c017d4b9b23c8b8508fe8daf1c986b26258fd8955fd55","s":"0x119bf7e6afdd5059720366b09db7bbbabda7abd4b87a2625443e6cce1da7edf6","to":"0x8c8822b07ff9f58376e6082ed60f0f99d31c9cf5","chainId":"0x539","accessList":[],"hash":"0x41e368ccfe6024c5fcbefb8a87c0d3144fd94add737d7f2290ed3823e54ff884"}`
const FEVM_PARSE_ERROR = "json: cannot unmarshal hex string \"0x\" into Go struct field txJSON.v of type *hexutil.Big"

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

	client, err := ethclient.Dial("https://filecoin-hyperspace.chainstacklabs.com/rpc/v0")

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
	// The key "f4d69c36885541f56f4728ddc002a6fa2fcb26c9f608910310a776c83b7fde47" is 0th account from the  WALLABY_MNEMONIC and WALLABY_HD_PATH
	// (But the hardhat deploy script actually ends up using account 0xE39dce95b1A924E2472E24C20C55eA3559a09251 or t410f4oo45fnrvesoerzoetbayvpkgvm2besropxbvxi)
	// It should be prefunded after every wallaby reset.
	naAddress := common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")
	caAddress := common.Address{}  // TODO use proper address
	vpaAddress := common.Address{} // TODO use proper address

	na, err := NitroAdjudicator.NewNitroAdjudicator(naAddress, client)
	if err != nil {
		t.Fatal(err)
	}

	wrappedClient := &FEVMWorkaroundChain{client, client}

	cs, err := NewEthChainService(wrappedClient, na, naAddress, caAddress, vpaAddress, txSubmitter, NoopLogger{})
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

// FEVMWorkaroundChain is an ethChain with a custom implementation of TransactionByHash
// as a workaround for https://github.com/filecoin-project/ref-fvm/issues/1158
type FEVMWorkaroundChain struct {
	ethChain
	tr ethereum.TransactionReader
}

func (fwc *FEVMWorkaroundChain) TransactionByHash(ctx context.Context, txHash common.Hash) (tx *ethTypes.Transaction, isPending bool, err error) {
	tx, pending, err := fwc.tr.TransactionByHash(context.Background(), txHash)
	if err != nil && strings.Contains(err.Error(), FEVM_PARSE_ERROR) {
		fmt.Printf("WARNING: Cannot parse transaction, using dummy tx json with asset address of 0x00...")
		tx := &ethTypes.Transaction{}
		err := tx.UnmarshalJSON([]byte(DUMMY_TX_DATA))
		if err != nil {

			panic(err)
		}
		return tx, false, nil
	}
	return tx, pending, err
}
