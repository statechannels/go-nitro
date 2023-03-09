// Package messageservice is a messaging service responsible for routing messages to peers and relaying messages received from peers.
package messageservice // import "github.com/statechannels/go-nitro/client/messageservice"

import "github.com/statechannels/go-nitro/protocols"

type MessageService interface {
	// Out returns a chan for receiving messages from the message service
	Out() <-chan protocols.Message
	// Send is for sending messages with the message service
	Send(protocols.Message)
	// Close closes the message service
	Close() error
}
