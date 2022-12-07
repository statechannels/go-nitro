package app

import (
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/types"
)

type App interface {
	Id() string

	HandleRequest(
		ch *consensus_channel.ConsensusChannel,
		from types.Address,
		ty string,
		data interface{},
	) error
}
