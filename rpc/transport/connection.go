package transport

type ConnectionType string

const (
	Nats ConnectionType = "nats"
	Ws   ConnectionType = "ws"
)

// Requester is a connection that can send requests and subscribe to notifications
type Requester interface {
	// Close closes the connection
	Close()

	// Request sends a blocking request and returns the response data or an error
	Request([]byte) ([]byte, error)
	// Subscribe provides a notification channel.
	// If subscription to notifications fails, it returns an error.
	Subscribe() (<-chan []byte, error)
}

// Responder is a connection that can respond to requests and send notifications
type Responder interface {
	// Close closes the connection
	Close()
	// Url returns the url that the responder is listening on
	Url() string

	// Respond listens for requests and calls the handler function when a request is received
	// It returns an error if the listener setup fails
	// The handler processes the incoming data and returns the response data
	Respond(func([]byte) []byte) error
	// Notify sends notification data without expecting a response
	Notify([]byte) error
}
