package p2pms

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
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
	"github.com/statechannels/go-nitro/crypto"
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
	DHT_PROTOCOL_PREFIX          protocol.ID = "/nitro" // use /nitro/kad/1.0.0 instead of /ipfs/kad/1.0.0
	PROTOCOL_ID                  protocol.ID = "/nitro/msg/1.0.0"
	PEER_EXCHANGE_PROTOCOL_ID    protocol.ID = "/nitro/peerinfo/1.0.0"
	DELIMITER                                = '\n'
	BUFFER_SIZE                              = 1_000
	NUM_CONNECT_ATTEMPTS                     = 20
	RETRY_SLEEP_DURATION                     = 5 * time.Second
	PEER_EXCHANGE_SLEEP_DURATION             = 10 * time.Second // how often we attempt FindPeers
	BOOTSTRAP_SLEEP_DURATION                 = 1 * time.Second  // how often we check for bootpeers in Peerstore
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

	MultiAddr string
}

// Id returns the libp2p peer ID of the message service.
func (ms *P2PMessageService) Id() peer.ID {
	id, _ := peer.IDFromPrivateKey(ms.key)
	return id
}

// NewMessageService returns a running P2PMessageService listening on the given ip, port and message key.
// If useMdnsPeerDiscovery is true, the message service will use mDNS to discover peers.
// Otherwise, peers must be added manually via `AddPeers`.
func NewMessageService(ip string, port int, pk []byte, useMdnsPeerDiscovery bool, logWriter io.Writer, bootPeers []string) *P2PMessageService {
	logging.ConfigureZeroLogger()

	me := crypto.GetAddressFromSecretKeyBytes(pk)
	ms := &P2PMessageService{
		toEngine:    make(chan protocols.Message, BUFFER_SIZE),
		newPeerInfo: make(chan basicPeerInfo, BUFFER_SIZE),
		peers:       &safesync.Map[peer.ID]{},
		me:          me,
		logger:      logging.WithAddress(zerolog.New(logWriter).With().Timestamp(), &me).Caller().Logger(),
	}

	messageKey, err := p2pcrypto.UnmarshalSecp256k1PrivateKey(pk)
	ms.checkError(err)

	ms.key = messageKey
	options := []libp2p.Option{
		libp2p.Identity(messageKey),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", ip, port)),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.NATPortMap(),
		libp2p.DefaultMuxers,
	}
	host, err := libp2p.New(options...)
	if err != nil {
		panic(err)
	}

	ms.p2pHost = host
	ms.p2pHost.SetStreamHandler(PROTOCOL_ID, ms.msgStreamHandler)

	// Print out my own peerInfo
	peerInfo := peer.AddrInfo{
		ID:    ms.p2pHost.ID(),
		Addrs: ms.p2pHost.Addrs(),
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	ms.checkError(err)
	ms.MultiAddr = addrs[0].String()
	ms.logger.Debug().Msgf("libp2p node multiaddrs: %v", addrs)

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
	ctx := context.Background()
	var options []dht.Option
	options = append(options, dht.BucketSize(20))
	options = append(options, dht.Mode(dht.ModeServer)) // allows other peers to connect to this node

	kademliaDHT, err := dht.New(ctx, ms.p2pHost, options...)
	ms.checkError(err)
	ms.dht = kademliaDHT

	// Setup notifications so we exchange nitro signing addresses when connected
	n := &network.NotifyBundle{}
	n.ConnectedF = func(n network.Network, conn network.Conn) {
		ms.logger.Debug().Msgf("notification: connected to peer %s", conn.RemotePeer().Pretty())
		ms.HandlePeerFound(peer.AddrInfo{ID: conn.RemotePeer(), Addrs: []multiaddr.Multiaddr{conn.RemoteMultiaddr()}})
	}
	n.DisconnectedF = func(n network.Network, conn network.Conn) {
		ms.logger.Debug().Msgf("notification: disconnected from peer: %s", conn.RemotePeer().Pretty())
	}

	ms.p2pHost.Network().Notify(n)

	ms.addBootPeers(bootPeers)

	expectedPeers := len(bootPeers)
	if expectedPeers > 0 {
		ms.logger.Debug().Msgf("waiting for %d bootpeer connections", expectedPeers)
		ticker := time.NewTicker(BOOTSTRAP_SLEEP_DURATION)
		for range ticker.C {
			peers := ms.p2pHost.Network().Peers()
			actualPeers := len(peers)
			ms.logger.Debug().Msgf("found peers: %v, expected peers: %d", actualPeers, expectedPeers)
			for _, peer := range peers {
				ms.logger.Debug().Msgf("peer info: %v", peer)
			}

			// Once we've discovered all bootPeers, stop the ticker
			if actualPeers >= expectedPeers {
				ms.logger.Debug().Msgf("discovered all expected bootPeers")
				ticker.Stop()
				break
			}
		}
	}

	err = ms.dht.Bootstrap(ctx) // Runs periodically to maintain a healthy routing table
	ms.checkError(err)

	ms.logger.Debug().Msgf("DHT setup complete")
}

// HandlePeerFound is called when a peer is discovered.
func (ms *P2PMessageService) HandlePeerFound(pi peer.AddrInfo) {
	ms.logger.Debug().Str("peer-id", pi.ID.Pretty()).Str("multiaddress", pi.Addrs[0].String()).Msgf("Attempting to add peer")
	ms.p2pHost.Peerstore().AddAddr(pi.ID, pi.Addrs[0], peerstore.PermanentAddrTTL)

	a, err := addressFromPeerID(pi.ID)
	ms.checkError(err)
	ms.logger.Debug().Str("peer-id", pi.ID.Pretty()).Str("sc-address", a.String()).Msgf("Derived SC address from peer ID")

	ms.peers.Store(a.String(), pi.ID)
	ms.newPeerInfo <- basicPeerInfo{Id: pi.ID, Address: a}
}

func (ms *P2PMessageService) msgStreamHandler(stream network.Stream) {
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

// Send sends messages to other participants.
// It blocks until the message is sent.
// It will retry establishing a stream NUM_CONNECT_ATTEMPTS times before giving up
func (ms *P2PMessageService) Send(msg protocols.Message) {
	raw, err := msg.Serialize()
	ms.checkError(err)

	peerId, ok := ms.peers.Load(msg.To.String())
	if !ok {
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
func (ms *P2PMessageService) checkError(err error) {
	if err == nil {
		return
	}
	ms.logger.Panic().Msg(err.Error())
}

// Out returns a channel that can be used to receive messages from the message service
func (ms *P2PMessageService) Out() <-chan protocols.Message {
	return ms.toEngine
}

// Close closes the P2PMessageService
func (ms *P2PMessageService) Close() error {
	// The mdns service is optional so we only close it if it exists
	if ms.mdns != nil {
		err := ms.mdns.Close()
		if err != nil {
			return err
		}
	}
	ms.p2pHost.RemoveStreamHandler(PROTOCOL_ID)
	return ms.p2pHost.Close()
}

// PeerInfoReceived returns a channel that receives a PeerInfo when a peer is discovered
func (ms *P2PMessageService) PeerInfoReceived() <-chan basicPeerInfo {
	return ms.newPeerInfo
}

func (ms *P2PMessageService) addBootPeers(peers []string) {
	for _, p := range peers {
		addr, err := multiaddr.NewMultiaddr(p)
		ms.checkError(err)

		peer, err := peer.AddrInfoFromP2pAddr(addr)
		ms.checkError(err)

		a, err := addressFromPeerID(peer.ID)
		ms.checkError(err)
		ms.peers.Store(a.String(), peer.ID)

		err = ms.p2pHost.Connect(context.Background(), *peer) // Adds peerInfo to local Peerstore
		ms.checkError(err)
		ms.logger.Debug().Msgf("connected to boot peer: %v (%s)", p, a.String())

	}
}

// addressFromPeerID attempts to extract a secp256k1 public key from a peer ID and returns the corresponding SC address
func addressFromPeerID(pId peer.ID) (types.Address, error) {
	p2pPub, err := pId.ExtractPublicKey()
	if err != nil {
		return types.Address{}, err
	}
	// raw will be a compressed 33 byte public key
	raw, err := p2pPub.Raw()
	if err != nil {
		return types.Address{}, err
	}

	pub, err := ethcrypto.DecompressPubkey(raw)
	if err != nil {
		return types.Address{}, err
	}

	return ethcrypto.PubkeyToAddress(*pub), nil
}
