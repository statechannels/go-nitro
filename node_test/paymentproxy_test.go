package node_test

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/rpc/transport"

	"github.com/statechannels/go-nitro/internal/logging"
	ta "github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/node/engine/chainservice"
	"github.com/statechannels/go-nitro/paymentproxy"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/types"
)

const (
	smallResponse = "Hello"

	parseErrorResponseBody     = "could not parse voucher"
	signatureErrorResponseBody = "error processing voucher"
	proxyAddress               = ":5511"
	bobRPCUrl                  = "127.0.0.1:4107/api/v1"
	destPort                   = 6622
	otherParam                 = "otherParam"
	otherParamValue            = "2"
	testFileContent            = "This a simple test file used in the payment proxy"
	testFileName               = "test_file.txt"
	serverReadyMaxWait         = 2 * time.Second
)

func setupTestFile(t *testing.T) func() {
	// Open the file for writing (create or truncate)
	file, err := os.Create(testFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(testFileContent)
	if err != nil {
		os.Remove(testFileName)
		t.Fatal(err)
	}
	return func() {
		err := os.Remove(testFileName)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestPaymentProxy(t *testing.T) {
	logFile := "payment_proxy.log"

	aliceClient, ireneClient, bobClient, cleanup := setupNitroClients(t, logFile)
	defer cleanup()

	paymentChannel := createChannelData(t, aliceClient, ireneClient, bobClient)

	// Startup a simple http server that will be used as the destination server
	// It serves a simple text response and on two endpoints `resourceWithParams` and `resource``
	destinationServerUrl, cleanupDestServer := runDestinationServer(t, destPort)
	defer cleanupDestServer()

	cleanupData := setupTestFile(t)
	defer cleanupData()

	// Create a PaymentProxy with the test destination server URL
	proxy := paymentproxy.NewPaymentProxy(
		proxyAddress,
		bobRPCUrl,
		destinationServerUrl,
		1)
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
	waitForServer(t, fmt.Sprintf("http://%s/", proxyAddress), serverReadyMaxWait)

	voucher := createVoucher(t, aliceClient, paymentChannel, 5)
	resp := performGetRequest(t, "", fmt.Sprintf("http://%s/resource?channelId=%s&amount=%d&signature=%s", proxyAddress, voucher.ChannelId, voucher.Amount.Int64(), voucher.Signature.ToHexString()))
	checkResponse(t, resp, smallResponse, http.StatusOK)

	// Using the same voucher again should result in a payment required response
	resp = performGetRequest(t, "", fmt.Sprintf("http://%s/resource?channelId=%s&amount=%d&signature=%s", proxyAddress, voucher.ChannelId, voucher.Amount.Int64(), voucher.Signature.ToHexString()))
	checkResponse(t, resp, expectedPaymentErrorMessage(5, 0), http.StatusPaymentRequired)

	// Not providing a voucher should result in a payment required response
	resp = performGetRequest(t, "", fmt.Sprintf("http://%s/resource", proxyAddress))
	checkResponse(t, resp, parseErrorResponseBody, http.StatusPaymentRequired)

	// A voucher less than 5 should be rejected
	voucher = createVoucher(t, aliceClient, paymentChannel, 4)
	resp = performGetRequest(t, "", fmt.Sprintf("http://%s/resource?channelId=%s&amount=%d&signature=%s", proxyAddress, voucher.ChannelId, voucher.Amount.Uint64(), voucher.Signature.ToHexString()))
	checkResponse(t, resp, expectedPaymentErrorMessage(5, 4), http.StatusPaymentRequired)

	// A voucher with a bad signature should be rejected
	voucher = createVoucher(t, aliceClient, paymentChannel, 5)
	// Manually modify some bytes in the signature to make it invalid
	voucher.Signature.S[3] = 0
	voucher.Signature.R[3] = 127

	resp = performGetRequest(t, "", fmt.Sprintf("http://%s/resource?channelId=%s&amount=%d&signature=%s", proxyAddress, voucher.ChannelId, voucher.Amount.Uint64(), voucher.Signature.ToHexString()))
	checkResponse(t, resp, signatureErrorResponseBody, http.StatusPaymentRequired)

	// Check that the proxy can handle non voucher params and pass them along to the destination server
	voucher = createVoucher(t, aliceClient, paymentChannel, 5)
	resp = performGetRequest(t, "", fmt.Sprintf("http://%s/resource/params?channelId=%s&amount=%d&signature=%s&otherParam=2", proxyAddress, voucher.ChannelId, voucher.Amount, voucher.Signature.ToHexString()))
	checkResponse(t, resp, smallResponse, http.StatusOK)

	// It should properly handle a request to a non existent endpoint
	voucher = createVoucher(t, aliceClient, paymentChannel, 5)
	resp = performGetRequest(t, "", fmt.Sprintf("http://%s/badpath?channelId=%s&amount=%d&signature=%s", proxyAddress, voucher.ChannelId, voucher.Amount.Uint64(), voucher.Signature.ToHexString()))
	checkResponse(t, resp, "", http.StatusNotFound)

	// It should return a larger response if the voucher is large enough
	voucher = createVoucher(t, aliceClient, paymentChannel, uint64(len(testFileContent)))
	resp = performGetRequest(t, "", fmt.Sprintf("http://%s/file?channelId=%s&amount=%d&signature=%s", proxyAddress, voucher.ChannelId, voucher.Amount.Int64(), voucher.Signature.ToHexString()))
	checkResponse(t, resp, testFileContent, http.StatusOK)

	// It should return a payment required response for a large response if the voucher is  not large enough
	voucher = createVoucher(t, aliceClient, paymentChannel, 5)
	resp = performGetRequest(t, "", fmt.Sprintf("http://%s/file?channelId=%s&amount=%d&signature=%s", proxyAddress, voucher.ChannelId, voucher.Amount.Int64(), voucher.Signature.ToHexString()))
	checkResponse(t, resp, expectedPaymentErrorMessage(len(testFileContent), 5), http.StatusPaymentRequired)

	// It should handle a simple range request
	voucher = createVoucher(t, aliceClient, paymentChannel, 5)
	resp = performGetRequest(t, "bytes=0-4", fmt.Sprintf("http://%s/file?channelId=%s&amount=%d&signature=%s", proxyAddress, voucher.ChannelId, voucher.Amount.Int64(), voucher.Signature.ToHexString()))
	// We want to make sure that the response body is only the first 4 bytes
	body, statusCode := getResponseInfo(t, resp)
	if body != testFileContent[0:5] {
		t.Fatalf("Expected response body to be %s, but got %s", testFileContent[0:5], body)
	}
	if statusCode != http.StatusPartialContent {
		t.Fatalf("Expected status code %d, but got %d", http.StatusPartialContent, statusCode)
	}

	// It should handle a complex range request
	// Multipart ranges include extra information about the ranges in the response body
	const multiPartResponseSize = 346
	voucher = createVoucher(t, aliceClient, paymentChannel, multiPartResponseSize)

	resp = performGetRequest(t, "bytes=0-1,3-4", fmt.Sprintf("http://%s/file?channelId=%s&amount=%d&signature=%s", proxyAddress, voucher.ChannelId, voucher.Amount.Int64(), voucher.Signature.ToHexString()))
	body, statusCode = getResponseInfo(t, resp)
	// The server response will also contain some extra information about the ranges
	if !strings.Contains(body, testFileContent[0:2]) || !strings.Contains(body, testFileContent[3:5]) {
		t.Fatalf("Expected response body to contain partial file contents")
	}
	if statusCode != http.StatusPartialContent {
		t.Fatalf("Expected status code %d, but got %d", http.StatusPartialContent, statusCode)
	}

	// It should reject a range request if the voucher is not large enough
	voucher = createVoucher(t, aliceClient, paymentChannel, 1)
	resp = performGetRequest(t, "bytes=0-1,3-4", fmt.Sprintf("http://%s/file?channelId=%s&amount=%d&signature=%s", proxyAddress, voucher.ChannelId, voucher.Amount.Int64(), voucher.Signature.ToHexString()))
	checkResponse(t, resp, expectedPaymentErrorMessage(multiPartResponseSize, 1), http.StatusPaymentRequired)
}

// createVoucher creates a voucher for the given channel and amount	using the given client
// If any error occurs it will fail the test
func createVoucher(t *testing.T, client rpc.RpcClientApi, channelId types.Destination, amount uint64) payments.Voucher {
	v, err := client.CreateVoucher(channelId, amount)
	if err != nil {
		t.Fatalf("Error creating voucher: %v", err)
	}
	return v
}

func expectedPaymentErrorMessage(total, numPaid int) string {
	return fmt.Sprintf("payment of %d required, the voucher only resulted in a payment of %d", total, numPaid)
}

// performGetRequest performs a GET request to the given url
// If any error occurs it will fail the test
func performGetRequest(t *testing.T, rangeVal string, url string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest("GET",
		url,
		nil)
	if err != nil {
		t.Fatalf("Error performing request: %v", err)
	}
	if rangeVal != "" {
		req.Header.Add("Range", rangeVal)
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
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, statusCode)
	}
}

// getResponseInfo reads the response body and returns it as a string
// If any error occurs it will fail the test
func getResponseInfo(t *testing.T, resp *http.Response) (body string, statusCode int) {
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading request data: %v", err)
	}

	return string(bodyText), resp.StatusCode
}

// setupNitroClients creates three nitro clients and connects them to each other
func setupNitroClients(t *testing.T, logFile string) (alice, irene, bob rpc.RpcClientApi, cleanup func()) {
	chain := chainservice.NewMockChain()

	logging.SetupDefaultFileLogger(logFile, slog.LevelDebug)

	aliceChainService := chainservice.NewMockChainService(chain, ta.Alice.Address())
	bobChainService := chainservice.NewMockChainService(chain, ta.Bob.Address())
	ireneChainService := chainservice.NewMockChainService(chain, ta.Irene.Address())
	ireneClient, msgIrene, ireneCleanup := setupNitroNodeWithRPCClient(t, ta.Irene.PrivateKey, 3106, 4106, ireneChainService, transport.Http, []string{})
	bootPeers := []string{msgIrene.MultiAddr}
	aliceClient, msgAlice, aliceCleanup := setupNitroNodeWithRPCClient(t, ta.Alice.PrivateKey, 3105, 4105, aliceChainService, transport.Http, bootPeers)

	bobClient, msgBob, bobCleanup := setupNitroNodeWithRPCClient(t, ta.Bob.PrivateKey, 3107, 4107, bobChainService, transport.Http, bootPeers)

	slog.Info("Clients created")

	waitForPeerInfoExchange(msgAlice, msgBob, msgIrene)
	slog.Info("Peer exchange complete")

	return aliceClient, ireneClient, bobClient, func() {
		aliceCleanup()
		ireneCleanup()
		bobCleanup()
		chain.Close()
	}
}

// createChannelData creates ledgers channels and a payment channel between Alice and Bob
func createChannelData(t *testing.T, aliceClient, ireneClient, bobClient rpc.RpcClientApi) (paymentChannelId types.Destination) {
	aliceLedgerRes, err := aliceClient.CreateLedgerChannel(ta.Irene.Address(), 100, simpleOutcome(ta.Alice.Address(), ta.Irene.Address(), 500, 500))
	if err != nil {
		t.Fatalf("Error creating channels: %v", err)
	}
	<-aliceClient.ObjectiveCompleteChan(aliceLedgerRes.Id)
	<-ireneClient.ObjectiveCompleteChan(aliceLedgerRes.Id)

	ireneLedgerRes, err := ireneClient.CreateLedgerChannel(ta.Bob.Address(), 100, simpleOutcome(ta.Irene.Address(), ta.Bob.Address(), 500, 500))
	if err != nil {
		t.Fatalf("Error creating channels: %v", err)
	}
	<-bobClient.ObjectiveCompleteChan(ireneLedgerRes.Id)
	<-ireneClient.ObjectiveCompleteChan(ireneLedgerRes.Id)

	initialOutcome := simpleOutcome(ta.Alice.Address(), ta.Bob.Address(), 500, 0)

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
		if r.URL.Path != "/resource" && r.URL.Path != "/resource/params" && r.URL.Path != "/file" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		params, err := url.ParseQuery(r.URL.RawQuery)
		checkError(err)
		// If this is a request to /resource/params, we check for the other param
		if checkForOtherParam := r.URL.Path == "/resource/params"; checkForOtherParam {
			if !params.Has(otherParam) {
				t.Fatalf("Did not find query param %s in url %s", otherParam, r.URL.RawQuery)
			}
			if params.Get(otherParam) != otherParamValue {
				t.Fatalf("Expected query param %s to have value %s, but got %s", otherParam, otherParamValue, params.Get(otherParam))
			}
		}

		// Always check that the voucher params were stripped out of every request
		for p := range params {
			if p == paymentproxy.AMOUNT_VOUCHER_PARAM || p == paymentproxy.CHANNEL_ID_VOUCHER_PARAM || p == paymentproxy.SIGNATURE_VOUCHER_PARAM {
				t.Fatalf("Expected no voucher information to be passed along, but got %s", p)
			}
		}

		if r.URL.Path == "/file" {
			http.ServeFile(w, r, testFileName)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/plain")
			_, err = w.Write([]byte(smallResponse))
			checkError(err)
		}
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handleRequest}

	go func() {
		err := server.ListenAndServe()
		checkError(err)
	}()

	url := fmt.Sprintf("http://localhost:%d", port)

	waitForServer(t, url, serverReadyMaxWait)

	return url, func() {
		server.Close()
	}
}

// waitForServer waits for the given url to be available by performing GET requests
func waitForServer(t *testing.T, url string, timeout time.Duration) {
	isReady := make(chan struct{})
	go func() {
		for {
			_, err := http.Get(url)
			if err == nil {
				close(isReady)
				return
			}

			time.Sleep(10 * time.Millisecond)
		}
	}()

	select {
	case <-isReady:
		return
	case <-time.After(timeout):
		t.Fatalf("server did not reply after %v", timeout)
	}
}
