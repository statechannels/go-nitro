package pingpong

import (
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/internal"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// TODO: Extract common errors into a common package
var ErrInvalidRequestType = internal.NewError("invalid request type")

type Balance struct {
	Remaining *big.Int
	Paid      *big.Int
}

type PingPongApp struct {
	engine *engine.Engine

	myAddress types.Address
}

func NewPingPongApp(engine *engine.Engine, myAddr types.Address) *PingPongApp {
	return &PingPongApp{
		engine: engine,

		myAddress: myAddr,
	}
}

func (a *PingPongApp) Id() string {
	return "pingpong"
}

func (a *PingPongApp) Ping(ch *consensus_channel.ConsensusChannel) error {
	for _, p := range ch.FixedPart().Participants {
		if p == a.myAddress {
			continue
		}

		fmt.Println("Sending ping to ", p.Hex())

		a.engine.SendMessages([]protocols.Message{
			{
				To: p,

				AppRequests: []types.AppRequest{
					{
						From: a.myAddress,

						AppId:       a.Id(),
						RequestType: RequestTypePing,
						ChannelId:   ch.Id,

						Data: nil,
					},
				},
			},
		})

	}

	return nil
}

func (a *PingPongApp) handlePing(
	ch *consensus_channel.ConsensusChannel,
	from types.Address,
	data interface{},
) {
	fmt.Println("Received ping")

	a.engine.SendMessages([]protocols.Message{
		{
			To: from,

			AppRequests: []types.AppRequest{
				{
					AppId:       a.Id(),
					RequestType: RequestTypePong,
					ChannelId:   ch.Id,
					Data:        nil,
				},
			},
		},
	})
}

func (a *PingPongApp) handlePong(
	ch *consensus_channel.ConsensusChannel,
	from types.Address,
	data interface{},
) {
	fmt.Println("Received pong")
}

func (a *PingPongApp) HandleRequest(
	ch *consensus_channel.ConsensusChannel,
	from types.Address,
	ty string,
	data interface{},
) error {
	switch ty {
	case RequestTypePing:
		a.handlePing(ch, from, data)

	case RequestTypePong:
		a.handlePong(ch, from, data)

	default:
		return ErrInvalidRequestType
	}

	return nil
}
