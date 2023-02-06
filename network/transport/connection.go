package transport

type Connection interface {
	Close()

	Request(string, []byte) ([]byte, error)
	Subscribe(string, func([]byte) []byte) error
}
