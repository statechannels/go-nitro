package netproto

var (
	// Message
	ErrUnexpectedMessage = NewError("unexpected message")
	ErrInvalidMessage    = NewError("invalid message")

	// Response
	ErrUnexpectedResponse = NewError("unexpected response")
	ErrInvalidResponse    = NewError("invalid response")
)
