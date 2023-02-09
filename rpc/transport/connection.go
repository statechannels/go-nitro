package transport

import "github.com/statechannels/go-nitro/rpc/serde"

type Requester interface {
	Close()
	Request(serde.RequestMethod, []byte) ([]byte, error)
}

type Subscriber interface {
	Close()
	Subscribe(serde.RequestMethod, func([]byte) []byte) error
}
