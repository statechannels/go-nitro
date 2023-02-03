package transport

type Transport interface {
	PollConnection() (Connection, error)

	Close()
}
