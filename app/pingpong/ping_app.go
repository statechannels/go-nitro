package pingpong

import (
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/internal"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// TODO: Extract common errors into a common package
var ErrInvalidRequestType = internal.NewError("invalid request type")

const AppId = "pingpong"

type Balance struct {
	Remaining *big.Int
	Paid      *big.Int
}

type PingPongApp struct {
	engine *engine.Engine
}

func NewPingPongApp(engine *engine.Engine) *PingPongApp {
	return &PingPongApp{engine}
}

func (a *PingPongApp) Type() string {
	return AppId
}

func (a *PingPongApp) Ping(ch *consensus_channel.ConsensusChannel) error {
	for i, p := range ch.Participants {
		if ch.MyIndex == uint(i) {
			continue
		}

		fmt.Println("Sending ping to ", p.Hex())
		a.engine.SendMessages([]protocols.Message{
			{
				To: p,
				AppRequests: []types.AppRequest{
					{
						AppId:       AppId,
						RequestType: "ping",
						ChannelId:   ch.Id,
						Data:        nil,
					},
				},
			},
		})

	}

	return nil
}

func (a *PingPongApp) HandlePing(ch *channel.Channel, data interface{}) {
	fmt.Println("Received ping")
	a.engine.SendMessages([]protocols.Message{
		{
			AppRequests: []types.AppRequest{
				{
					AppId:       AppId,
					RequestType: "pong",
					ChannelId:   ch.Id,
					Data:        nil,
				},
			},
		},
	})
}

func (a *PingPongApp) HandlePong(ch *channel.Channel, data interface{}) {
	fmt.Println("Received pong")
}

func (a *PingPongApp) HandleRequest(ch *channel.Channel, ty string, data interface{}) error {
	switch ty {

	case RequestTypePing:
		a.HandlePing(ch, data)

	case RequestTypePong:
		a.HandlePong(ch, data)

	default:
		return ErrInvalidRequestType
	}

	return nil
}
