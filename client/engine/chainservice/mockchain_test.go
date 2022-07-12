package chainservice

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func TestDeposit(t *testing.T) {
	// The MockChain should react to a deposit transaction for a given channel by sending an event with updated holdings for that channel

	// Construct MockChain
	var chain = NewMockChain()
	eventFeed := chain.EventFeed()

	// Prepare test data to trigger MockChainService
	testDeposit := types.Funds{
		common.HexToAddress("0x00"): big.NewInt(1),
	}
	testTx := protocols.NewDepositTransaction(types.Destination(common.HexToHash(`4ebd366d014a173765ba1e50f284c179ade31f20441bec41664712aac6cc461d`)), testDeposit)

	// Send one transaction and receive one event from it.
	err := chain.SendTransaction(testTx)
	if err != nil {
		t.Error(err)
	}
	event := <-eventFeed

	checkReceivedEventIsValid(t, event, testTx.Deposit, testTx.ChannelId())

	// Send the transaction again and receive another event
	err = chain.SendTransaction(testTx)
	if err != nil {
		t.Error(err)
	}

	event = <-eventFeed

	// The expectation is that the MockChainService remembered the previous deposit and added this one to it:
	expectedHoldings := testTx.Deposit.Add(testTx.Deposit)

	checkReceivedEventIsValid(t, event, expectedHoldings, testTx.ChannelId())
}

func checkReceivedEventIsValid(t *testing.T, receivedEvent Event, holdings types.Funds, channelId types.Destination) {
	if receivedEvent.ChannelID() != channelId {
		t.Fatalf(`channelId mismatch: expected %v but got %v`, channelId, receivedEvent.ChannelID())
	}

	depositEvent := receivedEvent.(DepositedEvent)
	if depositEvent.NowHeld.Cmp(holdings[depositEvent.AssetAddress]) != 0 {
		t.Fatalf(`holdings mismatch: expected %v but got %v`, holdings[depositEvent.AssetAddress], depositEvent.NowHeld)
	}
}
