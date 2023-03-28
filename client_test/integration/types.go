package integration_test

import (
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/internal/testactors"
)

const STORE_TEST_DATA_FOLDER = "../data/store_test"
const MAX_PARTICIPANTS = 4

const ledgerChannelDeposit = 5_000_000
const virtualChannelDeposit = 5000

type StoreType string

const MemStore StoreType = "MemStore"
const DurableStore StoreType = "DurableStore"

type ChainType string

const MockChain ChainType = "MockChain"
const SimulatedChain ChainType = "SimulatedChain"

type TestParticipant struct {
	StoreType StoreType
	Name      testactors.ActorName
}

type MessageService string

const TestMessageService MessageService = "TestMessageService"
const P2PMessageService MessageService = "P2PMessageService"

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
	mockChain      *chainservice.MockChain
	simulatedChain *chainservice.SimulatedChain
	bindings       *chainservice.Bindings
	ethAccounts    []*bind.TransactOpts
}
