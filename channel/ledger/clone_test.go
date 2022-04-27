package ledger

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

	cc, err := newLedgerChannel(fp(), Leader, 0, outcome, sigs)

	if err != nil {
		t.Fatal(err)
	}

	clone := cc.Clone()

	compareLedgerChannels := func(a, b LedgerChannel) string {
		return cmp.Diff(&a, &b,
			cmp.AllowUnexported(
				LedgerChannel{},
				Vars{},
				LedgerOutcome{},
				Guarantee{},
				Balance{},
				big.Int{},
				state.SignedState{}))
	}

	if diff := compareLedgerChannels(cc, *clone); diff != "" {
		t.Errorf("Clone: mismatch (-want +got):\n%s", diff)
	}

}
