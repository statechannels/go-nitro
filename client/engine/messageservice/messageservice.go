// Package messageservice is a messaging service responsible for routing messages to peers and relaying messages received from peers.
package messageservice // import "github.com/statechannels/go-nitro/client/messageservice"

import "github.com/statechannels/go-nitro/protocols"

type MessageService interface {
	// Inbox returns a chan for recieving messages from the message service
	Inbox() <-chan protocols.Message
	// Outbox returns a chan for sending messages to the message service
	Outbox() chan<- protocols.Message
}
