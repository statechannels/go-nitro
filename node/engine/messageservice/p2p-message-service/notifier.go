package p2pms

import (
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
)

// Triggered by p2pHost.Network events. When new peers are connected, the
// "Connected" notification handler will trigger nitro state channel address exchange
type NetworkNotifiee struct {
	ms *P2PMessageService
}

func (nn *NetworkNotifiee) Connected(n network.Network, conn network.Conn) {
	nn.ms.logger.Debug().Msgf("notification: connected to peer %s", conn.RemotePeer().String())
	go nn.ms.sendPeerInfo(conn.RemotePeer(), false)
}

func (nn NetworkNotifiee) Disconnected(n network.Network, conn network.Conn) {
	nn.ms.logger.Debug().Msgf("notification: disconnected from peer: %s", conn.RemotePeer().String())
}

func (nn NetworkNotifiee) Listen(network.Network, multiaddr.Multiaddr)      {}
func (nn NetworkNotifiee) ListenClose(network.Network, multiaddr.Multiaddr) {}
