package chainservice

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func TestNew(t *testing.T) {
	NewMockChainService()
}
func TestInstantDeposit(t *testing.T) {
	// MockChainService should react to a deposit transaction for a given channel by:
	// - immediately sending an event with updated holdings for that channel.

	// Construct Mock Chain Service and get references to chans
	mcs := NewMockChainService()
	in := mcs.In()
	out := mcs.Out()

	// Prepare test data to trigger MockChainService
	testDeposit := types.Funds{
		common.HexToAddress("0x00"): big.NewInt(1),
	}
	testTx := protocols.Transaction{
		ChannelId: types.Destination(common.HexToHash(`4ebd366d014a173765ba1e50f284c179ade31f20441bec41664712aac6cc461d`)),
		Deposit:   testDeposit,
	}

	// Send one transaction and recieve one event
	in <- testTx
	event := <-out

	if event.ChannelId != testTx.ChannelId {
		t.Error(`channelId mismatch`)
	}
	if !event.Holdings.Equal(testTx.Deposit) {
		t.Error(`holdings mismatch`)
	}

	// Send the transaction again and recieve another event
	in <- testTx
	event = <-out

	// The expectation is that the MockChainService remembered the previous deposit and added this one to it:
	expectedHoldings := types.Funds{
		common.HexToAddress("0x00"): big.NewInt(2),
	}

	if event.ChannelId != testTx.ChannelId {
		t.Error(`channelId mismatch`)
	}
	if !event.Holdings.Equal(expectedHoldings) {
		t.Error(`holdings mismatch`)
	}

}
