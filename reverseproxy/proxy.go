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

const (
	AMOUNT_VOUCHER_PARAM     = "amount"
	CHANNEL_ID_VOUCHER_PARAM = "channelId"
	SIGNATURE_VOUCHER_PARAM  = "signature"

	// HEADER_PREFIX is the prefix we attach to voucher parameters when we move them from query params to headers
	HEADER_PREFIX = "Nitro-"

	ErrPayment = types.ConstError("payment error")
)

// createPaymentError wraps an error with ErrPayment.
func createPaymentError(err error) error {
	return fmt.Errorf("%w: %w", ErrPayment, err)
}

// ReversePaymentProxy is an HTTP proxy that charges for HTTP requests.
type ReversePaymentProxy struct {
	server       *http.Server
	nitroClient  rpc.RpcClientApi
	costPerByte  uint64
	reverseProxy *httputil.ReverseProxy
	logger       zerolog.Logger
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

	// Creates a reverse proxy that will handle forwarding requests to the destination server
	proxy := &httputil.ReverseProxy{
		// Rewrite handles modifying the request before it is sent to the destination server
		// We override it to handle modifying the request URL (via SetURL) and moving the voucher from the query params to the header
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(destinationUrl)

			v, err := parseVoucher(r.In.URL.Query(), "")
			// If we can't parse the voucher we return and rely on ModifyResponse to throw an error when it doesn't find a voucher
			if err != nil {
				return
			}
			// We move the voucher from query params to the header, so it doesn't pollute the query params to the destination server
			removeVoucher(r.Out.URL.Query(), "")
			r.Out.URL.RawQuery = r.Out.URL.Query().Encode()
			addVoucher(v, r.Out.Header, HEADER_PREFIX)
		},
		// ModifyResponse handles modifying the response before it is sent back to the client
		// It attempts to parse the voucher from the request header and redeem it with the Nitro client
		// It will check the voucher amount against the cost (response size * cost per byte)
		ModifyResponse: func(r *http.Response) error {
			contentLength, err := strconv.ParseUint(r.Header.Get("Content-Length"), 10, 64)
			if err != nil {
				return err
			}

			v, err := parseVoucher(r.Request.Header, HEADER_PREFIX)
			if err != nil {
				return createPaymentError(fmt.Errorf("could not parse voucher: %w", err))
			}

			cost := costPerByte * contentLength

			logger.Debug().
				Uint64("costPerByte", costPerByte).
				Uint64("responseLength", contentLength).
				Uint64("cost", cost).
				Msg("Request cost")

			s, err := nitroClient.ReceiveVoucher(v)
			if err != nil {
				return createPaymentError(fmt.Errorf("error processing voucher %w", err))
			}

			logger.Debug().Msgf("Received voucher with delta %d", s.Delta.Uint64())

			// s.Delta is amount our balance increases by adding this voucher
			// AKA the payment amount we received in the request for this file
			if cost > s.Delta.Uint64() {
				return createPaymentError(fmt.Errorf("payment of %d required, the voucher only resulted in a payment of %d", cost, s.Delta.Uint64()))
			}

			logger.Debug().Msgf("Destination request URL %s", r.Request.URL.String())
			return nil
		},
		// ErrorHandler handles errors that occur during ModifyResponse
		// We use it to return a nice error message to the client depending on the error
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			if errors.Is(err, ErrPayment) {
				http.Error(w, err.Error(), http.StatusPaymentRequired)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			logger.Error().Err(err).Msgf("Error processing request")
		},
	}

	return &ReversePaymentProxy{
		server:       server,
		logger:       logger,
		nitroClient:  nitroClient,
		reverseProxy: proxy,
		costPerByte:  costPerByte,
	}
}

// Start starts the proxy server in a goroutine.
func (p *ReversePaymentProxy) Start() error {
	// Wire up our proxy to the http handler
	// This means that p.ServeHTTP will be called for every request
	p.server.Handler = p.reverseProxy

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

// keyValCollection is an interface that allows us to set, get and delete key value pairs.
// It is used to abstract away the differences between http.Header and url.Values.
type keyValCollection interface {
	Set(key, value string)
	Get(key string) string
	Del(key string)
}

// addVoucher takes in a voucher and adds it to the given keyValCollection.
// It prefixes the keys with the given prefix.
func addVoucher(v payments.Voucher, col keyValCollection, prefix string) {
	col.Set(prefix+CHANNEL_ID_VOUCHER_PARAM, v.ChannelId.String())
	col.Set(prefix+AMOUNT_VOUCHER_PARAM, v.Amount.String())
	col.Set(prefix+SIGNATURE_VOUCHER_PARAM, v.Signature.ToHexString())
}

// parseVoucher takes in an a keyValCollection  parses out a voucher.
func parseVoucher(col keyValCollection, prefix string) (payments.Voucher, error) {
	rawChId := col.Get(prefix + CHANNEL_ID_VOUCHER_PARAM)
	if rawChId == "" {
		return payments.Voucher{}, fmt.Errorf("missing channel ID")
	}
	rawAmt := col.Get(prefix + AMOUNT_VOUCHER_PARAM)
	if rawAmt == "" {
		return payments.Voucher{}, fmt.Errorf("missing amount")
	}
	rawSignature := col.Get(prefix + SIGNATURE_VOUCHER_PARAM)
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

// removeVoucherParams removes the voucher parameters from the request URL.
func removeVoucher(col keyValCollection, prefix string) {
	col.Del(prefix + CHANNEL_ID_VOUCHER_PARAM)
	col.Del(prefix + AMOUNT_VOUCHER_PARAM)
	col.Del(prefix + SIGNATURE_VOUCHER_PARAM)
}
