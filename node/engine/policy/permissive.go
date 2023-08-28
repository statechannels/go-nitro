package policy

import "github.com/statechannels/go-nitro/protocols"

// PermissivePolicy is a policy maker that decides to approve every unapproved objective
type PermissivePolicy struct{}

// ShouldApprove decides to approve o if it is currently unapproved
func (pp *PermissivePolicy) ShouldApprove(o protocols.Objective) bool {
	return o.GetStatus() == protocols.Unapproved
}
