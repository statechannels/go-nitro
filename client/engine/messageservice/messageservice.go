// Package messageservice is a messaging service responsible for routing messages to peers and relaying messages received from peers
package messageservice // import "github.com/statechannels/go-nitro/client/messageservice"

import "github.com/statechannels/go-nitro/protocols"

type MessageService interface {
	Out() <-chan protocols.Message // Returns a chan for recieving messages from the message service
	In() chan<- protocols.Message  // Returns a chan for sending messages to the message service
	Send(message protocols.Message)
}
