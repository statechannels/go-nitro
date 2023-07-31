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
	"strings"

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
	v, err := parseVoucher(r.URL.Query())
	if err != nil {
		p.handleError(w, r, createPaymentError(fmt.Errorf("could not parse voucher: %w", err)))
		return
	}

	removeVoucher(r)

	// We add the voucher to the request context so we can access it later if needed
	r = r.WithContext(context.WithValue(r.Context(), VOUCHER_CONTEXT_ARG, v))

	// If the request has a range header we can check the voucher before forwarding the request
	if hasRangeHeader(r.Header) {
		total, err := parseRangeHeader(r.Header)
		if err != nil {
			p.handleError(w, r, createPaymentError(fmt.Errorf("could not parse range: %w", err)))
			return
		}
		if err := p.checkVoucherAgainstCost(v, total); err != nil {
			p.handleError(w, r, createPaymentError(err))
			return
		}
	}
	p.reverseProxy.ServeHTTP(w, r)
}

// checkVoucherAgainstCost checks the voucher against the cost (content size * cost per byte)
// If the voucher amount is less than the cost, it will return an error
func (p *ReversePaymentProxy) checkVoucherAgainstCost(v payments.Voucher, contentLength uint64) error {
	s, err := p.nitroClient.ReceiveVoucher(v)
	if err != nil {
		return (fmt.Errorf("error processing voucher %w", err))
	}

	cost := contentLength * p.costPerByte

	// s.Delta is amount our balance increases by adding this voucher
	// AKA the payment amount we received in the request for this file
	if cost > s.Delta.Uint64() {
		return (fmt.Errorf("payment of %d required, the voucher only resulted in a payment of %d", cost, s.Delta.Uint64()))
	}
	return nil
}

// handleDestinationResponse modifies the response before it is sent back to the client
// It is responsible for parsing the voucher from the request header and redeeming it with the Nitro client
// It will check the voucher amount against the cost (response size * cost per byte)
// If the voucher amount is less than the cost, it will return a 402 Payment Required error instead of serving the content
func (p *ReversePaymentProxy) handleDestinationResponse(r *http.Response) error {
	// If the request has a range header the voucher was already checked in ServeHTTP
	if hasRangeHeader(r.Request.Header) {
		return nil
	}

	contentLength, err := strconv.ParseUint(r.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return err
	}

	v, ok := r.Request.Context().Value(VOUCHER_CONTEXT_ARG).(payments.Voucher)
	if !ok {
		return createPaymentError(fmt.Errorf("could not fetch voucher from context"))
	}

	if err := p.checkVoucherAgainstCost(v, contentLength); err != nil {
		return createPaymentError(err)
	}

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

const (
	RANGE_HEADER = "Range"
	UNIT_PREFIX  = "bytes="
)

// hasRangeHeader returns true if the given header has a range header value.
func hasRangeHeader(h http.Header) bool {
	return h.Get(RANGE_HEADER) != ""
}

// parseRangeHeader parses the range header and returns the total number of bytes requested.
func parseRangeHeader(h http.Header) (total uint64, err error) {
	if !hasRangeHeader(h) {
		return 0, fmt.Errorf("no range header")
	}
	rangeVal := h.Get(RANGE_HEADER)
	var hasPrefix bool
	rangeVal, hasPrefix = strings.CutPrefix(rangeVal, UNIT_PREFIX)
	if !hasPrefix {
		return 0, fmt.Errorf("range header value '%s' missing prefix 'byte='", rangeVal)
	}

	// It is valid to have multiple ranges, separated by commas. IE: bytes=0-499,1000-1499
	for _, r := range strings.Split(rangeVal, ",") {
		// Each range is a start and end separated by a dash. IE: 0-499
		numVals := strings.Split(r, "-")

		start, err := strconv.ParseUint(numVals[0], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("could not parse start: %w", err)
		}

		end, err := strconv.ParseUint(numVals[1], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("could not parse start: %w", err)
		}

		if end < start {
			return 0, fmt.Errorf("start cannot be greater than end")
		}

		total += end - start + 1 // +1 because the end is inclusive. IE: 0-9 would return 10 bytes.

	}
	return
}
