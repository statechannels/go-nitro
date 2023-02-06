package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/statechannels/go-nitro/network/transport"
)

type natsTransport struct {
	connection transport.Connection
}

var _ transport.Transport = (*natsTransport)(nil)

func NewNatsTransport(nc *nats.Conn) *natsTransport {
	// wouldn't it be better to get messaged directly into the channel?
	connection := NewNatsConnection(nc)
	natsTransport := &natsTransport{
		connection: connection,
	}

	return natsTransport
}

func (t *natsTransport) PollConnection() (transport.Connection, error) {
	return t.connection, nil
}

func (t *natsTransport) Close() {
	t.connection.Close()
}
