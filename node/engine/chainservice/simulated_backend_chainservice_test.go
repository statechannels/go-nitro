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
	NitroAdjudicator "github.com/statechannels/go-nitro/node/engine/chainservice/adjudicator"
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

func TestSimulatedBackendChainService(t *testing.T) {
	one := big.NewInt(1)
	three := big.NewInt(3)

	sim, bindings, ethAccounts, err := SetupSimulatedBackend(1)
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
		ChallengeDuration: 1000,
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

	receivedEvent := <-out
	crEvent := receivedEvent.(ChallengeRegisteredEvent)
	expectedChallengeRegisteredEvent := NewChallengeRegisteredEvent(concludeState.ChannelId(), 2, crEvent.candidate, crEvent.candidateSignatures)
	if diff := cmp.Diff(expectedChallengeRegisteredEvent, crEvent, cmp.AllowUnexported(ChallengeRegisteredEvent{}, commonEvent{}, big.Int{})); diff != "" {
		t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
	}

	testDeposit := types.Funds{
		common.HexToAddress("0x00"): three,
		bindings.Token.Address:      one,
	}
	testTx := protocols.NewDepositTransaction(concludeState.ChannelId(), testDeposit)

	// Submit transaction
	err = cs.SendTransaction(testTx)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the received events matches the expected event
	for i := 0; i < 2; i++ {
		receivedEvent := <-out
		dEvent := receivedEvent.(DepositedEvent)
		expectedDepositEvent := NewDepositedEvent(concludeState.ChannelId(), 5, dEvent.Asset, testDeposit[dEvent.Asset])
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
	concludeTx := protocols.NewWithdrawAllTransaction(cId, signedConcludeState)
	err = cs.SendTransaction(concludeTx)
	if err != nil {
		t.Fatal(err)
	}
	// Check that the recieved event matches the expected event
	concludedEvent := <-out
	expectedConcludeEvent := ConcludedEvent{commonEvent: commonEvent{channelID: cId, blockNum: 8}}
	if diff := cmp.Diff(expectedConcludeEvent, concludedEvent, cmp.AllowUnexported(ConcludedEvent{}, commonEvent{})); diff != "" {
		t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
	}

	// Check that the recieved event matches the expected event
	allocationUpdatedEvent := <-out
	expectedAllocationUpdatedEvent := NewAllocationUpdatedEvent(cId, 8, common.Address{}, new(big.Int).SetInt64(1))

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
