package policy

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

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
