package outcome

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestIsExternalDestination(t *testing.T) {

	e := common.HexToHash("0x00000000000000000000000096f7123E3A80C9813eF50213ADEd0e4511CB820f")
	i := common.HexToHash("0x6f7123E3A80C9813eF50213A96f7123E3A80C9813eF50213ADEd0e4511CB820f")

	if !IsExternalDestination(e) {
		t.Errorf("Received bytes %x was declared internal, when it is external", e)
	}

	if IsExternalDestination(i) {
		t.Errorf("Received bytes %x was declared external, when it is internal", i)
	}

}
