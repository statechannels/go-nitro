package chainservice

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var (
	Alice = testactors.Alice
	Bob   = testactors.Bob
)

var concludeOutcome = outcome.Exit{
	outcome.SingleAssetExit{
		Asset: types.Address{},
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: types.AddressToDestination(common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`)),
				Amount:      big.NewInt(1),
			},
			outcome.Allocation{
				Destination: types.AddressToDestination(common.HexToAddress(`0xEe18fF1575055691009aa246aE608132C57a422c`)),
				Amount:      big.NewInt(1),
			},
		},
	},
}

var concludeState = state.State{
	ChainId: big.NewInt(1337),
	Participants: []types.Address{
		Alice.Address(),
		Bob.Address(),
	},
	ChannelNonce:      big.NewInt(37140676580),
	AppDefinition:     types.Address{},
	ChallengeDuration: &big.Int{},
	AppData:           []byte{},
	Outcome:           concludeOutcome,
	TurnNum:           uint64(2),
	IsFinal:           true,
}

func TestDepositSimulatedBackendChainService(t *testing.T) {
	one := big.NewInt(1)
	sim, bindings, ethAccounts, err := SetupSimulatedBackend(1)
	if err != nil {
		t.Fatal(err)
	}

	cs := NewSimulatedBackendChainService(sim, bindings.Adjudicator.Contract, bindings.Adjudicator.Address, ethAccounts[0])

	_, err = bindings.Token.Contract.Approve(ethAccounts[0], bindings.Adjudicator.Address, one)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare test data to trigger EthChainService
	testDeposit := types.Funds{
		common.HexToAddress("0x00"): one,
		bindings.Token.Address:      one,
	}
	channelID := types.Destination(common.HexToHash(`4ebd366d014a173765ba1e50f284c179ade31f20441bec41664712aac6cc461d`))
	testTx := protocols.NewDepositTransaction(channelID, testDeposit)

	out := cs.SubscribeToEvents(ethAccounts[0].From)
	// Submit transactiom
	cs.SendTransaction(testTx)

	// Check that the recieved events matches the expected event
	for i := 0; i < 2; i++ {
		receivedEvent := <-out
		dEvent := receivedEvent.(DepositedEvent)
		expectedEvent := DepositedEvent{CommonEvent: CommonEvent{channelID: channelID, BlockNum: 2}, Asset: dEvent.Asset, NowHeld: testDeposit[dEvent.Asset], AmountDeposited: testDeposit[dEvent.Asset]}
		if diff := cmp.Diff(expectedEvent, dEvent, cmp.AllowUnexported(CommonEvent{}, big.Int{})); diff != "" {
			t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
		}
		delete(testDeposit, dEvent.Asset)
	}

	if len(testDeposit) != 0 {
		t.Fatalf("Mismatch between the deposit transaction and the received events")
	}

	sim.Close()
}

func TestConcludeSimulatedBackendChainService(t *testing.T) {
	// Generate Signatures
	aSig, _ := concludeState.Sign(Alice.PrivateKey)
	bSig, _ := concludeState.Sign(Bob.PrivateKey)

	sim, bindings, ethAccounts, err := SetupSimulatedBackend(1)
	if err != nil {
		t.Fatal(err)
	}
	cs := NewSimulatedBackendChainService(sim, bindings.Adjudicator.Contract, bindings.Adjudicator.Address, ethAccounts[0])
	out := cs.SubscribeToEvents(ethAccounts[0].From)

	// Fund channel
	testDeposit := types.Funds{
		common.HexToAddress("0x00"): big.NewInt(3),
	}
	cId := concludeState.ChannelId()

	depositTx := protocols.NewDepositTransaction(cId, testDeposit)
	cs.SendTransaction(depositTx)
	<-out

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
	cs.SendTransaction(concludeTx)

	// Check that the recieved event matches the expected event
	concludedEvent := <-out
	expectedEvent := ConcludedEvent{CommonEvent: CommonEvent{channelID: cId, BlockNum: 3}}
	if diff := cmp.Diff(expectedEvent, concludedEvent, cmp.AllowUnexported(CommonEvent{})); diff != "" {
		t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
	}

	// Check that the recieved event matches the expected event
	allocationUpdatedEvent := <-out
	expectedEvent2 := AllocationUpdatedEvent{
		CommonEvent:  CommonEvent{channelID: cId, BlockNum: 3},
		AssetAddress: common.Address{},
		AssetAmount:  new(big.Int).SetInt64(1)}

	if diff := cmp.Diff(expectedEvent2, allocationUpdatedEvent, cmp.AllowUnexported(CommonEvent{}, big.Int{})); diff != "" {
		t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
	}

	// Inspect state of chain (call StatusOf)
	statusOnChain, err := bindings.Adjudicator.Contract.StatusOf(&bind.CallOpts{}, cId)
	if err != nil {
		t.Fatal(err)
	}

	emptyBytes := [32]byte{}
	// Make assertion
	if !bytes.Equal(statusOnChain[:], emptyBytes[:]) {
		t.Fatalf("Adjudicator not updated as expected, got %v wanted %v", common.Bytes2Hex(statusOnChain[:]), common.Bytes2Hex(emptyBytes[:]))
	}

	// Not sure if this is necessary
	sim.Close()
}
