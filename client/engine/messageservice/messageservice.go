// Package messageservice is a messaging service responsible for routing messages to peers and relaying messages received from peers.
package messageservice // import "github.com/statechannels/go-nitro/client/messageservice"

import "github.com/statechannels/go-nitro/internal/protocols"

type MessageService interface {
	// Out returns a chan for recieving messages from the message service
	Out() <-chan protocols.Message
	// In returns a chan for sending messages to the message service
	In() chan<- protocols.Message
}
