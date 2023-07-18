package p2pms

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	p2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
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

type peerExchangeMessage struct {
	Id             peer.ID
	Address        types.Address
	ExpectResponse bool
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
	peers    *safesync.Map[peer.ID]

	me          types.Address
	key         p2pcrypto.PrivKey
	p2pHost     host.Host
	mdns        mdns.Service
	dht         *dht.IpfsDHT
	newPeerInfo chan basicPeerInfo
	logger      zerolog.Logger
}

// Id returns the libp2p peer ID of the message service.
func (ms *P2PMessageService) Id() peer.ID {
	id, _ := peer.IDFromPrivateKey(ms.key)
	return id
}

// NewMessageService returns a running P2PMessageService listening on the given ip, port and message key.
// If useMdnsPeerDiscovery is true, the message service will use mDNS to discover peers.
// Otherwise, peers must be added manually via `AddPeers`.
func NewMessageService(ip string, port int, me types.Address, pk []byte, useMdnsPeerDiscovery bool, logWriter io.Writer, bootPeers []string) *P2PMessageService {
	logging.ConfigureZeroLogger()

	ms := &P2PMessageService{
		toEngine:    make(chan protocols.Message, BUFFER_SIZE),
		newPeerInfo: make(chan basicPeerInfo, BUFFER_SIZE),
		peers:       &safesync.Map[peer.ID]{},
		me:          me,
		logger:      zerolog.New(logWriter).With().Timestamp().Str("message-service", me.String()[0:8]).Caller().Logger(),
	}

	messageKey, err := p2pcrypto.UnmarshalSecp256k1PrivateKey(pk)
	ms.checkError(err)

	ms.key = messageKey

	options := []libp2p.Option{
		libp2p.Identity(messageKey),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", ip, port)),
		// libp2p.Transport(tcp.NewTCPTransport),
		// libp2p.NoSecurity,
		// libp2p.DefaultMuxers,
	}
	host, err := libp2p.New(options...)
	if err != nil {
		panic(err)
	}

	ms.p2pHost = host
	ms.p2pHost.SetStreamHandler(PROTOCOL_ID, ms.msgStreamHandler)
	ms.p2pHost.SetStreamHandler(PEER_EXCHANGE_PROTOCOL_ID, ms.receivePeerInfo)

	if useMdnsPeerDiscovery {
		ms.setupMdns()
	} else {
		ms.setupDht(bootPeers)
	}

	return ms
}

func (ms *P2PMessageService) setupMdns() {
	// Since the mdns service could trigger a call to  `HandlePeerFound` at any time once started
	// We want to start mdns after the message service has been fully constructed
	ms.mdns = mdns.NewMdnsService(ms.p2pHost, "", ms)
	err := ms.mdns.Start()
	ms.checkError(err)
}

func (ms *P2PMessageService) setupDht(bootPeers []string) {
	log.SetAllLoggers(log.LevelInfo) // Set default log level for all loggers
	// Set specific log level for DHT module
	dhtKey := "dht"
	_ = log.SetLogLevel(dhtKey, "debug")

	ctx := context.Background()
	var options []dht.Option
	if len(bootPeers) == 0 {
		// Server mode: act as a bootstrapping node, allowing other peers to join this node
		options = append(options, dht.Mode(dht.ModeServer))
	}
	kademliaDHT, err := dht.New(ctx, ms.p2pHost, options...)
	ms.checkError(err)
	ms.dht = kademliaDHT

	// Print out my own peerInfo
	peerInfo := peer.AddrInfo{
		ID:    ms.p2pHost.ID(),
		Addrs: ms.p2pHost.Addrs(),
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	ms.checkError(err)
	fmt.Println("libp2p node address:", addrs[0])

	// Add bootstrap peers
	err = kademliaDHT.Bootstrap(ctx)
	ms.checkError(err)
	ms.addBootPeers(bootPeers)

	expectedPeers := len(bootPeers)
	if expectedPeers > 0 {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			peers := kademliaDHT.RoutingTable().ListPeers()
			fmt.Printf("Found peers: %v\n", len(peers))
			for _, peer := range peers {
				fmt.Println("Peer info: ", peer)
			}

			// Once we've discovered all bootPeers, stop the ticker
			if len(peers) >= expectedPeers {
				fmt.Println("Discovered all expected bootPeers.")
				ticker.Stop()
				break
			}
		}
		go ms.discoverPeers(ctx, string(PEER_EXCHANGE_PROTOCOL_ID)) // This will fail if DHT is empty (aka expectedPeers == 0)
	}

	ms.logger.Debug().Msgf("DHT setup complete")
}

// HandlePeerFound is called by the mDNS service when a peer is found.
func (ms *P2PMessageService) HandlePeerFound(pi peer.AddrInfo) {
	ms.logger.Debug().Msgf("Attempting to add mdns peer")
	ms.p2pHost.Peerstore().AddAddr(pi.ID, pi.Addrs[0], peerstore.PermanentAddrTTL)

	ms.sendPeerInfo(pi.ID, false)
}

func (ms *P2PMessageService) msgStreamHandler(stream network.Stream) {
	ms.logger.Debug().Msgf("received message")
	defer stream.Close()

	reader := bufio.NewReader(stream)
	// Create a buffer stream for non blocking read and write.
	raw, err := reader.ReadString(DELIMITER)

	// An EOF means the stream has been closed by the other side.
	if errors.Is(err, io.EOF) {
		return
	}
	ms.checkError(err)
	m, err := protocols.DeserializeMessage(raw)
	ms.checkError(err)
	ms.toEngine <- m
}

// sendPeerInfo sends our peer info to a given peerId
func (ms *P2PMessageService) sendPeerInfo(recipientId peer.ID, expectResponse bool) {
	stream, err := ms.p2pHost.NewStream(context.Background(), recipientId, PEER_EXCHANGE_PROTOCOL_ID)
	ms.checkError(err)
	defer stream.Close()

	raw, err := json.Marshal(peerExchangeMessage{
		Id:             ms.Id(),
		Address:        ms.me,
		ExpectResponse: expectResponse,
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
	ms.logger.Debug().Msgf("received peerInfo")
	defer stream.Close()

	reader := bufio.NewReader(stream)
	// Create a buffer stream for non blocking read and write.
	raw, err := reader.ReadString(DELIMITER)

	// An EOF means the stream has been closed by the other side.
	if errors.Is(err, io.EOF) {
		return
	}
	ms.checkError(err)

	var msg *peerExchangeMessage
	err = json.Unmarshal([]byte(raw), &msg)
	ms.checkError(err)

	peerInfo := basicPeerInfo{msg.Id, msg.Address}

	_, foundPeer := ms.peers.LoadOrStore(msg.Address.String(), msg.Id)
	if !foundPeer {
		ms.logger.Debug().Interface("peerInfo", peerInfo).Msgf("New peer found")
		ms.peers.Range(func(key string, value peer.ID) bool {
			fmt.Printf("peers map - key: %s, value: %+v\n", key, value)
			return true
		})
		ms.newPeerInfo <- peerInfo
	}

	if msg.ExpectResponse {
		ms.sendPeerInfo(msg.Id, false)
	}
}

// Send sends messages to other participants.
// It blocks until the message is sent.
// It will retry establishing a stream NUM_CONNECT_ATTEMPTS times before giving up
func (ms *P2PMessageService) Send(msg protocols.Message) {
	raw, err := msg.Serialize()
	ms.checkError(err)

	peerId, ok := ms.peers.Load(msg.To.String())
	if !ok {
		ms.peers.Range(func(key string, value peer.ID) bool {
			fmt.Printf("key: %s, value: %+v\n", key, value)
			return true
		})
		panic(fmt.Errorf("could not load peer %s", msg.To.String()))
	}

	for i := 0; i < NUM_CONNECT_ATTEMPTS; i++ {
		s, err := ms.p2pHost.NewStream(context.Background(), peerId, PROTOCOL_ID)
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
	// The mdns service is optional so we only close it if it exists
	if s.mdns != nil {
		err := s.mdns.Close()
		if err != nil {
			return err
		}
	}
	s.p2pHost.RemoveStreamHandler(PROTOCOL_ID)
	return s.p2pHost.Close()
}

// PeerInfoReceived returns a channel that receives a PeerInfo when a peer is discovered
func (s *P2PMessageService) PeerInfoReceived() <-chan basicPeerInfo {
	return s.newPeerInfo
}

func (ms *P2PMessageService) addBootPeers(peers []string) {
	for _, p := range peers {
		addr, err := multiaddr.NewMultiaddr(p)
		ms.checkError(err)

		peer, err := peer.AddrInfoFromP2pAddr(addr)
		ms.checkError(err)

		err = ms.p2pHost.Connect(context.Background(), *peer)
		ms.checkError(err)
		ms.logger.Debug().Msgf("connected to boot peer: %v", p)

		ms.sendPeerInfo(peer.ID, true)
	}
}

func (ms *P2PMessageService) discoverPeers(ctx context.Context, topic string) {
	routingDiscovery := routing.NewRoutingDiscovery(ms.dht)
	_, err := routingDiscovery.Advertise(ctx, topic) // fires every 3 hours by default
	ms.checkError(err)

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

			peersChan, err := routingDiscovery.FindPeers(ctx, topic)
			ms.checkError(err)

			for p := range peersChan {
				if p.ID == ms.p2pHost.ID() {
					continue
				}
				ms.logger.Debug().Msgf("inspecting peer found through discovery")
				if ms.p2pHost.Network().Connectedness(p.ID) != network.Connected {

					_, err = ms.p2pHost.Network().DialPeer(ctx, p.ID)
					ms.checkError(err)
					ms.p2pHost.Peerstore().AddAddr(p.ID, p.Addrs[0], peerstore.PermanentAddrTTL)
					ms.logger.Debug().Msgf("connected to new peer: %+v", p)

					ms.sendPeerInfo(p.ID, true)
				}
			}
		}
	}
}
