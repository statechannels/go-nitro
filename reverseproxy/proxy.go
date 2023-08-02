package reverseproxy

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog"
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

// ReversePaymentProxy is an HTTP proxy that charges for HTTP requests.
type ReversePaymentProxy struct {
	server         *http.Server
	nitroClient    rpc.RpcClientApi
	costPerByte    uint64
	reverseProxy   *httputil.ReverseProxy
	logger         zerolog.Logger
	destinationUrl *url.URL
}

// NewReversePaymentProxy creates a new ReversePaymentProxy.
func NewReversePaymentProxy(proxyAddress string, nitroEndpoint string, destinationURL string, costPerByte uint64, logger zerolog.Logger) *ReversePaymentProxy {
	server := &http.Server{Addr: proxyAddress}
	nitroClient, err := rpc.NewHttpRpcClient(nitroEndpoint)
	if err != nil {
		panic(err)
	}
	destinationUrl, err := url.Parse(destinationURL)
	if err != nil {
		panic(err)
	}

	p := &ReversePaymentProxy{
		server:         server,
		logger:         logger,
		nitroClient:    nitroClient,
		costPerByte:    costPerByte,
		destinationUrl: destinationUrl,
		reverseProxy:   &httputil.ReverseProxy{},
	}

	// Wire up our handlers to the reverse proxy
	p.reverseProxy.Rewrite = func(pr *httputil.ProxyRequest) { pr.SetURL(p.destinationUrl) }
	p.reverseProxy.ModifyResponse = p.handleDestinationResponse
	p.reverseProxy.ErrorHandler = p.handleError
	// Wire up our handler to the server
	p.server.Handler = p

	return p
}

// ServeHTTP is the main entry point for the reverse payment proxy server.
// It is responsible for parsing the voucher from the query params and moving it to the request header
// It then delegates to the reverse proxy to handle rewriting the request and sending it to the destination
func (p *ReversePaymentProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
func (p *ReversePaymentProxy) handleDestinationResponse(r *http.Response) error {
	contentLength, err := strconv.ParseUint(r.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return err
	}

	v, ok := r.Request.Context().Value(VOUCHER_CONTEXT_ARG).(payments.Voucher)
	if !ok {
		return createPaymentError(fmt.Errorf("could not fetch voucher from context"))
	}
	cost := p.costPerByte * contentLength

	p.logger.Debug().
		Uint64("costPerByte", p.costPerByte).
		Uint64("responseLength", contentLength).
		Uint64("cost", cost).
		Msg("Request cost")

	s, err := p.nitroClient.ReceiveVoucher(v)
	if err != nil {
		return createPaymentError(fmt.Errorf("error processing voucher %w", err))
	}

	p.logger.Debug().Msgf("Received voucher with delta %d", s.Delta.Uint64())

	// s.Delta is amount our balance increases by adding this voucher
	// AKA the payment amount we received in the request for this file
	if cost > s.Delta.Uint64() {
		return createPaymentError(fmt.Errorf("payment of %d required, the voucher only resulted in a payment of %d", cost, s.Delta.Uint64()))
	}

	p.logger.Debug().Msgf("Destination request URL %s", r.Request.URL.String())
	return nil
}

// handleError is responsible for logging the error and returning the appropriate HTTP status code
func (p *ReversePaymentProxy) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, ErrPayment) {
		http.Error(w, err.Error(), http.StatusPaymentRequired)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	p.logger.Error().Err(err).Msgf("Error processing request")
}

// Start starts the proxy server in a goroutine.
func (p *ReversePaymentProxy) Start() error {
	go func() {
		p.logger.Info().Msgf("Starting reverse payment proxy listening on %s.", p.server.Addr)

		if err := p.server.ListenAndServe(); err != http.ErrServerClosed {
			p.logger.Err(err).Msg("ListenAndServe()")
		}
	}()

	return nil
}

// Stop stops the proxy server and closes everything.
func (p *ReversePaymentProxy) Stop() error {
	p.logger.Info().Msgf("Stopping reverse payment proxy listening on %s", p.server.Addr)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	// Check if the request is an OPTIONS preflight request.
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}
