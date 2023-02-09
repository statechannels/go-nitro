package transport

import "github.com/statechannels/go-nitro/rpc/serde"

type Requester interface {
	Close()

	// Request sends data for a topic and returns the response data or an error
	Request(serde.RequestMethod, []byte) ([]byte, error)
	// Subscribe listens for notifications for a topic and does not respond to notifications
	Subscribe(serde.NotificationMethod, func([]byte)) error
}

type Subscriber interface {
	Close()

	// Respond listens for requests for a topic and calls the handler function when a request is received
	Respond(serde.RequestMethod, func([]byte) []byte) error
	// Notify sends data for a topic without expecting a response
	Notify(serde.NotificationMethod, []byte) error
}
