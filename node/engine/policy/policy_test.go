package policy

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/types"
)

func TestDenyListPolicy(t *testing.T) {
	badBob := crypto.GetAddressFromSecretKeyBytes(common.Hex2Bytes(`92df90ad792e7987539f2bcbae9e4d3e539fd6d919eb185cc2f2beb81adb473c`))
	policy := NewDenyListPolicy([]types.Address{badBob}, Ledger)

	df := testdata.Objectives.Directfund.GenericDFO()
	if !policy.ShouldApprove(&df) {
		t.Fatal("Policy should approve objective since bad bob is not a participant")
	}

	df.C.Participants = append(df.C.Participants, badBob)
	if policy.ShouldApprove(&df) {
		t.Fatal("Policy should deny objective since bad bob is a participant")
	}

	policy = NewDenyListPolicy([]types.Address{badBob}, Payment)

	vf := testdata.Objectives.Virtualfund.GenericVFO()
	if !policy.ShouldApprove(&vf) {
		t.Fatal("Policy should approve objective since bad bob is not a participant")
	}

	vf.V.Participants = append(vf.V.Participants, badBob)
	if policy.ShouldApprove(&vf) {
		t.Fatal("Policy should deny objective since bad bob is a participant")
	}

	ddf := testdata.Objectives.Directfund.GenericDFO()
	if !policy.ShouldApprove(&ddf) {
		t.Fatal("Policy should always approve direct defund objectives")
	}
}

func TestAllowListPolicy(t *testing.T) {
	policy := NewAllowListPolicy([]types.Address{testactors.Alice.Address(), testactors.Bob.Address()}, Ledger)

	df := testdata.Objectives.Directfund.GenericDFO()
	if !policy.ShouldApprove(&df) {
		t.Fatal("Policy should approve objective since only alice and bob are participants")
	}

	df.C.Participants = append(df.C.Participants, testactors.Irene.Address())
	if policy.ShouldApprove(&df) {
		t.Fatal("Policy should deny objective since irene is not on the allow list")
	}

	policy = NewAllowListPolicy([]types.Address{testactors.Alice.Address(), testactors.Bob.Address()}, Payment)

	vf := testdata.Objectives.Virtualfund.GenericVFO()
	if policy.ShouldApprove(&vf) {
		t.Fatal("Policy should approve objective since irene is not a participant")
	}
	policy = NewAllowListPolicy([]types.Address{testactors.Alice.Address(), testactors.Bob.Address(), testactors.Irene.Address()}, Payment)

	if !policy.ShouldApprove(&vf) {
		t.Fatal("Policy should approve objective since all participants are on the allow list")
	}

	ddf := testdata.Objectives.Directfund.GenericDFO()
	if !policy.ShouldApprove(&ddf) {
		t.Fatal("Policy should always approve direct defund objectives")
	}
}

func TestFairOutcomePolicy(t *testing.T) {
	policy := NewFairOutcomePolicy(testactors.Alice.Address())

	df := testdata.Objectives.Directfund.GenericDFO()
	if !policy.ShouldApprove(&df) {
		t.Fatal("Policy should approve objective since the outcome is fair")
	}

	df = testdata.GenerateDFOFromOutcome(testdata.Outcomes.Create(testactors.Alice.Address(), testactors.Bob.Address(), 1, 5, common.Address{}))
	if policy.ShouldApprove(&df) {
		t.Fatal("Policy should reject the outcome as unfair")
	}

	vf := testdata.GenerateVFOFromOutcome(testdata.Outcomes.Create(testactors.Alice.Address(), testactors.Bob.Address(), 10, 0, common.Address{}))
	if !policy.ShouldApprove(&vf) {
		t.Fatal("Policy should approve objective since the outcome is fair")
	}
	vf = testdata.GenerateVFOFromOutcome(testdata.Outcomes.Create(testactors.Alice.Address(), testactors.Bob.Address(), 5, 1, common.Address{}))
	if policy.ShouldApprove(&vf) {
		t.Fatal("Policy should reject objective since the outcome is unfair")
	}

	ddf := testdata.Objectives.Directfund.GenericDFO()
	if !policy.ShouldApprove(&ddf) {
		t.Fatal("Policy should always approve direct defund objectives")
	}
}

func TestMaxSpendPolicy(t *testing.T) {
	policy := NewLedgerChannelMaxSpendPolicy(testactors.Alice.Address(), types.Funds{types.Address{}: big.NewInt(5)})

	df := testdata.GenerateDFOFromOutcome(testdata.Outcomes.Create(testactors.Alice.Address(), testactors.Bob.Address(), 5, 5, common.Address{}))
	if !policy.ShouldApprove(&df) {
		t.Fatal("Policy should approve objective since the outcome required by the objective is less than the max spend")
	}

	df = testdata.GenerateDFOFromOutcome(testdata.Outcomes.Create(testactors.Alice.Address(), testactors.Bob.Address(), 10, 5, common.Address{}))
	if policy.ShouldApprove(&df) {
		t.Fatal("Policy should reject objective since the outcome requires more than the max spend")
	}

	policy = NewPaymentChannelMaxSpendPolicy(testactors.Alice.Address(), types.Funds{types.Address{}: big.NewInt(5)})
	vf := testdata.GenerateVFOFromOutcome(testdata.Outcomes.Create(testactors.Alice.Address(), testactors.Bob.Address(), 5, 0, common.Address{}))
	if !policy.ShouldApprove(&vf) {
		t.Fatal("Policy should approve objective since the outcome required by the objective is less than the max spend")
	}
	vf = testdata.GenerateVFOFromOutcome(testdata.Outcomes.Create(testactors.Alice.Address(), testactors.Bob.Address(), 10, 0, common.Address{}))
	if policy.ShouldApprove(&vf) {
		t.Fatal("Policy should reject objective since the outcome requires more than the max spend")
	}

	ddf := testdata.Objectives.Directfund.GenericDFO()
	if !policy.ShouldApprove(&ddf) {
		t.Fatal("Policy should always approve direct defund objectives")
	}
}

func TestPolicies(t *testing.T) {
	policies := NewPolicies(
		NewFairOutcomePolicy(testactors.Alice.Address()),
		NewAllowListPolicy([]types.Address{testactors.Alice.Address(), testactors.Bob.Address()}, Ledger),
		NewLedgerChannelMaxSpendPolicy(testactors.Alice.Address(), types.Funds{types.Address{}: big.NewInt(5)}),
		NewPaymentChannelMaxSpendPolicy(testactors.Alice.Address(), types.Funds{types.Address{}: big.NewInt(5)}),
	)

	df := testdata.GenerateDFOFromOutcome(testdata.Outcomes.Create(testactors.Alice.Address(), testactors.Bob.Address(), 5, 5, common.Address{}))
	if !policies.ShouldApprove(&df) {
		t.Fatal("Policies should approve objective because it is fair, alice and bob are on the allow list, and the outcome is less than the max spend")
	}

	df = testdata.GenerateDFOFromOutcome(testdata.Outcomes.Create(testactors.Alice.Address(), testactors.Bob.Address(), 10, 5, common.Address{}))
	if policies.ShouldApprove(&df) {
		t.Fatal("Policies should reject objective because of the ledger max spend policy")
	}
}
