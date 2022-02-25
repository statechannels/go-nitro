package virtualfund

import (
	"bytes"
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

func TestSingleHopVirtualFund(t *testing.T) {

	type actor struct {
		address     types.Address
		destination types.Destination
		privateKey  []byte
		role        uint
	}

	////////////
	// ACTORS //
	////////////

	alice := actor{
		address:     common.HexToAddress(`0xD9995BAE12FEe327256FFec1e3184d492bD94C31`),
		destination: types.AddressToDestination(common.HexToAddress(`0xD9995BAE12FEe327256FFec1e3184d492bD94C31`)),
		privateKey:  common.Hex2Bytes(`7ab741b57e8d94dd7e1a29055646bafde7010f38a900f55bbd7647880faa6ee8`),
		role:        0,
	}

	p1 := actor{ // Aliases: The Hub, Irene
		address:     common.HexToAddress(`0xd4Fa489Eacc52BA59438993f37Be9fcC20090E39`),
		destination: types.AddressToDestination(common.HexToAddress(`0xd4Fa489Eacc52BA59438993f37Be9fcC20090E39`)),
		privateKey:  common.Hex2Bytes(`2030b463177db2da82908ef90fa55ddfcef56e8183caf60db464bc398e736e6f`),
		role:        1,
	}

	bob := actor{
		address:     common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`),
		destination: types.AddressToDestination(common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`)),
		privateKey:  common.Hex2Bytes(`62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2`),
		role:        2,
	}

	// assertSideEffectsContainsMessageWith fails the test instantly if the supplied side effects does not contain a message for the supplied actor with the supplied expected signed state.
	assertSideEffectsContainsMessageWith := func(ses protocols.SideEffects, expectedSignedState state.SignedState, to actor, t *testing.T) {
		for _, msg := range ses.MessagesToSend {
			for _, ss := range msg.SignedStates {
				if ss.Equal(expectedSignedState) && bytes.Equal(msg.To[:], to.address[:]) {
					return
				}
			}
		}
		t.Fatalf("side effects %v do not contain signed state %v for %v", ses, expectedSignedState, to)
	}

	// assertSideEffectsContainsMessageWith calls assertSideEffectsContainsMessageWith for all peers of the actor with role myRole.
	assertSideEffectsContainsMessagesForPeersWith := func(ses protocols.SideEffects, expectedSignedState state.SignedState, myRole uint, t *testing.T) {
		if myRole != alice.role {
			assertSideEffectsContainsMessageWith(ses, expectedSignedState, alice, t)
		}
		if myRole != p1.role {
			assertSideEffectsContainsMessageWith(ses, expectedSignedState, p1, t)
		}
		if myRole != bob.role {
			assertSideEffectsContainsMessageWith(ses, expectedSignedState, bob, t)
		}
	}

	collectPeerSignaturesOnSetupState := func(V *channel.SingleHopVirtualChannel, myRole uint, prefund bool) {
		var state state.State
		if prefund {
			state = V.PreFundState()
		} else {
			state = V.PostFundState()
		}

		if myRole != alice.role {
			aliceSig, _ := state.Sign(alice.privateKey)
			V.AddStateWithSignature(state, aliceSig)
		}
		if myRole != p1.role {
			p1Sig, _ := state.Sign(p1.privateKey)
			V.AddStateWithSignature(state, p1Sig)
		}
		if myRole != bob.role {
			bobSig, _ := state.Sign(bob.privateKey)
			V.AddStateWithSignature(state, bobSig)
		}
	}

	/////////////////////
	// VIRTUAL CHANNEL //
	/////////////////////

	// Virtual Channel
	vPreFund := state.State{
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
	vPostFund := vPreFund.Clone()
	vPostFund.TurnNum = 1

	TestAs := func(my actor, t *testing.T) {

		// prepareLedgerChannels constructs and return two new ledger channels. Depending on the role, one of them may be nil.
		prepareLedgerChannels := func(role uint) (*channel.TwoPartyLedger, *channel.TwoPartyLedger) {
			var l *channel.TwoPartyLedger
			var r *channel.TwoPartyLedger
			switch my.role {
			case 0:
				{
					r, _ = ledger.CreateTestLedger(
						outcome.Allocation{Destination: my.destination, Amount: big.NewInt(5)},
						outcome.Allocation{Destination: p1.destination, Amount: big.NewInt(5)},
						&my.privateKey, 0, big.NewInt(0))
					ledger.SignPreAndPostFundingStates(r, []*[]byte{&alice.privateKey, &p1.privateKey}) // TODO these steps could be absorbed into CreateTestLedger

				}
			case 1:
				{
					l, _ = ledger.CreateTestLedger(
						outcome.Allocation{Destination: alice.destination, Amount: big.NewInt(5)},
						outcome.Allocation{Destination: my.destination, Amount: big.NewInt(5)},
						&alice.privateKey, 1, big.NewInt(0))
					r, _ = ledger.CreateTestLedger(
						outcome.Allocation{Destination: my.destination, Amount: big.NewInt(5)},
						outcome.Allocation{Destination: bob.destination, Amount: big.NewInt(5)},
						&alice.privateKey, 0, big.NewInt(0))
					ledger.SignPreAndPostFundingStates(l, []*[]byte{&alice.privateKey, &p1.privateKey})
					ledger.SignPreAndPostFundingStates(r, []*[]byte{&p1.privateKey, &bob.privateKey})
				}
			case 2:
				{
					l, _ = ledger.CreateTestLedger(
						outcome.Allocation{Destination: p1.destination, Amount: big.NewInt(5)},
						outcome.Allocation{Destination: my.destination, Amount: big.NewInt(5)},
						&alice.privateKey, 1, big.NewInt(0))
					ledger.SignPreAndPostFundingStates(l, []*[]byte{&bob.privateKey, &p1.privateKey})
				}
			default:
				{
					panic(`invalid role`)
				}

			}
			return l, r
		}

		testNew := func(t *testing.T) {
			ledgerChannelToMyLeft, ledgerChannelToMyRight := prepareLedgerChannels(my.role)
			// Assert that a valid set of constructor args does not result in an error
			o, err := New(false, vPreFund, my.address, ledgerChannelToMyLeft, ledgerChannelToMyRight)
			if err != nil {
				t.Error(err)
			}

			if my.role != alice.role {
				// Check left guarantee
				got := o.ToMyLeft.ExpectedGuarantees[types.Address{}] // VState only has one (native) asset represented by the zero address
				expectedGuaranteeMetadata := outcome.GuaranteeMetadata{Left: ledgerChannelToMyLeft.TheirDestination(), Right: ledgerChannelToMyLeft.MyDestination()}
				expectedEncodedGuaranteeMetadata, _ := expectedGuaranteeMetadata.Encode()
				want := outcome.Allocation{
					Destination:    o.V.Id,
					Amount:         big.NewInt(0).Set(vPreFund.VariablePart().Outcome[0].TotalAllocated()),
					AllocationType: outcome.GuaranteeAllocationType,
					Metadata:       expectedEncodedGuaranteeMetadata,
				}
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("TestNew: expectedGuarantee mismatch (-want +got):\n%s", diff)
				}
			}

			if my.role != bob.role {
				// Check right guarantee
				got := o.ToMyRight.ExpectedGuarantees[types.Address{}] // VState only has one (native) asset represented by the zero address
				expectedGuaranteeMetadata := outcome.GuaranteeMetadata{Left: ledgerChannelToMyRight.MyDestination(), Right: ledgerChannelToMyRight.TheirDestination()}
				expectedEncodedGuaranteeMetadata, _ := expectedGuaranteeMetadata.Encode()
				want := outcome.Allocation{
					Destination:    o.V.Id,
					Amount:         big.NewInt(0).Set(vPreFund.VariablePart().Outcome[0].TotalAllocated()),
					AllocationType: outcome.GuaranteeAllocationType,
					Metadata:       expectedEncodedGuaranteeMetadata,
				}
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("TestNew: expectedGuarantee mismatch (-want +got):\n%s", diff)
				}
			}
		}

		testCrank := func(t *testing.T) {
			ledgerChannelToMyLeft, ledgerChannelToMyRight := prepareLedgerChannels(my.role)
			var s, _ = New(false, vPreFund, my.address, ledgerChannelToMyLeft, ledgerChannelToMyRight)
			// Assert that cranking an unapproved objective returns an error
			if _, _, _, _, err := s.Crank(&my.privateKey); err == nil {
				t.Error(`Expected error when cranking unapproved objective, but got nil`)
			}

			// Approve the objective, so that the rest of the test cases can run.
			o := s.Approve().(VirtualFundObjective)
			// To test the finite state progression, we are going to progressively mutate o
			// And then crank it to see which "pause point" (WaitingFor) we end up at.

			// Initial Crank
			oObj, got, waitingFor, _, err := o.Crank(&my.privateKey)
			o = oObj.(VirtualFundObjective)
			if err != nil {
				t.Error(err)
			}
			if waitingFor != WaitingForCompletePrefund {
				t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForCompletePrefund, waitingFor)
			}

			expectedSignedState := state.NewSignedState(o.V.PreFundState())
			mySig, _ := o.V.PreFundState().Sign(my.privateKey)
			_ = expectedSignedState.AddSignature(mySig)
			assertSideEffectsContainsMessagesForPeersWith(got, expectedSignedState, my.role, t)

			// Manually progress the extended state by collecting prefund signatures
			collectPeerSignaturesOnSetupState(o.V, my.role, true)

			// Cranking should move us to the next waiting point, generate ledger requests as a side effect, and alter the extended state to reflect that
			var gotRequests []protocols.LedgerRequest
			oObj, _, waitingFor, gotRequests, err = o.Crank(&my.privateKey)
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

			wantRequests := []protocols.LedgerRequest{}
			switch my.role {
			case 0:
				{
					wantRequests = append(wantRequests, protocols.LedgerRequest{
						ObjectiveId: o.Id(),
						LedgerId:    ledgerChannelToMyRight.Id,
						Destination: s.V.Id,
						Left:        my.destination, Right: p1.destination,
						LeftAmount:  types.Funds{types.Address{}: big.NewInt(5)},
						RightAmount: types.Funds{types.Address{}: big.NewInt(5)},
					})
				}
			case 1:
				{
					wantRequests = append(wantRequests, protocols.LedgerRequest{
						ObjectiveId: o.Id(),
						LedgerId:    ledgerChannelToMyLeft.Id,
						Destination: s.V.Id,
						Left:        alice.destination, Right: my.destination,
						LeftAmount:  types.Funds{types.Address{}: big.NewInt(5)},
						RightAmount: types.Funds{types.Address{}: big.NewInt(5)},
					})
					wantRequests = append(wantRequests, protocols.LedgerRequest{
						ObjectiveId: o.Id(),
						LedgerId:    ledgerChannelToMyRight.Id,
						Destination: s.V.Id,
						Left:        my.destination, Right: bob.destination,
						LeftAmount:  types.Funds{types.Address{}: big.NewInt(5)},
						RightAmount: types.Funds{types.Address{}: big.NewInt(5)},
					})
				}
			case 2:
				{
					wantRequests = append(wantRequests, protocols.LedgerRequest{
						ObjectiveId: o.Id(),
						LedgerId:    ledgerChannelToMyLeft.Id,
						Destination: s.V.Id,
						Left:        p1.destination, Right: my.destination,
						LeftAmount:  types.Funds{types.Address{}: big.NewInt(5)},
						RightAmount: types.Funds{types.Address{}: big.NewInt(5)},
					})
				}
			}

			if diff := cmp.Diff(wantRequests, gotRequests, cmp.Comparer(types.Equal)); diff != "" {
				t.Errorf("TestCrank: ledger requests mismatch (-want +got):\n%s", diff)
			}

			ledgerManager := ledger.NewLedgerManager()
			switch my.role {
			case 0:
				{
					_, _ = ledgerManager.HandleRequest(o.ToMyRight.Channel, gotRequests[0], &my.privateKey)
					ledger.SignLatest(o.ToMyRight.Channel, [][]byte{p1.privateKey})
				}
			case 1:
				{
					_, _ = ledgerManager.HandleRequest(o.ToMyLeft.Channel, gotRequests[0], &my.privateKey)
					ledger.SignLatest(o.ToMyLeft.Channel, [][]byte{alice.privateKey})
					_, _ = ledgerManager.HandleRequest(o.ToMyRight.Channel, gotRequests[1], &my.privateKey)
					ledger.SignLatest(o.ToMyRight.Channel, [][]byte{bob.privateKey})

				}
			case 2:
				{
					_, _ = ledgerManager.HandleRequest(o.ToMyLeft.Channel, gotRequests[0], &my.privateKey)
					ledger.SignLatest(o.ToMyLeft.Channel, [][]byte{p1.privateKey})

				}
			}

			// Cranking now should not generate side effects, because we already did that
			oObj, got, waitingFor, _, err = o.Crank(&my.privateKey)
			o = oObj.(VirtualFundObjective)
			if err != nil {
				t.Error(err)
			}
			if waitingFor != WaitingForCompletePostFund {
				t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForCompletePostFund, waitingFor)
			}
			expectedSignedState = state.NewSignedState(o.V.PostFundState())
			mySig, _ = o.V.PostFundState().Sign(my.privateKey)
			_ = expectedSignedState.AddSignature(mySig)
			assertSideEffectsContainsMessagesForPeersWith(got, expectedSignedState, my.role, t)

			// Manually progress the extended state by collecting postfund signatures
			collectPeerSignaturesOnSetupState(o.V, my.role, false)

			// This should be the final crank...
			_, _, waitingFor, _, err = o.Crank(&my.privateKey)
			if err != nil {
				t.Error(err)
			}
			if waitingFor != WaitingForNothing {
				t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForNothing, waitingFor)
			}

		}

		testUpdate := func(t *testing.T) {
			ledgerChannelToMyLeft, ledgerChannelToMyRight := prepareLedgerChannels(my.role)
			var s, _ = New(false, vPreFund, my.address, ledgerChannelToMyLeft, ledgerChannelToMyRight)
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

			vPostFund := s.V.PostFundState()
			ss := state.NewSignedState(vPostFund)

			switch my.role {
			case 0:
				{
					_ = ss.Sign(&p1.privateKey)

				}
			case 1:
				{
					_ = ss.Sign(&alice.privateKey)

				}
			case 2:
				{
					_ = ss.Sign(&p1.privateKey)

				}
			}
			e.SignedStates = append(e.SignedStates, ss)

			updatedObj, err := s.Update(e)
			updated := updatedObj.(VirtualFundObjective)
			if err != nil {
				t.Error(err)
			}

			switch my.role {
			case 0:
				{
					if !updated.V.SignedStateForTurnNum[1].HasSignatureForParticipant(p1.role) {
						t.Error(`Objective data not updated as expected`)
					}
				}
			case 1:
				{
					if !updated.V.SignedStateForTurnNum[1].HasSignatureForParticipant(alice.role) {
						t.Error(`Objective data not updated as expected`)
					}
				}
			case 2:
				{
					if !updated.V.SignedStateForTurnNum[1].HasSignatureForParticipant(p1.role) {
						t.Error(`Objective data not updated as expected`)
					}
				}
			}

			// Part 2: a signature on a relevant ledger channel
			f := protocols.ObjectiveEvent{
				ObjectiveId: s.Id(),
			}
			f.SignedStates = make([]state.SignedState, 0)
			someTurnNum := uint64(99)
			switch my.role {
			case 0:
				{
					s := ledgerChannelToMyRight.PreFundState().Clone()
					s.TurnNum = someTurnNum
					ss = state.NewSignedState(s)
					_ = ss.Sign(&p1.privateKey)
				}
			case 1:
				{
					s := ledgerChannelToMyRight.PreFundState().Clone()
					s.TurnNum = someTurnNum
					ss = state.NewSignedState(s)
					_ = ss.Sign(&bob.privateKey)
				}
			case 2:
				{
					s := ledgerChannelToMyLeft.PreFundState().Clone()
					s.TurnNum = someTurnNum
					ss = state.NewSignedState(s)
					_ = ss.Sign(&p1.privateKey)
				}
			}
			f.SignedStates = append(f.SignedStates, ss)

			updatedObj, err = s.Update(f)
			updated = updatedObj.(VirtualFundObjective)
			if err != nil {
				t.Error(err)
			}

			switch my.role {
			case 0:
				{
					if !updated.ToMyRight.Channel.SignedStateForTurnNum[someTurnNum].HasSignatureForParticipant((updated.ToMyRight.Channel.MyIndex + 1) % 2) {
						t.Error(`Objective data not updated as expected`)
					}
				}
			case 1:
				{
					if !updated.ToMyRight.Channel.SignedStateForTurnNum[someTurnNum].HasSignatureForParticipant((updated.ToMyRight.Channel.MyIndex + 1) % 2) {
						t.Error(`Objective data not updated as expected`)
					}
				}
			case 2:
				{
					if !updated.ToMyLeft.Channel.SignedStateForTurnNum[someTurnNum].HasSignatureForParticipant((updated.ToMyLeft.Channel.MyIndex + 1) % 2) {
						t.Error(`Objective data not updated as expected`)
					}
				}
			}

		}
		t.Run(`New`, testNew)
		t.Run(`Crank`, testCrank)
		t.Run(`Update`, testUpdate)

	}

	t.Run(`AsAlice`, func(t *testing.T) { TestAs(alice, t) })
	t.Run(`AsBob`, func(t *testing.T) { TestAs(bob, t) })
	t.Run(`AsP1`, func(t *testing.T) { TestAs(p1, t) })
}
