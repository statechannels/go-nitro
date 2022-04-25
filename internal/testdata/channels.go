package testdata

import (
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

type channelCollection struct {
	// MockConsensusChannel constructs and returns a ledger channel
	MockConsensusChannel virtualfund.GetLedgerFunction
}

var Channels channelCollection = channelCollection{
	MockConsensusChannel: mockConsensusChannel,
}

func mockTwoPartyLedger(firstParty, secondParty types.Address) (ledger *channel.TwoPartyLedger, ok bool) {
	ledger, err := channel.NewTwoPartyLedger(createLedgerState(
		firstParty,
		secondParty,
		100,
		100,
	), 0) // todo: make myIndex configurable
	if err != nil {
		panic(fmt.Errorf("error mocking a twoPartyLedger: %w", err))
	}
	return ledger, true
}

func mockConsensusChannel(counterparty types.Address) (ledger *consensus_channel.ConsensusChannel, ok bool) {
	ts := testState.Clone()
	request := directfund.ObjectiveRequest{
		MyAddress:         ts.Participants[0],
		CounterParty:      ts.Participants[1],
		AppData:           ts.AppData,
		AppDefinition:     ts.AppDefinition,
		ChallengeDuration: ts.ChallengeDuration,
		Nonce:             ts.ChannelNonce.Int64(),
		Outcome:           ts.Outcome,
	}
	testObj, _ := directfund.NewObjective(request, false)
	cc, _ := testObj.CreateLedgerChannel()
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
func (l LedgerNetwork) GetLedgerLookup(seeker types.Address) virtualfund.GetLedgerFunction {
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

// createLedgerNetwork returns active, funded consensus_channels connecting the supplied
// actors according the the supplied edge list.
//
// Edges specify actors via their indices in the actors slice,
// and each edge is ordered [leader, follower].
func createLedgerNetwork(actors []testactors.Actor, edges [][2]int) LedgerNetwork {
	// naive connectedness check: does not detect, eg
	// a--b  c--d (no path from a to c or d, etc)
	for i, a := range actors {
		connected := false
		for _, edge := range edges {
			if edge[0] == i || edge[1] == i {
				connected = true
			}
		}
		if !connected {
			panic(fmt.Sprintf("actor %v is not connected in the test network", a))
		}
	}

	var ret LedgerNetwork

	for _, edge := range edges {
		if edge[0] > len(actors)-1 ||
			edge[1] > len(actors)-1 ||
			edge[0] == edge[1] {
			panic(fmt.Sprintf("malformed ledger network edge list: %v", edge))
		}
		leader := actors[edge[0]]
		follower := actors[edge[1]]

		testLedger := createTestLedger(leader, follower)

		ret.ledgers = append(ret.ledgers, testLedger)
	}

	return ret
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
	fp.Participants[0] = leader.Address
	fp.Participants[1] = follower.Address

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
		fp,
		0,
		*outcome,
		sigs,
	)
	if err != nil {
		panic(fmt.Sprintf("error creating leader channel in testLedger: %v", err))
	}
	followCh, err := consensus_channel.NewFollowerChannel(
		fp,
		0,
		*outcome,
		sigs,
	)
	if err != nil {
		panic(fmt.Sprintf("error creating follwer channel in testLedger: %v", err))
	}

	return TestLedger{leaderCh, followCh}
}
