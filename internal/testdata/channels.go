package testdata

import (
	"fmt"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

type channelCollection struct {
	// MockTwoPartyLedger constructs and returns a ledger channel
	MockTwoPartyLedger virtualfund.GetTwoPartyLedgerFunction
}

var Channels channelCollection = channelCollection{
	MockTwoPartyLedger: mockTwoPartyLedger,
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
