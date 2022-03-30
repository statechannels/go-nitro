package virtualfund

import (
	"bytes"
	"fmt"
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

// signPreAndPostFundingStates is a test utility function which applies signatures from
// multiple participants to pre and post fund states
func signPreAndPostFundingStates(ledger *channel.TwoPartyLedger, secretKeys []*[]byte) {
	for _, sk := range secretKeys {
		_, _ = ledger.SignAndAddPrefund(sk)
		_, _ = ledger.SignAndAddPostfund(sk)
	}
}

// signLatest is a test utility function which applies signatures from
// multiple participants to the latest recorded state
func signLatest(ledger *channel.TwoPartyLedger, secretKeys [][]byte) {

	// Find the largest turn num and therefore the latest state
	turnNum := uint64(0)
	for t := range ledger.SignedStateForTurnNum {
		if t > turnNum {
			turnNum = t
		}
	}
	// Sign it
	toSign := ledger.SignedStateForTurnNum[turnNum]
	for _, secretKey := range secretKeys {
		_ = toSign.Sign(&secretKey)
	}
	ledger.Channel.AddSignedState(toSign)
}

// addLedgerProposal calculates the ledger proposal state, signs it and adds it to the ledger.
func addLedgerProposal(
	ledger *channel.TwoPartyLedger,
	left types.Destination,
	right types.Destination,
	guaranteeDestination types.Destination,
	secretKey *[]byte,
) {

	supported, _ := ledger.LatestSupportedState()
	nextState := constructLedgerProposal(supported, left, right, guaranteeDestination)
	_, _ = ledger.SignAndAddState(nextState, secretKey)
}

// constructLedgerProposal returns a new ledger state with an updated outcome that includes the proposal
func constructLedgerProposal(
	supported state.State,
	left types.Destination,
	right types.Destination,
	guaranteeDestination types.Destination,
) state.State {
	leftAmount := types.Funds{types.Address{}: big.NewInt(6)}
	rightAmount := types.Funds{types.Address{}: big.NewInt(4)}
	nextState := supported.Clone()

	nextState.TurnNum = nextState.TurnNum + 1
	nextState.Outcome, _ = nextState.Outcome.DivertToGuarantee(left, right, leftAmount, rightAmount, guaranteeDestination)
	return nextState
}

// newTestTwoPartyLedger creates a new two party ledger channel based on the provided allocations. The channel will appear to be fully-funded on chain.
func newTestTwoPartyLedger(allocations []outcome.Allocation, myAddress types.Address, nonce *big.Int) (*channel.TwoPartyLedger, error) {

	initialState := state.State{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{},
		ChannelNonce:      nonce,
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
		AppData:           []byte{},
		Outcome: outcome.Exit{outcome.SingleAssetExit{
			Allocations: outcome.Allocations{},
		}},
		TurnNum: 0,
		IsFinal: false,
	}
	foundMyAddress := false
	var myIndex uint
	for i, alloc := range allocations {
		a, err := alloc.Destination.ToAddress()
		if err != nil {
			ntpl := channel.TwoPartyLedger{}
			return &ntpl, fmt.Errorf("could not extract address: %w", err)
		}
		initialState.Participants = append(initialState.Participants, a)
		initialState.Outcome[0].Allocations = append(initialState.Outcome[0].Allocations, alloc)
		if a == myAddress {
			foundMyAddress = true
			myIndex = uint(i)
		}
	}

	if !foundMyAddress {
		panic("Destination corresponding to myAddress not found in outcome")
	}

	ledger, lErr := channel.NewTwoPartyLedger(initialState, myIndex)
	if lErr != nil {
		return ledger, fmt.Errorf("error creating ledger: %w", lErr)
	}
	ledger.OnChainFunding = ledger.PreFundState().Outcome.TotalAllocated()

	return ledger, nil
}

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
					Amount:      big.NewInt(6),
				},
				outcome.Allocation{
					Destination: bob.destination,
					Amount:      big.NewInt(4),
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
					r, _ = newTestTwoPartyLedger(
						outcome.Allocations{
							{Destination: my.destination, Amount: big.NewInt(6)},
							{Destination: p1.destination, Amount: big.NewInt(4)},
						},
						my.address, big.NewInt(0))
					signPreAndPostFundingStates(r, []*[]byte{&alice.privateKey, &p1.privateKey}) // TODO these steps could be absorbed into CreateTestLedger

				}
			case 1:
				{
					l, _ = newTestTwoPartyLedger(
						outcome.Allocations{
							{Destination: alice.destination, Amount: big.NewInt(6)},
							{Destination: my.destination, Amount: big.NewInt(4)},
						},
						my.address, big.NewInt(0))
					r, _ = newTestTwoPartyLedger(
						outcome.Allocations{
							{Destination: my.destination, Amount: big.NewInt(6)},
							{Destination: bob.destination, Amount: big.NewInt(4)},
						},
						my.address, big.NewInt(0))
					signPreAndPostFundingStates(l, []*[]byte{&alice.privateKey, &p1.privateKey})
					signPreAndPostFundingStates(r, []*[]byte{&p1.privateKey, &bob.privateKey})
				}
			case 2:
				{
					l, _ = newTestTwoPartyLedger(
						outcome.Allocations{
							{Destination: p1.destination, Amount: big.NewInt(6)},
							{Destination: my.destination, Amount: big.NewInt(4)},
						},
						my.address, big.NewInt(0))
					signPreAndPostFundingStates(l, []*[]byte{&bob.privateKey, &p1.privateKey})
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
			o, err := constructFromState(false, vPreFund, my.address, ledgerChannelToMyLeft, nil, ledgerChannelToMyRight, nil) // todo: #420 deprecate TwoPartyLedgers
			if err != nil {
				t.Fatal(err)
			}

			var expectedGuaranteeMetadataLeft outcome.GuaranteeMetadata
			var expectedGuaranteeMetadataRight outcome.GuaranteeMetadata
			switch my.role {
			case alice.role:
				{
					expectedGuaranteeMetadataRight = outcome.GuaranteeMetadata{Left: alice.destination, Right: p1.destination}
				}
			case p1.role:
				{
					expectedGuaranteeMetadataLeft = outcome.GuaranteeMetadata{Left: alice.destination, Right: p1.destination}
					expectedGuaranteeMetadataRight = outcome.GuaranteeMetadata{Left: p1.destination, Right: bob.destination}
				}
			case bob.role:
				{
					expectedGuaranteeMetadataLeft = outcome.GuaranteeMetadata{Left: p1.destination, Right: bob.destination}
				}
			}
			if (expectedGuaranteeMetadataLeft != outcome.GuaranteeMetadata{}) {
				gotLeft := o.ToMyLeft.getExpectedGuarantees()[types.Address{}] // VState only has one (native) asset represented by the zero address
				expectedEncodedGuaranteeMetadataLeft, _ := expectedGuaranteeMetadataLeft.Encode()
				wantLeft := outcome.Allocation{
					Destination:    o.V.Id,
					Amount:         big.NewInt(0).Set(vPreFund.VariablePart().Outcome[0].TotalAllocated()),
					AllocationType: outcome.GuaranteeAllocationType,
					Metadata:       expectedEncodedGuaranteeMetadataLeft,
				}
				if diff := cmp.Diff(wantLeft, gotLeft); diff != "" {
					t.Fatalf("TestNew: expectedGuarantee mismatch (-want +got):\n%s", diff)
				}
			}
			if (expectedGuaranteeMetadataRight != outcome.GuaranteeMetadata{}) {
				gotRight := o.ToMyRight.getExpectedGuarantees()[types.Address{}] // VState only has one (native) asset represented by the zero address
				expectedEncodedGuaranteeMetadataRight, _ := expectedGuaranteeMetadataRight.Encode()
				wantRight := outcome.Allocation{
					Destination:    o.V.Id,
					Amount:         big.NewInt(0).Set(vPreFund.VariablePart().Outcome[0].TotalAllocated()),
					AllocationType: outcome.GuaranteeAllocationType,
					Metadata:       expectedEncodedGuaranteeMetadataRight,
				}
				if diff := cmp.Diff(wantRight, gotRight); diff != "" {
					t.Fatalf("TestNew: expectedGuarantee mismatch (-want +got):\n%s", diff)
				}
			}

		}

		testclone := func(t *testing.T) {
			ledgerChannelToMyLeft, ledgerChannelToMyRight := prepareLedgerChannels(my.role)

			o, _ := constructFromState(false, vPreFund, my.address, ledgerChannelToMyLeft, nil, ledgerChannelToMyRight, nil) // todo: #420 deprecate TwoPartyLedgers

			clone := o.clone()

			if diff := cmp.Diff(o, clone); diff != "" {
				t.Fatalf("Clone: mismatch (-want +got):\n%s", diff)
			}
		}

		testCrank := func(t *testing.T) {
			ledgerChannelToMyLeft, ledgerChannelToMyRight := prepareLedgerChannels(my.role)
			var s, _ = constructFromState(false, vPreFund, my.address, ledgerChannelToMyLeft, nil, ledgerChannelToMyRight, nil) // todo: #420 deprecate TwoPartyLedgers
			// Assert that cranking an unapproved objective returns an error
			if _, _, _, err := s.Crank(&my.privateKey); err == nil {
				t.Fatal(`Expected error when cranking unapproved objective, but got nil`)
			}

			// Approve the objective, so that the rest of the test cases can run.
			o := s.Approve().(*Objective)
			// To test the finite state progression, we are going to progressively mutate o
			// And then crank it to see which "pause point" (WaitingFor) we end up at.

			// Initial Crank
			oObj, got, waitingFor, err := o.Crank(&my.privateKey)
			o = oObj.(*Objective)
			if err != nil {
				t.Fatal(err)
			}
			if waitingFor != WaitingForCompletePrefund {
				t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForCompletePrefund, waitingFor)
			}

			expectedSignedState := state.NewSignedState(o.V.PreFundState())
			mySig, _ := o.V.PreFundState().Sign(my.privateKey)
			_ = expectedSignedState.AddSignature(mySig)
			assertSideEffectsContainsMessagesForPeersWith(got, expectedSignedState, my.role, t)

			// Manually progress the extended state by collecting prefund signatures
			collectPeerSignaturesOnSetupState(o.V, my.role, true)

			// Cranking should move us to the next waiting point, update the ledger channel, and alter the extended state to reflect that
			// TODO: Check that ledger channel is updated as expected
			oObj, got, waitingFor, _ = o.Crank(&my.privateKey)

			// Check that the messsages contain the expected ledger proposals
			// We only expect a proposal in the right ledger channel, as we will be the leader in that ledger channel
			switch my.role {
			case 0:
				{
					supported, _ := o.ToMyRight.Channel.LatestSupportedState()
					expectedSignedState := state.NewSignedState(constructLedgerProposal(supported, types.AddressToDestination(alice.address), types.AddressToDestination(p1.address), o.V.Id))
					_ = expectedSignedState.Sign(&my.privateKey)

					assertSideEffectsContainsMessageWith(got, expectedSignedState, p1, t)

				}
			case 1:
				{
					supported, _ := o.ToMyRight.Channel.LatestSupportedState()
					expectedSignedState := state.NewSignedState(constructLedgerProposal(supported, types.AddressToDestination(p1.address), types.AddressToDestination(bob.address), o.V.Id))
					_ = expectedSignedState.Sign(&my.privateKey)

					assertSideEffectsContainsMessageWith(got, expectedSignedState, bob, t)
				}
			}
			if waitingFor != WaitingForCompleteFunding {
				t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForCompleteFunding, waitingFor)
			}

			o = oObj.(*Objective)

			//Update the ledger funding by mimicing other participants either proposing an update or accepting our update
			switch my.role {
			case 0:
				{
					signLatest(o.ToMyRight.Channel, [][]byte{p1.privateKey})
				}
			case 1:
				{
					// If we are P1 we mimic Alice proposing the update to the ledger channel
					addLedgerProposal(o.ToMyLeft.Channel, types.AddressToDestination(alice.address), types.AddressToDestination(p1.address), o.V.Id, &alice.privateKey)
					// We mimic Bob accepting the proposal on the right
					signLatest(o.ToMyRight.Channel, [][]byte{bob.privateKey})

				}
			case 2:
				{
					// If we are Bob we mimic P1 proposing the update to the ledger channel
					addLedgerProposal(o.ToMyLeft.Channel, types.AddressToDestination(p1.address), types.AddressToDestination(bob.address), o.V.Id, &p1.privateKey)

				}
			}

			// Cranking now should not generate side effects, because we already did that
			oObj, got, waitingFor, err = o.Crank(&my.privateKey)
			o = oObj.(*Objective)
			if err != nil {
				t.Fatal(err)
			}
			if waitingFor != WaitingForCompletePostFund {
				t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForCompletePostFund, waitingFor)
			}

			// Check that the messsages contain the expected ledger acceptances
			// We only expect an acceptance in the left ledger channel as we will be the follower in that ledger channel
			switch my.role {
			case 1:
				{
					supported, _ := o.ToMyLeft.Channel.LatestSupportedState()
					expectedSignedState := state.NewSignedState(supported)
					_ = expectedSignedState.Sign(&my.privateKey)

					assertSideEffectsContainsMessageWith(got, expectedSignedState, alice, t)

				}
			case 2:
				{
					supported, _ := o.ToMyLeft.Channel.LatestSupportedState()
					expectedSignedState := state.NewSignedState(supported)
					_ = expectedSignedState.Sign(&my.privateKey)

					assertSideEffectsContainsMessageWith(got, expectedSignedState, p1, t)
				}
			}

			expectedSignedState = state.NewSignedState(o.V.PostFundState())
			mySig, _ = o.V.PostFundState().Sign(my.privateKey)
			_ = expectedSignedState.AddSignature(mySig)
			assertSideEffectsContainsMessagesForPeersWith(got, expectedSignedState, my.role, t)

			// Manually progress the extended state by collecting postfund signatures
			collectPeerSignaturesOnSetupState(o.V, my.role, false)

			// This should be the final crank...
			_, _, waitingFor, err = o.Crank(&my.privateKey)
			if err != nil {
				t.Fatal(err)
			}
			if waitingFor != WaitingForNothing {
				t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForNothing, waitingFor)
			}

		}

		testUpdate := func(t *testing.T) {
			ledgerChannelToMyLeft, ledgerChannelToMyRight := prepareLedgerChannels(my.role)
			var s, _ = constructFromState(false, vPreFund, my.address, ledgerChannelToMyLeft, nil, ledgerChannelToMyRight, nil) // todo: #420 deprecate TwoPartyLedgers
			// Prepare an event with a mismatched objectiveId
			e := protocols.ObjectiveEvent{
				ObjectiveId: "some-other-id",
			}
			// Assert that Updating the objective with such an event returns an error
			// TODO is this the behaviour we want? Below with the signatures, we prefer a log + NOOP (no error)
			if _, err := s.Update(e); err == nil {
				t.Fatal(`Objective ID mismatch -- expected an error but did not get one`)
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
			updated := updatedObj.(*Objective)
			if err != nil {
				t.Fatal(err)
			}

			switch my.role {
			case 0:
				{
					if !updated.V.SignedStateForTurnNum[1].HasSignatureForParticipant(p1.role) {
						t.Fatal(`Objective data not updated as expected`)
					}
				}
			case 1:
				{
					if !updated.V.SignedStateForTurnNum[1].HasSignatureForParticipant(alice.role) {
						t.Fatal(`Objective data not updated as expected`)
					}
				}
			case 2:
				{
					if !updated.V.SignedStateForTurnNum[1].HasSignatureForParticipant(p1.role) {
						t.Fatal(`Objective data not updated as expected`)
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
			updated = updatedObj.(*Objective)
			if err != nil {
				t.Fatal(err)
			}

			switch my.role {
			case 0:
				{
					if !updated.ToMyRight.Channel.SignedStateForTurnNum[someTurnNum].HasSignatureForParticipant((updated.ToMyRight.Channel.MyIndex + 1) % 2) {
						t.Fatal(`Objective data not updated as expected`)
					}
				}
			case 1:
				{
					if !updated.ToMyRight.Channel.SignedStateForTurnNum[someTurnNum].HasSignatureForParticipant((updated.ToMyRight.Channel.MyIndex + 1) % 2) {
						t.Fatal(`Objective data not updated as expected`)
					}
				}
			case 2:
				{
					if !updated.ToMyLeft.Channel.SignedStateForTurnNum[someTurnNum].HasSignatureForParticipant((updated.ToMyLeft.Channel.MyIndex + 1) % 2) {
						t.Fatal(`Objective data not updated as expected`)
					}
				}
			}

		}
		t.Run(`New`, testNew)
		t.Run(`clone`, testclone)
		t.Run(`Crank`, testCrank)
		t.Run(`Update`, testUpdate)

	}

	t.Run(`AsAlice`, func(t *testing.T) { TestAs(alice, t) })
	t.Run(`AsBob`, func(t *testing.T) { TestAs(bob, t) })
	t.Run(`AsP1`, func(t *testing.T) { TestAs(p1, t) })
}
