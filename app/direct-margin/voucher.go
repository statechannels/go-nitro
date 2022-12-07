package directmargin

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

type MaginState struct {
	ChannelId      types.Destination
	LeaderAmount   *big.Int
	FollowerAmount *big.Int
	LeaderSig      state.Signature
	FollowerSig    state.Signature
}

// NOTE: maybe rename "Request" with "Message", since not only requests are sent (for example: voucher signed by the leader or follower)

// Margin App request types
const RequestTypeFunding = "funding"
const RequestTypeMarginProposal = "margin-proposal"
const RequestTypeMarginAccept = "margin-accept"
const RequestTypeMarginReject = "margin-reject"
const RequestTypeDefunding = "defunding"

