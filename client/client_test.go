package client

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/types"
)

func TestNew(t *testing.T) {

	aKey, a := GeneratePrivateKeyAndAddress()
	bKey, b := GeneratePrivateKeyAndAddress()
	chain := chainservice.NewMockChain([]types.Address{a, b})

	chainservA := chainservice.NewSimpleChainService(chain, a)
	messageserviceA := messageservice.NewTestMessageService(a)
	storeA := store.NewMockStore()
	clientA := New(messageserviceA, chainservA, storeA, aKey, a)

	chainservB := chainservice.NewSimpleChainService(chain, b)
	messageserviceB := messageservice.NewTestMessageService(a)
	storeB := store.NewMockStore()
	New(messageserviceB, chainservB, storeB, bKey, b)

	messageserviceA.Connect(messageserviceB)
	messageserviceB.Connect(messageserviceA)

	clientA.CreateDirectChannel(b, types.Address{}, types.Bytes{}, outcome.Exit{}, big.NewInt(0))
}
