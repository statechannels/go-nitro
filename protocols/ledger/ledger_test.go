package ledger

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type actor struct {
	address     types.Address
	destination types.Destination
	privateKey  []byte
}

var alice = actor{
	address:     common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`),
	destination: types.AddressToDestination(common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`)),
	privateKey:  common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`),
}

var bob = actor{
	address:     common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`),
	destination: types.AddressToDestination(common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`)),
	privateKey:  common.Hex2Bytes(`62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2`),
}

func TestCreateLedger(t *testing.T) {

	left := outcome.Allocation{Destination: alice.destination, Amount: big.NewInt(3)}
	right := outcome.Allocation{Destination: bob.destination, Amount: big.NewInt(2)}

	ledger, err := CreateTestLedger(left, right, &alice.privateKey, 0, big.NewInt(0))
	if err != nil {
		t.Error(err)
	}

	if ledger.ChannelNonce.Cmp(big.NewInt(0)) != 0 {
		t.Error("TestCreateLedger: initial ledger channel should use the 0 nonce")
	}

}

func TestHandleLedgerRequest(t *testing.T) {
	ledgerManager := NewLedgerManager()
	left := outcome.Allocation{Destination: alice.destination, Amount: big.NewInt(3)}
	right := outcome.Allocation{Destination: bob.destination, Amount: big.NewInt(2)}

	ledger, _ := CreateTestLedger(left, right, &alice.privateKey, 0, big.NewInt(0))

	destination := types.AddressToDestination(common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`))
	asset := types.Address{}
	oId := protocols.ObjectiveId("Test")

	validRequest := protocols.LedgerRequest{
		ObjectiveId: oId,
		LedgerId:    ledger.Id,
		Left:        left.Destination,
		Right:       right.Destination,
		Destination: destination,
		Amount:      types.Funds{asset: big.NewInt(4)},
	}
	invalidRequest := protocols.LedgerRequest{
		ObjectiveId: oId,
		LedgerId:    ledger.Id,
		Left:        left.Destination,
		Right:       right.Destination,
		Destination: destination,
		Amount:      types.Funds{asset: big.NewInt(10)},
	}

	_, err := ledgerManager.HandleRequest(ledger, validRequest, &alice.privateKey)
	if err == nil {
		t.Errorf("TestHandleLedgerRequest: expected request to fail as there is no supported state")
	}

	SignPreAndPostFundingStates(ledger, []*[]byte{&alice.privateKey, &bob.privateKey})

	_, err = ledgerManager.HandleRequest(ledger, invalidRequest, &alice.privateKey)
	if err == nil {
		t.Errorf("TestHandleLedgerRequest: expected request to fail as the ledger does not have enough funds")
	}

	sideEffects, err := ledgerManager.HandleRequest(ledger, validRequest, &alice.privateKey)
	if err != nil {
		t.Error(err)
	}
	guarantee, _ := outcome.GuaranteeMetadata{
		Left:  alice.destination,
		Right: bob.destination,
	}.Encode()

	expectedState := state.State{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{alice.address, bob.address},
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
		AppData:           []byte{},
		Outcome: outcome.Exit{outcome.SingleAssetExit{
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: alice.destination,
					Amount:      big.NewInt(1),
				},
				outcome.Allocation{
					Destination: bob.destination,
					Amount:      big.NewInt(0),
				},
				outcome.Allocation{
					Destination:    destination,
					Amount:         big.NewInt(4),
					AllocationType: outcome.GuaranteeAllocationType,
					Metadata:       guarantee,
				},
			}}},
		TurnNum: 2,
		IsFinal: false}

	expectedSigned := state.NewSignedState(expectedState)
	err = expectedSigned.Sign(&alice.privateKey)
	if err != nil {
		t.Error(err)
	}

	expectedMessage := protocols.Message{To: bob.address, ObjectiveId: oId, SignedStates: []state.SignedState{expectedSigned}}

	if diff := cmp.Diff(sideEffects.MessagesToSend[0], expectedMessage); diff != "" {
		t.Errorf("TestHandleRequest: ledger message mismatch (-want +got):\n%s", diff)
	}
}
