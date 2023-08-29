package policy

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

type ChannelType string

const (
	Ledger  ChannelType = "ledger"
	Payment ChannelType = "payment"
)

func NewDenyListPolicy(denied []types.Address, channelType ChannelType) PolicyMaker {
	deniedMap := make(map[types.Address]struct{})
	for _, a := range denied {
		deniedMap[a] = struct{}{}
	}
	return &listPolicy{participants: deniedMap, mode: deny, channelType: channelType}
}

func NewAllowListPolicy(allowed []types.Address, channelType ChannelType) PolicyMaker {
	allowedMap := make(map[types.Address]struct{})
	for _, a := range allowed {
		allowedMap[a] = struct{}{}
	}
	return &listPolicy{participants: allowedMap, mode: allow, channelType: channelType}
}

func (lp *listPolicy) shouldCheckObjective(o protocols.Objective) bool {
	_, isDf := o.(*directfund.Objective)
	_, isVf := o.(*virtualfund.Objective)
	_, isDdf := o.(*directdefund.Objective)
	_, isVdf := o.(*virtualdefund.Objective)

	if (isDdf || isDf) && lp.channelType == Ledger {
		return true
	}

	if (isVf || isVdf) && lp.channelType == Payment {
		return true
	}
	return false
}

// ShouldApprove decides to approve o if it is currently unapproved
func (lp *listPolicy) ShouldApprove(o protocols.Objective) bool {
	if lp.shouldCheckObjective(o) {
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
	channelType  ChannelType
}
