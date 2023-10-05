// Package messageservice is a messaging service responsible for routing messages to peers and relaying messages received from peers.
package messageservice // import "github.com/statechannels/go-nitro/node/messageservice"

import (
	p2pms "github.com/statechannels/go-nitro/node/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/protocols"
)

type MessageService interface {
	// P2PMessages returns a chan for receiving messages from the message service
	P2PMessages() <-chan protocols.Message
	// SignRequests returns a chan for receiving signature requests from the message service
	SignRequests() <-chan p2pms.SignatureRequest
	// Send is for sending messages with the message service
	Send(protocols.Message) error
	// Close closes the message service
	Close() error
}
