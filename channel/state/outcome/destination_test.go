package outcome

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestIsExternal(t *testing.T) {

	e := common.HexToHash("0x00000000000100000000000096f7123E3A80C9813eF50213ADEd0e4511CB820f")
	i := common.HexToHash("0x6f7123E3A80C9813eF50213A96f7123E3A80C9813eF50213ADEd0e4511CB820f")

	if !IsExternal(e) {
		t.Errorf("Received bytes %x was declared internal, when it is external", e)
	}

	if IsExternal(i) {
		t.Errorf("Received bytes %x was declared external, when it is internal", i)
	}

}

func TestToAddress(t *testing.T) {
	areExternal := []common.Hash{
		common.HexToHash("0x000000000000000000000000aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000096f7123E3A80C9813eF50213ADEd0e4511CB820f"),
	}

	areNotExternal := []common.Hash{
		common.HexToHash("0x000000000000000000000001aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		common.HexToHash("0x100000000000000000000000aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		common.HexToHash("0x0000000000b0000000000000aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
	}

	for _, extAddress := range areExternal {
		if _, err := ToAddress(extAddress); err != nil {
			t.Errorf("expected to convert %x to an external address, but failed", extAddress)
		}
	}

	for _, notExtAddress := range areNotExternal {
		if _, err := ToAddress(notExtAddress); err == nil {
			t.Errorf("expected to fail when converting %x to an external address, but succeeded", notExtAddress)
		}
	}
}
