package policy

import (
	"fmt"
	"log/slog"
	"math/big"

	"github.com/statechannels/go-nitro/internal/logging"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

// PolicyMaker is used to decide whether to approve or reject an objective
type PolicyMaker interface {
	ShouldApprove(o protocols.Objective) bool
}

type Policies struct {
	policies []PolicyMaker
}

func NewPolicies(policies ...PolicyMaker) *Policies {
	return &Policies{policies: policies}
}

// ShouldApprove decides to approve o if all policies decide to approve it
func (p *Policies) ShouldApprove(o protocols.Objective) bool {
	for _, policy := range p.policies {

		approves := policy.ShouldApprove(o)
		slog.Debug("Policymaker decision", "policy-maker", fmt.Sprintf("%T", policy), "approves", approves, logging.WithObjectiveIdAttribute(o.Id()))

		if !approves {
			return false
		}
	}

	return true
}

// PermissivePolicy is a policy maker that decides to approve every unapproved objective
type PermissivePolicy struct{}

// ShouldApprove decides to approve o if it is currently unapproved
func (pp *PermissivePolicy) ShouldApprove(o protocols.Objective) bool {
	return o.GetStatus() == protocols.Unapproved
}

type AllowListPolicy struct {
	allowed map[types.Address]bool
}

func NewAllowListPolicy(allowed []types.Address) *AllowListPolicy {
	allowedMap := make(map[types.Address]bool)
	for _, a := range allowed {
		allowedMap[a] = true
	}
	return &AllowListPolicy{allowed: allowedMap}
}

// ShouldApprove decides to approve o if it is currently unapproved
func (ap *AllowListPolicy) ShouldApprove(o protocols.Objective) bool {
	for _, p := range o.GetParticipants() {
		if !ap.allowed[p] {
			return false
		}
	}
	return true
}

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
					return false
				} else if i > 0 && a.Amount.Cmp(big.NewInt(0)) != 0 {
					return false
				}
			}
		}
	}

	return true
}
