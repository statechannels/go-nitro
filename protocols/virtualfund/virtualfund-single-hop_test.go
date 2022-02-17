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
	"github.com/statechannels/go-nitro/protocols/ledger"
	"github.com/statechannels/go-nitro/types"
)

// TODO move these package-level symbols inside the scope of the test
type actor struct {
	address     types.Address
	destination types.Destination
	privateKey  []byte
	role        uint
}

func signState(s state.State, a actor) state.SignedState {
	ss := state.NewSignedState(s)
	sig, err := s.Sign(a.privateKey)
	if err != nil {
		panic(err)
	}
	err = ss.AddSignature(sig)
	if err != nil {
		panic(err)
	}
	return ss
}
func TestSingleHopVirtualFund(t *testing.T) {
	var n = uint(1) // number of intermediaries

	var n = uint(2) // number of ledger channels (num_hops + 1)

	// In general
	// Alice = P_0 <=L_0=> P_1 <=L_1=> ... P_n <=L_n>= P_n+1 = Bob

	// For these tests
	// Alice <=L_0=> P_1 <=L_1=> Bob

	////////////
	// ACTORS //
	////////////

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

	// }
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

	var AsAlice = func(t *testing.T) {

		/////////////////////
		// BEGIN test data //
		/////////////////////

		// In this test, we play Alice
		my := alice

		// Alice plays role 0 so has no ledger channel on her left
		var ledgerChannelToMyLeft *channel.TwoPartyLedger

		ledgerManager := ledger.NewLedgerManager()
		left := outcome.Allocation{Destination: alice.destination, Amount: big.NewInt(5)}
		right := outcome.Allocation{Destination: p1.destination, Amount: big.NewInt(5)}
		myIndex := uint(0) // because Alice is in the "left" slot

		// She has a single ledger channel L_0 connecting her to P_1
		var ledgerChannelToMyRight, _ = ledger.CreateTestLedger(left, right, &alice.privateKey, myIndex, big.NewInt(0))
		// Ensure this channel is fully funded on chain
		ledgerChannelToMyRight.OnChainFunding = ledgerChannelToMyRight.PreFundState().Outcome.TotalAllocated()

		///////////////////
		// END test data //
		///////////////////

		testNew := func(t *testing.T) {

			// Assert that a valid set of constructor args does not result in an error
			o, err := New(false, vPreFund, my.address, n, my.role, ledgerChannelToMyLeft, ledgerChannelToMyRight)
			if err != nil {
				t.Error(err)
			}

			got := o.ToMyRight.ExpectedGuarantees[types.Address{}] // VState only has one (native) asset represented by the zero address
			var expectedGuaranteeMetadata = outcome.GuaranteeMetadata{Left: ledgerChannelToMyRight.MyDestination(), Right: ledgerChannelToMyRight.TheirDestination()}
			var expectedEncodedGuaranteeMetadata, _ = expectedGuaranteeMetadata.Encode()
			var expectedGuarantee outcome.Allocation = outcome.Allocation{
				Destination:    o.V.Id,
				Amount:         big.NewInt(0).Set(vPreFund.VariablePart().Outcome[0].TotalAllocated()),
				AllocationType: outcome.GuaranteeAllocationType,
				Metadata:       expectedEncodedGuaranteeMetadata,
			}
			want := expectedGuarantee

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("TestNew: expectedGuarantee mismatch (-want +got):\n%s", diff)
			}
		}

		testCrank := func(t *testing.T) {
			var s, _ = New(false, vPreFund, my.address, n, my.role, ledgerChannelToMyLeft, ledgerChannelToMyRight)
			// Assert that cranking an unapproved objective returns an error
			if _, _, _, err := s.Crank(&my.privateKey); err == nil {
				t.Error(`Expected error when cranking unapproved objective, but got nil`)
			}

			// Approve the objective, so that the rest of the test cases can run.
			o := s.Approve().(VirtualFundObjective)
			// To test the finite state progression, we are going to progressively mutate o
			// And then crank it to see which "pause point" (WaitingFor) we end up at.

			// Initial Crank
			oObj, got, waitingFor, err := o.Crank(&my.privateKey)
			o = oObj.(VirtualFundObjective)
			if err != nil {
				t.Error(err)
			}
			if waitingFor != WaitingForCompletePrefund {
				t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompletePrefund, waitingFor)
			}

			expectedSignedState := state.NewSignedState(o.V.PreFundState())
			aliceSig, _ := o.V.PreFundState().Sign(my.privateKey)
			_ = expectedSignedState.AddSignature(aliceSig)

			forBob := protocols.Message{To: bob.address, ObjectiveId: o.Id(), SignedStates: []state.SignedState{expectedSignedState}}
			forIrene := protocols.Message{To: p1.address, ObjectiveId: o.Id(), SignedStates: []state.SignedState{expectedSignedState}}
			want := protocols.SideEffects{MessagesToSend: []protocols.Message{forIrene, forBob}}

			if diff := cmp.Diff(want, got); diff != "" {
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("TestCrank: side effects mismatch (-want +got):\n%s", diff)
				}

			}
			// Manually progress the extended state by collecting prefund signatures
			bobSig, _ := vPreFund.Sign(bob.privateKey)
			p1Sig, _ := vPreFund.Sign(p1.privateKey)
			o.V.AddStateWithSignature(vPreFund, bobSig)
			o.V.AddStateWithSignature(vPreFund, p1Sig)

			// Cranking should move us to the next waiting point, generate ledger requests as a side effect, and alter the extended state to reflect that
			oObj, got, waitingFor, err = o.Crank(&my.privateKey)
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
			var expectedLedgerRequests = []protocols.LedgerRequest{{
				ObjectiveId: o.Id(),
				LedgerId:    ledgerChannelToMyRight.Id,
				Destination: s.V.Id,

				Left: ledgerChannelToMyRight.MyDestination(), Right: ledgerChannelToMyRight.TheirDestination(),
				LeftAmount:  types.Funds{types.Address{}: big.NewInt(5)},
				RightAmount: types.Funds{types.Address{}: big.NewInt(5)},
			}}
			want = protocols.SideEffects{LedgerRequests: expectedLedgerRequests}

			if diff := cmp.Diff(want, got, cmp.Comparer(types.Equal)); diff != "" {
				t.Errorf("TestCrank: side effects mismatch (-want +got):\n%s", diff)
			}

			ledger.SignPreAndPostFundingStates(o.ToMyRight.Channel, []*[]byte{&alice.privateKey, &p1.privateKey})

			_, _ = ledgerManager.HandleRequest(o.ToMyRight.Channel, got.LedgerRequests[0], &alice.privateKey)
			ledger.SignLatest(o.ToMyRight.Channel, [][]byte{p1.privateKey})
			// Cranking now should not generate side effects, because we already did that
			oObj, got, waitingFor, err = o.Crank(&my.privateKey)
			o = oObj.(VirtualFundObjective)
			if err != nil {
				t.Error(err)
			}
			if waitingFor != WaitingForCompletePostFund {
				t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompletePostFund, waitingFor)
			}

			expectedSignedState = signState(o.V.PostFundState(), my)

			forBob = protocols.Message{To: bob.address, ObjectiveId: o.Id(), SignedStates: []state.SignedState{expectedSignedState}}
			forIrene = protocols.Message{To: p1.address, ObjectiveId: o.Id(), SignedStates: []state.SignedState{expectedSignedState}}
			want = protocols.SideEffects{MessagesToSend: []protocols.Message{forIrene, forBob}}

			if diff := cmp.Diff(want, got); diff != "" {
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("TestCrank: side effects mismatch (-want +got):\n%s", diff)
				}

			}

			// Manually progress the extended state by collecting postfund signatures
			bobPost := signState(o.V.PostFundState(), bob)
			p1Post := signState(o.V.PostFundState(), p1)
			o.V.AddSignedStates([]state.SignedState{bobPost, p1Post})
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
			var s, _ = New(false, vPreFund, my.address, n, my.role, ledgerChannelToMyLeft, ledgerChannelToMyRight)
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
			prefundsignedstate := signState(s.V.PreFundState(), bob)
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

			ledger, _ := ledger.CreateTestLedger(left, right, &alice.privateKey, 0, big.NewInt(0))
			ss := signState(ledger.PreFundState(), alice)

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

	var AsBob = func(t *testing.T) {

		/////////////////////
		// BEGIN test data //
		/////////////////////

		// In this test, we play Bob
		my := bob

		// Bob plays role 2 so has no ledger channel on his right
		var ledgerChannelToMyRight *channel.TwoPartyLedger

		ledgerManager := ledger.NewLedgerManager()
		left := outcome.Allocation{Destination: p1.destination, Amount: big.NewInt(5)}
		right := outcome.Allocation{Destination: bob.destination, Amount: big.NewInt(5)}
		myIndex := uint(1) // because Bob is in the "right" slot

		// He has a single ledger channel L_1 connecting him to P_1
		var ledgerChannelToMyLeft, _ = ledger.CreateTestLedger(left, right, &bob.privateKey, myIndex, big.NewInt(0))
		// Ensure this channel is fully funded on chain
		ledgerChannelToMyLeft.OnChainFunding = ledgerChannelToMyLeft.PreFundState().Outcome.TotalAllocated()

		///////////////////
		// END test data //
		///////////////////

		testNew := func(t *testing.T) {

			// Assert that a valid set of constructor args does not result in an error
			o, err := New(false, vPreFund, my.address, n, my.role, ledgerChannelToMyLeft, ledgerChannelToMyRight)
			if err != nil {
				t.Error(err)
			}

			got := o.ToMyLeft.ExpectedGuarantees[types.Address{}] // VState only has one (native) asset represented by the zero address
			var expectedGuaranteeMetadata = outcome.GuaranteeMetadata{Left: ledgerChannelToMyLeft.TheirDestination(), Right: ledgerChannelToMyLeft.MyDestination()}
			var expectedEncodedGuaranteeMetadata, _ = expectedGuaranteeMetadata.Encode()
			var expectedGuarantee outcome.Allocation = outcome.Allocation{
				Destination:    o.V.Id,
				Amount:         big.NewInt(0).Set(vPreFund.VariablePart().Outcome[0].TotalAllocated()),
				AllocationType: outcome.GuaranteeAllocationType,
				Metadata:       expectedEncodedGuaranteeMetadata,
			}
			want := expectedGuarantee

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("TestNew: expectedGuarantee mismatch (-want +got):\n%s", diff)
			}
		}

		testCrank := func(t *testing.T) {
			var s, _ = New(false, vPreFund, my.address, n, my.role, ledgerChannelToMyLeft, ledgerChannelToMyRight)
			// Assert that cranking an unapproved objective returns an error
			if _, _, _, err := s.Crank(&my.privateKey); err == nil {
				t.Error(`Expected error when cranking unapproved objective, but got nil`)
			}

			// Approve the objective, so that the rest of the test cases can run.
			o := s.Approve().(VirtualFundObjective)
			// To test the finite state progression, we are going to progressively mutate o
			// And then crank it to see which "pause point" (WaitingFor) we end up at.

			// Initial Crank
			oObj, got, waitingFor, err := o.Crank(&my.privateKey)
			o = oObj.(VirtualFundObjective)
			if err != nil {
				t.Error(err)
			}
			if waitingFor != WaitingForCompletePrefund {
				t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompletePrefund, waitingFor)
			}

			expectedSignedState := state.NewSignedState(o.V.PreFundState())
			bobSig, _ := o.V.PreFundState().Sign(my.privateKey)
			_ = expectedSignedState.AddSignature(bobSig)

			forAlice := protocols.Message{To: alice.address, ObjectiveId: o.Id(), SignedStates: []state.SignedState{expectedSignedState}}
			forIrene := protocols.Message{To: p1.address, ObjectiveId: o.Id(), SignedStates: []state.SignedState{expectedSignedState}}
			want := protocols.SideEffects{MessagesToSend: []protocols.Message{forAlice, forIrene}}
			// TODO ^^^ The test is currently sensitive to the order of the messages. It should not be.

			if diff := cmp.Diff(want, got); diff != "" {
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("TestCrank: side effects mismatch (-want +got):\n%s", diff)
				}

			}
			// Manually progress the extended state by collecting prefund signatures
			aliceSig, _ := vPreFund.Sign(alice.privateKey)
			o.V.AddStateWithSignature(vPreFund, aliceSig)
			p1Sig, _ := vPreFund.Sign(p1.privateKey)
			o.V.AddStateWithSignature(vPreFund, p1Sig)

			// Cranking should move us to the next waiting point, generate ledger requests as a side effect, and alter the extended state to reflect that
			oObj, got, waitingFor, err = o.Crank(&my.privateKey)
			o = oObj.(VirtualFundObjective)
			if err != nil {
				t.Error(err)
			}
			if waitingFor != WaitingForCompleteFunding {
				t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForCompleteFunding, waitingFor)
			}
			if o.requestedLedgerUpdates != true {
				t.Error(`Expected ledger update idempotency flag to be raised, but it wasn't`)
			}
			var expectedLedgerRequests = []protocols.LedgerRequest{{
				ObjectiveId: o.Id(),
				LedgerId:    ledgerChannelToMyLeft.Id,
				Destination: s.V.Id,
				Amount:      types.Funds{types.Address{}: s.V.PreFundState().VariablePart().Outcome[0].Allocations.Total()},
				Left:        ledgerChannelToMyLeft.TheirDestination(), Right: ledgerChannelToMyLeft.MyDestination(),
			}}
			want = protocols.SideEffects{LedgerRequests: expectedLedgerRequests}

			if diff := cmp.Diff(want, got, cmp.Comparer(types.Equal)); diff != "" {
				t.Errorf("TestCrank: side effects mismatch (-want +got):\n%s", diff)
			}

			ledger.SignPreAndPostFundingStates(o.ToMyLeft.Channel, []*[]byte{&bob.privateKey, &p1.privateKey})

			_, _ = ledgerManager.HandleRequest(o.ToMyLeft.Channel, got.LedgerRequests[0], &bob.privateKey)
			ledger.SignLatest(o.ToMyLeft.Channel, [][]byte{p1.privateKey})
			// Cranking now should not generate side effects, because we already did that
			oObj, got, waitingFor, err = o.Crank(&my.privateKey)
			o = oObj.(VirtualFundObjective)
			if err != nil {
				t.Error(err)
			}
			if waitingFor != WaitingForCompletePostFund {
				t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompletePostFund, waitingFor)
			}

			expectedSignedState = signState(o.V.PostFundState(), my)

			forAlice = protocols.Message{To: alice.address, ObjectiveId: o.Id(), SignedStates: []state.SignedState{expectedSignedState}}
			forIrene = protocols.Message{To: p1.address, ObjectiveId: o.Id(), SignedStates: []state.SignedState{expectedSignedState}}
			want = protocols.SideEffects{MessagesToSend: []protocols.Message{forAlice, forIrene}}
			// TODO ^^^ The test is currently sensitive to the order of the messages. It should not be.

			if diff := cmp.Diff(want, got); diff != "" {
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("TestCrank: side effects mismatch (-want +got):\n%s", diff)
				}

			}

			// Manually progress the extended state by collecting postfund signatures
			alicePost := signState(o.V.PostFundState(), alice)
			p1Post := signState(o.V.PostFundState(), p1)
			o.V.AddSignedStates([]state.SignedState{alicePost, p1Post})

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
			var s, _ = New(false, vPreFund, my.address, n, my.role, ledgerChannelToMyLeft, ledgerChannelToMyRight)
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
			prefundsignedstate := signState(s.V.PreFundState(), alice)
			e.SignedStates = append(e.SignedStates, prefundsignedstate)

			updatedObj, err := s.Update(e)
			updated := updatedObj.(VirtualFundObjective)
			if err != nil {
				t.Error(err)
			}
			if updated.V.SignedStateForTurnNum[0].HasSignatureForParticipant(alice.role) != true {
				t.Error(`Objective data not updated as expected`)
			}

			// Part 2: a signature on Bob's ledger channel (on his left)
			f := protocols.ObjectiveEvent{
				ObjectiveId: s.Id(),
			}
			f.SignedStates = make([]state.SignedState, 0)

			ledger, _ := ledger.CreateTestLedger(left, right, &bob.privateKey, 0, big.NewInt(0))
			ss := signState(ledger.PreFundState(), bob)

			f.SignedStates = append(f.SignedStates, ss)

			updatedObj, err = s.Update(f)
			updated = updatedObj.(VirtualFundObjective)
			if err != nil {
				t.Error(err)
			}
			if !updated.ToMyLeft.ledgerChannelAffordsExpectedGuarantees() != true {
				t.Error(`Objective data not updated as expected`)
			}

		}

		t.Run(`New`, testNew)
		t.Run(`Update`, testUpdate)
		t.Run(`Crank`, testCrank)

	}

	t.Run(`AsAlice`, AsAlice)
	t.Run(`AsBob`, AsBob)

}
