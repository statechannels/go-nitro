package chainservice

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func TestDeposit(t *testing.T) {
	// The MockChain should react to a deposit transaction for a given channel by sending an event with updated holdings for that channel to all subsribers

	var a = types.Address(common.HexToAddress(`a`))
	var b = types.Address(common.HexToAddress(`b`))

	// Construct MockChain and tell it the subscriber addresses.
	// This is not super elegant but gets around data races -- the constructor will make channels and then run a listener which will send on them.
	var chain = NewMockChain()
	chain.SubscribeToEvents(a)
	chain.SubscribeToEvents(b)

	eventFeedA, err := chain.EventFeed(a)
	if err != nil {
		t.Fatalf("subscription for address a failed")
	}

	eventFeedB, err := chain.EventFeed(b)
	if err != nil {
		t.Fatalf("subscription for address b failed")
	}

	// Prepare test data to trigger MockChainService
	testDeposit := types.Funds{
		common.HexToAddress("0x00"): big.NewInt(1),
	}
	testTx := protocols.NewDepositTransaction(types.Destination(common.HexToHash(`4ebd366d014a173765ba1e50f284c179ade31f20441bec41664712aac6cc461d`)), testDeposit)

	// Send one transaction and receive one event from it.
	chain.SendTransaction(testTx)
	event := <-eventFeedA

	receivedEventIsValid(t, event, testTx.Deposit, testTx.ChannelId())

	// Send the transaction again and receive another event
	chain.SendTransaction(testTx)
	event = <-eventFeedA

	// The expectation is that the MockChainService remembered the previous deposit and added this one to it:
	expectedHoldings := testTx.Deposit.Add(testTx.Deposit)

	receivedEventIsValid(t, event, expectedHoldings, testTx.ChannelId())

	// Pull an event out of the other mock chain service and check that
	eventB := <-eventFeedB
	receivedEventIsValid(t, eventB, testTx.Deposit, testTx.ChannelId())

	// Pull another event out of the other mock chain service and check that
	eventB = <-eventFeedB
	receivedEventIsValid(t, eventB, expectedHoldings, testTx.ChannelId())
}

func receivedEventIsValid(t *testing.T, receivedEvent Event, holdings types.Funds, channelId types.Destination) {
	if receivedEvent.ChannelID() != channelId {
		t.Fatalf(`channelId mismatch: expected %v but got %v`, channelId, receivedEvent.ChannelID())
	}

	depositEvent := receivedEvent.(DepositedEvent)
	if depositEvent.NowHeld.Cmp(holdings[depositEvent.Asset]) != 0 {
		t.Fatalf(`holdings mismatch: expected %v but got %v`, holdings[depositEvent.Asset], depositEvent.NowHeld)
	}
}
