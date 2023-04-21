package testdata

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

type channelCollection struct {
	// MockConsensusChannel constructs and returns a ledger channel
	MockConsensusChannel virtualfund.GetTwoPartyConsensusLedgerFunction
}

var Channels channelCollection = channelCollection{
	MockConsensusChannel: mockConsensusChannel,
}

func mockConsensusChannel(counterparty types.Address) (ledger *consensus_channel.ConsensusChannel, ok bool) {
	ts := testState.Clone()
	ts.TurnNum = 0
	ss := state.NewSignedState(ts)
	id := protocols.ObjectiveId(directfund.ObjectivePrefix + testState.ChannelId().String())
	op := protocols.CreateObjectivePayload(id, directfund.SignedStatePayload, ss)
	testObj, err := directfund.ConstructFromPayload(true, op, ts.Participants[0])
	if err != nil {
		return &consensus_channel.ConsensusChannel{}, false
	}

	// Manually progress the extended state by collecting postfund signatures
	correctSignatureByAliceOnPostFund, _ := testObj.C.PostFundState().Sign(testactors.Alice.PrivateKey)
	correctSignatureByBobOnPostFund, _ := testObj.C.PostFundState().Sign(testactors.Bob.PrivateKey)
	testObj.C.AddStateWithSignature(testObj.C.PostFundState(), correctSignatureByAliceOnPostFund)
	testObj.C.AddStateWithSignature(testObj.C.PostFundState(), correctSignatureByBobOnPostFund)

	cc, err := testObj.CreateConsensusChannel()
	cc.OnChainFunding = types.Funds{
		common.HexToAddress("0x00"): big.NewInt(2),
	}

	if err != nil {
		panic(err)
	}

	return cc, true
}

// LedgerNetwork is a collection of in-memory consensus_channel ledgers
// which expose both the leader and follower perspective on the ledger
type LedgerNetwork struct {
	ledgers []TestLedger
}

// TestLedger wraps the leader and follower views of a consensus_channel
type TestLedger struct {
	LeaderView   consensus_channel.ConsensusChannel
	FollowerView consensus_channel.ConsensusChannel
}

// GetLedgerLookup returns a ledger-lookup function for the given ledger seeker.
//
// The returned function inspects the ledgers from the ledger set, and returns
// the ledger between the seeker and given counterparty from the seeker's perspective.
func (l LedgerNetwork) GetLedgerLookup(seeker types.Address) virtualfund.GetTwoPartyConsensusLedgerFunction {
	var myLedgers []TestLedger

	// package all of seeker's ledgers for the closure
	for _, ledger := range l.ledgers {
		if ledger.FollowerView.Leader() == seeker ||
			ledger.FollowerView.Follower() == seeker {
			myLedgers = append(myLedgers, ledger)
		}
	}

	return func(counterparty types.Address) (ledger *consensus_channel.ConsensusChannel, ok bool) {
		for _, ledger := range myLedgers {
			if ledger.FollowerView.Follower() == seeker &&
				ledger.FollowerView.Leader() == counterparty {
				return &ledger.FollowerView, true
			}

			if ledger.LeaderView.Leader() == seeker &&
				ledger.LeaderView.Follower() == counterparty {
				return &ledger.LeaderView, true
			}
		}
		return nil, false
	}
}

// createLedgerPath returns active, funded consensus_channels connecting the supplied
// actors in left-to-right fashion from actors[0] to actors[len(actors)-1]. The
// leftmost actor in each consensuschannel is the channel's leader.
//
// Constructed and returned channels can be accessed from either "perspective"
// of leader or follower
func createLedgerPath(actors []testactors.Actor) LedgerNetwork {
	var ret LedgerNetwork

	for i, leader := range actors {
		if i+1 >= len(actors) {
			break
		}
		follower := actors[i+1]
		testLedger := createTestLedger(leader, follower)

		ret.ledgers = append(ret.ledgers, testLedger)
	}
	return ret
}

func createTestLedger(leader, follower testactors.Actor) TestLedger {
	fp := testState.Clone().FixedPart()
	fp.Participants[0] = leader.Address()
	fp.Participants[1] = follower.Address()

	outcome := consensus_channel.NewLedgerOutcome(
		types.Address{}, // the zero asset
		consensus_channel.NewBalance(leader.Destination(), big.NewInt(100)),
		consensus_channel.NewBalance(follower.Destination(), big.NewInt(200)),
		[]consensus_channel.Guarantee{},
	)

	initVars := consensus_channel.Vars{Outcome: *outcome, TurnNum: 0}
	leaderSig, _ := initVars.AsState(fp).Sign(leader.PrivateKey)
	followSig, _ := initVars.AsState(fp).Sign(follower.PrivateKey)

	sigs := [2]state.Signature{leaderSig, followSig}

	leaderCh, err := consensus_channel.NewLeaderChannel(
		initVars.AsState(fp),
		0,
		*outcome,
		sigs,
	)
	if err != nil {
		panic(fmt.Sprintf("error creating leader channel in testLedger: %v", err))
	}
	followCh, err := consensus_channel.NewFollowerChannel(
		initVars.AsState(fp),
		0,
		*outcome,
		sigs,
	)
	if err != nil {
		panic(fmt.Sprintf("error creating follwer channel in testLedger: %v", err))
	}

	return TestLedger{leaderCh, followCh}
}
