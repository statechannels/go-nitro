package directdefund

import (
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var alice, bob testactors.Actor = testactors.Alice, testactors.Bob

var testState = state.State{
	ChainId:           big.NewInt(9001),
	Participants:      []types.Address{alice.Address, bob.Address},
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

// newChannelFromSignedState constructs a new Channel from the signed state.
func newChannelFromSignedState(ss state.SignedState, myIndex uint) (*channel.Channel, error) {
	s := ss.State()
	prefund := s.Clone()
	prefund.TurnNum = 0
	c, err := channel.New(prefund, myIndex)
	if err != nil {
		return c, err
	}

	sss := make([]state.SignedState, 1)
	sss[0] = ss
	allOk := c.AddSignedStates(sss)
	if !allOk {
		return c, errors.New("Unable to add a state to channel")
	}
	c.OnChainFunding = map[common.Address]*big.Int{types.Address{}: big.NewInt(1)}
	return c, nil
}

func newTestObjective(signByBob bool) (Objective, error) {
	o := Objective{}

	toSign := []bool{true, true}
	if !signByBob {
		toSign = []bool{true, false}
	}
	ss, err := signedTestState(testState, toSign)
	if err != nil {
		return o, err
	}

	testChannel, err := newChannelFromSignedState(ss, 0)
	if err != nil {
		return o, err
	}

	// Assert that valid constructor args do not result in error
	o, err = NewObjective(true, testChannel)
	if err != nil {
		return o, err
	}
	return o, nil
}

// TestNew tests the constructor using a TestState fixture
func TestNew(t *testing.T) {
	if _, err := newTestObjective(false); err == nil {
		t.Error("expected an error constructing the defund objective from a state without signatures from all participant, but got nil")
	}

	if _, err := newTestObjective(true); err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	o, _ := newTestObjective(true)

	// Prepare an event with a mismatched channelId
	e := protocols.ObjectiveEvent{
		ObjectiveId: "some-id",
	}
	// Assert that Updating the objective with such an event returns an error
	if _, err := o.Update(e); err == nil {
		t.Error(`ChannelId mismatch -- expected an error but did not get one`)
	}

	// Try updating the objective with a non-final state. An error is expected.
	s := testState.Clone()
	s.TurnNum = 3
	e.ObjectiveId = o.Id()
	ss, _ := signedTestState(s, []bool{true, false})
	e.SignedStates = []state.SignedState{ss}

	if _, err := o.Update(e); err.Error() != "direct defund objective can only be updated with final states" {
		t.Error(err)
	}

	// Try updating the objective with a final state with the wrong turn number. An error is expected.
	s.TurnNum = 4
	s.IsFinal = true
	ss, _ = signedTestState(s, []bool{true, false})
	e.SignedStates = []state.SignedState{ss}

	if _, err := o.Update(e); err.Error() != "expected state with turn number 3, received turn number 4" {
		t.Error(err)
	}
}

func compareSideEffect(a, b protocols.SideEffects) string {
	return cmp.Diff(a, b, cmp.AllowUnexported(a, state.SignedState{}))
}

func TestCrankAlice(t *testing.T) {
	// The starting channel state is:
	//  - Channel has a non-final consensus state
	//  - Channel has funds
	o, _ := newTestObjective(true)

	// The first crank. Alice is expected to create and sign a final state
	updated, se, wf, err := o.Crank(&alice.PrivateKey)

	if err != nil {
		t.Error(err)
	}

	if wf != WaitingForFinalization {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForFinalization, wf)
	}

	// Create the state we expect Alice to send
	finalState := testState.Clone()
	finalState.TurnNum = 3
	finalState.IsFinal = true
	finalStateSignedByAlice, _ := signedTestState(finalState, []bool{true, false})

	expectedSE := protocols.SideEffects{
		MessagesToSend: []protocols.Message{{
			To:          bob.Address,
			ObjectiveId: o.Id(),
			SignedStates: []state.SignedState{
				finalStateSignedByAlice,
			},
			SignedProposals: []consensus_channel.SignedProposal{},
		},
		}}

	if diff := compareSideEffect(expectedSE, se); diff != "" {
		t.Errorf("Side effects mismatch (-want +got):\n%s", diff)
	}

	// The second update and crank. Alice is expected to create a withdrawAll transaction
	finalStateSignedByAliceBob, _ := signedTestState(finalState, []bool{true, true})
	e := protocols.ObjectiveEvent{ObjectiveId: o.Id(), SignedStates: []state.SignedState{finalStateSignedByAliceBob}}

	updated, err = updated.Update(e)
	if err != nil {
		t.Error(err)
	}
	_, se, wf, err = updated.Crank(&alice.PrivateKey)
	if err != nil {
		t.Error(err)
	}

	if wf != WaitingForWithdraw {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForWithdraw, wf)
	}

	expectedSE = protocols.SideEffects{TransactionsToSubmit: []protocols.ChainTransaction{{
		Type:      protocols.WithdrawAllTransactionType,
		ChannelId: o.C.Id,
		Deposit:   types.Funds{},
	}}}

	if diff := cmp.Diff(expectedSE, se); diff != "" {
		t.Fatalf("Side effects mismatch (-want +got):\n%s", diff)
	}

	// The third crank. Alice is expected to enter the terminal state of the defunding protocol.
	updated.C.OnChainFunding = types.Funds{}
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

	o, _ := newTestObjective(true)
	o.C.MyIndex = 1

	// Update the objective with Alice's final state
	finalState := testState.Clone()
	finalState.TurnNum = 3
	finalState.IsFinal = true
	finalStateSignedByAlice, _ := signedTestState(finalState, []bool{true, false})
	e := protocols.ObjectiveEvent{ObjectiveId: o.Id(), SignedStates: []state.SignedState{finalStateSignedByAlice}}
	o, err := o.Update(e)
	if err != nil {
		t.Error(err)
	}

	// The first crank. Bob is expected to create and sign a final state
	o, se, wf, err := o.Crank(&bob.PrivateKey)

	if err != nil {
		t.Error(err)
	}

	if wf != WaitingForWithdraw {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForWithdraw, wf)
	}

	// Create the state we expect Bob to send
	finalStateSignedByBob, _ := signedTestState(finalState, []bool{false, true})
	expectedSE := protocols.SideEffects{
		MessagesToSend: []protocols.Message{{
			To:          alice.Address,
			ObjectiveId: o.Id(),
			SignedStates: []state.SignedState{
				finalStateSignedByBob,
			},
			SignedProposals: []consensus_channel.SignedProposal{},
		},
		}}

	if diff := compareSideEffect(expectedSE, se); diff != "" {
		t.Errorf("Side effects mismatch (-want +got):\n%s", diff)
	}

	// The second update and crank. Bob is expected to NOT create any transactions or side effects
	o, err = o.Update(e)
	if err != nil {
		t.Error(err)
	}
	_, se, wf, err = o.Crank(&bob.PrivateKey)
	if err != nil {
		t.Error(err)
	}

	if wf != WaitingForWithdraw {
		t.Fatalf(`WaitingFor: expected %v, got %v`, WaitingForWithdraw, wf)
	}

	expectedSE = protocols.SideEffects{}

	if diff := cmp.Diff(expectedSE, se); diff != "" {
		t.Fatalf("Side effects mismatch (-want +got):\n%s", diff)
	}

	// The third crank. Bob is expected to enter the terminal state of the defunding protocol.
	o, err = o.UpdateWithChainEvent(chainservice.DepositedEvent{Holdings: types.Funds{}})

	if err != nil {
		t.Error(err)
	}

	_, se, wf, err = o.Crank(&bob.PrivateKey)
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
	ddfo, _ := newTestObjective(true)

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
