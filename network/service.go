package network

import (
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/network/serde"
	"github.com/statechannels/go-nitro/network/transport"
)

type NetworkService struct {
	Logger     zerolog.Logger
	Connection transport.Connection

	errorHander func(uint64, []byte)
}

func NewNetworkService(con transport.Connection) *NetworkService {
	p := &NetworkService{
		Connection: con,
	}

	return p
}

func (p *NetworkService) Subscribe(topic string, handler func([]byte) []byte) {
	err := p.Connection.Subscribe(topic, handler)
	if err != nil {
		panic(err)
	}
}

func (p *NetworkService) RegisterErrorHandler(handler func(uint64, []byte)) {
	p.errorHander = handler
	p.Logger.Trace().Msg("registered error handler")
}

func (p *NetworkService) UnregisterErrorHandler(method serde.RequestMethod) {
	p.errorHander = nil
	p.Logger.Trace().Msg("unregistered error handler")
}

func (p *NetworkService) Close() {
	p.Connection.Close()
}
