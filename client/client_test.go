package client

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/types"
)

func TestNew(t *testing.T) {
	a := types.Address(common.HexToAddress(`a`))
	b := types.Address(common.HexToAddress(`b`))
	chainservice := chainservice.NewMockChain([]types.Address{a, b})

	messageservice := messageservice.NewTestMessageService(a)

	client := New(messageservice, chainservice)
}
