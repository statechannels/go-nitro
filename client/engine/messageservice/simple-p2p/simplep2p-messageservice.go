package simplep2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const (
	MESSAGE_ADDRESS = "/messages/1.0.0"
	DELIMETER       = '\n'
	BUFFER_SIZE     = 1_000_000
)

// P2PMessageService is a rudimentary message service that uses TCP to send and receive messages
type P2PMessageService struct {
	in  chan protocols.Message // for receiving messages from engine
	out chan protocols.Message // for sending message to engine

	peers map[types.Address]PeerInfo

	quit chan struct{} // quit is used to signal the goroutine to stop

	me      PeerInfo
	p2pHost host.Host
}

type PeerInfo struct {
	Address    types.Address
	Port       int64
	Id         peer.ID
	MessageKey crypto.PrivKey
}

// MultiAddress returns the multiaddress of the peer based on their port and Id
func (p PeerInfo) MultiAddress() multiaddr.Multiaddr {

	a, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/p2p/%s", p.Port, p.Id))
	if err != nil {
		panic(err)
	}

	return a
}

// GeneratePeerInfo generates a random  message key/ peer id and returns a PeerInfo
func GeneratePeerInfo(add types.Address, port int64) PeerInfo {
	messageKey, _, err := crypto.GenerateECDSAKeyPair(rand.Reader)
	if err != nil {
		panic(err)
	}
	id, err := peer.IDFromPrivateKey(messageKey)
	if err != nil {
		panic(err)
	}
	return PeerInfo{Id: id, Address: add, Port: port, MessageKey: messageKey}
}

// NewTestMessageService returns a running SimpleTcpMessageService listening on the given url
func NewP2PMessageService(me PeerInfo, peers map[types.Address]PeerInfo) *P2PMessageService {

	options := []libp2p.Option{libp2p.Identity(me.MessageKey),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", me.Port)),
		libp2p.DefaultTransports,
		libp2p.NoSecurity,
		libp2p.DefaultMuxers,
	}
	host, err := libp2p.New(options...)
	if err != nil {

		panic(err)
	}

	h := &P2PMessageService{
		in:      make(chan protocols.Message, BUFFER_SIZE),
		out:     make(chan protocols.Message, BUFFER_SIZE),
		peers:   peers,
		p2pHost: host,
		quit:    make(chan struct{}),
		me:      me,
	}
	h.p2pHost.SetStreamHandler(MESSAGE_ADDRESS, func(stream network.Stream) {

		reader := bufio.NewReader(stream)
		for {
			select {
			case <-h.quit:
				stream.Close()
				return
			default:

				// Create a buffer stream for non blocking read and write.
				raw, err := reader.ReadString(DELIMETER)

				h.checkError(err)
				m, err := protocols.DeserializeMessage(raw)
				h.checkError(err)
				h.out <- m
			}
		}
	})

	return h

}

// DialPeers dials all peers in the peer list and establishs a connection with them.
// This should be called once all the message services are running.
// TODO: The message service should handle this internally
func (s *P2PMessageService) DialPeers() {
	go s.connectToPeers()
}

// connectToPeers establishes a stream with all our peers and uses that stream to send messages
func (s *P2PMessageService) connectToPeers() {
	// create a map with streams to all peers
	peerStreams := make(map[types.Address]network.Stream)
	for _, p := range s.peers {
		if p.Address == s.me.Address {
			continue
		}
		// Extract the peer ID from the multiaddr.
		info, err := peer.AddrInfoFromP2pAddr(p.MultiAddress())
		s.checkError(err)

		// Add the destination's peer multiaddress in the peerstore.
		// This will be used during connection and stream creation by libp2p.
		s.p2pHost.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

		stream, err := s.p2pHost.NewStream(context.Background(), info.ID, MESSAGE_ADDRESS)
		s.checkError(err)
		peerStreams[p.Address] = stream
	}
	for {
		select {
		case <-s.quit:

			for _, writer := range peerStreams {
				writer.Close()
			}
			return
		case m := <-s.in:
			raw, err := m.Serialize()
			s.checkError(err)
			writer := bufio.NewWriter(peerStreams[m.To])
			_, err = writer.WriteString(raw)
			s.checkError(err)
			err = writer.WriteByte(DELIMETER)
			s.checkError(err)
			writer.Flush()

		}
	}

}

// checkError panics if the SimpleTCPMessageService is running, otherwise it just returns
func (s *P2PMessageService) checkError(err error) {
	if err == nil {
		return
	}
	select {

	case <-s.quit: // If we are quitting we can ignore the error
		return
	default:
		panic(err)
	}
}

func (s *P2PMessageService) Out() <-chan protocols.Message {
	return s.out
}

func (s *P2PMessageService) In() chan<- protocols.Message {
	return s.in
}

// Close closes the SimpleTCPMessageService
func (s *P2PMessageService) Close() {
	close(s.quit)

}
