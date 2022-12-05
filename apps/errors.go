package apps

import "github.com/statechannels/go-nitro/internal"

var (
	ErrAppNotRegistered = internal.NewError("app not registered")
	ErrChannelNotFound  = internal.NewError("channel not found")
)
