package reverseproxy

import (
	"context"
	"fmt"
	"log"
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

// ReversePaymentProxy is an HTTP proxy that charges for HTTP requests.
type ReversePaymentProxy struct {
	server      *http.Server
	path        string
	nitroClient *rpc.RpcClient
	ctx         context.Context
	cancel      context.CancelFunc

	reverseProxy *httputil.ReverseProxy
}

// NewReversePaymentProxy creates a new ReversePaymentProxy.
func NewReversePaymentProxy(proxyPort uint, nitroEndpoint string, destination string) *ReversePaymentProxy {
	server := &http.Server{}
	nitroClient, err := rpc.NewHttpRpcClient(nitroEndpoint)
	if err != nil {
		panic(err)
	}
	destinationUrl, err := url.Parse(destination)
	if err != nil {
		panic(err)
	}
	// Creates a reverse proxy that will handle forwarding requests to the destination server
	proxy := httputil.NewSingleHostReverseProxy(destinationUrl)

	return &ReversePaymentProxy{
		path:   fmt.Sprintf("localhost:%d", proxyPort),
		server: server,

		nitroClient:  nitroClient,
		reverseProxy: proxy,
	}
}

// Start starts the proxy server in a goroutine.
func (p *ReversePaymentProxy) Start(ctx context.Context) error {
	p.ctx, p.cancel = context.WithCancel(ctx)

	go func() {
		fmt.Printf("Starting reverse payment proxy listening on %s\n", p.path)
		if err := http.ListenAndServe(p.path, p); err != http.ErrServerClosed {
			log.Fatalf("http.ListenAndServe(): %v", err)
		}
	}()

	return nil
}

// Stop stops the proxy server and closes everything.
func (p *ReversePaymentProxy) Stop() error {
	p.cancel()
	err := p.nitroClient.Close()
	if err != nil {
		return err
	}
	return p.server.Close()
}

// ServeHTTP is the main entry point for the proxy.
// It looks for voucher parameters in the request to construct a voucher.
// It then passes the voucher to the nitro client to process.
// Based on the amount added by the voucher, it either forwards the request to the destination server or returns an error.
func (p *ReversePaymentProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// This the payment we expect to receive for the file.
	const expectedPayment = int64(5)

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		webError(w, fmt.Errorf("could not query params: %w", err), http.StatusBadRequest)
		return
	}

	v, err := parseVoucher(params)
	if err != nil {
		webError(w, fmt.Errorf("could not parse voucher: %w", err), http.StatusBadRequest)
		return
	}

	s, err := p.nitroClient.ReceiveVoucher(v)
	if err != nil {
		webError(w, fmt.Errorf("error processing voucher %w", err), http.StatusBadRequest)
		return
	}

	// s.Delta is amount our balance increases by adding this voucher
	// AKA the payment amount we received in the request for this file
	if s.Delta.Cmp(big.NewInt(expectedPayment)) < 0 {
		webError(w, fmt.Errorf("payment of %d required, the voucher only resulted in a payment of %d", expectedPayment, s.Delta.Uint64()), http.StatusPaymentRequired)
		return
	}

	// Strip out the voucher params so the destination server doesn't need to handle them
	removeVoucherParams(r.URL)

	// Forward the request to the destination server
	p.reverseProxy.ServeHTTP(w, r)
}

// parseVoucher takes in an a collection of query params and parses out a voucher.
func parseVoucher(params url.Values) (payments.Voucher, error) {
	if !params.Has("channelId") {
		return payments.Voucher{}, fmt.Errorf("a valid channel id must be provided")
	}
	if !params.Has("amount") {
		return payments.Voucher{}, fmt.Errorf("a valid amount must be provided")
	}
	if !params.Has("signature") {
		return payments.Voucher{}, fmt.Errorf("a valid signature must be provided")
	}
	rawChId := params.Get("channelId")
	rawAmt := params.Get("amount")
	amount := big.NewInt(0)
	amount.SetString(rawAmt, 10)
	rawSignature := params.Get("signature")

	v := payments.Voucher{
		ChannelId: types.Destination(common.HexToHash(rawChId)),
		Amount:    amount,
		Signature: crypto.SplitSignature(hexutil.MustDecode(rawSignature)),
	}
	return v, nil
}

// removeVoucherParams removes the voucher parameters from the request URL.
func removeVoucherParams(u *url.URL) {
	queryParams := u.Query()
	delete(queryParams, "channelId")
	delete(queryParams, "signature")
	delete(queryParams, "amount")
	// Update the request URL without the voucher parameters
	u.RawQuery = queryParams.Encode()
}

// webError is a helper function to return an http error.
func webError(w http.ResponseWriter, err error, code int) {
	// TODO: This is a hack to allow CORS requests to the gateway for the boost integration demo.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	fmt.Printf("ERROR CODE %d\n", code)
	http.Error(w, err.Error(), code)
}
