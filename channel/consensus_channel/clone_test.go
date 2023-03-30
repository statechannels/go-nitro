package consensus_channel

import (
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

func TestClone(t *testing.T) {
	outcome := makeOutcome(
		allocation(alice, aBal),
		allocation(bob, bBal),
		guarantee(vAmount, types.Destination{7}, alice, bob),
	)

	initialVars := Vars{Outcome: outcome, TurnNum: 0}
	aliceSig, _ := initialVars.AsState(fp()).Sign(alice.PrivateKey)
	bobsSig, _ := initialVars.AsState(fp()).Sign(bob.PrivateKey)
	sigs := [2]state.Signature{aliceSig, bobsSig}

	cc, err := newConsensusChannel(fp(), Leader, 0, outcome, sigs)
	if err != nil {
		t.Fatal(err)
	}

	clone := cc.Clone()

	compareConsensusChannels := func(a, b ConsensusChannel) string {
		return cmp.Diff(&a, &b,
			cmp.AllowUnexported(
				ConsensusChannel{},
				Vars{},
				LedgerOutcome{},
				Guarantee{},
				Balance{},
				big.Int{},
				state.SignedState{}))
	}

	if diff := compareConsensusChannels(cc, *clone); diff != "" {
		t.Errorf("Clone: mismatch (-want +got):\n%s", diff)
	}
}
