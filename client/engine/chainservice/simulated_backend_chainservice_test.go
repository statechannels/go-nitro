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
	sim, na, naAddress, ethAccounts, err := SetupSimulatedBackend(1)
	if err != nil {
		t.Fatal(err)
	}

	cs := NewSimulatedBackendChainService(sim, na, naAddress, ethAccounts[0])

	// Prepare test data to trigger EthChainService
	testDeposit := types.Funds{
		common.HexToAddress("0x00"): big.NewInt(1),
	}
	channelID := types.Destination(common.HexToHash(`4ebd366d014a173765ba1e50f284c179ade31f20441bec41664712aac6cc461d`))
	testTx := protocols.NewDepositTransaction(channelID, testDeposit)

	out := cs.SubscribeToEvents(ethAccounts[0].From)
	// Submit transactiom
	cs.SendTransaction(testTx)

	// Check that the recieved event matches the expected event
	receivedEvent := <-out
	expectedEvent := DepositedEvent{CommonEvent: CommonEvent{channelID: channelID, BlockNum: 2}, Holdings: testDeposit}
	if diff := cmp.Diff(expectedEvent, receivedEvent, cmp.AllowUnexported(CommonEvent{})); diff != "" {
		t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
	}

	sim.Close()
}

func TestConcludeSimulatedBackendChainService(t *testing.T) {
	// Generate Signatures
	aSig, _ := concludeState.Sign(Alice.PrivateKey)
	bSig, _ := concludeState.Sign(Bob.PrivateKey)

	sim, na, naAddress, ethAccounts, err := SetupSimulatedBackend(1)
	if err != nil {
		t.Fatal(err)
	}
	cs := NewSimulatedBackendChainService(sim, na, naAddress, ethAccounts[0])
	out := cs.SubscribeToEvents(ethAccounts[0].From)

	// Fund channel
	testDeposit := types.Funds{
		common.HexToAddress("0x00"): big.NewInt(2),
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
	expectedEvent2 := AllocationUpdatedEvent{CommonEvent: CommonEvent{channelID: cId, BlockNum: 3}, Holdings: types.Funds{}}
	if diff := cmp.Diff(expectedEvent2, allocationUpdatedEvent, cmp.AllowUnexported(CommonEvent{})); diff != "" {
		t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
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

	// Not sure if this is necessary
	sim.Close()
}
