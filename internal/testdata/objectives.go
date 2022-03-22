package td

import "github.com/statechannels/go-nitro/protocols/directfund"

// objectiveCollection namespaces literal objectives, precomputed objectives, and
// procedural objective generators for consumption
type objectiveCollection struct {
	Directfund dfoCollection
	// todo
	// virtualfund  vfoCollection
	// todo
	// directdefund ddfoCollection
}

type dfoCollection struct {
	// GenericDFO returns a non-specific directfund.Objective with nonzero data.
	GenericDFO func() directfund.Objective
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
}

func genericDFO() directfund.Objective {
	ts := testState
	request := directfund.ObjectiveRequest{
		MyAddress:         ts.Participants[0],
		CounterParty:      ts.Participants[1],
		AppData:           ts.AppData,
		AppDefinition:     ts.AppDefinition,
		ChallengeDuration: ts.ChallengeDuration,
		Nonce:             ts.ChannelNonce.Int64(),
		Outcome:           ts.Outcome,
	}
	testObj, _ := directfund.NewObjective(request, false)
	return testObj
}
