package transport

import "github.com/statechannels/go-nitro/rpc/serde"

type Connection interface {
	// Request sends data for a topic and returns the response data or an error
	Request(serde.RequestMethod, []byte) ([]byte, error)
	// Respond listens for requests for a topic and calls the handler function when a request is received
	Respond(serde.RequestMethod, func([]byte) []byte) error

	// Close shuts down the connection
	Close()
}
