package p2pms

import (
	"bufio"
	"context"
	"encoding/json"
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
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/internal/logging"
	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// basicPeerInfo contains the basic information about a peer
type basicPeerInfo struct {
	Id      peer.ID
	Address types.Address
}

const (
	PROTOCOL_ID                  protocol.ID = "/go-nitro/msg/1.0.0"
	PEER_EXCHANGE_PROTOCOL_ID    protocol.ID = "/go-nitro/peerinfo/1.0.0"
	DELIMITER                                = '\n'
	BUFFER_SIZE                              = 1_000
	NUM_CONNECT_ATTEMPTS                     = 20
	RETRY_SLEEP_DURATION                     = 5 * time.Second
	PEER_EXCHANGE_SLEEP_DURATION             = 100 * time.Millisecond
)

// P2PMessageService is a rudimentary message service that uses TCP to send and receive messages.
type P2PMessageService struct {
	toEngine chan protocols.Message // for forwarding processed messages to the engine
	peers    *safesync.Map[basicPeerInfo]

	me          types.Address
	key         p2pcrypto.PrivKey
	p2pHost     host.Host
	mdns        mdns.Service
	newPeerInfo chan basicPeerInfo
	logger      zerolog.Logger
}

// Id returns the libp2p peer ID of the message service.
func (ms *P2PMessageService) Id() peer.ID {
	id, _ := peer.IDFromPrivateKey(ms.key)
	return id
}

// NewMessageService returns a running P2PMessageService listening on the given ip, port and message key.
func NewMessageService(ip string, port int, me types.Address, pk []byte, logWriter io.Writer) *P2PMessageService {
	logging.ConfigureZeroLogger()

	ms := &P2PMessageService{
		toEngine:    make(chan protocols.Message, BUFFER_SIZE),
		newPeerInfo: make(chan basicPeerInfo, BUFFER_SIZE),
		peers:       &safesync.Map[basicPeerInfo]{},
		me:          me,
		logger:      zerolog.New(logWriter).With().Timestamp().Str("message-service", me.String()[0:8]).Caller().Logger(),
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

	mdns := mdns.NewMdnsService(host, "", ms)
	err = mdns.Start()
	ms.checkError(err)
	ms.mdns = mdns
	ms.p2pHost = host

	ms.p2pHost.SetStreamHandler(PROTOCOL_ID, ms.msgStreamHandler)

	ms.p2pHost.SetStreamHandler(PEER_EXCHANGE_PROTOCOL_ID, func(stream network.Stream) {
		ms.receivePeerInfo(stream)
		stream.Close()
	})

	return ms
}

// HandlePeerFound is called by the mDNS service when a peer is found.
func (ms *P2PMessageService) HandlePeerFound(pi peer.AddrInfo) {
	ms.p2pHost.Peerstore().AddAddr(pi.ID, pi.Addrs[0], peerstore.PermanentAddrTTL)
	stream, err := ms.p2pHost.NewStream(context.Background(), pi.ID, PEER_EXCHANGE_PROTOCOL_ID)
	ms.checkError(err)
	ms.sendPeerInfo(stream)
	stream.Close()
}

func (ms *P2PMessageService) msgStreamHandler(stream network.Stream) {
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
}

// sendPeerInfo sends our peer info over the given stream
func (ms *P2PMessageService) sendPeerInfo(stream network.Stream) {
	raw, err := json.Marshal(basicPeerInfo{
		Id:      ms.Id(),
		Address: ms.me,
	})

	ms.checkError(err)
	writer := bufio.NewWriter(stream)
	// We don't care about the number of bytes written
	_, err = writer.WriteString(string(raw) + string(DELIMITER))
	ms.checkError(err)
	writer.Flush()
}

// receivePeerInfo receives peer info from the given stream
func (ms *P2PMessageService) receivePeerInfo(stream network.Stream) {
	reader := bufio.NewReader(stream)
	// Create a buffer stream for non blocking read and write.
	raw, err := reader.ReadString(DELIMITER)

	// An EOF means the stream has been closed by the other side.
	if errors.Is(err, io.EOF) {
		stream.Close()
		return
	}
	ms.checkError(err)

	var peerInfo *basicPeerInfo
	err = json.Unmarshal([]byte(raw), &peerInfo)
	ms.checkError(err)

	_, foundPeer := ms.peers.LoadOrStore(peerInfo.Address.String(), *peerInfo)
	if !foundPeer {
		ms.logger.Debug().Interface("peerInfo", peerInfo).Msgf("New peer found")
		ms.newPeerInfo <- *peerInfo
	}
}

// Send sends messages to other participants.
// It blocks until the message is sent.
// It will retry establishing a stream NUM_CONNECT_ATTEMPTS times before giving up
func (ms *P2PMessageService) Send(msg protocols.Message) {
	raw, err := msg.Serialize()
	ms.checkError(err)

	peerInfo, ok := ms.peers.Load(msg.To.String())
	if !ok {
		panic(fmt.Errorf("could not load peer %s", msg.To.String()))
	}

	for i := 0; i < NUM_CONNECT_ATTEMPTS; i++ {
		s, err := ms.p2pHost.NewStream(context.Background(), peerInfo.Id, PROTOCOL_ID)
		if err == nil {

			writer := bufio.NewWriter(s)

			// We don't care about the number of bytes written
			_, err = writer.WriteString(raw + string(DELIMITER))

			ms.checkError(err)

			writer.Flush()
			s.Close()

			return
		}

		ms.logger.Info().Int("attempt", i).Str("to", msg.To.String()).Msg("Could not open stream")
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
	s.mdns.Close()
	s.p2pHost.RemoveStreamHandler(PROTOCOL_ID)
	return s.p2pHost.Close()
}

// PeerInfoReceived returns a channel that receives a PeerInfo when a peer is discovered
func (s *P2PMessageService) PeerInfoReceived() <-chan basicPeerInfo {
	return s.newPeerInfo
}

// PeerInfo contains peer information and the ip address/port
type PeerInfo struct {
	Port      int
	Id        peer.ID
	Address   types.Address
	IpAddress string
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
		ms.p2pHost.Peerstore().AddAddr(p.Id, multi, peerstore.PermanentAddrTTL)

		ms.peers.Store(p.Address.String(), basicPeerInfo{p.Id, p.Address})
	}
}
