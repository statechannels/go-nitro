package network

import (
	"github.com/statechannels/go-nitro/internal"
)

var (
	// Service
	ErrServiceClosed = internal.NewError("service closed")

	// NetworkServiceConnection
	ErrPeerClosed = internal.NewError("peer closed")

	// Protocol
	ErrRequestError  = internal.NewError("request error")
	ErrResponseError = internal.NewError("response error")
)
