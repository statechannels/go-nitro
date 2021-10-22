package protocols

import (
	"crypto/rand"
	"fmt"

	"github.com/statechannels/go-nitro/channel/state"
)

// TestObjectives simulates some computationally intensive task when cranked
type TestObjective struct {
}

func (t TestObjective) Id() ObjectiveId {
	return `test`
}
func (t TestObjective) Approve() Objective {
	return t
}
func (t TestObjective) Reject() Objective                     { return t }
func (t TestObjective) Update(event ObjectiveEvent) Objective { return t }
func (t TestObjective) Crank(secretKey *[]byte) (Objective, SideEffects, WaitingFor, error) {
	b := make([]byte, 100)
	rand.Read(b)
	state.SignEthereumMessage(b, *secretKey)
	fmt.Println(`signed a message`)
	return t, SideEffects{}, ``, nil
}
