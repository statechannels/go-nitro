package policy

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func NewDenyListPolicy(denied []types.Address) PolicyMaker {
	deniedMap := make(map[types.Address]struct{})
	for _, a := range denied {
		deniedMap[a] = struct{}{}
	}
	return &listPolicy{participants: deniedMap, mode: deny}
}

func NewAllowListPolicy(allowed []types.Address) PolicyMaker {
	allowedMap := make(map[types.Address]struct{})
	for _, a := range allowed {
		allowedMap[a] = struct{}{}
	}
	return &listPolicy{participants: allowedMap, mode: allow}
}

// ShouldApprove decides to approve o if it is currently unapproved
func (lp *listPolicy) ShouldApprove(o protocols.Objective) bool {
	for _, p := range o.GetParticipants() {
		_, found := lp.participants[p]

		// If this is an allow list then we reject any objective that has a participant not on the list
		if lp.mode == allow && !found {
			return false
		}
		// If this is a deny list then we reject any objective that has a participant on the list
		if lp.mode == deny && found {
			return false
		}

	}
	return true
}

type listMode string

const (
	allow listMode = "allow"
	deny  listMode = "deny"
)

type listPolicy struct {
	participants map[types.Address]struct{}
	mode         listMode
}
