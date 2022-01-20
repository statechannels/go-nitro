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
	_, b := GeneratePrivateKeyAndAddress()
	chain := chainservice.MockChain{}

	chainserv := chainservice.NewSimpleChainService(chain, a)
	messageservice := messageservice.NewTestMessageService(a)
	store := store.NewMockStore()
	client := New(messageservice, chainserv, store, aKey, a)

	chain = chainservice.NewMockChain([]types.Address{a, b})

	client.CreateDirectChannel(b, types.Address{}, types.Bytes{}, outcome.Exit{}, big.NewInt(0))
}
