package testdata

import (
	"fmt"

	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
)

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
//     testdata.Objectives.twopartyledgers.irene_ivan
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
	ts.TurnNum = 0
	testObj, err := directfund.ConstructFromState(false, ts, ts.Participants[0])
	if err != nil {
		panic(fmt.Errorf("error constructing genericDFO: %w", err))
	}
	return testObj
}

func genericVFO() virtualfund.Objective {
	ts := testVirtualState.Clone()
	ts.Participants[0] = testactors.Alice.Address()
	ts.Participants[1] = testactors.Irene.Address()
	ts.Participants[2] = testactors.Bob.Address()

	request := virtualfund.ObjectiveRequest{
		ts.Participants[0],
		ts.Participants[1],
		ts.Participants[2],
		ts.AppDefinition,
		ts.AppData,
		ts.ChallengeDuration,
		ts.Outcome,
		ts.ChannelNonce.Int64(),
	}

	ledgerPath := createLedgerPath([]testactors.Actor{
		testactors.Alice,
		testactors.Irene,
		testactors.Bob,
	})
	lookup := ledgerPath.GetLedgerLookup(testactors.Alice.Address())

	testVFO, err := virtualfund.NewObjective(request, lookup)
	if err != nil {
		panic(fmt.Errorf("error constructing genericVFO: %w", err))
	}
	return testVFO
}
