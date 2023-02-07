package network

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/rs/zerolog"

	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/rpc/serde"
	"github.com/statechannels/go-nitro/rpc/transport"
)

// Response includes a payload or an error.
type Response[T serde.ResponsePayload] struct {
	Payload T
	Error   error
}

// Request uses the supplied connection and payload to send a non-blocking JSONRPC request.
// It returns a channel that sends a response payload. If the request fails to send, an error is returned.
func Request[T serde.RequestPayload, U serde.ResponsePayload](connection transport.Connection, request T, logger zerolog.Logger) (<-chan Response[U], error) {
	returnChan := make(chan Response[U], 1)

	var method serde.RequestMethod
	switch any(request).(type) {
	case directfund.ObjectiveRequest:
		method = serde.DirectFundRequestMethod
	case directdefund.ObjectiveRequest:
		method = serde.DirectDefundRequestMethod
	case virtualfund.ObjectiveRequest:
		method = serde.VirtualFundRequestMethod
	case virtualdefund.ObjectiveRequest:
		method = serde.VirtualDefundRequestMethod
	case serde.PaymentRequest:
		method = serde.PayRequestMethod
	default:
		return nil, fmt.Errorf("unknown request type %v", request)
	}
	requestId := rand.Uint64()
	message := serde.NewJsonRpcRequest(requestId, method, request)
	data, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	logger.Trace().
		Str("method", string(method)).
		Msg("sent message")

	go func() {
		responseData, err := connection.Request(method, data)
		if err != nil {
			returnChan <- Response[U]{Error: err}
		}

		logger.Trace().Msgf("Rpc client received response: %+v", responseData)

		jsonResponse := serde.JsonRpcResponse[U]{}
		err = json.Unmarshal(responseData, &jsonResponse)
		if err != nil {
			returnChan <- Response[U]{Error: err}
		}

		returnChan <- Response[U]{jsonResponse.Result, nil}
	}()

	return returnChan, nil
}
