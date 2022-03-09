package directfund

import (
	"errors"
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

type actor struct {
	address     types.Address
	destination types.Destination
	privateKey  []byte
}

var alicePK = common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`)
var bobPK = common.Hex2Bytes(`62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2`)

var alice = actor{
	address:     common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`),
	destination: types.AddressToDestination(common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`)),
	privateKey:  alicePK,
}

var bob = actor{
	address:     common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`),
	destination: types.AddressToDestination(common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`)),
	privateKey:  bobPK,
}

var testState = state.State{
	ChainId:           big.NewInt(9001),
	Participants:      []types.Address{alice.address, bob.address},
	ChannelNonce:      big.NewInt(37140676580),
	AppDefinition:     common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`),
	ChallengeDuration: big.NewInt(60),
	AppData:           []byte{},
	Outcome: outcome.Exit{
		outcome.SingleAssetExit{
			Asset: types.Address{},
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: bob.destination, // Bob is first so we can easily test WaitingForMyTurnToFund
					Amount:      big.NewInt(5),
				},
				outcome.Allocation{
					Destination: alice.destination,
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
	pks := [2][]byte{alicePK, bobPK}
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
	c.OnChainFunding = map[common.Address]*big.Int{common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`): big.NewInt(1)}
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
		if err != nil {
			t.Error(err)
		}
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

	s := testState.Clone()
	s.TurnNum = 3
	e.ObjectiveId = o.Id()
	ss, _ := signedTestState(s, []bool{true, false})
	e.SignedStates = []state.SignedState{ss}

	if _, err := o.Update(e); err == nil {
		t.Error("expected an error when updating with a non-final state")
	}
}

func TestCrankAlice(t *testing.T) {
	o, _ := newTestObjective(true)

	// The first crank. Alice is expected to create and sign a final state
	updated, se, wf, _, err := o.Crank(&alicePK)

	if err != nil {
		t.Error(err)
	}

	if wf != WaitingForFinalization {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForFinalization, wf)
	}

	// Create the state we expect Alice to send
	finalState := testState.Clone()
	finalState.TurnNum = 3
	finalState.IsFinal = true
	finalStateSignedByAlice, _ := signedTestState(finalState, []bool{true, false})

	expectedSE := protocols.SideEffects{
		MessagesToSend: []protocols.Message{{
			To:          bob.address,
			ObjectiveId: o.Id(),
			SignedStates: []state.SignedState{
				finalStateSignedByAlice,
			},
		},
		}}

	if diff := cmp.Diff(expectedSE, se); diff != "" {
		t.Errorf("Side effects mismatch (-want +got):\n%s", diff)
	}

	// The second update and crank. Alice is expected to create a withdrawAll transaction
	finalStateSignedByAliceBob, _ := signedTestState(finalState, []bool{true, true})
	e := protocols.ObjectiveEvent{ObjectiveId: o.Id(), SignedStates: []state.SignedState{finalStateSignedByAliceBob}}

	updated, err = updated.Update(e)
	if err != nil {
		t.Error(err)
	}
	_, se, wf, _, err = updated.Crank(&alicePK)
	if err != nil {
		t.Error(err)
	}

	if wf != WaitingForWithdraw {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForWithdraw, wf)
	}

	expectedSE = protocols.SideEffects{TransactionsToSubmit: []protocols.ChainTransaction{{
		Type:      protocols.WithdrawAllTransactionType,
		ChannelId: o.C.Id,
		Deposit:   types.Funds{},
	}}}

	if diff := cmp.Diff(expectedSE, se); diff != "" {
		t.Errorf("Side effects mismatch (-want +got):\n%s", diff)
	}

	// The third crank. Alice is expected to enter the terminal state of the defunding protocol.
	updated.C.OnChainFunding = types.Funds{}
	_, se, wf, _, err = updated.Crank(&alicePK)
	if err != nil {
		t.Error(err)
	}
	if wf != WaitingForNothing {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForWithdraw, wf)
	}

	expectedSE = protocols.SideEffects{}
	if diff := cmp.Diff(expectedSE, se); diff != "" {
		t.Errorf("Side effects mismatch (-want +got):\n%s", diff)
	}
}

func TestCrankBob(t *testing.T) {
	o, _ := newTestObjective(true)
	o.C.MyIndex = 1

	// The first crank. Bob is expected to create and sign a final state
	updated, se, wf, _, err := o.Crank(&bobPK)

	if err != nil {
		t.Error(err)
	}

	if wf != WaitingForFinalization {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForFinalization, wf)
	}

	// Create the state we expect Bob to send
	finalState := testState.Clone()
	finalState.TurnNum = 3
	finalState.IsFinal = true
	finalStateSignedByBob, _ := signedTestState(finalState, []bool{false, true})

	expectedSE := protocols.SideEffects{
		MessagesToSend: []protocols.Message{{
			To:          alice.address,
			ObjectiveId: o.Id(),
			SignedStates: []state.SignedState{
				finalStateSignedByBob,
			},
		},
		}}

	if diff := cmp.Diff(expectedSE, se); diff != "" {
		t.Errorf("Side effects mismatch (-want +got):\n%s", diff)
	}

	// The second update and crank. Bob is expected to NOT create any transactions
	finalStateSignedByAliceBob, _ := signedTestState(finalState, []bool{true, true})
	e := protocols.ObjectiveEvent{ObjectiveId: o.Id(), SignedStates: []state.SignedState{finalStateSignedByAliceBob}}

	updated, err = updated.Update(e)
	if err != nil {
		t.Error(err)
	}
	_, se, wf, _, err = updated.Crank(&alicePK)
	if err != nil {
		t.Error(err)
	}

	if wf != WaitingForWithdraw {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForWithdraw, wf)
	}

	expectedSE = protocols.SideEffects{}

	if diff := cmp.Diff(expectedSE, se); diff != "" {
		t.Errorf("Side effects mismatch (-want +got):\n%s", diff)
	}

	// The third crank. Bob is expected to enter the terminal state of the defunding protocol.
	updated.C.OnChainFunding = types.Funds{}
	_, se, wf, _, err = updated.Crank(&alicePK)
	if err != nil {
		t.Error(err)
	}
	if wf != WaitingForNothing {
		t.Errorf(`WaitingFor: expected %v, got %v`, WaitingForWithdraw, wf)
	}

	expectedSE = protocols.SideEffects{}
	if diff := cmp.Diff(expectedSE, se); diff != "" {
		t.Errorf("Side effects mismatch (-want +got):\n%s", diff)
	}
}
