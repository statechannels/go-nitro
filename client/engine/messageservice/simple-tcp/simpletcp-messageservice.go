package simpletcp

import (
	"bufio"
	"fmt"
	"net"

	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const (
	CONN_TYPE = "tcp"
	DELIMETER = '\n'
)

// SimpleTCPMessageService is a rudimentary message service that uses TCP to send and receive messages
type SimpleTCPMessageService struct {
	out chan protocols.Message // for sending message to engine

	peers map[types.Address]string

	listener net.Listener // The listener for incoming connections on our port

	quit chan struct{} // quit is used to signal the goroutine to stop

}

// NewTestMessageService returns a running SimpleTcpMessageService listening on the given url
func NewSimpleTCPMessageService(myUrl string, peers map[types.Address]string) *SimpleTCPMessageService {

	l, err := net.Listen(CONN_TYPE, myUrl)
	if err != nil {
		panic(err)
	}
	h := &SimpleTCPMessageService{
		out:   make(chan protocols.Message, 5),
		peers: peers,

		listener: l,
		quit:     make(chan struct{}),
	}

	go h.listenForIncoming()

	return h

}

// Send dispatches messages
func (s *SimpleTCPMessageService) Send(msg protocols.Message) {
	peer, ok := s.peers[msg.To]

	if !ok {
		panic(fmt.Errorf("no peer port for %s", msg.To))
	}

	raw, err := msg.Serialize()

	// Append the delimiter to the message to indicate the end of the message
	raw = fmt.Sprintf("%s%c", raw, DELIMETER)

	if err != nil {
		s.panicIfRunning(err)
		return
	}

	conn, err := net.Dial(CONN_TYPE, peer)
	if err != nil {
		s.panicIfRunning(err)
		return
	}
	_, err = conn.Write([]byte(raw))
	if err != nil {
		s.panicIfRunning(err)
		return
	}
	conn.Close()

}

// listenForIncoming listens for any incoming messages from other peers
func (s *SimpleTCPMessageService) listenForIncoming() {
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			s.panicIfRunning(err)
			return
		}

		raw, err := bufio.NewReader(conn).ReadString(DELIMETER)

		if err != nil {
			s.panicIfRunning(err)
			return

		}
		m, err := protocols.DeserializeMessage(raw)
		if err != nil {
			s.panicIfRunning(err)
			return
		}
		s.out <- m

	}

}

// panicIfRunning panics if the SimpleTCPMessageService is running, otherwise it just returns
func (s *SimpleTCPMessageService) panicIfRunning(err error) {
	select {
	case <-s.quit: // If we are quitting we can ignore the error
		return
	default:
		panic(err)
	}
}

func (s *SimpleTCPMessageService) Out() <-chan protocols.Message {
	return s.out
}

// Close closes the SimpleTCPMessageService
func (s *SimpleTCPMessageService) Close() {
	close(s.quit)
	s.listener.Close()

}
