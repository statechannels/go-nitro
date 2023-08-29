package policy

import (
	"fmt"
	"log/slog"

	"github.com/statechannels/go-nitro/internal/logging"
	"github.com/statechannels/go-nitro/protocols"
)

// PolicyMaker is used to decide whether to approve or reject an objective
type PolicyMaker interface {
	ShouldApprove(o protocols.Objective) bool
}

type Policies []PolicyMaker

func NewPolicies(policies ...PolicyMaker) *Policies {
	p := Policies(policies)
	return &p
}

// ShouldApprove decides to approve o if all policies decide to approve it
func (p *Policies) ShouldApprove(o protocols.Objective) bool {
	for _, policy := range *p {

		approves := policy.ShouldApprove(o)
		slog.Debug("Policymaker decision", "policy-maker", fmt.Sprintf("%T", policy), "approves", approves, logging.WithObjectiveIdAttribute(o.Id()))

		if !approves {
			return false
		}
	}

	return true
}
