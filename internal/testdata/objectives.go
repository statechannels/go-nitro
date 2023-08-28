package testdata

import (
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

const TEST_CHAIN_ID = 1337

// objectiveCollection namespaces literal objectives, precomputed objectives, and
// procedural objective generators for consumption
type objectiveCollection struct {
	Directfund  dfoCollection
	Virtualfund vfoCollection
	// todo
	// directdefund ddfoCollection
}

type dfoCollection struct {
	// GenericDFO returns a non-specific directfund.Objective with nonzero data.
	GenericDFO func() directfund.Objective
}

type vfoCollection struct {
	// GenericVFO returns a non-specific virtualfund.Objective with nonzero data.
	GenericVFO func() virtualfund.Objective
}

// Objectives is the endpoint for tests to consume constructed objectives or
// objective generating utility functions
//
// eg, a test wanting an Irene-Ivan ledger creation objective could import via
//
//	testdata.Objectives.twopartyledgers.irene_ivan
var Objectives objectiveCollection = objectiveCollection{
	Directfund: dfoCollection{
		GenericDFO: genericDFO,
	},
	Virtualfund: vfoCollection{
		GenericVFO: genericVFO,
	},
}

func genericDFO() directfund.Objective {
	ts := testState.Clone()
	return GenerateDFOFromOutcome(ts.Outcome)
}

func GenerateDFOFromOutcome(o outcome.Exit) directfund.Objective {
	ts := testState.Clone()
	ts.TurnNum = 0
	ts.Outcome = o.Clone()
	ss := state.NewSignedState(ts)
	id := protocols.ObjectiveId(directfund.ObjectivePrefix + testState.ChannelId().String())
	op, err := protocols.CreateObjectivePayload(id, directfund.SignedStatePayload, ss)
	if err != nil {
		panic(fmt.Errorf("error constructing objective payload: %w", err))
	}
	testObj, err := directfund.ConstructFromPayload(false, op, ts.Participants[0])
	if err != nil {
		panic(fmt.Errorf("error constructing genericDFO: %w", err))
	}
	return testObj
}

func GenerateVFOFromOutcome(o outcome.Exit) virtualfund.Objective {
	ts := testVirtualState.Clone()
	ts.Outcome = o.Clone()
	ts.Participants[0] = testactors.Alice.Address()
	ts.Participants[1] = testactors.Irene.Address()
	ts.Participants[2] = testactors.Bob.Address()

	request := virtualfund.NewObjectiveRequest(
		[]types.Address{ts.Participants[1]},
		ts.Participants[2],
		ts.ChallengeDuration,
		o,
		ts.ChannelNonce,
		ts.AppDefinition,
	)
	ledgerPath := createLedgerPath([]testactors.Actor{
		testactors.Alice,
		testactors.Irene,
		testactors.Bob,
	})
	lookup := ledgerPath.GetLedgerLookup(testactors.Alice.Address())

	testVFO, err := virtualfund.NewObjective(request, true, ts.Participants[0], big.NewInt(TEST_CHAIN_ID), lookup)
	if err != nil {
		panic(fmt.Errorf("error constructing genericVFO: %w", err))
	}
	return testVFO
}

func genericVFO() virtualfund.Objective {
	ts := testVirtualState.Clone()
	return GenerateVFOFromOutcome(ts.Outcome)
}
