package rpc

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/statechannels/go-nitro/node/query"
	"github.com/statechannels/go-nitro/types"
)

type LedgerChannelView struct {
	ID           types.Destination
	Counterparty types.Address
	AssetAddress types.Address
	MyBalance    *hexutil.Big
	TheirBalance *hexutil.Big
	// todo: running guarantees?
}

// NewLedgerChannelView creates a "view with perspective" into a given ledger channel.
func NewLedgerChannelView(info query.LedgerChannelInfo, myAddress types.Address) (LedgerChannelView, error) {
	lcv := LedgerChannelView{}

	lcv.ID = info.ID
	lcv.AssetAddress = info.Balance.AssetAddress

	if info.Balance.Leader == myAddress {
		lcv.Counterparty = info.Balance.Follower
		lcv.MyBalance = info.Balance.LeaderBalance
		lcv.TheirBalance = info.Balance.FollowerBalance
	} else if info.Balance.Follower == myAddress {
		lcv.Counterparty = info.Balance.Leader
		lcv.MyBalance = info.Balance.FollowerBalance
		lcv.TheirBalance = info.Balance.LeaderBalance
	} else {
		return lcv, fmt.Errorf("%s is not a participant in ledger channel %s", myAddress, info.ID)
	}

	return lcv, nil
}
