package transport

import (
	"fmt"
)

var (
	// Transport
	ErrTransportClosed = fmt.Errorf("transport closed")

	// Connection
	ErrConnectionClosed = fmt.Errorf("connection closed")
)
