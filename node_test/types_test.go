package node_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/node/engine/chainservice"
	"github.com/statechannels/go-nitro/node/engine/messageservice"
)

const (
	MAX_PARTICIPANTS       = 4
	STORE_TEST_DATA_FOLDER = "../data/store_test"
	ledgerChannelDeposit   = 5_000_000
	defaultTimeout         = 10 * time.Second
	virtualChannelDeposit  = 5000
	DURABLE_STORE_FOLDER   = "../data/node_test"
)

type StoreType string

const (
	MemStore     StoreType = "MemStore"
	DurableStore StoreType = "DurableStore"
)

type ChainType string

const (
	MockChain      ChainType = "MockChain"
	SimulatedChain ChainType = "SimulatedChain"
)

type TestParticipant struct {
	testactors.Actor
	StoreType StoreType
}

type MessageService string

const (
	TestMessageService MessageService = "TestMessageService"
	P2PMessageService  MessageService = "P2PMessageService"
)

// TestCase is a test case for the node integration test.
type TestCase struct {
	Description    string
	Chain          ChainType
	MessageService MessageService
	NumOfChannels  uint
	NumOfPayments  uint
	MessageDelay   time.Duration
	LogName        string
	NumOfHops      uint
	Participants   []TestParticipant
}

// Validate validates the test case and makes sure that the current test supports the test case.
func (tc *TestCase) Validate() error {
	if tc.NumOfHops < 1 || tc.NumOfHops > 2 {
		return fmt.Errorf("NumOfHops must be 1 or 2")
	}
	if tc.NumOfHops == 1 && len(tc.Participants) != 3 ||
		tc.NumOfHops == 2 && len(tc.Participants) != 4 {
		return fmt.Errorf("NumOfHops is %d, but there are %d participants", tc.NumOfHops, len(tc.Participants))
	}
	if tc.NumOfChannels < 1 || tc.NumOfChannels > 9 {
		return fmt.Errorf("NumOfChannels must be greater than 0 and less than 10. Supplied %d", tc.NumOfChannels)
	}
	if tc.MessageDelay > 5*time.Second {
		return fmt.Errorf("MessageDelay must be smaller than 5s")
	}
	return nil
}

// sharedTestInfrastructure is a struct that contains shared information liker the message broker, the simulated chain, and the ethereum accounts.
type sharedTestInfrastructure struct {
	broker         *messageservice.Broker
	mockChain      *chainservice.MockChain
	simulatedChain chainservice.SimulatedChain
	bindings       *chainservice.Bindings
	ethAccounts    []*bind.TransactOpts
}

func (sti *sharedTestInfrastructure) Close(t *testing.T) {
	if sti.simulatedChain != nil {
		closeSimulatedChain(t, sti.simulatedChain)
	}
}
