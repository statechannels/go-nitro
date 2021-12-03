package virtualfund

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

/////////////////////
// BEGIN test data //
/////////////////////

// General
var Alice = struct {
	address     types.Address
	destination types.Destination
	privateKey  []byte
}{
	address:     common.HexToAddress(`0xD9995BAE12FEe327256FFec1e3184d492bD94C31`),
	destination: types.AdddressToDestination(common.HexToAddress(`0xD9995BAE12FEe327256FFec1e3184d492bD94C31`)),
	privateKey:  common.Hex2Bytes(`7ab741b57e8d94dd7e1a29055646bafde7010f38a900f55bbd7647880faa6ee8`),
}
var my = Alice

var P_0 = struct {
	address     types.Address
	destination types.Destination
	privateKey  []byte
}{
	address:     common.HexToAddress(`0xd4Fa489Eacc52BA59438993f37Be9fcC20090E39`),
	destination: types.AdddressToDestination(common.HexToAddress(`0xd4Fa489Eacc52BA59438993f37Be9fcC20090E39`)),
	privateKey:  common.Hex2Bytes(`2030b463177db2da82908ef90fa55ddfcef56e8183caf60db464bc398e736e6f`),
}

var Bob = struct {
	address     types.Address
	destination types.Destination
	privateKey  []byte
}{
	address:     common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`),
	destination: types.AdddressToDestination(common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`)),
	privateKey:  common.Hex2Bytes(`62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2`),
}

// Ledger Channel(s)

// Alice = P_0 <=L_0=> P_1 <=L_1=>...P_n <=L_n>= P_n+1 = Bob
// ^^^^^

// Alice plays role 0 so has no ledger channel on her left
var ledgerChannelToMyLeft channel.Channel

// She has a single ledger channel L_0 connecting her to P_1
var L_0state = state.State{
	ChainId:           big.NewInt(9001),
	Participants:      []types.Address{my.address, P_0.address},
	ChannelNonce:      big.NewInt(0),
	AppDefinition:     types.Address{},
	ChallengeDuration: big.NewInt(45),
	AppData:           []byte{},
	Outcome: outcome.Exit{outcome.SingleAssetExit{
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: my.destination,
				Amount:      big.NewInt(5),
			},
			outcome.Allocation{
				Destination: P_0.destination,
				Amount:      big.NewInt(5),
			},
		},
	}},
	TurnNum: big.NewInt(1),
	IsFinal: false,
}

var ledgerChannelToMyRight channel.Channel = channel.New(
	L_0state,
	true,
	my.destination,
	P_0.destination,
)

// Virtual Channel
var VState = state.State{
	ChainId:           big.NewInt(9001),
	Participants:      []types.Address{my.address, P_0.address, Bob.address}, // A single hop virtual channel
	ChannelNonce:      big.NewInt(0),
	AppDefinition:     types.Address{},
	ChallengeDuration: big.NewInt(45),
	AppData:           []byte{},
	Outcome: outcome.Exit{outcome.SingleAssetExit{
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: my.destination,
				Amount:      big.NewInt(5),
			},
			outcome.Allocation{
				Destination: Bob.destination,
				Amount:      big.NewInt(5),
			},
		},
	}},
	TurnNum: big.NewInt(0),
	IsFinal: false,
}

// Objective
var myRole = uint(0) // In this test, we play Alice
var s, _ = New(VState, my.address, myRole, ledgerChannelToMyLeft, ledgerChannelToMyRight)
var expectedEncodedGuaranteeMetadata, _ = outcome.GuaranteeMetadata{Left: ledgerChannelToMyRight.MyDestination, Right: ledgerChannelToMyRight.TheirDestination}.Encode()
var expectedGuarantee outcome.Allocation = outcome.Allocation{
	Destination:    s.V.Id,
	Amount:         big.NewInt(0).Set(VState.VariablePart().Outcome[0].TotalAllocated()),
	AllocationType: outcome.GuaranteeAllocationType,
	Metadata:       expectedEncodedGuaranteeMetadata,
}

///////////////////
// END test data //
///////////////////

// TestNew tests the constructor using a TestState fixture
func TestNew(t *testing.T) {
	// Assert that a valid set of constructor args does not result in an error
	o, err := New(VState, my.address, myRole, ledgerChannelToMyLeft, ledgerChannelToMyRight)
	if err != nil {
		t.Error(err)
	}

	got := o.ExpectedGuarantees[0][types.Address{}] // VState only has one (native) asset represented by the zero address
	want := expectedGuarantee

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestNew: expectedGuarantee mismatch (-want +got):\n%s", diff)
	}
}

func TestCrank(t *testing.T) {
	// Assert that cranking an unapproved objective returns an error
	if _, _, _, err := s.Crank(&my.privateKey); err == nil {
		t.Error(`Expected error when cranking unapproved objective, but got nil`)
	}

	// Approve the objective, so that the rest of the test cases can run.
	o := s.Approve()

	// To test the finite state progression, we are going to progressively mutate o
	// And then crank it to see which "pause point" (WaitingFor) we end up at.

	// Initial Crank
	_, _, waitingFor, err := o.Crank(&my.privateKey)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForCompletePrefund {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompletePrefund, waitingFor)
	}

	// Manually progress the extended state by collecting prefund signatures
	o.(VirtualFundObjective).preFundSigned[0] = true
	o.(VirtualFundObjective).preFundSigned[1] = true
	o.(VirtualFundObjective).preFundSigned[2] = true

	// Cranking should move us to the next waiting point, generate ledger requests as a side effect, and alter the extended state to reflect that
	o, _, waitingFor, err = o.Crank(&my.privateKey)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForCompleteFunding {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompleteFunding, waitingFor)
	}
	if o.(VirtualFundObjective).requestedLedgerUpdates != true { // TODO && sideeffects are as expected
		t.Error(`Expected ledger updates to be requested, but they weren't`)
	}

	// Manually progress the extended state by "completing funding" from this wallet's point of view
	var UpdatedL0Outcome = outcome.Exit{
		outcome.SingleAssetExit{ // TODO this is not realistic as it does not contain allocations for either Alice (P_0) nor P_1
			Asset: types.Address{},
			Allocations: outcome.Allocations{
				expectedGuarantee,
			},
		},
	}
	var UpdatedL0State = o.(VirtualFundObjective).L[0].LatestSupportedState
	UpdatedL0State.Outcome = UpdatedL0Outcome
	var UpdatedL0Channel = o.(VirtualFundObjective).L[0]
	UpdatedL0Channel.LatestSupportedState = UpdatedL0State
	o.(VirtualFundObjective).L[0] = UpdatedL0Channel
	o.(VirtualFundObjective).L[0].OnChainFunding[types.Address{}] = UpdatedL0Outcome[0].Allocations.Total() // Make this channel fully funded

	// Cranking now should not generate side effects, because we already did that
	o, _, waitingFor, err = o.Crank(&my.privateKey)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != protocols.WaitingForCompletePostFund {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompletePostFund, waitingFor)
	}

	// Manually progress the extended state by collecting postfund signatures
	o.(VirtualFundObjective).postFundSigned[0] = true
	o.(VirtualFundObjective).postFundSigned[1] = true
	o.(VirtualFundObjective).postFundSigned[2] = true

	// This should be the final crank...
	_, _, waitingFor, err = o.Crank(&my.privateKey)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForNothing {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForNothing, waitingFor)
	}

}
