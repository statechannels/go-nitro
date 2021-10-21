package msg

import "github.com/statechannels/go-nitro/protocols"

type Msg interface {
	GetRecieveChan() chan protocols.Message
	Send(message protocols.Message)
}
