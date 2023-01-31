package network

import (
	"errors"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/network/serde"
	"github.com/statechannels/go-nitro/network/transport"
)

const (
	DirectFundRequestMethod    = "direct_fund"
	DirectDefundRequestMethod  = "direct_defund"
	VirtualFundRequestMethod   = "virtual_fund"
	VirtualDefundRequestMethod = "virtual_defund"
)

type NetworkService struct {
	Logger     zerolog.Logger
	Connection transport.Connection

	handlerRequest sync.Map
	handlerError   sync.Map
	responseHander func([]byte)
}

func NewNetworkService(con transport.Connection) *NetworkService {
	p := &NetworkService{
		Connection: con,
	}

	go p.handleMessages()

	return p
}

func (p *NetworkService) RegisterRequestHandler(method string, handler func([]byte)) {
	p.handlerRequest.Store(method, handler)
	p.Logger.Trace().Str("method", method).Msg("registered request handler")
}

func (p *NetworkService) UnregisterRequestHandler(method string) {
	p.handlerRequest.Delete(method)
	p.Logger.Trace().Str("method", method).Msg("unregistered request handler")
}

func (p *NetworkService) RegisterErrorHandler(method string, handler func([]byte)) {
	p.handlerError.Store(method, handler)
	p.Logger.Trace().Str("method", method).Msg("registered error handler")
}

func (p *NetworkService) UnregisterErrorHandler(method string) {
	p.handlerError.Delete(method)
	p.Logger.Trace().Str("method", method).Msg("unregistered error handler")
}

func (p *NetworkService) RegisterResponseHandler(handler func([]byte)) {
	p.responseHander = handler
	p.Logger.Trace().Msg("registered response handler")
}

func (p *NetworkService) UnregisterResponseHandler() {
	p.responseHander = nil
	p.Logger.Trace().Msg("unregistered response handler")
}

// TODO: implement (un)registerPublicEventHandler
// TODO: implement (un)registerPrivateEventHandler

func (p *NetworkService) handleMessages() {
	for {
		data, err := p.Connection.Recv()
		if err != nil {
			if errors.Is(err, transport.ErrConnectionClosed) {
				p.Logger.Info().Msg("connection closed")
				break
			}

			// TODO: handle error
			p.Logger.Fatal().Err(err).Msg("failed to receive message")
		}

		msg, messageType, err := serde.Deserialize(data)

		if err != nil {
			p.Logger.Error().Err(err).Msg("failed to deserialize message")
			return
		}

		// NOTE: we do not hande messages in a separate goroutine
		// to ensure that messages are handled in the order they are received
		// and to avoid inconsistencies in the state of the peer
		p.handleMessage(msg.Method, messageType, data)
	}
}

func (p *NetworkService) SendMessage(method string, data []byte) {
	// FIXME: we can use one topic per app, but they have to be different
	topic := fmt.Sprintf("nitro.%s", method)
	p.Connection.Send(topic, data)

	p.Logger.Trace().
		Str("method", method).
		Msg("sent message")
}

func (p *NetworkService) getHandler(method string, messageType serde.MessageType) func([]byte) {
	switch messageType {
	case serde.TypeRequest:
		function, ok := p.handlerRequest.Load(method)
		if ok {
			return function.(func([]byte))
		}

	case serde.TypeResponse:
		return p.responseHander

	case serde.TypeError:
		function, ok := p.handlerError.Load(method)
		if ok {
			return function.(func([]byte))
		}
	}
	// TODO: case handlerPublicEvent
	// TODO: case handlerPrivateEvent

	return nil
}

func (p *NetworkService) handleMessage(method string, messageType serde.MessageType, data []byte) {
	p.Logger.Trace().
		Str("method", method).
		Msg("received message")

	h := p.getHandler(method, messageType)

	if h == nil {
		p.Logger.Error().
			Str("method", method).
			Msg("missing handler")
		return
	}

	h(data)
}

func (p *NetworkService) Close() {
	p.Connection.Close()
}
