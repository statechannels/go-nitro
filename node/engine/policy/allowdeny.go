package policy

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func NewDenyListPolicy(denied []types.Address) PolicyMaker {
	deniedMap := make(map[types.Address]bool)
	for _, a := range denied {
		deniedMap[a] = false
	}
	return &listPolicy{denied: deniedMap}
}

func NewAllowListPolicy(allowed []types.Address) PolicyMaker {
	allowedMap := make(map[types.Address]bool)
	for _, a := range allowed {
		allowedMap[a] = false
	}
	return &listPolicy{allowed: allowedMap}
}

// ShouldApprove decides to approve o if it is currently unapproved
func (lp *listPolicy) ShouldApprove(o protocols.Objective) bool {
	for _, p := range o.GetParticipants() {
		if !lp.allowed[p] {
			return false
		}
	}
	return true
}

type listPolicy struct {
	allowed map[types.Address]bool
	denied  map[types.Address]bool
}
