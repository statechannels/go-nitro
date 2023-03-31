package integration_test

import (
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

type sharedInra struct {
	broker         *messageservice.Broker
	peers          []p2pms.PeerInfo
	mockChain      *chainservice.MockChain
	simulatedChain *chainservice.SimulatedChain
	bindings       *chainservice.Bindings
	ethAccounts    []*bind.TransactOpts
}
