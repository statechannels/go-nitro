package transport

type Connection interface {
	Send(string, []byte)
	Recv() ([]byte, error)

	Close()

	Request(string, []byte) ([]byte, error)
	Subscribe(string, func([]byte) []byte) error
}
