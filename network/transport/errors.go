package transport

import (
	"github.com/statechannels/go-nitro/internal"
)

var (
	// Transport
	ErrTransportClosed = internal.NewError("transport closed")

	// Connection
	ErrConnectionClosed = internal.NewError("connection closed")
)
