package transport

import "github.com/statechannels/go-nitro/rpc/serde"

type Connection interface {
	Close()

	Request(serde.RequestMethod, []byte) ([]byte, error)
	Subscribe(serde.RequestMethod, func([]byte) []byte) error
}
