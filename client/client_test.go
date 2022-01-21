package client

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/types"
)

func TestNew(t *testing.T) {

	skA := common.Hex2Bytes("c417e8a75ebe2bfe16fe108e1e04802c324974eef6ea2cc3d55194fa38677b5e")
	a := types.Address(common.HexToAddress(`0xaaa3D879df547333a9ac87341C92f11e5FB79CD4`))
	b := types.Address(common.HexToAddress(`b`))
	chain := chainservice.NewMockChain([]types.Address{a, b})
	chainservice := chainservice.NewSimpleChainService(chain, a)
	messageservice := messageservice.NewTestMessageService(a)
	store := store.NewMockStore(skA)
	New(messageservice, chainservice, store) // TODO reinstate client:=
}
