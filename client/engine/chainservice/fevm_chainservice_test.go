package chainservice

import (
	"context"
	"log"
	"testing"
)

func TestFevmChainService(t *testing.T) {

	fcs, err := NewFevmChainService("https://wallaby.node.glif.io/rpc/v0", "9182b5bf5b9c966e001934ebaf008f65516290cef6e3069d11e718cbd4336aae", log.Default().Writer())
	if err != nil {
		t.Fatal(fcs)
	}
	blockNumber, err := fcs.chain.BlockNumber(context.Background())
	if err != nil {
		t.Fatal(fcs)
	}
	t.Log("Block number is", blockNumber)
	nonce, err := fcs.filecoinNonce()

	if err != nil {
		t.Fatal(fcs)
	}
	t.Log("Filecoin.MpoolGetNonce call returned", nonce)
}
