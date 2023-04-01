package integration_test

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/internal/testactors"
)

const (
	MAX_PARTICIPANTS       = 4
	STORE_TEST_DATA_FOLDER = "../data/store_test"
	ledgerChannelDeposit   = 5_000_000
	defaultTimeout         = 10 * time.Second
	virtualChannelDeposit  = 5000
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
	StoreType StoreType
	Name      testactors.ActorName
}

type MessageService string

const (
	TestMessageService MessageService = "TestMessageService"
	P2PMessageService  MessageService = "P2PMessageService"
)

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
	if tc.NumOfChannels < 1 {
		return fmt.Errorf("NumOfChannels must be greater than 0")
	}
	if tc.NumOfChannels > 10 {
		return fmt.Errorf("NumOfChannels must be smaller than 10")
	}
	if tc.MessageDelay > 5*time.Second {
		return fmt.Errorf("MessageDelay must be smaller than 5s")
	}
	return nil
}

type sharedInra struct {
	broker         *messageservice.Broker
	peers          []p2pms.PeerInfo
	mockChain      *chainservice.MockChain
	simulatedChain *chainservice.SimulatedChain
	bindings       *chainservice.Bindings
	ethAccounts    []*bind.TransactOpts
}
