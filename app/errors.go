package app

import "github.com/statechannels/go-nitro/internal"

var (
	ErrAppNotRegistered = internal.NewError("app not registered")
	ErrChannelNotFound  = internal.NewError("channel not found")

	ErrUnknownRequestType = internal.NewError("unknown request type")
)
