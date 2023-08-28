package policy

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/types"
)

func TestDenyListPolicy(t *testing.T) {
	badBob := crypto.GetAddressFromSecretKeyBytes(common.Hex2Bytes(`92df90ad792e7987539f2bcbae9e4d3e539fd6d919eb185cc2f2beb81adb473c`))
	policy := NewDenyListPolicy([]types.Address{badBob})

	df := testdata.Objectives.Directfund.GenericDFO()
	if !policy.ShouldApprove(&df) {
		t.Fatal("Policy should approve objective since bad bob is not a participant")
	}

	df.C.Participants = append(df.C.Participants, badBob)
	if policy.ShouldApprove(&df) {
		t.Fatal("Policy should deny objective since bad bob is a participant")
	}

	vf := testdata.Objectives.Virtualfund.GenericVFO()
	if !policy.ShouldApprove(&vf) {
		t.Fatal("Policy should approve objective since bad bob is not a participant")
	}

	vf.V.Participants = append(vf.V.Participants, badBob)
	if policy.ShouldApprove(&vf) {
		t.Fatal("Policy should deny objective since bad bob is a participant")
	}
}

func TestAllowListPolicy(t *testing.T) {
	policy := NewAllowListPolicy([]types.Address{testactors.Alice.Address(), testactors.Bob.Address()})

	df := testdata.Objectives.Directfund.GenericDFO()
	if !policy.ShouldApprove(&df) {
		t.Fatal("Policy should approve objective since only alice and bob are participants")
	}

	df.C.Participants = append(df.C.Participants, testactors.Irene.Address())
	if policy.ShouldApprove(&df) {
		t.Fatal("Policy should deny objective since irene is not on the allow list")
	}

	vf := testdata.Objectives.Virtualfund.GenericVFO()
	if policy.ShouldApprove(&vf) {
		t.Fatal("Policy should approve objective since irene is not a participant")
	}

	policy = NewAllowListPolicy([]types.Address{testactors.Alice.Address(), testactors.Bob.Address(), testactors.Irene.Address()})

	if !policy.ShouldApprove(&vf) {
		t.Fatal("Policy should approve objective since all participants are on the allow list")
	}
}
