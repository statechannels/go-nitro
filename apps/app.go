package apps

import "github.com/statechannels/go-nitro/channel"

type App interface {
	Type() string

	HandleRequest(ch *channel.Channel, ty string, data interface{}) error
}
