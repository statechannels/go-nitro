package policy

import "github.com/statechannels/go-nitro/protocols"

// permissivePolicy is a policy maker that decides to approve every unapproved objective
type permissivePolicy struct{}

// ShouldApprove decides to approve o if it is currently unapproved
func (pp *permissivePolicy) ShouldApprove(o protocols.Objective) bool {
	return o.GetStatus() == protocols.Unapproved
}

// NewPermissivePolicy returns a new PermissivePolicy
func NewPermissivePolicy() PolicyMaker {
	return &permissivePolicy{}
}
