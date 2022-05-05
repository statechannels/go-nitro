package directfund

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/internal/channel"
	"github.com/statechannels/go-nitro/internal/channel/consensus_channel"
	"github.com/statechannels/go-nitro/internal/channel/state"
	"github.com/statechannels/go-nitro/internal/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/protocols"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/types"
)

var alice, bob testactors.Actor = testactors.Alice, testactors.Bob

var testState = state.State{
	ChainId:           big.NewInt(9001),
	Participants:      []types.Address{alice.Address(), bob.Address()},
	ChannelNonce:      big.NewInt(37140676580),
	AppDefinition:     common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`),
	ChallengeDuration: big.NewInt(60),
	AppData:           []byte{},
	Outcome: outcome.Exit{
		outcome.SingleAssetExit{
			Asset: types.Address{},
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: bob.Destination(), // Bob is first so we can easily test WaitingForMyTurnToFund
					Amount:      big.NewInt(5),
				},
				outcome.Allocation{
					Destination: alice.Destination(),
					Amount:      big.NewInt(5),
				},
			},
		},
	},
	TurnNum: 0,
	IsFinal: false,
}

// TestNew tests the constructor using a TestState fixture
func TestNew(t *testing.T) {
	// Assert that valid constructor args do not result in error
	if _, err := ConstructFromState(false, testState, testState.Participants[0]); err != nil {
		t.Error(err)
	}

	// Construct a final state
	finalState := testState.Clone()
	finalState.IsFinal = true

	if _, err := ConstructFromState(false, finalState, testState.Participants[0]); err == nil {
		t.Error("expected an error when constructing with an intial state marked final, but got nil")
	}

	nonParticipant := common.HexToAddress("0x5b53f71453aeCb03D837bfe170570d40aE736CB4")
	if _, err := ConstructFromState(false, testState, nonParticipant); err == nil {
		t.Error("expected an error when constructing with a participant not in the channel, but got nil")
	}
}

func TestConstructFromState(t *testing.T) {
	// Assert that valid constructor args do not result in error
	if _, err := ConstructFromState(false, testState, testState.Participants[0]); err != nil {
		t.Error(err)
	}

	// Construct a final state
	finalState := testState.Clone()
	finalState.IsFinal = true

	if _, err := ConstructFromState(false, finalState, testState.Participants[0]); err == nil {
		t.Error("expected an error when constructing with an intial state marked final, but got nil")
	}

	nonParticipant := common.HexToAddress("0x5b53f71453aeCb03D837bfe170570d40aE736CB4")
	if _, err := ConstructFromState(false, testState, nonParticipant); err == nil {
		t.Error("expected an error when constructing with a participant not in the channel, but got nil")
	}
}
func TestUpdate(t *testing.T) {
	// Construct various variables for use in TestUpdate
	var s, _ = ConstructFromState(false, testState, testState.Participants[0])

	var stateToSign state.State = s.C.PreFundState()
	var correctSignatureByParticipant, _ = stateToSign.Sign(alice.PrivateKey)
	// Prepare an event with a mismatched channelId
	e := protocols.ObjectiveEvent{
		ObjectiveId: "some-id",
	}
	// Assert that Updating the objective with such an event returns an error
	// TODO is this the behaviour we want? Below with the signatures, we prefer a log + NOOP (no error)
	if _, err := s.Update(e); err == nil {
		t.Error(`ChannelId mismatch -- expected an error but did not get one`)
	}

	// Now modify the event to give it the "correct" objective id,
	// and make a new Sigs map.
	// This prepares us for the rest of the test. We will reuse the same event multiple times
	e.ObjectiveId = s.Id()

	// Next, attempt to update the objective with correct signature by a participant on a relevant state
	// Assert that this results in an appropriate change in the extended state of the objective
	ss := state.NewSignedState(stateToSign)
	err := ss.AddSignature(correctSignatureByParticipant)
	if err != nil {
		t.Error(err)
	}
	e.SignedState = ss
	updatedObjective, err := s.Update(e)
	if err != nil {
		t.Error(err)
	}
	updated := updatedObjective.(*Objective)
	if updated.C.PreFundSignedByMe() != true {
		t.Error(`Objective data not updated as expected`)
	}

	// Finally, add some Holdings information to the event
	// Updating the objective with this event should overwrite the holdings that are stored
	newFunding := types.Funds{
		common.Address{}: big.NewInt(3),
	}
	highBlockNum := uint64(200)
	updatedObjective, err = s.UpdateWithChainEvent(chainservice.DepositedEvent{Holdings: newFunding, CommonEvent: chainservice.CommonEvent{BlockNum: highBlockNum}})
	if err != nil {
		t.Error(err)
	}
	updated = updatedObjective.(*Objective)
	if !updated.C.OnChainFunding.Equal(newFunding) {
		t.Error(`Objective data not updated as expected`, updated.C.OnChainFunding, newFunding)
	}
	if updated.latestBlockNumber != uint64(highBlockNum) {
		t.Error("Latest block number not updated as expected", updated.latestBlockNumber, highBlockNum)
	}

	// Update with stale funding information should be ignored
	staleFunding := types.Funds{}
	staleFunding[common.Address{}] = big.NewInt(2)
	lowBlockNum := uint64(100)

	updatedObjective, _ = updated.UpdateWithChainEvent(chainservice.DepositedEvent{Holdings: staleFunding, CommonEvent: chainservice.CommonEvent{BlockNum: uint64(lowBlockNum)}})

	updated = updatedObjective.(*Objective)

	if updated.C.OnChainFunding.Equal(staleFunding) {
		t.Error("OnChainFunding was updated to stale funding information", updated.C.OnChainFunding, staleFunding)
	}
	if updated.latestBlockNumber == uint64(lowBlockNum) {
		t.Error("latestBlockNumber was updated to stale block number", updated.latestBlockNumber, lowBlockNum)
	}

}

func compareSideEffect(a, b protocols.SideEffects) string {
	return cmp.Diff(a, b, cmp.AllowUnexported(a, state.SignedState{}, consensus_channel.Add{}, consensus_channel.Guarantee{}, consensus_channel.Remove{}))
}

func TestCrank(t *testing.T) {

	// BEGIN test data preparation
	var s, _ = ConstructFromState(false, testState, testState.Participants[0])
	var correctSignatureByAliceOnPreFund, _ = s.C.PreFundState().Sign(alice.PrivateKey)
	var correctSignatureByBobOnPreFund, _ = s.C.PreFundState().Sign(bob.PrivateKey)

	var correctSignatureByAliceOnPostFund, _ = s.C.PostFundState().Sign(alice.PrivateKey)
	var correctSignatureByBobOnPostFund, _ = s.C.PostFundState().Sign(bob.PrivateKey)

	// Prepare expected side effects
	preFundSS := state.NewSignedState(s.C.PreFundState())
	_ = preFundSS.AddSignature(correctSignatureByAliceOnPreFund)
	expectedPreFundSideEffects := protocols.SideEffects{
		MessagesToSend: []protocols.Message{
			{
				To: bob.Address(),
				Payloads: []protocols.MessagePayload{{
					ObjectiveId: s.Id(),
					SignedState: preFundSS,
				}},
			}}}

	postFundSS := state.NewSignedState(s.C.PostFundState())
	_ = postFundSS.AddSignature(correctSignatureByAliceOnPostFund)
	expectedPostFundSideEffects := protocols.SideEffects{
		MessagesToSend: []protocols.Message{
			{
				To: bob.Address(),
				Payloads: []protocols.MessagePayload{{
					ObjectiveId: s.Id(),
					SignedState: postFundSS,
				}},
			},
		}}
	expectedFundingSideEffects := protocols.SideEffects{
		TransactionsToSubmit: []protocols.ChainTransaction{{
			Type:      protocols.DepositTransactionType,
			ChannelId: s.C.Id,
			Deposit: types.Funds{
				testState.Outcome[0].Asset: testState.Outcome[0].Allocations[0].Amount,
			},
		}},
	}
	// END test data preparation

	// Assert that cranking an unapproved objective returns an error
	if _, _, _, err := s.Crank(&alice.PrivateKey); err == nil {
		t.Error(`Expected error when cranking unapproved objective, but got nil`)
	}

	// Approve the objective, so that the rest of the test cases can run.
	o := s.Approve().(*Objective)

	// To test the finite state progression, we are going to progressively mutate o
	// And then crank it to see
	//  - which "pause point" (WaitingFor) we end up at,
	//  - what side effects are declared.

	// Initial Crank
	_, sideEffects, waitingFor, err := o.Crank(&alice.PrivateKey)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForCompletePrefund {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForCompletePrefund, waitingFor)
	}

	if diff := compareSideEffect(expectedPreFundSideEffects, sideEffects); diff != "" {
		t.Errorf("Side effects mismatch (-want +got):\n%s", diff)
	}

	// Manually progress the extended state by collecting prefund signatures
	o.C.AddStateWithSignature(o.C.PreFundState(), correctSignatureByAliceOnPreFund)
	o.C.AddStateWithSignature(o.C.PreFundState(), correctSignatureByBobOnPreFund)

	// Cranking should move us to the next waiting point
	_, _, waitingFor, err = o.Crank(&alice.PrivateKey)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForMyTurnToFund {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForMyTurnToFund, waitingFor)
	}

	// Manually make the first "deposit"
	o.C.OnChainFunding[testState.Outcome[0].Asset] = testState.Outcome[0].Allocations[0].Amount
	_, sideEffects, waitingFor, err = o.Crank(&alice.PrivateKey)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForCompleteFunding {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForCompleteFunding, waitingFor)
	}

	if diff := cmp.Diff(expectedFundingSideEffects, sideEffects); diff != "" {
		t.Fatalf("Side effects mismatch (-want +got):\n%s", diff)
	}

	// Manually make the second "deposit"
	totalAmountAllocated := testState.Outcome[0].TotalAllocated()
	o.C.OnChainFunding[testState.Outcome[0].Asset] = totalAmountAllocated
	_, sideEffects, waitingFor, err = o.Crank(&alice.PrivateKey)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForCompletePostFund {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForCompletePostFund, waitingFor)
	}
	if diff := compareSideEffect(expectedPostFundSideEffects, sideEffects); diff != "" {
		t.Errorf("Side effects mismatch (-want +got):\n%s", diff)
	}

	// Manually progress the extended state by collecting postfund signatures
	o.C.AddStateWithSignature(o.C.PostFundState(), correctSignatureByAliceOnPostFund)
	o.C.AddStateWithSignature(o.C.PostFundState(), correctSignatureByBobOnPostFund)

	// This should be the final crank
	o.C.OnChainFunding[testState.Outcome[0].Asset] = totalAmountAllocated
	_, _, waitingFor, err = o.Crank(&alice.PrivateKey)
	if err != nil {
		t.Error(err)
	}
	if waitingFor != WaitingForNothing {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForNothing, waitingFor)
	}
}

func TestClone(t *testing.T) {
	compareObjectives := func(a, b protocols.Objective) string {
		return cmp.Diff(&a, &b, cmp.AllowUnexported(Objective{}, channel.Channel{}, big.Int{}, state.SignedState{}))
	}

	var s, _ = ConstructFromState(false, testState, testState.Participants[0])

	clone := s.clone()

	if diff := compareObjectives(&s, &clone); diff != "" {
		t.Errorf("Clone: mismatch (-want +got):\n%s", diff)
	}
}

func TestMarshalJSON(t *testing.T) {
	dfo, _ := ConstructFromState(false, testState, testState.Participants[0])

	encodedDfo, err := json.Marshal(dfo)

	if err != nil {
		t.Fatalf("error encoding direct-fund objective %v", dfo)
	}

	got := Objective{}
	if err := got.UnmarshalJSON(encodedDfo); err != nil {
		t.Fatalf("error unmarshaling test direct fund objective: %s", err.Error())
	}

	if !got.myDepositSafetyThreshold.Equal(dfo.myDepositSafetyThreshold) {
		t.Fatalf("expected myDepositSafetyThreshhold %v but got %v",
			dfo.myDepositSafetyThreshold, got.myDepositSafetyThreshold)
	}
	if !got.myDepositTarget.Equal(dfo.myDepositTarget) {
		t.Fatalf("expected myDepositTarget %v but got %v",
			dfo.myDepositTarget, got.myDepositTarget)
	}
	if !got.fullyFundedThreshold.Equal(dfo.fullyFundedThreshold) {
		t.Fatalf("expected fullyFundedThreshold %v but got %v",
			dfo.fullyFundedThreshold, got.fullyFundedThreshold)
	}
	if !(got.Status == dfo.Status) {
		t.Fatalf("expected Status %v but got %v", dfo.Status, got.Status)
	}
	if got.C.Id != dfo.C.Id {
		t.Fatalf("expected channel Id %s but got %s", dfo.C.Id, got.C.Id)
	}
}
