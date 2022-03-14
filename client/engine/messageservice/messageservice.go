// Package messageservice is a messaging service responsible for routing messages to peers and relaying messages received from peers.
package messageservice // import "github.com/statechannels/go-nitro/client/messageservice"

import "github.com/statechannels/go-nitro/protocols"

type MessageService interface {
	// Out returns a chan for sending messages out of the MessageService
	Out() <-chan protocols.Message
	// In returns a chan for receiving messages into the MessageService
	In() chan<- protocols.Message
	Send(message protocols.Message)
}
