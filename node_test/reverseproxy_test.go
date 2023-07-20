package node_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
	parseErrorResponseBody        = "could not parse voucher"
	signatureErrorResponseBody    = "error processing voucher"
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

	resp := performGetRequest(t, fmt.Sprintf("http://localhost:%d/resource?channelId=%s&amount=%d&signature=%s", proxyPort, paymentChannel, 5, v.Signature.ToHexString()))
	checkResponse(t, resp, destinationServerResponseBody, http.StatusOK)

	// Using the same voucher again should result in a payment required response
	resp = performGetRequest(t, fmt.Sprintf("http://localhost:%d/resource?channelId=%s&amount=%d&signature=%s", proxyPort, paymentChannel, 5, v.Signature.ToHexString()))
	checkResponse(t, resp, expectedPaymentErrorMessage(0), http.StatusPaymentRequired)

	// Not providing a voucher should result in a payment required response
	resp = performGetRequest(t, fmt.Sprintf("http://localhost:%d/resource", proxyPort))
	checkResponse(t, resp, parseErrorResponseBody, http.StatusPaymentRequired)

	// A voucher less than 5 should be rejected
	v, err = aliceClient.CreateVoucher(paymentChannel, 4)
	if err != nil {
		t.Fatalf("Error creating voucher: %v", err)
	}
	resp = performGetRequest(t, fmt.Sprintf("http://localhost:%d/resource?channelId=%s&amount=%d&signature=%s", proxyPort, v.ChannelId, v.Amount.Uint64(), v.Signature.ToHexString()))
	checkResponse(t, resp, expectedPaymentErrorMessage(4), http.StatusPaymentRequired)

	// A voucher with a bad signature should be rejected
	v, err = aliceClient.CreateVoucher(paymentChannel, 5)
	if err != nil {
		t.Fatalf("Error creating voucher: %v", err)
	}

	// Manually modify some bytes in the signature to make it invalid
	v.Signature.S[3] = 0
	v.Signature.R[3] = 127

	resp = performGetRequest(t,
		fmt.Sprintf("http://localhost:%d/resource?channelId=%s&amount=%d&signature=%s", proxyPort, v.ChannelId, v.Amount.Uint64(), v.Signature.ToHexString()))
	checkResponse(t, resp, signatureErrorResponseBody, http.StatusPaymentRequired)
}

func expectedPaymentErrorMessage(numPaid int) string {
	return fmt.Sprintf("payment of 5 required, the voucher only resulted in a payment of %d", numPaid)
}

// performGetRequest performs a GET request to the given url
// If any error occurs it will fail the test
func performGetRequest(t *testing.T, url string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest("GET",
		url,
		nil)
	if err != nil {
		t.Fatalf("Error performing request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error performing request: %v", err)
	}
	return resp
}

func checkResponse(t *testing.T, resp *http.Response, expectedBody string, expectedStatusCode int) {
	responseBodyText, statusCode := getResponseInfo(t, resp)
	if !strings.Contains(responseBodyText, expectedBody) {
		t.Errorf("The body of the response %s did not contain the expected text %s ", responseBodyText, expectedBody)
	}
	if statusCode != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, statusCode)
	}
}

// getResponseInfo reads the response body and returns it as a string
// If any error occurs it will fail the test
func getResponseInfo(t *testing.T, resp *http.Response) (body string, statusCode int) {
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading request data: %v", err)
	}
	resp.Body.Close()

	return string(bodyText), resp.StatusCode
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
