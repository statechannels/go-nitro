package app

import (
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type RequestHandler func(
	ch *consensus_channel.ConsensusChannel,
	from types.Address,
	data interface{},
) error

type App struct {
	id string

	MyAddress types.Address

	requestHandlers map[string]RequestHandler

	MessageCh chan protocols.Message
}

func NewApp(id string, myAddr types.Address) *App {
	return &App{
		id: id,

		MyAddress: myAddr,

		requestHandlers: make(map[string]RequestHandler),

		MessageCh: make(chan protocols.Message),
	}
}

// NOTE: i would revert to "Type" here, since there is now way for the sender to know what
// will be the app id of the receiver...
func (a *App) Id() string {
	return a.id
}

func (a *App) RegisterRequestHandler(ty string, handler RequestHandler) {
	a.requestHandlers[ty] = handler
}

func (a *App) HandleRequest(
	ch *consensus_channel.ConsensusChannel,
	from types.Address,
	ty string,
	data interface{},
) error {
	handler, ok := a.requestHandlers[ty]
	if !ok {
		return ErrUnknownRequestType
	}

	return handler(ch, from, data)
}

func (a *App) SendRequest(
	ch *consensus_channel.ConsensusChannel,
	to types.Address,
	ty string,
	data interface{},
) {
	a.MessageCh <- protocols.Message{
		To: to,

		AppRequests: []types.AppRequest{
			{
				From: a.MyAddress,

				AppId:       a.Id(),
				RequestType: ty,
				ChannelId:   ch.Id,

				Data: data,
			},
		},
	}
}
