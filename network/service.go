package network

import (
	"errors"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/network/serde"
	"github.com/statechannels/go-nitro/network/transport"
)

type NetworkService struct {
	Logger     zerolog.Logger
	Connection transport.Connection

	handlerRequest  sync.Map
	responseHandler func(uint64, []byte)
	errorHander     func(uint64, []byte)
}

func NewNetworkService(con transport.Connection) *NetworkService {
	p := &NetworkService{
		Connection: con,
	}

	go p.handleMessages()

	return p
}

func (p *NetworkService) RegisterRequestHandler(method serde.RequestMethod, handler func(uint64, []byte)) {
	p.handlerRequest.Store(string(method), handler)
	p.Logger.Trace().Str("method", string(method)).Msg("registered request handler")
}

func (p *NetworkService) UnregisterRequestHandler(method serde.RequestMethod) {
	p.handlerRequest.Delete(string(method))
	p.Logger.Trace().Str("method", string(method)).Msg("unregistered request handler")
}

func (p *NetworkService) RegisterErrorHandler(handler func(uint64, []byte)) {
	p.errorHander = handler
	p.Logger.Trace().Msg("registered error handler")
}

func (p *NetworkService) UnregisterErrorHandler(method serde.RequestMethod) {
	p.errorHander = nil
	p.Logger.Trace().Msg("unregistered error handler")
}

func (p *NetworkService) RegisterResponseHandler(handler func(uint64, []byte)) {
	p.responseHandler = handler
	p.Logger.Trace().Msg("registered response handler")
}
func (p *NetworkService) UnregisterResponseHandler() {
	p.responseHandler = nil
	p.Logger.Trace().Msg("unregistered response handler")
}

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
		p.handleMessage(msg.Id, msg.Method, messageType, data)
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

func (p *NetworkService) getHandler(method string, messageType serde.MessageType) func(uint64, []byte) {
	switch messageType {
	case serde.TypeRequest:
		function, ok := p.handlerRequest.Load(method)
		if ok {
			return function.(func(uint64, []byte))
		}
	case serde.TypeResponse:
		return p.responseHandler
	case serde.TypeError:
		return p.errorHander
	}

	return nil
}

func (p *NetworkService) handleMessage(id uint64, method serde.RequestMethod, messageType serde.MessageType, data []byte) {
	p.Logger.Trace().
		Str("method", string(method)).
		Uint64("id", id).
		Msg("received message")

	h := p.getHandler(string(method), messageType)

	if h == nil {
		p.Logger.Error().
			Str("method", string(method)).
			Uint64("id", id).
			Msg("missing handler")
		return
	}

	h(id, data)
}

func (p *NetworkService) Close() {
	p.Connection.Close()
}
