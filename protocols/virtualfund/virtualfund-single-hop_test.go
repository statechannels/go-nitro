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

func TestSingleHopVirtualFund(t *testing.T) {

	// In general
	// Alice = P_0 <=L_0=> P_1 <=L_1=> ... P_n <=L_n>= P_n+1 = Bob

	// For these tests
	// Alice <=L_0=> P_1 <=L_1=> Bob

	////////////
	// ACTORS //
	////////////
	type actor struct {
		address     types.Address
		destination types.Destination
		privateKey  []byte
		role        uint
	}

	var alice = actor{
		address:     common.HexToAddress(`0xD9995BAE12FEe327256FFec1e3184d492bD94C31`),
		destination: types.AddressToDestination(common.HexToAddress(`0xD9995BAE12FEe327256FFec1e3184d492bD94C31`)),
		privateKey:  common.Hex2Bytes(`7ab741b57e8d94dd7e1a29055646bafde7010f38a900f55bbd7647880faa6ee8`),
		role:        0,
	}

	var p1 = actor{ // Aliases: The Hub, Irene
		address:     common.HexToAddress(`0xd4Fa489Eacc52BA59438993f37Be9fcC20090E39`),
		destination: types.AddressToDestination(common.HexToAddress(`0xd4Fa489Eacc52BA59438993f37Be9fcC20090E39`)),
		privateKey:  common.Hex2Bytes(`2030b463177db2da82908ef90fa55ddfcef56e8183caf60db464bc398e736e6f`),
		role:        1,
	}

	var bob = actor{
		address:     common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`),
		destination: types.AddressToDestination(common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`)),
		privateKey:  common.Hex2Bytes(`62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2`),
		role:        2,
	}

	/////////////////////
	// VIRTUAL CHANNEL //
	/////////////////////

	// Virtual Channel
	var vPreFund = state.State{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{alice.address, p1.address, bob.address}, // A single hop virtual channel
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
		AppData:           []byte{},
		Outcome: outcome.Exit{outcome.SingleAssetExit{
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: alice.destination,
					Amount:      big.NewInt(5),
				},
				outcome.Allocation{
					Destination: bob.destination,
					Amount:      big.NewInt(5),
				},
			},
		}},
		TurnNum: 0,
		IsFinal: false,
	}

	/////////////////////
	// LEDGER CHANNELS //
	/////////////////////

	var l0state = state.State{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{alice.address, p1.address},
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
		AppData:           []byte{},
		Outcome: outcome.Exit{outcome.SingleAssetExit{
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: alice.destination,
					Amount:      big.NewInt(5),
				},
				outcome.Allocation{
					Destination: p1.destination,
					Amount:      big.NewInt(5),
				},
			},
		}},
		TurnNum: 0, // We use turnNum 0 so that we can use github.com/statechannels/go-nitro/channel.New().
		// It would be more realistic to have a higher TurnNum, but that would involve more boilerplate code.
		IsFinal: false,
	}

	var vId, _ = vPreFund.ChannelId()

	var l0guaranteemetadataemcoded, _ = outcome.GuaranteeMetadata{
		Left:  alice.destination,
		Right: p1.destination,
	}.Encode()

	var l0updatedstate = state.State{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{alice.address, p1.address},
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
		AppData:           []byte{},
		Outcome: outcome.Exit{outcome.SingleAssetExit{
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: alice.destination,
					Amount:      big.NewInt(0),
				},
				outcome.Allocation{
					Destination: p1.destination,
					Amount:      big.NewInt(0),
				},
				outcome.Allocation{
					Destination:    vId,
					Amount:         big.NewInt(10),
					AllocationType: outcome.GuaranteeAllocationType,
					Metadata:       l0guaranteemetadataemcoded,
				},
			},
		}},
		TurnNum: 2, // This needs to be greater than the previous state else it will be rejected by Channel.AddSignedState
		IsFinal: false,
	}

	var AsAlice = func(t *testing.T) {

		/////////////////////
		// BEGIN test data //
		/////////////////////

		// In this test, we play Alice
		my := alice

		// Alice plays role 0 so has no ledger channel on her left
		var ledgerChannelToMyLeft channel.TwoPartyLedger

		// She has a single ledger channel L_0 connecting her to P_1
		var ledgerChannelToMyRight, _ = channel.NewTwoPartyLedger(
			l0state,
			0,
		)

		// Ensure this channel is fully funded on chain
		ledgerChannelToMyRight.OnChainFunding = ledgerChannelToMyRight.PreFundState().Outcome.TotalAllocated()

		// Objective
		var n = uint(2) // number of ledger channels (num_hops + 1)
		var s, _ = New(false, vPreFund, my.address, n, my.role, &ledgerChannelToMyLeft, ledgerChannelToMyRight)
		var expectedGuaranteeMetadata = outcome.GuaranteeMetadata{Left: ledgerChannelToMyRight.MyDestination(), Right: ledgerChannelToMyRight.TheirDestination()}
		var expectedEncodedGuaranteeMetadata, _ = expectedGuaranteeMetadata.Encode()
		var expectedGuarantee outcome.Allocation = outcome.Allocation{
			Destination:    s.V.Id,
			Amount:         big.NewInt(0).Set(vPreFund.VariablePart().Outcome[0].TotalAllocated()),
			AllocationType: outcome.GuaranteeAllocationType,
			Metadata:       expectedEncodedGuaranteeMetadata,
		}
		var expectedLedgerRequests = []protocols.LedgerRequest{{
			LedgerId:    ledgerChannelToMyRight.Id,
			Destination: s.V.Id,
			Amount:      types.Funds{types.Address{}: s.V.PreFundState().VariablePart().Outcome[0].Allocations.Total()},
			Left:        ledgerChannelToMyRight.MyDestination(), Right: ledgerChannelToMyRight.TheirDestination(),
		}}

		var correctSignatureByAliceOnVPreFund, _ = s.V.PreFundState().Sign(alice.privateKey)
		var correctSignatureByP_1OnVPreFund, _ = s.V.PreFundState().Sign(p1.privateKey)
		var correctSignatureByBobOnVPreFund, _ = s.V.PreFundState().Sign(bob.privateKey)

		var correctSignatureByAliceOnVPostFund, _ = s.V.PostFundState().Sign(alice.privateKey)
		var correctSignatureByP_1OnVPostFund, _ = s.V.PostFundState().Sign(p1.privateKey)
		var correctSignatureByBobOnVPostFund, _ = s.V.PostFundState().Sign(bob.privateKey)

		var correctSignatureByAliceOnL_0updatedsate, _ = l0updatedstate.Sign(alice.privateKey)
		var correctSignatureByP_1OnL_0updatedsate, _ = l0updatedstate.Sign(p1.privateKey)

		///////////////////
		// END test data //
		///////////////////

		testNew := func(t *testing.T) {
			// Assert that a valid set of constructor args does not result in an error
			o, err := New(false, vPreFund, my.address, 2, my.role, &ledgerChannelToMyLeft, ledgerChannelToMyRight)
			if err != nil {
				t.Error(err)
			}

			got := o.ToMyRight.ExpectedGuarantees[types.Address{}] // VState only has one (native) asset represented by the zero address
			want := expectedGuarantee

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("TestNew: expectedGuarantee mismatch (-want +got):\n%s", diff)
			}
		}

		testCrank := func(t *testing.T) {

			// Assert that cranking an unapproved objective returns an error
			if _, _, _, err := s.Crank(&my.privateKey); err == nil {
				t.Error(`Expected error when cranking unapproved objective, but got nil`)
			}

			// Approve the objective, so that the rest of the test cases can run.
			o := s.Approve().(VirtualFundObjective)
			// To test the finite state progression, we are going to progressively mutate o
			// And then crank it to see which "pause point" (WaitingFor) we end up at.

			// Initial Crank
			_, got, waitingFor, err := o.Crank(&my.privateKey)
			if err != nil {
				t.Error(err)
			}
			if waitingFor != WaitingForCompletePrefund {
				t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompletePrefund, waitingFor)
			}

			expectedSignedState := state.NewSignedState(o.V.PreFundState())
			_ = expectedSignedState.AddSignature(correctSignatureByAliceOnVPreFund)

			forBob := protocols.Message{To: bob.address, ObjectiveId: o.Id(), SignedStates: []state.SignedState{expectedSignedState}}
			forIrene := protocols.Message{To: p1.address, ObjectiveId: o.Id(), SignedStates: []state.SignedState{expectedSignedState}}
			want := protocols.SideEffects{MessagesToSend: []protocols.Message{forIrene, forBob}}

			if diff := cmp.Diff(want, got); diff != "" {
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("TestCrank: side effects mismatch (-want +got):\n%s", diff)
				}

			}
			// Manually progress the extended state by collecting prefund signatures

			o.V.AddStateWithSignature(vPreFund, correctSignatureByBobOnVPreFund)
			o.V.AddStateWithSignature(vPreFund, correctSignatureByP_1OnVPreFund)

			// Cranking should move us to the next waiting point, generate ledger requests as a side effect, and alter the extended state to reflect that
			oObj, got, waitingFor, err := o.Crank(&my.privateKey)
			o = oObj.(VirtualFundObjective)
			if err != nil {
				t.Error(err)
			}
			if waitingFor != WaitingForCompleteFunding {
				t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompleteFunding, waitingFor)
			}
			if o.requestedLedgerUpdates != true {
				t.Error(`Expected ledger update idempotency flag to be raised, but it wasn't`)
			}

			want = protocols.SideEffects{LedgerRequests: expectedLedgerRequests}

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("TestCrank: side effects mismatch (-want +got):\n%s", diff)
			}

			// Manually progress the extended state by "completing funding" from this wallet's point of view
			o.ToMyRight.Channel.AddStateWithSignature(l0updatedstate, correctSignatureByAliceOnL_0updatedsate)
			o.ToMyRight.Channel.AddStateWithSignature(l0updatedstate, correctSignatureByP_1OnL_0updatedsate)
			o.ToMyRight.Channel.OnChainFunding[types.Address{}] = l0state.Outcome[0].Allocations.Total() // Make this channel fully funded

			// Cranking now should generate signed post fund messages
			oObj, got, waitingFor, err = o.Crank(&my.privateKey)
			o = oObj.(VirtualFundObjective)
			if err != nil {
				t.Error(err)
			}
			if waitingFor != WaitingForCompletePostFund {
				t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompletePostFund, waitingFor)
			}

			expectedSignedState = state.NewSignedState(o.V.PostFundState())
			_ = expectedSignedState.AddSignature(correctSignatureByAliceOnVPostFund)

			forBob = protocols.Message{To: bob.address, ObjectiveId: o.Id(), SignedStates: []state.SignedState{expectedSignedState}}
			forIrene = protocols.Message{To: p1.address, ObjectiveId: o.Id(), SignedStates: []state.SignedState{expectedSignedState}}
			want = protocols.SideEffects{MessagesToSend: []protocols.Message{forIrene, forBob}}

			if diff := cmp.Diff(want, got); diff != "" {
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("TestCrank: side effects mismatch (-want +got):\n%s", diff)
				}

			}
			// Manually progress the extended state by collecting postfund signatures

			o.V.AddStateWithSignature(o.V.PostFundState(), correctSignatureByBobOnVPostFund)
			o.V.AddStateWithSignature(o.V.PostFundState(), correctSignatureByP_1OnVPostFund)

			// This should be the final crank...
			_, _, waitingFor, err = o.Crank(&my.privateKey)
			if err != nil {
				t.Error(err)
			}
			if waitingFor != WaitingForNothing {
				t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForNothing, waitingFor)
			}

		}

		testUpdate := func(t *testing.T) {

			// Prepare an event with a mismatched objectiveId
			e := protocols.ObjectiveEvent{
				ObjectiveId: "some-other-id",
			}
			// Assert that Updating the objective with such an event returns an error
			// TODO is this the behaviour we want? Below with the signatures, we prefer a log + NOOP (no error)
			if _, err := s.Update(e); err == nil {
				t.Error(`Objective ID mismatch -- expected an error but did not get one`)
			}

			// Now modify the event to give it the "correct" channelId (matching the objective),
			// and make a new Sigs map.
			// This prepares us for the rest of the test. We will reuse the same event multiple times
			e.ObjectiveId = s.Id()
			e.SignedStates = make([]state.SignedState, 0)

			// Next, attempt to update the objective with correct signature by a participant on a relevant state
			// Assert that this results in an appropriate change in the extended state of the objective
			// Part 1: a signature on a state in channel V
			prefundsignedstate := state.NewSignedState(s.V.PreFundState())
			err := prefundsignedstate.AddSignature(correctSignatureByBobOnVPreFund)
			if err != nil {
				t.Error(err)
			}
			e.SignedStates = append(e.SignedStates, prefundsignedstate)

			updatedObj, err := s.Update(e)
			updated := updatedObj.(VirtualFundObjective)
			if err != nil {
				t.Error(err)
			}
			if updated.V.SignedStateForTurnNum[0].HasSignatureForParticipant(bob.role) != true {
				t.Error(`Objective data not updated as expected`)
			}

			// Part 2: a signature on Alice's ledger channel (on her right)
			f := protocols.ObjectiveEvent{
				ObjectiveId: s.Id(),
			}
			f.SignedStates = make([]state.SignedState, 0)
			ss := state.NewSignedState(l0updatedstate)
			err = ss.AddSignature(correctSignatureByAliceOnL_0updatedsate)
			if err != nil {
				t.Error(err)
			}
			f.SignedStates = append(f.SignedStates, ss)

			updatedObj, err = s.Update(f)
			updated = updatedObj.(VirtualFundObjective)
			if err != nil {
				t.Error(err)
			}
			if !updated.ToMyRight.ledgerChannelAffordsExpectedGuarantees() != true {
				t.Error(`Objective data not updated as expected`)
			}

		}

		t.Run(`New`, testNew)
		t.Run(`Update`, testUpdate)
		t.Run(`Crank`, testCrank)

	}
	t.Run(`AsAlice`, AsAlice)
}
