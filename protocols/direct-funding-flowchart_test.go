// TODO:
// - prepare a bunch of initial states
// - crank them
// - assert on the sideffects and waitingFor string
package protocols

func assertObjective(o Objective) {
	var dfObjective DirectFundingObjectiveState

	approved := dfObjective.Approve()

	assertObjective(approved)
}
