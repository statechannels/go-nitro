package policy

import (
	"log/slog"
	"math/big"

	"github.com/statechannels/go-nitro/internal/logging"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

type FairOutcomePolicy struct {
	me types.Address
}

func NewFairOutcomePolicy(me types.Address) *FairOutcomePolicy {
	return &FairOutcomePolicy{me: me}
}

// ShouldApprove decides to approve o if it is currently unapproved
func (fp *FairOutcomePolicy) ShouldApprove(o protocols.Objective) bool {
	df, isDf := o.(*directfund.Objective)
	if isDf {
		for _, e := range df.C.PreFundState().Outcome {
			forMe := e.TotalAllocatedFor(types.AddressToDestination(fp.me))
			for _, a := range e.Allocations {
				if a.Amount.Cmp(forMe) != 0 {
					slog.Warn("FairOutcomePolicy: rejecting directfund objective with unequal allocations", "objective", o.Id(), "allocations", e.Allocations, logging.WithObjectiveIdAttribute(o.Id()))
					return false
				}
			}
		}
	}
	vf, isVf := o.(*virtualfund.Objective)
	if isVf {
		for _, e := range vf.V.PreFundState().Outcome {

			total := e.TotalAllocated()
			for i, a := range e.Allocations {
				if i == 0 && a.Amount.Cmp(total) != 0 {
					slog.Warn("FairOutcomePolicy: rejecting virtualfund objective, expected payer to start with full amount", "allocations", e.Allocations, logging.WithObjectiveIdAttribute(o.Id()))
					return false
				} else if i > 0 && a.Amount.Cmp(big.NewInt(0)) != 0 {
					slog.Warn("FairOutcomePolicy: rejecting virtualfund objective, expected payee to start with 0", "allocations", e.Allocations, logging.WithObjectiveIdAttribute(o.Id()))
					return false
				}
			}
		}
	}

	return true
}
