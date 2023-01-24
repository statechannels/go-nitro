package transport

type Connection interface {
	Send(string, []byte)
	Recv() ([]byte, error)

	Close()
}
