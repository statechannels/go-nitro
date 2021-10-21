package msg

import "github.com/statechannels/go-nitro/protocols"

type Msg interface {
	Send(message protocols.Message)
}
