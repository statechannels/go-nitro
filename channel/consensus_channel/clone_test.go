package consensus_channel

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
)

func TestClone(t *testing.T) {

	cc, _ := newConsensusChannel(state.TestState.FixedPart(), 0, 1, ledgerOutcome(), [2]state.Signature{
		{
			R: common.Hex2Bytes(`704b3afcc6e702102ca1af3f73cf3b37f3007f368c40e8b81ca823a65740a053`),
			S: common.Hex2Bytes(`14040ad4c598dbb055a50430142a13518e1330b79d24eed86fcbdff1a7a95589`),
			V: byte(27),
		},
		{
			R: common.Hex2Bytes(`704b3afcc6e702102ca1af3f73cf3b37f3007f368c40e8b81ca823a65740a053`),
			S: common.Hex2Bytes(`14040ad4c598dbb055a50430142a13518e1330b79d24eed86fcbdff1a7a95589`),
			V: byte(27),
		},
	})

	clone := cc.Clone()

	compareConsensusChannels := func(a, b ConsensusChannel) string {
		return cmp.Diff(&a, &b, cmp.AllowUnexported(ConsensusChannel{}, channel.Channel{}, big.Int{}, state.SignedState{}))
	}

	if diff := compareConsensusChannels(cc, *clone); diff != "" {
		t.Errorf("Clone: mismatch (-want +got):\n%s", diff)
	}

}
