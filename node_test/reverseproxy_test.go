package node_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/internal/logging"

	ta "github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/node/engine/chainservice"
	"github.com/statechannels/go-nitro/reverseproxy"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/types"
)

const (
	destinationServerResponseBody = "Hello! This is from the destination server"
	proxyPort                     = 5511
	bobRPCUrl                     = ":4107/api/v1"
	destPort                      = 6622
)

func TestReversePaymentProxy(t *testing.T) {
	logFile := "reverse_payment_proxy.log"

	logDestination := logging.NewLogWriter("../artifacts", logFile)

	aliceClient, ireneClient, bobClient, cleanup := setupNitroClients(t, logDestination)
	defer cleanup()

	paymentChannel := createChannelData(t, aliceClient, ireneClient, bobClient)

	// Start up a test http server that acts as the destination server
	// It will return a simple response
	destinationServerUrl, cleanupDestServer := runDestinationServer(t, destPort)
	defer cleanupDestServer()

	// Create a ReversePaymentProxy with the test destination server URL
	proxy := reverseproxy.NewReversePaymentProxy(proxyPort, bobRPCUrl, destinationServerUrl)
	defer func() {
		err := proxy.Stop()
		if err != nil {
			t.Fatalf("Error stopping proxy: %v", err)
		}
	}()

	err := proxy.Start()
	if err != nil {
		t.Fatalf("Error starting proxy: %v", err)
	}

	v, err := aliceClient.CreateVoucher(paymentChannel, 5)
	if err != nil {
		t.Fatalf("Error creating voucher: %v", err)
	}

	// Create a request to the proxy server for some resource
	req, err := http.NewRequest("GET",
		fmt.Sprintf("http://localhost:%d/resource?channelId=%s&amount=%d&signature=%s", proxyPort, paymentChannel, 5, v.Signature.ToHexString()),
		nil)
	if err != nil {
		t.Fatalf("Error creating test request: %v", err)
	}

	// Make the request to the proxy server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error making request to proxy server: %v", err)
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading request data: %v", err)
	}
	// Check if the response from the destination server is correct

	if string(bodyText) != (destinationServerResponseBody) {
		t.Errorf("Expected response %q, but got %q", destinationServerResponseBody, bodyText)
	}
}

// setupNitroClients creates three nitro clients and connects them to each other
func setupNitroClients(t *testing.T, logDestination *os.File) (alice, irene, bob *rpc.RpcClient, cleanup func()) {
	chain := chainservice.NewMockChain()
	logger := testLogger(logDestination)

	aliceChainService := chainservice.NewMockChainService(chain, ta.Alice.Address())
	bobChainService := chainservice.NewMockChainService(chain, ta.Bob.Address())
	ireneChainService := chainservice.NewMockChainService(chain, ta.Irene.Address())

	aliceClient, msgAlice, aliceCleanup := setupNitroNodeWithRPCClient(t, ta.Alice.PrivateKey, 3105, 4105, aliceChainService, logDestination, "ws")
	ireneClient, msgIrene, ireneCleanup := setupNitroNodeWithRPCClient(t, ta.Irene.PrivateKey, 3106, 4106, ireneChainService, logDestination, "ws")
	bobClient, msgBob, bobCleanup := setupNitroNodeWithRPCClient(t, ta.Bob.PrivateKey, 3107, 4107, bobChainService, logDestination, "ws")

	logger.Info().Msg("Clients created")

	waitForPeerInfoExchange(msgAlice, msgBob, msgIrene)
	logger.Info().Msg("Peer exchange complete")

	return aliceClient, ireneClient, bobClient, func() {
		aliceCleanup()
		ireneCleanup()
		bobCleanup()
		chain.Close()
	}
}

// createChannelData creates ledgers channels and a payment channel between Alice and Bob
func createChannelData(t *testing.T, aliceClient, ireneClient, bobClient *rpc.RpcClient) (paymentChannelId types.Destination) {
	aliceLedgerRes, err := aliceClient.CreateLedgerChannel(ta.Irene.Address(), 100, simpleOutcome(ta.Alice.Address(), ta.Irene.Address(), 100, 100))
	if err != nil {
		t.Fatalf("Error creating channels: %v", err)
	}
	<-aliceClient.ObjectiveCompleteChan(aliceLedgerRes.Id)
	<-ireneClient.ObjectiveCompleteChan(aliceLedgerRes.Id)

	ireneLedgerRes, err := ireneClient.CreateLedgerChannel(ta.Bob.Address(), 100, simpleOutcome(ta.Irene.Address(), ta.Bob.Address(), 100, 100))
	if err != nil {
		t.Fatalf("Error creating channels: %v", err)
	}
	<-bobClient.ObjectiveCompleteChan(ireneLedgerRes.Id)
	<-ireneClient.ObjectiveCompleteChan(ireneLedgerRes.Id)

	initialOutcome := simpleOutcome(ta.Alice.Address(), ta.Bob.Address(), 100, 0)

	createPayCh, err := aliceClient.CreatePaymentChannel(
		[]common.Address{ta.Irene.Address()},
		ta.Bob.Address(),
		100,
		initialOutcome,
	)
	if err != nil {
		t.Fatalf("Error creating channels: %v", err)
	}
	<-aliceClient.ObjectiveCompleteChan(createPayCh.Id)
	<-bobClient.ObjectiveCompleteChan(createPayCh.Id)
	return createPayCh.ChannelId
}

// runDestinationServer runs a simple http server that returns a simple response
// It performs a basic check to make sure no voucher information was passed along
func runDestinationServer(t *testing.T, port uint) (destUrl string, cleanup func()) {
	checkError := func(err error) {
		if err == http.ErrServerClosed {
			return
		}
		if err != nil {
			t.Fatalf("Error running the destination server: %+v", err)
		}
	}

	handleRequest := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/resource" {
			t.Fatalf("Expected voucher information to be stripped off got %s instead", r.URL.String())
		}

		// Simulate the destination server's response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write([]byte(destinationServerResponseBody))
		checkError(err)
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handleRequest}

	go func() {
		err := server.ListenAndServe()
		checkError(err)
	}()
	return fmt.Sprintf("http://localhost:%d", port), func() {
		server.Close()
	}
}
