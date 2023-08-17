package directdefund

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/internal/testhelpers"
	"github.com/statechannels/go-nitro/node/engine/chainservice"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var alice, bob testactors.Actor = testactors.Alice, testactors.Bob

var testState = state.State{
	Participants:      []types.Address{alice.Address(), bob.Address()},
	ChannelNonce:      37140676580,
	AppDefinition:     common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`),
	ChallengeDuration: 60,
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
	TurnNum: 2,
	IsFinal: false,
}

// signedTestState returns a signed state with signatures requested in toSign
func signedTestState(s state.State, toSign []bool) (state.SignedState, error) {
	ss := state.NewSignedState(s)
	pks := [2][]byte{alice.PrivateKey, bob.PrivateKey}
	for i, pk := range pks {
		if !toSign[i] {
			continue
		}

		sig, err := ss.State().Sign(pk)
		if err != nil {
			return ss, err
		}
		err = ss.AddSignature(sig)
		if err != nil {
			return ss, err
		}
	}
	return ss, nil
}

// newTestObjective returns a directdefund Objective constructed with a MockConsensusChannel.
func newTestObjective() (Objective, error) {
	cc, _ := testdata.Channels.MockConsensusChannel(alice.Address())

	getConsensusChannel := func(id types.Destination) (channel *consensus_channel.ConsensusChannel, err error) {
		return cc, nil
	}
	request := NewObjectiveRequest(cc.Id)
	// Assert that valid constructor args do not result in error
	o, err := NewObjective(request, true, getConsensusChannel)
	if err != nil {
		return Objective{}, err
	}
	return o, nil
}

// TestNew tests the constructor using a TestState fixture
func TestNew(t *testing.T) {
	if _, err := newTestObjective(); err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	o, _ := newTestObjective()

	// Prepare an event with a mismatched channelId
	op, err := protocols.CreateObjectivePayload("some-id", SignedStatePayload, testState)
	testhelpers.Ok(t, err)

	// Assert that Updating the objective with such an event returns an error
	if _, err := o.Update(op); err == nil {
		t.Error(`ChannelId mismatch -- expected an error but did not get one`)
	}

	// Try updating the objective with a non-final state. An error is expected.
	s := testState.Clone()
	s.TurnNum = 3

	ss, _ := signedTestState(s, []bool{true, false})
	op, err = protocols.CreateObjectivePayload(o.Id(), SignedStatePayload, ss)
	testhelpers.Ok(t, err)

	if _, err := o.Update(op); err.Error() != "direct defund objective can only be updated with final states" {
		t.Error(err)
	}

	// Try updating the objective with a final state with the wrong turn number. An error is expected.
	s.TurnNum = 4
	s.IsFinal = true
	ss, _ = signedTestState(s, []bool{true, false})
	op, err = protocols.CreateObjectivePayload(o.Id(), SignedStatePayload, ss)
	testhelpers.Ok(t, err)

	if _, err := o.Update(op); err.Error() != "expected state with turn number 2, received turn number 4" {
		t.Error(err)
	}
}

func compareSideEffect(a, b protocols.SideEffects) string {
	return cmp.Diff(a, b, cmp.AllowUnexported(a, state.SignedState{}, consensus_channel.Add{}, consensus_channel.Remove{}, consensus_channel.Guarantee{}, protocols.Message{}, payments.Voucher{}))
}

func TestCrankAlice(t *testing.T) {
	// The starting channel state is:
	//  - Channel has a non-final consensus state
	//  - Channel has funds
	o, _ := newTestObjective()

	// The first crank. Alice is expected to create and sign a final state
	updated, se, wf, err := o.Crank(&alice.PrivateKey)
	if err != nil {
		t.Error(err)
	}

	if wf != WaitingForFinalization {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForFinalization, wf)
	}

	// Create the state we expect Alice to send
	finalState, err := o.C.LatestSupportedState()
	if err != nil {
		t.Fatal(err)
	}
	finalState.TurnNum = 2
	finalState.IsFinal = true
	finalStateSignedByAlice, _ := signedTestState(finalState, []bool{true, false})

	msgs, err := protocols.CreateObjectivePayloadMessage(o.Id(), finalStateSignedByAlice, SignedStatePayload, o.otherParticipants()...)
	testhelpers.Ok(t, err)

	expectedSE := protocols.SideEffects{MessagesToSend: msgs}

	if diff := compareSideEffect(expectedSE, se); diff != "" {
		t.Errorf("Side effects mismatch (-want +got):\n%s", diff)
	}

	// The second update and crank. Alice is expected to create a withdrawAll transaction
	finalStateSignedByAliceBob, _ := signedTestState(finalState, []bool{true, true})
	op, err := protocols.CreateObjectivePayload(o.Id(), SignedStatePayload, finalStateSignedByAliceBob)
	testhelpers.Ok(t, err)

	updated, err = updated.Update(op)
	if err != nil {
		t.Error(err)
	}
	updated, se, wf, err = updated.Crank(&alice.PrivateKey)
	if err != nil {
		t.Error(err)
	}
	if !updated.(*Objective).withdrawTransactionSubmitted {
		t.Fatalf("Expected transactionSubmitted flag to be set to true")
	}

	if wf != WaitingForWithdraw {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForWithdraw, wf)
	}

	expectedSE = protocols.SideEffects{TransactionsToSubmit: []protocols.ChainTransaction{protocols.NewWithdrawAllTransaction(o.C.Id, finalStateSignedByAliceBob)}}

	if diff := cmp.Diff(expectedSE, se, cmp.AllowUnexported(expectedSE, state.SignedState{}, protocols.ChainTransactionBase{})); diff != "" {
		t.Fatalf("Side effects mismatch (-want +got):\n%s", diff)
	}

	// The third crank. Alice is expected to enter the terminal state of the defunding protocol.
	updated.(*Objective).C.OnChainFunding = types.Funds{}
	_, se, wf, err = updated.Crank(&alice.PrivateKey)
	if err != nil {
		t.Error(err)
	}
	if wf != WaitingForNothing {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForNothing, wf)
	}

	expectedSE = protocols.SideEffects{}
	if diff := cmp.Diff(expectedSE, se); diff != "" {
		t.Fatalf("Side effects mismatch (-want +got):\n%s", diff)
	}
}

func TestCrankBob(t *testing.T) {
	// The starting channel state is:
	//  - Channel has a non-final non-consensus state
	//  - Channel has funds

	o, _ := newTestObjective()
	o.C.MyIndex = 1

	// Update the objective with Alice's final state
	finalState := testState.Clone()
	finalState.TurnNum = 2
	finalState.IsFinal = true
	finalStateSignedByAlice, _ := signedTestState(finalState, []bool{true, false})

	op, err := protocols.CreateObjectivePayload(o.Id(), SignedStatePayload, finalStateSignedByAlice)
	testhelpers.Ok(t, err)

	updated, err := o.Update(op)
	if err != nil {
		t.Fatal(err)
	}

	// The first crank. Bob is expected to create and sign a final state
	updated, se, wf, err := updated.Crank(&bob.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	if wf != WaitingForWithdraw {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForWithdraw, wf)
	}

	// Create the state we expect Bob to send
	finalStateSignedByBob, _ := signedTestState(finalState, []bool{false, true})

	msgs, err := protocols.CreateObjectivePayloadMessage(updated.Id(), finalStateSignedByBob, SignedStatePayload, o.otherParticipants()...)
	testhelpers.Ok(t, err)
	expectedSE := protocols.SideEffects{MessagesToSend: msgs}

	if diff := compareSideEffect(expectedSE, se); diff != "" {
		t.Errorf("Side effects mismatch (-want +got):\n%s", diff)
	}

	// The second update and crank. Bob is expected to NOT create any transactions or side effects
	updated, err = updated.Update(op)
	if err != nil {
		t.Error(err)
	}
	updated, se, wf, err = updated.Crank(&bob.PrivateKey)
	if err != nil {
		t.Error(err)
	}

	if updated.(*Objective).withdrawTransactionSubmitted {
		t.Fatalf("Expected transactionSubmitted flag to be set to false")
	}

	if wf != WaitingForWithdraw {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForWithdraw, wf)
	}

	expectedSE = protocols.SideEffects{}

	if diff := cmp.Diff(expectedSE, se); diff != "" {
		t.Fatalf("Side effects mismatch (-want +got):\n%s", diff)
	}

	// The third crank. Bob is expected to enter the terminal state of the defunding protocol.
	_, err = updated.(*Objective).C.UpdateWithChainEvent(chainservice.NewAllocationUpdatedEvent(types.Destination{}, 1, common.Address{}, common.Big0))

	if err != nil {
		t.Error(err)
	}

	_, se, wf, err = updated.Crank(&bob.PrivateKey)
	if err != nil {
		t.Error(err)
	}
	if wf != WaitingForNothing {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForNothing, wf)
	}

	expectedSE = protocols.SideEffects{}
	if diff := cmp.Diff(expectedSE, se); diff != "" {
		t.Fatalf("Side effects mismatch (-want +got):\n%s", diff)
	}
}

func TestMarshalJSON(t *testing.T) {
	ddfo, _ := newTestObjective()

	encodedDdfo, err := json.Marshal(ddfo)
	if err != nil {
		t.Fatalf("error encoding directdefund objective %v", ddfo)
	}

	got := Objective{}
	if err := got.UnmarshalJSON(encodedDdfo); err != nil {
		t.Fatalf("error unmarshaling test directdefund objective: %s", err.Error())
	}

	if got.finalTurnNum != ddfo.finalTurnNum {
		t.Fatalf("expected finalTurnNum %d but got %d",
			ddfo.finalTurnNum, got.finalTurnNum)
	}
	if !(got.Status == ddfo.Status) {
		t.Fatalf("expected Status %v but got %v", ddfo.Status, got.Status)
	}
	if got.C.Id != ddfo.C.Id {
		t.Fatalf("expected channel Id %s but got %s", ddfo.C.Id, got.C.Id)
	}
}

func TestApproveReject(t *testing.T) {
	o, err := newTestObjective()
	testhelpers.Ok(t, err)

	approved := o.Approve()
	if approved.GetStatus() != protocols.Approved {
		t.Errorf("Expected approved status, got %v", approved.GetStatus())
	}
	rejected, sideEffects := o.Reject()
	if rejected.GetStatus() != protocols.Rejected {
		t.Errorf("Expected rejceted status, got %v", approved.GetStatus())
	}
	if len(sideEffects.MessagesToSend) != 1 {
		t.Errorf("Expected to send one message")
	}
}
