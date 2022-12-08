package virtualmargin

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/client/engine/store/safesync"
	"github.com/statechannels/go-nitro/types"
)

type marginStatus struct {
	leaderAddr      common.Address
	followerAddr    common.Address
	startingBalance *big.Int
	latestMargin    MarginApp
	currentBalance  Balance
}

type MarginAppManager struct {
	channels *safesync.Map[*marginStatus]
	me       common.Address
}

func NewMarginAppManager(me types.Address) *MarginAppManager {
	channels := safesync.Map[*marginStatus]{}
	return &MarginAppManager{&channels, me}
}

func (ma *MarginAppManager) Register(channelId types.Destination, leader common.Address, follower common.Address, startingBalance *big.Int) error {
	// lets have our initial balance, equal by both sides
	// if starting balance is 100, then we will start to propose our margin from
	// LeaderAmount = 50, FollowerAmount = 50
	// TODO think about the case one of the malicious part can challenge this state

	halfStartingBalance := startingBalance.Div(startingBalance, big.NewInt(2))
	balance := Balance{halfStartingBalance, halfStartingBalance}
	marginApp := MarginApp{ChannelId: channelId, LeaderAmount: halfStartingBalance, FollowerAmount: halfStartingBalance, Version: big.NewInt(1)}

	data := &marginStatus{leader, follower, big.NewInt(0).Set(startingBalance), marginApp, balance}
	if _, ok := ma.channels.Load(channelId.String()); ok {
		return fmt.Errorf("channel already registered")
	}

	ma.channels.Store(channelId.String(), data)

	return nil
}

func (ma *MarginAppManager) Remove(channelId types.Destination) {
	ma.channels.Delete(channelId.String())
}

func (ma *MarginAppManager) ProposeAndSign(channelId types.Destination, leaderAmount *big.Int, followerAmount *big.Int, pk []byte) (MarginApp, error) {
	pStatus, ok := ma.channels.Load(channelId.String())
	if !ok {
		return MarginApp{}, fmt.Errorf("channel not found")
	}

	margin := MarginApp{LeaderAmount: &big.Int{}, FollowerAmount: &big.Int{}}
	newBalance := big.NewInt(0).Add(followerAmount, leaderAmount)
	if types.Gt(newBalance, pStatus.startingBalance) {
		return MarginApp{}, fmt.Errorf("sum of proposed funds should be equal to starting balance")
	}

	if pStatus.leaderAddr != ma.me {
		return MarginApp{}, fmt.Errorf("can only propose vouchers if leader")
	}

	margin.LeaderAmount = leaderAmount
	margin.FollowerAmount = followerAmount
	margin.ChannelId = channelId
	margin.Version = big.NewInt(0).Add(pStatus.latestMargin.Version, big.NewInt(0))

	if err := margin.LeaderSign(pk); err != nil {
		return margin, err
	}

	// Update margin app manager
	pStatus.currentBalance = Balance{Leader: leaderAmount, Follower: followerAmount}
	pStatus.latestMargin = margin

	return margin, nil
}

func (ma *MarginAppManager) HandleProposal(margin MarginApp) error {
	// Here All checks if valid proposal
	status, ok := ma.channels.Load(margin.ChannelId.String())
	if !ok {
		return fmt.Errorf("channel not registered")
	}

	if status.followerAddr != ma.me {
		return fmt.Errorf("only follower could handle proposal")
	}

	marginBalance := big.NewInt(0).Add(margin.FollowerAmount, margin.LeaderAmount)
	if !types.Gt(marginBalance, status.startingBalance) {
		return fmt.Errorf("sum of proposed funds should be equal to starting balance")
	}

	// Verify that margin was signed by Leader
	signer, err := margin.RecoverLeaderSigner()
	if err != nil {
		return fmt.Errorf("couldn't recover leader signature")
	}

	if signer != status.leaderAddr {
		return fmt.Errorf("margin was not signed by leader")
	}

	return nil
}

func (ma *MarginAppManager) SignProposal(margin MarginApp, pk []byte) (MarginApp, error) {
	status, ok := ma.channels.Load(margin.ChannelId.String())
	if !ok {
		return MarginApp{}, fmt.Errorf("channel not registered")
	}

	if err := margin.FollowerSign(pk); err != nil {
		return MarginApp{}, fmt.Errorf("couldn't sign latest margin")
	}

	status.latestMargin = margin
	status.currentBalance.Leader = margin.LeaderAmount
	status.currentBalance.Follower = margin.FollowerAmount

	//TODO
	// Make channel to push signed proposal

	return MarginApp{}, nil
}

func (ma *MarginAppManager) RejectProposal() {
	//TODO
	// Make channel to push rejected proposal
}

func (ma *MarginAppManager) ChannelRegistered(channelId types.Destination) bool {
	_, ok := ma.channels.Load(channelId.String())
	return ok

}

// For Leader - ProposeAndSign
// For Follower - HandleProposal, SignProposl/RejectProposal
