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
	"github.com/statechannels/go-nitro/protocols"
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
		Actors.Alice.Address,
		Actors.Bob.Address,
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
	aSig, _ := concludeState.Sign(Actors.Alice.PrivateKey)
	bSig, _ := concludeState.Sign(Actors.Bob.PrivateKey)

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
