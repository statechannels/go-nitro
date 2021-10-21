package chain

import "github.com/statechannels/go-nitro/protocols"

type Chain interface {
	Submit(tx protocols.Transaction)
}
