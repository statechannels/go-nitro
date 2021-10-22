package protocols

import (
	"crypto/rand"
	"time"

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
	time.Sleep(200 * time.Millisecond) // TODO consider choosing a time that roughly matches ECDSA hash signature time
	b := make([]byte, 100)
	rand.Read(b)
	state.SignEthereumMessage(b, *secretKey)
	return t, SideEffects{}, ``, nil
}
