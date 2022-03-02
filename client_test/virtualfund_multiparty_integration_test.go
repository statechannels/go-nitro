package client_test

import (
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// TestMultiPartyVirtualFundIntegration tests the scenario where Alice creates virtual channels with Bob and Brian using Irene as the intermediary.
func TestMultiPartyVirtualFundIntegration(t *testing.T) {
	t.Skip()

	// Set up logging
	logDestination, err := os.OpenFile("virtualfund_multiparty_client_test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Reset log destination file
	err = logDestination.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}

	brianKey := common.Hex2Bytes("0aca28ba64679f63d71e671ab4dbb32aaa212d4789988e6ca47da47601c18fe2")
	brian := common.HexToAddress("0xB2B22ec3889d11f2ddb1A1Db11e80D20EF367c01")

	aliceKey := common.Hex2Bytes(`2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d`)
	alice := common.HexToAddress(`0xAAA6628Ec44A8a742987EF3A114dDFE2D4F7aDCE`)

	bobKey := common.Hex2Bytes(`0279651921cd800ac560c21ceea27aab0107b67daf436cdd25ce84cad30159b4`)
	bob := common.HexToAddress(`0xBBB676f9cFF8D242e9eaC39D063848807d3D1D94`)

	ireneKey := common.Hex2Bytes(`febb3b74b0b52d0976f6571d555f4ac8b91c308dfa25c7b58d1e6a7c3f50c781`)
	irene := common.HexToAddress(`0x111A00868581f73AB42FEEF67D235Ca09ca1E8db`)

	chain := chainservice.NewMockChain([]types.Address{alice, bob, irene, brian})

	clientAlice, aliceMS := setupClient(aliceKey, chain, logDestination)
	clientBob, bobMS := setupClient(bobKey, chain, logDestination)

	clientBrian, brianMS := setupClient(brianKey, chain, logDestination)

	clientIrene, ireneMS := setupClient(ireneKey, chain, logDestination)
	connectMessageServices([]messageservice.TestMessageService{aliceMS, bobMS, ireneMS, brianMS})

	directlyFundALedgerChannel(clientAlice, clientIrene)
	directlyFundALedgerChannel(clientIrene, clientBob)
	directlyFundALedgerChannel(clientIrene, clientBrian)

	id := clientAlice.CreateVirtualChannel(bob, irene, types.Address{}, types.Bytes{}, createVirtualOutcome(alice, bob), big.NewInt(0))
	id2 := clientAlice.CreateVirtualChannel(brian, irene, types.Address{}, types.Bytes{}, createVirtualOutcome(alice, brian), big.NewInt(0))

	waitForCompletedObjectiveId(id, &clientBob)
	waitForCompletedObjectiveId(id2, &clientBrian)

	waitForCompletedObjectiveIds([]protocols.ObjectiveId{id, id2}, &clientAlice)
	waitForCompletedObjectiveIds([]protocols.ObjectiveId{id, id2}, &clientIrene)

}
