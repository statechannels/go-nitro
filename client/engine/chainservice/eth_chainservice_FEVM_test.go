package chainservice

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func TestEthChainServiceFEVM(t *testing.T) {
	// Since this is hitting a contract on a test chain we only want to run it selectively
	t.Skip()
	// This is funded key on wallaby based on the test test ... fake mnemoic
	pkString := "6b65fdf763faebfbcf9a43d5ab3dd2fb639a3d69c10df99eddc0a6eb30a99ba7"
	// Due to https://github.com/filecoin-project/ref-fvm/issues/1182
	// the on-chain chainid() function returns the incorrect chain id (31415926)
	// To work around this we use ths incorrect chain id in the state
	// So the on-chain check passes
	workaroundChainId := big.NewInt(31415926)
	pk, err := crypto.HexToECDSA(pkString)
	if err != nil {
		t.Fatal(err)
	}

	client, err := ethclient.Dial("https://wallaby.node.glif.io/rpc/v0")

	if err != nil {
		t.Fatal(err)
	}
	wallabyChainId := big.NewInt(31415)
	// When submitting a transaction it's signed against a specific chain id
	// To get the correct signature we need to use the correct chain id that wallaby is expecting
	txSubmitter, err := bind.NewKeyedTransactorWithChainID(pk, wallabyChainId)
	if err != nil {
		log.Fatal(err)
	}
	// By setting the GasTipCap we signal this is a type 2 transaction
	// FEVM does NOT support type 1 transactions
	txSubmitter.GasTipCap = big.NewInt(300000)

	// This is the deployed contract on wallaby
	naAddress := common.HexToAddress("0xab5c7Ff206Ed23180DbF9c6F3b98Ec984D0b0aB8")
	caAddress := common.Address{}  // TODO use proper address
	vpaAddress := common.Address{} // TODO use proper address

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
		ChainId: workaroundChainId,
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
