package pingpong

import (
	"fmt"

	"github.com/statechannels/go-nitro/app"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/types"
)

type PingPongApp struct {
	*app.App
}

func NewPingPongApp(engine *engine.Engine, myAddr types.Address) *PingPongApp {
	a := &PingPongApp{
		App: app.NewApp("pingpong", myAddr),
	}

	a.RegisterRequestHandler(
		RequestTypePing,
		func(ch *consensus_channel.ConsensusChannel, from types.Address, data interface{}) error {
			a.handlePing(ch, from, data)

			return nil
		},
	)

	a.RegisterRequestHandler(
		RequestTypePong,
		func(ch *consensus_channel.ConsensusChannel, from types.Address, data interface{}) error {
			a.handlePong(ch, from, data)

			return nil
		},
	)

	return a
}

func (a *PingPongApp) Id() string {
	return "pingpong"
}

func (a *PingPongApp) Ping(ch *consensus_channel.ConsensusChannel) error {
	for _, p := range ch.FixedPart().Participants {
		if p == a.MyAddress {
			continue
		}

		fmt.Println("Sending ping to", p, "with data", 42)

		a.SendRequest(ch, p, RequestTypePing, 42)

	}

	return nil
}

func (a *PingPongApp) handlePing(
	ch *consensus_channel.ConsensusChannel,
	from types.Address,
	data interface{},
) {
	fmt.Println("Received ping from", from, "with data", data)

	fmt.Println("Sending pong to", from, "with data", data.(float64)+1)

	a.SendRequest(ch, from, RequestTypePong, data.(float64)+1)
}

func (a *PingPongApp) handlePong(
	ch *consensus_channel.ConsensusChannel,
	from types.Address,
	data interface{},
) {
	fmt.Println("Received pong from", from, "with data", data)
}
