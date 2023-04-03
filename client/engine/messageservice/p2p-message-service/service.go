package p2pms

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/libp2p/go-libp2p"
	p2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"
	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const (
	PROTOCOL_ID          protocol.ID = "/go-nitro/msg/1.0.0"
	DELIMITER                        = '\n'
	BUFFER_SIZE                      = 1_000
	NUM_CONNECT_ATTEMPTS             = 20
	RETRY_SLEEP_DURATION             = 5 * time.Second
)

// P2PMessageService is a rudimentary message service that uses TCP to send and receive messages.
type P2PMessageService struct {
	toEngine chan protocols.Message // for forwarding processed messages to the engine

	peers *safesync.Map[peer.ID]

	me      types.Address
	key     p2pcrypto.PrivKey
	p2pHost host.Host
}

// Id returns the libp2p peer ID of the message service.
func (ms *P2PMessageService) Id() peer.ID {
	id, _ := peer.IDFromPrivateKey(ms.key)
	return id
}

// AddPeers adds the peers to the message service.
// We ignore peers that are ourselves.
func (ms *P2PMessageService) AddPeers(peers []PeerInfo) {
	for _, p := range peers {
		// Ignore ourselves
		if p.Address == ms.me {
			continue
		}
		multi, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s", p.IpAddress, p.Port, p.Id))

		// Extract the peer ID from the multiaddr.
		info, _ := peer.AddrInfoFromP2pAddr(multi)
		// Add the destination's peer multiaddress in the peerstore.
		// This will be used during connection and stream creation by libp2p.
		ms.p2pHost.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

		ms.peers.Store(p.Address.String(), info.ID)
	}
}

// NewMessageService returns a running P2PMessageService listening on the given ip, port and message key.
func NewMessageService(ip string, port int, me types.Address, pk []byte) *P2PMessageService {
	ms := &P2PMessageService{
		toEngine: make(chan protocols.Message, BUFFER_SIZE),
		peers:    &safesync.Map[peer.ID]{},
		me:       me,
	}

	messageKey, err := p2pcrypto.UnmarshalSecp256k1PrivateKey(pk)
	ms.checkError(err)

	ms.key = messageKey

	options := []libp2p.Option{
		libp2p.Identity(messageKey),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", ip, port)),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.NoSecurity,
		libp2p.DefaultMuxers,
	}
	host, err := libp2p.New(options...)
	if err != nil {
		panic(err)
	}
	ms.p2pHost = host

	ms.p2pHost.SetStreamHandler(PROTOCOL_ID, func(stream network.Stream) {
		reader := bufio.NewReader(stream)
		// Create a buffer stream for non blocking read and write.
		raw, err := reader.ReadString(DELIMITER)

		// An EOF means the stream has been closed by the other side.
		if errors.Is(err, io.EOF) {
			stream.Close()
			return
		}
		ms.checkError(err)
		m, err := protocols.DeserializeMessage(raw)
		ms.checkError(err)
		ms.toEngine <- m
		stream.Close()
	})

	return ms
}

// Send sends messages to other participants.
// It blocks until the message is sent.
// It will retry establishing a stream NUM_CONNECT_ATTEMPTS times before giving up
func (ms *P2PMessageService) Send(msg protocols.Message) {
	raw, err := msg.Serialize()
	ms.checkError(err)

	id, ok := ms.peers.Load(msg.To.String())
	if !ok {
		panic(fmt.Errorf("could not load peer %s", msg.To.String()))
	}

	for i := 0; i < NUM_CONNECT_ATTEMPTS; i++ {
		s, err := ms.p2pHost.NewStream(context.Background(), id, PROTOCOL_ID)
		if err == nil {

			writer := bufio.NewWriter(s)

			// We don't care about the number of bytes written
			_, err = writer.WriteString(raw + string(DELIMITER))

			ms.checkError(err)

			writer.Flush()
			s.Close()

			return
		}

		// TODO: Hook up to a logger
		fmt.Printf("attempt %d: could not open stream to %s, retrying in %s\n", i, msg.To.String(), RETRY_SLEEP_DURATION.String())
		time.Sleep(RETRY_SLEEP_DURATION)

	}
}

// checkError panics if the message service is running and there is an error, otherwise it just returns
func (s *P2PMessageService) checkError(err error) {
	if err == nil {
		return
	}
	panic(err)
}

// Out returns a channel that can be used to receive messages from the message service
func (s *P2PMessageService) Out() <-chan protocols.Message {
	return s.toEngine
}

// Close closes the P2PMessageService
func (s *P2PMessageService) Close() error {
	s.p2pHost.RemoveStreamHandler(PROTOCOL_ID)
	return s.p2pHost.Close()
}
