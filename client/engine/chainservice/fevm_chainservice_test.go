package chainservice

import (
	"context"
	"log"
	"testing"
)

// Run this using
// go test ./client/engine/chainservice -run TestFevmChainService -v -p=1
func TestFevmChainService(t *testing.T) {
	t.Skip() // should not run in CI because it depends on wallaby testnet
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
		t.Fatal(err)
	}
	t.Log("Filecoin.MpoolGetNonce call returned", nonce)

	err = fcs.DeployAdjudicator()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("chain service stored adjudicator address:", fcs.naAddress)
}
