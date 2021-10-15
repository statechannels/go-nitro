package protocols

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
)

// TODO these are placeholders for now
type SideEffects []interface{}
type WaitingFor string

// Protocol is the interface for off-chain protocols
type Protocol interface {
	Initialize(initialState state.State) Protocol // returns the initial Protocol object, does not declare effects

	SignatureRecieved(sig state.Signature) Protocol // returns an updated Protocol (a copy, no mutation allowed), does not declare effects
	DepositObserved(holding *big.Int) Protocol      // returns an updated Protocol (a copy, no mutation allowed), does not declare effects
	// other protocol specific events: returns an updated Protocol (a copy, no mutation allowed), does not declare effects

	Crank() (SideEffects, WaitingFor, error) // does *not* accept an event, but *does* declare side effects, does *not* return an updated Protocol

}
