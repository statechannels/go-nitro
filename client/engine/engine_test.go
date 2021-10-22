package engine

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// testObjective simulates a computationally intensive task when cranked
type testObjective struct {
}

func (t testObjective) Id() protocols.ObjectiveId {
	return `test`
}
func (t testObjective) Approve() protocols.Objective {
	return t
}
func (t testObjective) Reject() protocols.Objective                               { return t }
func (t testObjective) Update(event protocols.ObjectiveEvent) protocols.Objective { return t }
func (t testObjective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	b := make([]byte, 100)
	rand.Read(b)
	state.SignEthereumMessage(b, *secretKey)
	fmt.Println(`signed a message`)
	return t, protocols.SideEffects{}, ``, nil
}

// testStore holds a specific secret key and returns a testObjective
type testStore struct{}

func (testStore) GetChannelSecretKey() *[]byte {
	k := common.Hex2Bytes(`187bb12e927c1652377405f81d93ce948a593f7d66cfba383ee761858b05921a`)
	return &k
}

func (testStore) GetObjectiveById(protocols.ObjectiveId) protocols.Objective {
	return testObjective{}
}
func (testStore) GetObjectiveByChannelId(types.Bytes32) protocols.Objective {
	return testObjective{}
}
func (testStore) SetObjective(protocols.Objective) error {
	return nil
}
func (testStore) UpdateProgressLastMadeAt(protocols.ObjectiveId, protocols.WaitingFor) {}

// testMessageService has a pair of chans, that's it
type testMessageService struct{}

var recieveChan chan protocols.Message = make(chan protocols.Message)
var sendChan chan protocols.Message = make(chan protocols.Message)

func (testMessageService) GetRecieveChan() chan protocols.Message { return recieveChan }
func (testMessageService) GetSendChan() chan protocols.Message    { return sendChan }
func (testMessageService) Send(message protocols.Message)         {}

type testChainService struct{}

var recieveChanChain chan chainservice.Event = make(chan chainservice.Event)
var sendChanChain chan protocols.Transaction = make(chan protocols.Transaction)

func (testChainService) GetRecieveChan() chan chainservice.Event { return recieveChanChain }
func (testChainService) GetSendChan() chan protocols.Transaction { return sendChanChain }
func (testChainService) Submit(tx protocols.Transaction)         {}

// TestRun stresses the engine by sending it a large number of messages
func TestRun(t *testing.T) {
	fmt.Println(`constructing engine`)

	// construct an engine with a test store and test messaging and chain services

	e := New(testMessageService{}, testChainService{}, testStore{})

	go e.Run()
	msg := protocols.Message{
		ObjectiveId: `test`,
	}

	for j := 1; j <= 100_000; j++ {
		// hit the API repeatedly in such a way as to trigger the testobjective's "expensive" crank
		e.fromMsg <- msg
	}

}
