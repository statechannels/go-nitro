package node_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/internal/testactors"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/internal/testhelpers"
	"github.com/statechannels/go-nitro/node"
	"github.com/statechannels/go-nitro/node/engine/messageservice"
	p2pms "github.com/statechannels/go-nitro/node/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/node/query"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const maxHops = 12

func BenchmarkVirtualFund(b *testing.B) {
	dataFolder, cleanup := testhelpers.GenerateTempStoreFolder()
	defer cleanup()

	particpants := make([]TestParticipant, 2+maxHops)

	particpants[0] = TestParticipant{StoreType: MemStore, Actor: testactors.Alice}
	particpants[1] = TestParticipant{StoreType: MemStore, Actor: testactors.Bob}
	for i := 2; i < maxHops+2; i++ {
		// Generate a new private key
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			b.Fatal(err)
		}
		privateKeyBytes := crypto.FromECDSA(privateKey)

		particpants[i] = TestParticipant{StoreType: MemStore, Actor: testactors.Actor{
			PrivateKey: privateKeyBytes,
			Role:       2, // UNUSED
			Name:       testactors.ActorName("intermediary-" + fmt.Sprint(i)),
			Port:       0, // UNUSED
		}}
	}

	tc := TestCase{
		Description:    "Bench",
		Chain:          MockChain,
		MessageService: TestMessageService,
		NumOfChannels:  1,
		MessageDelay:   time.Millisecond * 100,
		LogName:        "bench",
		NumOfHops:      maxHops,
		NumOfPayments:  0,
		Participants:   particpants,
	}

	infra := setupSharedInfra(tc)

	msgServices := make([]messageservice.MessageService, 0)

	// Setup clients
	b.Log("Initalizing intermediary node(s)...")
	intermediaries := make([]node.Node, 0, tc.NumOfHops)
	bootPeers := make([]string, 0)
	for _, intermediary := range tc.Participants[2:] {
		clientI, msgI, multiAddr := setupIntegrationNode(tc, intermediary, infra, []string{}, dataFolder)

		intermediaries = append(intermediaries, clientI)
		msgServices = append(msgServices, msgI)
		bootPeers = append(bootPeers, multiAddr)
	}

	defer func() {
		for i := range intermediaries {
			intermediaries[i].Close()
		}
	}()
	b.Log("Intermediary node(s) setup complete")

	clientA, msgA, _ := setupIntegrationNode(tc, tc.Participants[0], infra, bootPeers, dataFolder)
	defer clientA.Close()
	msgServices = append(msgServices, msgA)

	clientB, msgB, _ := setupIntegrationNode(tc, tc.Participants[1], infra, bootPeers, dataFolder)
	defer clientB.Close()
	msgServices = append(msgServices, msgB)

	if tc.MessageService != TestMessageService {
		p2pServices := make([]*p2pms.P2PMessageService, len(tc.Participants))
		for i, msgService := range msgServices {
			p2pServices[i] = msgService.(*p2pms.P2PMessageService)
		}

		b.Log("Waiting for peer info exchange...")
		waitForPeerInfoExchange(p2pServices...)
		b.Log("Peer info exchange complete")
	}

	asset := common.Address{}

	// connect Alice to first intermediary
	aliceLedger := openLedgerChannel(b, clientA, intermediaries[0], asset)
	checkLedgerChannel(b, aliceLedger, initialLedgerOutcome(*clientA.Address, *intermediaries[0].Address, asset), query.Open, clientA)

	// connect Bob to every intermediary
	for _, intermediary := range intermediaries {
		bobLedger := openLedgerChannel(b, intermediary, clientB, asset)
		checkLedgerChannel(b, bobLedger, initialLedgerOutcome(*intermediary.Address, *clientB.Address, asset), query.Open, clientB)
	}

	// So that all of these paths exist:
	// A - I - B
	// A - I - I - B
	// A - I - I - I - B

	// connect intermediaries in a linear chain
	for i := 0; i+1 < len(intermediaries); i++ {
		openLedgerChannel(b, intermediaries[i], intermediaries[i+1], asset)
	}

	benchmarkVirtualfund := func(numHops int, b *testing.B) {
		path := clientAddresses(intermediaries[0:numHops])
		// b.Log("Path:" + fmt.Sprint(path))
		outcome := td.Outcomes.Create(testactors.Alice.Address(), testactors.Bob.Address(), virtualChannelDeposit, 0, types.Address{})
		response, err := clientA.CreatePaymentChannel(
			path,
			testactors.Bob.Address(),
			0,
			outcome,
		)
		if err != nil {
			b.Fatal(err)
		}
		objectiveId := response.Id
		// Wait for all the virtual channels to be ready
		waitForObjectives(b, clientA, clientB, intermediaries[0:numHops-1], []protocols.ObjectiveId{objectiveId})
	}

	for j := 1; j < int(tc.NumOfHops); j++ {
		b.Run("benchmark "+fmt.Sprint(j)+" hop virtual fund",
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					benchmarkVirtualfund(j, b)
				}
			},
		)
	}
}
