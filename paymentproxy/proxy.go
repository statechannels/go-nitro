package paymentproxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/types"
)

type contextKey string

const (
	AMOUNT_VOUCHER_PARAM     = "amount"
	CHANNEL_ID_VOUCHER_PARAM = "channelId"
	SIGNATURE_VOUCHER_PARAM  = "signature"

	VOUCHER_CONTEXT_ARG contextKey = "voucher"

	ErrPayment = types.ConstError("payment error")
)

// createPaymentError wraps an error with ErrPayment.
func createPaymentError(err error) error {
	return fmt.Errorf("%w: %w", ErrPayment, err)
}

// PaymentProxy is an HTTP proxy that charges for HTTP requests.
type PaymentProxy struct {
	server       *http.Server
	nitroClient  rpc.RpcClientApi
	costPerByte  uint64
	reverseProxy *httputil.ReverseProxy

	destinationUrl            *url.URL
	certFilePath, certKeyPath string
}

// NewPaymentProxy creates a new PaymentProxy.
func NewPaymentProxy(proxyAddress string, nitroEndpoint string, destinationURL string, costPerByte uint64, certFilePath, certKeyPath string) *PaymentProxy {
	server := &http.Server{Addr: proxyAddress}

	nitroClient, err := rpc.NewHttpRpcClient(nitroEndpoint)
	if err != nil {
		panic(err)
	}
	destinationUrl, err := url.Parse(destinationURL)
	if err != nil {
		panic(err)
	}

	p := &PaymentProxy{
		server:         server,
		nitroClient:    nitroClient,
		costPerByte:    costPerByte,
		destinationUrl: destinationUrl,
		reverseProxy:   &httputil.ReverseProxy{},
		certFilePath:   certFilePath,
		certKeyPath:    certKeyPath,
	}
	// Wire up our handlers to the reverse proxy
	p.reverseProxy.Rewrite = func(pr *httputil.ProxyRequest) { pr.SetURL(p.destinationUrl) }
	p.reverseProxy.ModifyResponse = p.handleDestinationResponse
	p.reverseProxy.ErrorHandler = p.handleError
	// Wire up our handler to the server
	p.server.Handler = p

	return p
}

// ServeHTTP is the main entry point for the payment proxy server.
// It is responsible for parsing the voucher from the query params and moving it to the request header
// It then delegates to the reverse proxy to handle rewriting the request and sending it to the destination
func (p *PaymentProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If the request is a health check, return a 200 OK
	if r.URL.Path == "/health" {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Proxy is healthy"))
		if err != nil {
			p.handleError(w, r, err)
			return
		}
		return
	}
	enableCORS(w, r)
	v, err := parseVoucher(r.URL.Query())
	if err != nil {
		p.handleError(w, r, createPaymentError(fmt.Errorf("could not parse voucher: %w", err)))
		return
	}

	removeVoucher(r)

	// We add the voucher to the request context so we can access it in the response handler
	r = r.WithContext(context.WithValue(r.Context(), VOUCHER_CONTEXT_ARG, v))

	p.reverseProxy.ServeHTTP(w, r)
}

// handleDestinationResponse modifies the response before it is sent back to the client
// It is responsible for parsing the voucher from the request header and redeeming it with the Nitro client
// It will check the voucher amount against the cost (response size * cost per byte)
// If the voucher amount is less than the cost, it will return a 402 Payment Required error instead of serving the content
func (p *PaymentProxy) handleDestinationResponse(r *http.Response) error {
	// Ignore OPTIONS requests as they are preflight requests
	if r.Request.Method == "OPTIONS" {
		return nil
	}

	contentLength := uint64(0)
	// If the Content-Length header is set, use that
	// Otherwise, read the body to get the length
	if r.ContentLength != -1 {
		contentLength = uint64(r.ContentLength)
	} else {
		var err error
		contentLength, err = readBodyLength(r.Body)
		if err != nil {
			return createPaymentError(err)
		}
	}

	v, ok := r.Request.Context().Value(VOUCHER_CONTEXT_ARG).(payments.Voucher)
	if !ok {
		return createPaymentError(fmt.Errorf("could not fetch voucher from context"))
	}
	cost := p.costPerByte * contentLength

	slog.Debug("Request cost", "cost-per-byte", p.costPerByte, "response-length", contentLength, "cost", cost)

	s, err := p.nitroClient.ReceiveVoucher(v)
	if err != nil {
		return createPaymentError(fmt.Errorf("error processing voucher %w", err))
	}
	slog.Debug("Received voucher", "delta", s.Delta.Uint64())

	// s.Delta is amount our balance increases by adding this voucher
	// AKA the payment amount we received in the request for this file
	if cost > s.Delta.Uint64() {
		return createPaymentError(fmt.Errorf("payment of %d required, the voucher only resulted in a payment of %d", cost, s.Delta.Uint64()))
	}
	slog.Debug("Destination request", "url", r.Request.URL.String())

	return nil
}

// handleError is responsible for logging the error and returning the appropriate HTTP status code
func (p *PaymentProxy) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, ErrPayment) {
		http.Error(w, err.Error(), http.StatusPaymentRequired)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	slog.Error("Error processing request", "error", err)
}

// Start starts the proxy server in a goroutine.
func (p *PaymentProxy) Start() error {
	go func() {
		if p.certFilePath != "" && p.certKeyPath != "" {
			if err := p.server.ListenAndServeTLS(p.certFilePath, p.certKeyPath); err != http.ErrServerClosed {
				slog.Error("Error while listening", "error", err)
			}
		} else {
			slog.Info("Starting a payment proxy", "address", p.server.Addr)
			if err := p.server.ListenAndServe(); err != http.ErrServerClosed {
				slog.Error("Error while listening", "error", err)
			}
		}
	}()

	return nil
}

// Stop stops the proxy server and closes everything.
func (p *PaymentProxy) Stop() error {
	slog.Info("Stopping a payment proxy", "address", p.server.Addr)

	err := p.server.Shutdown(context.Background())
	if err != nil {
		return err
	}

	return p.nitroClient.Close()
}

// parseVoucher takes in an a collection of query params and parses out a voucher.
func parseVoucher(params url.Values) (payments.Voucher, error) {
	rawChId := params.Get(CHANNEL_ID_VOUCHER_PARAM)
	if rawChId == "" {
		return payments.Voucher{}, fmt.Errorf("missing channel ID")
	}
	rawAmt := params.Get(AMOUNT_VOUCHER_PARAM)
	if rawAmt == "" {
		return payments.Voucher{}, fmt.Errorf("missing amount")
	}
	rawSignature := params.Get(SIGNATURE_VOUCHER_PARAM)
	if rawSignature == "" {
		return payments.Voucher{}, fmt.Errorf("missing signature")
	}

	amount := big.NewInt(0)
	amount.SetString(rawAmt, 10)

	v := payments.Voucher{
		ChannelId: types.Destination(common.HexToHash(rawChId)),
		Amount:    amount,
		Signature: crypto.SplitSignature(hexutil.MustDecode(rawSignature)),
	}
	return v, nil
}

// removeVoucherParams removes the voucher parameters from the request URL
func removeVoucher(r *http.Request) {
	queryParams := r.URL.Query()

	queryParams.Del(CHANNEL_ID_VOUCHER_PARAM)
	queryParams.Del(AMOUNT_VOUCHER_PARAM)
	queryParams.Del(SIGNATURE_VOUCHER_PARAM)

	r.URL.RawQuery = queryParams.Encode()
}

// enableCORS enables CORS headers in the response.
func enableCORS(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers to allow all origins (*).
	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	// Check if the request is an OPTIONS preflight request.
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func readBodyLength(b io.ReadCloser) (uint64, error) {
	var byteCount uint64
	buffer := make([]byte, 1024)

	// Read the response body and count the bytes
	for {
		n, err := b.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break // Reached the end of the response body
			}
			return 0, err
		}
		byteCount += uint64(n)
	}

	return byteCount, nil
}
