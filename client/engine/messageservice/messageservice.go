// Package messageservice is a messaging service responsible for routing messages to peers and relaying messages received from peers
package messageservice // import "github.com/statechannels/go-nitro/client/messageservice"

import "github.com/statechannels/go-nitro/protocols"

type MessageService interface {
	GetReceiveChan() <-chan protocols.Message
	GetSendChan() chan<- protocols.Message
	Send(message protocols.Message)
}
