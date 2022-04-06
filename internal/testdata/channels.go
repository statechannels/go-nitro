package testdata

import (
	"fmt"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

type channelCollection struct {
	// MockTwoPartyLedger constructs and returns a ledger channel
	MockTwoPartyLedger   virtualfund.GetTwoPartyLedgerFunction
	MockConsensusChannel virtualfund.GetTwoPartyConsensusLedgerFunction
}

var Channels channelCollection = channelCollection{
	MockTwoPartyLedger:   mockTwoPartyLedger,
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
	cc, _ := testObj.CreateConsensusChannel()
	return cc, true
}
