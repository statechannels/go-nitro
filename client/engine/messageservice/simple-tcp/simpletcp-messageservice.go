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
	in  chan protocols.Message // for receiving messages from engine
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
		in:    make(chan protocols.Message, 5),
		out:   make(chan protocols.Message, 5),
		peers: peers,

		listener: l,
		quit:     make(chan struct{}),
	}

	go h.listenForIncoming()
	go h.listenForOutgoing()

	return h

}

// listenForOutgoing listens for any messages that should be sent out to peers and sends them
func (s *SimpleTCPMessageService) listenForOutgoing() {
	for {
		select {
		case <-s.quit:
			return

		case m := <-s.in:
			{
				peer, ok := s.peers[m.To]

				if !ok {
					panic(fmt.Errorf("no peer port for %s", m.To))
				}

				raw, err := m.Serialize()

				// Append the delimiter to the message to indicate the end of the message
				raw = fmt.Sprintf("%s%c", raw, DELIMETER)

				if err != nil {
					panic(err)
				}

				conn, err := net.Dial(CONN_TYPE, peer)
				if err != nil {
					select {
					case <-s.quit: // If we are quitting we can ignore the error
						return
					default:
						panic(err)
					}
				}
				_, err = conn.Write([]byte(raw))
				if err != nil {
					select {
					case <-s.quit: // If we are quitting we can ignore the error
						return
					default:
						panic(err)
					}
				}
				conn.Close()

			}

		}
	}

}

// listenForIncoming listens for any incoming messages from other peers
func (s *SimpleTCPMessageService) listenForIncoming() {
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				panic(err)
			}
		}

		raw, err := bufio.NewReader(conn).ReadString(DELIMETER)

		if err != nil {
			select {
			case <-s.quit: // If we are quitting we can ignore the error
				return
			default:
				panic(err)
			}

		}
		m, err := protocols.DeserializeMessage(raw)
		if err != nil {
			panic(err)
		}
		s.out <- m

		conn.Close()

	}

}

func (s *SimpleTCPMessageService) Out() <-chan protocols.Message {
	return s.out
}

func (s *SimpleTCPMessageService) In() chan<- protocols.Message {
	return s.in
}

// Close closes the SimpleTCPMessageService
func (s *SimpleTCPMessageService) Close() {
	close(s.quit)
	s.listener.Close()

}
