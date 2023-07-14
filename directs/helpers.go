package directs

import (
	"bytes"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
)

// GetLatestProposedStateByAppData returns the latest proposed state
// (state that signed by one party and has bigger turn num than supported one)
// with the given appData
func GetLatestProposedStateByAppData(ch *channel.Channel, appData []byte) *state.State {
	latestSupportedState, err := ch.LatestSupportedState()
	if err != nil {
		return nil
	}

	latestSignedState, err := ch.LatestSignedState()
	if err != nil {
		return nil
	}

	for i := latestSupportedState.TurnNum + 1; i <= latestSignedState.State().TurnNum; i++ {
		s := ch.SignedStateForTurnNum[i].State()
		if bytes.Equal(s.AppData, appData) {
			return &s
		}
	}

	return nil
}
