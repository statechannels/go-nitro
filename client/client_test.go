package client

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

func TestNew(t *testing.T) {

	aKey, a := crypto.GeneratePrivateKeyAndAddress()
	bKey, b := crypto.GeneratePrivateKeyAndAddress()
	chain := chainservice.NewMockChain([]types.Address{a, b})

	chainservA := chainservice.NewSimpleChainService(chain, a)
	messageserviceA := messageservice.NewTestMessageService(a)
	storeA := store.NewMockStore(aKey)
	clientA := New(messageserviceA, chainservA, storeA)

	chainservB := chainservice.NewSimpleChainService(chain, b)
	messageserviceB := messageservice.NewTestMessageService(a)
	storeB := store.NewMockStore(bKey)
	New(messageserviceB, chainservB, storeB)

	messageserviceA.Connect(messageserviceB)
	messageserviceB.Connect(messageserviceA)

	clientA.CreateDirectChannel(b, types.Address{}, types.Bytes{}, outcome.Exit{}, big.NewInt(0))

}
