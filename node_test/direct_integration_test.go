package node_test

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/utils"
	"github.com/statechannels/go-nitro/node"
	"github.com/statechannels/go-nitro/node/engine/messageservice"
	p2pms "github.com/statechannels/go-nitro/node/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/node/query"
	"github.com/statechannels/go-nitro/types"
	"github.com/stretchr/testify/require"
)

func TestDirectIntegrationScenario(t *testing.T) {
	directCase := TestCase{
		Description:    "Direct test",
		Chain:          MockChain,
		MessageService: TestMessageService,
		MessageDelay:   0,
		LogName:        "direct_integration_run.log",
		Participants: []TestParticipant{
			{StoreType: MemStore, Actor: testactors.Alice},
			{StoreType: MemStore, Actor: testactors.Bob},
		},
		Actions: []DirectTestCaseAction{
			// Alice proposes a channel update, Bob agrees
			{Initiator: 0, Responder: 1, ProposalAppData: types.Bytes{0}, IsResponderAgree: true, LatestSupportedAppData: types.Bytes{0}},
			// Bob proposes a channel update, Alice agrees
			{Initiator: 1, Responder: 0, ProposalAppData: types.Bytes{1}, IsResponderAgree: true, LatestSupportedAppData: types.Bytes{1}},
			// Alice proposes a channel update, Bob disagrees
			{Initiator: 0, Responder: 1, ProposalAppData: types.Bytes{2}, IsResponderAgree: false, LatestSupportedAppData: types.Bytes{1}},
			// Alice proposes more channel updates, Bob continue to disagree
			{Initiator: 0, Responder: 1, ProposalAppData: types.Bytes{3}, IsResponderAgree: false, LatestSupportedAppData: types.Bytes{1}},
			{Initiator: 0, Responder: 1, ProposalAppData: types.Bytes{4}, IsResponderAgree: false, LatestSupportedAppData: types.Bytes{1}},
			// Bob agrees on one of Alice's channel updates, Alice's agreement is not needed
			{Initiator: 1, Responder: 0, ProposalAppData: types.Bytes{3}, IsResponderAgree: false, LatestSupportedAppData: types.Bytes{3}},
			// Bob agrees on one of Alice's expired channel updates, so he creates a new one. Alice disagrees
			{Initiator: 1, Responder: 0, ProposalAppData: types.Bytes{2}, IsResponderAgree: false, LatestSupportedAppData: types.Bytes{3}},
			// Bob proposes another channel update, and Alice agrees, so all previous updates expire
			{Initiator: 1, Responder: 0, ProposalAppData: types.Bytes{5}, IsResponderAgree: true, LatestSupportedAppData: types.Bytes{5}},
		},
	}

	RunDirectIntegrationTestCase(directCase, t)
}

// RunDirectIntegrationTestCase runs the integration test case.
func RunDirectIntegrationTestCase(tc TestCase, t *testing.T) {
	// Clean up all the test data we create at the end of the test
	defer os.RemoveAll(STORE_TEST_DATA_FOLDER)

	t.Run(tc.Description, func(t *testing.T) {
		err := tc.ValidateDirect()
		if err != nil {
			t.Fatal(err)
		}
		infra := setupSharedInfra(tc)
		defer infra.Close(t)

		msgServices := make([]messageservice.MessageService, 0)

		// Setup clients
		clientA, msgA, storeA := setupIntegrationNode(tc, tc.Participants[0], infra)
		defer clientA.Close()
		msgServices = append(msgServices, msgA)

		clientB, msgB, storeB := setupIntegrationNode(tc, tc.Participants[1], infra)
		defer clientB.Close()
		msgServices = append(msgServices, msgB)

		clients := []*node.Node{&clientA, &clientB}

		if tc.MessageService == P2PMessageService {
			p2pServices := make([]*p2pms.P2PMessageService, len(tc.Participants))
			for i, msgService := range msgServices {
				p2pServices[i] = msgService.(*p2pms.P2PMessageService)
			}

			utils.WaitForPeerInfoExchange(p2pServices...)
		}

		// Create direct channel between peers
		asset := common.Address{}
		customAppDef := common.HexToAddress("0x8D033747C268ef77460055Ff86CBB2643330D1A1") // Should not be empty
		outcome := initialLedgerOutcome(*clientA.Address, *clientB.Address, asset)

		channel := setupDirectChannel(t, clientA, clientB, asset, customAppDef, types.Bytes{})
		checkLedgerChannel(t, channel, outcome, query.Open, clientA)

		for _, act := range tc.Actions {
			clients[act.Initiator].UpdateChannel(channel, act.ProposalAppData)

			update := <-clients[act.Responder].ReceivedChannelUpdates()
			if act.IsResponderAgree {
				clients[act.Responder].UpdateChannel(update.ChannelId, update.AppData)
				<-clients[act.Initiator].ReceivedChannelUpdates()
			}

			chA, _ := storeA.GetChannelById(channel)
			lssA, _ := chA.LatestSupportedState()
			require.Equal(t, act.LatestSupportedAppData, lssA.AppData)

			chB, _ := storeB.GetChannelById(channel)
			lssB, _ := chB.LatestSupportedState()
			require.Equal(t, act.LatestSupportedAppData, lssB.AppData)
		}

		// TODO: finish closing channel after ledger channel refactoring
		// closeDirectChannel(t, clientA, clientB, channel)
	})
}
