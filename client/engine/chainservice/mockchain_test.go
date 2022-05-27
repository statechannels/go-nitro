package chainservice

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func TestDeposit(t *testing.T) {
	// The MockChain and SimpleChainService should work together to react to a deposit transaction for a given channel by:
	//  - sending an event with updated holdings for that channel to all SimpleChainServices which are subscribed

	var a = types.Address(common.HexToAddress(`a`))
	var b = types.Address(common.HexToAddress(`b`))

	// Construct MockChain and tell it the addresses of the SimpleChainServices which will subscribe to it.
	// This is not super elegant but gets around data races -- the constructor will make channels and then run a listener which will send on them.
	var chain = NewMockChain()
	chain.Subscribe(a)
	chain.Subscribe(b)

	// Construct SimpleChainServices
	mcsA := NewSimpleChainService(chain, a)
	mcsB := NewSimpleChainService(chain, b)

	inA := mcsA.In()
	outA := mcsA.Out()

	// Prepare test data to trigger MockChainService
	testDeposit := types.Funds{
		common.HexToAddress("0x00"): big.NewInt(1),
	}
	testTx := protocols.ChainTransaction{
		ChannelId: types.Destination(common.HexToHash(`4ebd366d014a173765ba1e50f284c179ade31f20441bec41664712aac6cc461d`)),
		Deposit:   testDeposit,
		Type:      protocols.DepositTransactionType,
	}

	// Send one transaction into one of the SimpleChainServices and receive one event from it.
	inA <- testTx
	event := <-outA

	if event.ChannelID() != testTx.ChannelId {
		t.Fatalf(`channelId mismatch: expected %v but got %v`, testTx.ChannelId, event.ChannelID())
	}
	if !event.(DepositedEvent).Holdings.Equal(testTx.Deposit) {
		t.Fatalf(`holdings mismatch: expected %v but got %v`, testTx.Deposit, event.(DepositedEvent).Holdings)
	}

	// Send the transaction again and receive another event
	inA <- testTx
	event = <-outA

	// The expectation is that the MockChainService remembered the previous deposit and added this one to it:
	expectedHoldings := testTx.Deposit.Add(testTx.Deposit)

	if event.ChannelID() != testTx.ChannelId {
		t.Fatalf(`channelId mismatch: expected %v but got %v`, testTx.ChannelId, event.ChannelID())
	}
	if !event.(DepositedEvent).Holdings.Equal(expectedHoldings) {
		t.Fatalf(`holdings mismatch: expected %v but got %v`, expectedHoldings, event.(DepositedEvent).Holdings)
	}

	// Pull an event out of the other mock chain service and check that
	eventB := <-mcsB.Out()

	if eventB.ChannelID() != testTx.ChannelId {
		t.Fatalf(`channelId mismatch: expected %v but got %v`, testTx.ChannelId, eventB.ChannelID())
	}
	if !eventB.(DepositedEvent).Holdings.Equal(testTx.Deposit) {
		t.Fatalf(`holdings mismatch: expected %v but got %v`, testTx.Deposit, eventB.(DepositedEvent).Holdings)
	}

	// Pull another event out of the other mock chain service and check that
	eventB = <-mcsB.Out()

	if eventB.ChannelID() != testTx.ChannelId {
		t.Fatalf(`channelId mismatch: expected %v but got %v`, testTx.ChannelId, eventB.ChannelID())
	}
	if !eventB.(DepositedEvent).Holdings.Equal(expectedHoldings) {
		t.Fatalf(`holdings mismatch: expected %v but got %v`, expectedHoldings, eventB.(DepositedEvent).Holdings)
	}

}
