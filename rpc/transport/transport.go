package transport

type TransportType string

const (
	Nats TransportType = "nats"
	Ws   TransportType = "ws"
)

// Requester is a transport that can send requests and subscribe to notifications
type Requester interface {
	// Close closes the connection
	Close() error

	// Request sends a blocking request and returns the response data or an error
	Request([]byte) ([]byte, error)
	// Subscribe provides a notification channel.
	// If subscription to notifications fails, it returns an error.
	Subscribe() (<-chan []byte, error)
}

// Responder is a transport that can respond to requests and send notifications
type Responder interface {
	// Close closes the connection
	Close() error
	// Url returns the url that the responder is listening on
	Url() string

	// RegisterRequestHandler registers a handler that accepts a request and returns a response.
	// It returns an error if the registration setup fails
	RegisterRequestHandler(string, func([]byte) []byte) error
	// Notify sends notification data without expecting a response
	Notify([]byte) error
}
