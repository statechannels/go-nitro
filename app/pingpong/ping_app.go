package pingpong

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/statechannels/go-nitro/app"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/types"
)

type PingPongApp struct {
	*app.App
	random   *rand.Rand
	reqInfos map[string]RequestInfo
}

type RequestInfo struct {
	sendAt   int64
	callback chan int64
}

func NewPingPongApp(engine *engine.Engine, myAddr types.Address) *PingPongApp {
	a := &PingPongApp{
		App:      app.NewApp("pingpong", myAddr),
		random:   rand.New(rand.NewSource(time.Now().UnixNano())),
		reqInfos: map[string]RequestInfo{},
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

func (a *PingPongApp) Ping(ch *consensus_channel.ConsensusChannel, out chan int64) error {
	reqId := fmt.Sprintf("%d", a.random.Uint64())
	ts := time.Now().UnixNano()
	a.reqInfos[reqId] = RequestInfo{ts, out}

	for _, p := range ch.FixedPart().Participants {
		if p == a.MyAddress {
			continue
		}

		fmt.Println("Sending ping to", p, "with reqId", reqId, "and timestamp", ts)

		a.SendRequest(ch, p, RequestTypePing, reqId)
	}

	return nil
}

func (a *PingPongApp) handlePing(
	ch *consensus_channel.ConsensusChannel,
	from types.Address,
	data interface{},
) {
	reqId := data.(string)

	fmt.Println("Received ping from", from, "with reqId", reqId, "and timestamp")

	a.SendRequest(ch, from, RequestTypePong, reqId)
}

func (a *PingPongApp) handlePong(
	ch *consensus_channel.ConsensusChannel,
	from types.Address,
	data interface{},
) {
	ts := time.Now().UnixNano()
	reqId := data.(string)
	reqInfos, ok := a.reqInfos[reqId]
	if !ok {
		fmt.Printf("ERR: Received pong with unknown id %s\n", reqId)
		return
	}
	delete(a.reqInfos, reqId)
	rtt := ts - reqInfos.sendAt
	reqInfos.callback <- rtt
	fmt.Println("Received pong from", from, "with reqId", reqId, "RTT", rtt, "ns")
}
