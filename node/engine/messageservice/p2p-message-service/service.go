package p2pms

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	p2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"
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
	DHT_PROTOCOL_PREFIX     protocol.ID = "/nitro" // use /nitro/kad/1.0.0 instead of /ipfs/kad/1.0.0
	GENERAL_MSG_PROTOCOL_ID protocol.ID = "/nitro/msg/1.0.0"

	DELIMITER                = '\n'
	BUFFER_SIZE              = 1_000
	NUM_CONNECT_ATTEMPTS     = 10
	RETRY_SLEEP_DURATION     = 5 * time.Second
	BOOTSTRAP_SLEEP_DURATION = 100 * time.Millisecond // how often we check for bootpeers in Peerstore
)

// P2PMessageService is a rudimentary message service that uses TCP to send and receive messages.
type P2PMessageService struct {
	initComplete chan struct{}
	toEngine     chan protocols.Message // for forwarding processed messages to the engine
	peers        *safesync.Map[peer.ID]

	me          types.Address
	key         p2pcrypto.PrivKey
	p2pHost     host.Host
	dht         *dht.IpfsDHT
	newPeerInfo chan basicPeerInfo
	logger      *slog.Logger

	MultiAddr string
}

// Id returns the libp2p peer ID of the message service.
func (ms *P2PMessageService) Id() peer.ID {
	id, _ := peer.IDFromPrivateKey(ms.key)
	return id
}

// NewMessageService returns a running P2PMessageService listening on the given ip, port and message key.
func NewMessageService(ip string, port int, me types.Address, pk []byte, bootPeers []string) *P2PMessageService {
	ms := &P2PMessageService{
		initComplete: make(chan struct{}, 1),
		toEngine:     make(chan protocols.Message, BUFFER_SIZE),
		newPeerInfo:  make(chan basicPeerInfo, BUFFER_SIZE),
		peers:        &safesync.Map[peer.ID]{},
		me:           me,
		logger:       logging.LoggerWithAddress(slog.Default(), me),
	}

	messageKey, err := p2pcrypto.UnmarshalSecp256k1PrivateKey(pk)
	if err != nil {
		panic(err)
	}

	ms.key = messageKey
	options := []libp2p.Option{
		libp2p.Identity(messageKey),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", "0.0.0.0", port)),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.NATPortMap(),
		libp2p.EnableNATService(),
		libp2p.DefaultMuxers,
	}
	host, err := libp2p.New(options...)
	if err != nil {
		panic(err)
	}

	ms.p2pHost = host
	ms.p2pHost.SetStreamHandler(GENERAL_MSG_PROTOCOL_ID, ms.msgStreamHandler)

	// Print out my own peerInfo
	peerInfo := peer.AddrInfo{
		ID:    ms.p2pHost.ID(),
		Addrs: ms.p2pHost.Addrs(),
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	ms.MultiAddr = addrs[0].String()
	ms.logger.Info("libp2p node", "multiaddrs", addrs)

	err = ms.setupDht(bootPeers)

	if err != nil {
		panic(err)
	}

	return ms
}

func (ms *P2PMessageService) setupDht(bootPeers []string) error {
	ctx := context.Background()
	var options []dht.Option
	options = append(options, dht.BucketSize(20))
	options = append(options, dht.Mode(dht.ModeServer))                                                    // allows other peers to connect to this node
	options = append(options, dht.ProtocolPrefix(DHT_PROTOCOL_PREFIX))                                     // need this to allow custom NamespacedValidator
	options = append(options, dht.NamespacedValidator(DHT_NAMESPACE, stateChannelAddrToPeerIDValidator{})) // all records prefixed with /scaddr/ will use this custom validator

	kademliaDHT, err := dht.New(ctx, ms.p2pHost, options...)
	if err != nil {
		return err
	}
	ms.dht = kademliaDHT

	// Setup network connection notifications
	n := &network.NotifyBundle{}
	n.ConnectedF = func(n network.Network, conn network.Conn) {
		ms.logger.Debug("notification: connected to peer", "peerId", conn.RemotePeer().String(), "peerCount", len(ms.p2pHost.Network().Peers()))

		peerInfo := basicPeerInfo{Id: conn.RemotePeer()}
		ms.newPeerInfo <- peerInfo
	}
	n.DisconnectedF = func(n network.Network, conn network.Conn) {
		ms.logger.Debug("notification: disconnected from peer", "peerId", conn.RemotePeer().String(), "peerCount", len(ms.p2pHost.Network().Peers()))
	}
	ms.p2pHost.Network().Notify(n)
	ms.connectBootPeers(bootPeers)

	err = ms.dht.Bootstrap(ctx) // Sends FIND_NODE queries periodically to populate dht routing table
	if err != nil {
		return err
	}

	// Must wait until dht RoutingTable has an entry before adding custom dht record
	// This is a restriction enforced by the libp2p library. When we try to put a value
	// into the DHT, the node is not storing it locally. Instead its telling other peers
	// to store it. The key-value pairs are stored on nodes with IDs closest to the key.
	// If the RoutingTable is empty, the node has no peers to propagate this information to.
	go func() {
		ticker := time.NewTicker(BOOTSTRAP_SLEEP_DURATION)
		for range ticker.C {
			if ms.dht.RoutingTable().Size() > 0 {
				ms.addScaddrDhtRecord(ctx)
				ticker.Stop()
				close(ms.initComplete)
				return
			}
		}
	}()

	ms.logger.Info("DHT setup complete")
	return nil
}

// InitComplete returns a chan that gets closed once the message service is initalized
func (ms *P2PMessageService) InitComplete() <-chan struct{} {
	return ms.initComplete
}

// addScaddrDhtRecord adds this node's state channel address to the custom dht namespace
func (ms *P2PMessageService) addScaddrDhtRecord(ctx context.Context) {
	ms.logger.Debug("Adding state channel address to dht")

	recordData := &dhtData{
		SCAddr:    ms.me.String(),
		PeerID:    ms.Id().String(),
		Timestamp: time.Time.Unix(time.Now()),
	}
	recordDataBytes, err := json.Marshal(recordData)
	ms.checkError(err)

	signature, err := ms.key.Sign(recordDataBytes)
	ms.checkError(err)

	fullRecord := &dhtRecord{
		Data:      *recordData,
		Signature: signature,
	}
	fullRecordBytes, err := json.Marshal(fullRecord)
	ms.checkError(err)

	key := DHT_RECORD_PREFIX + ms.me.String()
	err = ms.dht.PutValue(ctx, key, fullRecordBytes)
	ms.checkError(err)
	ms.logger.Info("Added state channel address to dht")
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
	if err != nil {
		ms.logger.Error("error reading from stream", "err", err)
		return
	}
	m, err := protocols.DeserializeMessage(raw)
	if err != nil {
		ms.logger.Error("error deserializing message", "err", err)
		return
	}
	ms.toEngine <- m
}

func (ms *P2PMessageService) getPeerIdFromDht(scaddr string) (peer.ID, error) {
	recordBytes, err := ms.dht.GetValue(context.Background(), DHT_RECORD_PREFIX+scaddr)
	if err != nil {
		return "", err
	}

	recordData := &dhtRecord{}
	err = json.Unmarshal(recordBytes, recordData)
	if err != nil {
		return "", err
	}

	peerId, err := peer.Decode(recordData.Data.PeerID)
	if err != nil {
		return "", err
	}
	ms.logger.Debug("found address in dht", "scaddr", scaddr, "peerId", peerId.String())

	ms.peers.Store(scaddr, peerId) // Cache this info locally for use next time
	return peerId, nil
}

// Send sends messages to other participants.
// It blocks until the message is sent.
// It will retry establishing a stream NUM_CONNECT_ATTEMPTS times before giving up
func (ms *P2PMessageService) Send(msg protocols.Message) error {
	raw, err := msg.Serialize()
	if err != nil {
		return err
	}

	// First try to get peerId from local "peers" map. If the address is not found there,
	// query the dht to retrieve the peerId, then store in local map for next time
	peerId, ok := ms.peers.Load(msg.To.String())
	if !ok {
		peerId, err = ms.getPeerIdFromDht(msg.To.String())
		if err != nil {
			return err
		}
	}

	for i := 0; i < NUM_CONNECT_ATTEMPTS; i++ {
		s, err := ms.p2pHost.NewStream(context.Background(), peerId, GENERAL_MSG_PROTOCOL_ID)
		if err == nil {
			writer := bufio.NewWriter(s)
			_, err = writer.WriteString(raw + string(DELIMITER)) // We don't care about the number of bytes written
			if err != nil {
				return err
			}

			writer.Flush()
			s.Close()
			return nil
		}

		ms.logger.Warn("error opening stream", "err", err, "attempt", i, "to", msg.To.String())
		time.Sleep(RETRY_SLEEP_DURATION)
	}
	return nil
}

// checkError panics if the message service is running and there is an error, otherwise it just returns
func (ms *P2PMessageService) checkError(err error) {
	if err == nil {
		return
	}
	ms.logger.Error("error in message service", "err", err)
	panic(err)
}

// Out returns a channel that can be used to receive messages from the message service
func (ms *P2PMessageService) Out() <-chan protocols.Message {
	return ms.toEngine
}

// Close closes the P2PMessageService
func (ms *P2PMessageService) Close() error {
	ms.p2pHost.RemoveStreamHandler(GENERAL_MSG_PROTOCOL_ID)
	return ms.p2pHost.Close()
}

// PeerInfoReceived returns a channel that receives a PeerInfo when a peer is discovered
func (ms *P2PMessageService) PeerInfoReceived() <-chan basicPeerInfo {
	return ms.newPeerInfo
}

// connectBootPeers connects to the given boot peers
func (ms *P2PMessageService) connectBootPeers(bootPeers []string) {
	expectedPeers := len(bootPeers)
	if expectedPeers == 0 {
		return
	}

	for _, p := range bootPeers {
		addr, err := multiaddr.NewMultiaddr(p)
		ms.checkError(err)

		peer, err := peer.AddrInfoFromP2pAddr(addr)
		ms.checkError(err)

		err = ms.p2pHost.Connect(context.Background(), *peer) // Adds peerInfo to local Peerstore
		ms.checkError(err)

		ms.logger.Debug("connected to boot peer", "peer", p)

	}

	// Add bootpeers and wait for connections before proceeding

	ms.logger.Info("waiting for bootpeer connections", "expectedPeers", expectedPeers)

	ticker := time.NewTicker(BOOTSTRAP_SLEEP_DURATION)
	for range ticker.C {
		peers := ms.p2pHost.Network().Peers()
		actualPeers := len(peers)
		ms.logger.Debug("peers found", "found-peers", actualPeers, "expected-peers", expectedPeers)

		for _, peer := range peers {
			ms.logger.Debug("peer info", "peer", peer.String())
		}

		// Once we've connected to enough peers, stop the ticker
		if actualPeers >= expectedPeers {
			ms.logger.Info("initial threshold for peer connections has been met")
			ticker.Stop()
			return
		}
	}
}
