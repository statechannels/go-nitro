package engine

import "github.com/statechannels/go-nitro/protocols"

// PolicyMaker is used to decide whether to approve or reject an objective
type PolicyMaker interface {
	ShouldApprove(o protocols.Objective) bool
}

// PermissivePolicy is a policy maker that decides to approve every unapproved objective
type PermissivePolicy struct{}

// ShouldApprove decides to approve o if it is currently unapproved
func (pp *PermissivePolicy) ShouldApprove(o protocols.Objective) bool {
	return o.GetStatus() == protocols.Unapproved
}
