package nats

import (
	"testing"

	"github.com/statechannels/go-nitro/network/transport"
)

func TestNatsConnectionType(t *testing.T) {
	var _ transport.Connection = (*natsConnection)(nil)
}
