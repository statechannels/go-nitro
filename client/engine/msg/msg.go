// Package msg is a messaging service responsible for routing messages to peers and relaying messages recieved from peers
package msg // import "github.com/statechannels/go-nitro/client/msg"

import "github.com/statechannels/go-nitro/protocols"

type Msg interface {
	GetRecieveChan() chan protocols.Message
	GetSendChan() chan protocols.Message
	Send(message protocols.Message)
}
