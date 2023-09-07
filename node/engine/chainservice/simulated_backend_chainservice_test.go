package chainservice

import (
	"bytes"
	"log/slog"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/logging"
	"github.com/statechannels/go-nitro/internal/testactors"
	NitroAdjudicator "github.com/statechannels/go-nitro/node/engine/chainservice/adjudicator"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var (
	CHALLENGE_DURATION = uint32(1000) // 1000 seconds. Much longer than the duration of the test
	Alice              = testactors.Alice
	Bob                = testactors.Bob
	challengeBlockNum  = uint64(2)
	depositBlockNum    = uint64(5)
	concludeBlockNum   = uint64(8)
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

func TestSimulatedBackendChainService(t *testing.T) {
	logging.SetupDefaultFileLogger("simulatedBackendChainService.log", slog.LevelDebug)

	one := big.NewInt(1)
	three := big.NewInt(3)

	var receivedEvent Event

	sim, bindings, ethAccounts, err := SetupSimulatedBackend(2)
	defer closeSimulatedChain(t, sim)
	if err != nil {
		t.Fatal(err)
	}

	cs, err := NewSimulatedBackendChainService(sim, bindings, ethAccounts[0])
	defer closeChainService(t, cs)
	if err != nil {
		t.Fatal(err)
	}

	concludeState := state.State{
		Participants: []types.Address{
			Alice.Address(),
			Bob.Address(),
		},
		ChannelNonce:      37140676580,
		AppDefinition:     bindings.ConsensusApp.Address,
		ChallengeDuration: CHALLENGE_DURATION,
		AppData:           []byte{},
		Outcome:           concludeOutcome,
		TurnNum:           uint64(2),
		IsFinal:           true,
	}

	challengerSig, err := NitroAdjudicator.SignChallengeMessage(concludeState, Alice.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	concludeSignedState := state.NewSignedState(concludeState)
	aSig, err := concludeState.Sign(Alice.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	bSig, err := concludeState.Sign(Bob.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	_ = concludeSignedState.AddSignature(aSig)
	_ = concludeSignedState.AddSignature(bSig)

	challengeTx := protocols.NewChallengeTransaction(concludeState.ChannelId(), concludeSignedState, make([]state.SignedState, 0), challengerSig)

	out := cs.EventFeed()
	err = cs.SendTransaction(challengeTx)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the received events matches the expected event
	receivedEvent = <-out
	crEvent := receivedEvent.(ChallengeRegisteredEvent)
	expectedChallengeRegisteredEvent := NewChallengeRegisteredEvent(concludeState.ChannelId(), challengeBlockNum, crEvent.candidate, crEvent.candidateSignatures)
	if diff := cmp.Diff(expectedChallengeRegisteredEvent, crEvent, cmp.AllowUnexported(ChallengeRegisteredEvent{}, commonEvent{}, big.Int{})); diff != "" {
		t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
	}

	testDeposit := types.Funds{
		common.HexToAddress("0x00"): three,
		bindings.Token.Address:      one,
	}
	depositTx := protocols.NewDepositTransaction(concludeState.ChannelId(), testDeposit)

	// Submit transaction
	err = cs.SendTransaction(depositTx)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the received events match the expected events
	for i := 0; i < 2; i++ {
		receivedEvent = <-out
		dEvent := receivedEvent.(DepositedEvent)
		expectedDepositEvent := NewDepositedEvent(concludeState.ChannelId(), depositBlockNum, dEvent.Asset, testDeposit[dEvent.Asset])
		if diff := cmp.Diff(expectedDepositEvent, dEvent, cmp.AllowUnexported(DepositedEvent{}, commonEvent{}, big.Int{})); diff != "" {
			t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
		}
		delete(testDeposit, dEvent.Asset)
	}

	if len(testDeposit) != 0 {
		t.Fatalf("Mismatch between the deposit transaction and the received events")
	}

	cId := concludeState.ChannelId()

	signedConcludeState := state.NewSignedState(concludeState)
	err = signedConcludeState.AddSignature(aSig)
	if err != nil {
		t.Fatal(err)
	}
	err = signedConcludeState.AddSignature(bSig)
	if err != nil {
		t.Fatal(err)
	}

	// Start new chain service. It should detect old chain events that were emitted while it was offline
	cs2, err := NewSimulatedBackendChainService(sim, bindings, ethAccounts[1])
	defer closeChainService(t, cs2)
	if err != nil {
		t.Fatal(err)
	}

	concludeTx := protocols.NewWithdrawAllTransaction(cId, signedConcludeState)
	err = cs.SendTransaction(concludeTx)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the recieved event matches the expected event
	concludedEvent := <-out
	expectedConcludedEvent := ConcludedEvent{commonEvent: commonEvent{channelID: cId, blockNum: concludeBlockNum}}
	if diff := cmp.Diff(expectedConcludedEvent, concludedEvent, cmp.AllowUnexported(ConcludedEvent{}, commonEvent{})); diff != "" {
		t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
	}

	// Check that the recieved event matches the expected event
	allocationUpdatedEvent := <-out
	expectedAllocationUpdatedEvent := NewAllocationUpdatedEvent(cId, concludeBlockNum, common.Address{}, new(big.Int).SetInt64(1))
	if diff := cmp.Diff(expectedAllocationUpdatedEvent, allocationUpdatedEvent, cmp.AllowUnexported(AllocationUpdatedEvent{}, commonEvent{}, big.Int{})); diff != "" {
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

	// Check latest confirmed block number recognized by each chainservice
	blockNum := cs.GetLatestConfirmedBlockNum()
	if blockNum != concludeBlockNum {
		t.Fatalf("cs.GetLatestConfirmedBlockNum does not match expected: got %v wanted %v", blockNum, concludeBlockNum)
	}
	blockNum2 := cs2.GetLatestConfirmedBlockNum()
	if blockNum2 != concludeBlockNum {
		t.Fatalf("cs2.GetLatestConfirmedBlockNum does not match expected: got %v wanted %v", blockNum2, concludeBlockNum)
	}

	// Check events from cs2 to ensure they match the expected values
	receivedEvent = <-cs2.EventFeed()
	crEvent = receivedEvent.(ChallengeRegisteredEvent)
	expectedChallengeRegisteredEvent = NewChallengeRegisteredEvent(concludeState.ChannelId(), challengeBlockNum, crEvent.candidate, crEvent.candidateSignatures)
	if diff := cmp.Diff(expectedChallengeRegisteredEvent, crEvent, cmp.AllowUnexported(ChallengeRegisteredEvent{}, commonEvent{}, big.Int{})); diff != "" {
		t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
	}

	for i := 0; i < 2; i++ {
		receivedEvent = <-cs2.EventFeed()
		_, ok := receivedEvent.(DepositedEvent)
		if !ok {
			t.Fatalf("Expected chain event to be DepositedEvent")
		}
	}

	receivedEvent = <-cs2.EventFeed()
	if diff := cmp.Diff(expectedConcludedEvent, receivedEvent, cmp.AllowUnexported(ConcludedEvent{}, commonEvent{})); diff != "" {
		t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
	}

	receivedEvent = <-cs2.EventFeed()
	if diff := cmp.Diff(expectedAllocationUpdatedEvent, receivedEvent, cmp.AllowUnexported(AllocationUpdatedEvent{}, commonEvent{}, big.Int{})); diff != "" {
		t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
	}
}

func closeChainService(t *testing.T, cs ChainService) {
	if err := cs.Close(); err != nil {
		t.Fatal(err)
	}
}

func closeSimulatedChain(t *testing.T, chain SimulatedChain) {
	if err := chain.Close(); err != nil {
		t.Fatal(err)
	}
}
