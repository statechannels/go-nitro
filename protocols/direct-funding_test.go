// TODO:
// - prepare a bunch of initial states
// - crank them
package protocols

func assertObjective(o Objective) {
	var dfObjective DirectFundingObjectiveState

	approved := dfObjective.Approve()

	assertObjective(approved)
}
