package network

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/network/serde"
	"github.com/statechannels/go-nitro/network/transport"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
)

// Response includes JSONRPC response data or an error.
// todo: Response should include an unmarshalled payload instead of the entire JSONRPC response.
type Response struct {
	Data  any
	Error error
}

// Request uses the supplied connection and payload to send a non-blocking JSONRPC request.
// It returns a channel that sends a JSONRPC response. If the request fails to send, an error is returned.
func Request[T serde.RequestPayload](connection transport.Connection, request T, logger zerolog.Logger) (<-chan Response, error) {
	returnChan := make(chan Response, 1)

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
			returnChan <- Response{nil, err}
		}

		logger.Trace().Msgf("Rpc client received response: %+v", responseData)
		switch any(request).(type) {
		case directfund.ObjectiveRequest:
			unmarshalAndSend(responseData, directfund.ObjectiveResponse{}, returnChan)
		case directdefund.ObjectiveRequest, virtualdefund.ObjectiveRequest:
			unmarshalAndSend(responseData, protocols.ObjectiveId(""), returnChan)
		case virtualfund.ObjectiveRequest:
			unmarshalAndSend(responseData, virtualfund.ObjectiveResponse{}, returnChan)
		case serde.PaymentRequest:
			unmarshalAndSend(responseData, serde.PaymentRequest{}, returnChan)
		default:
			returnChan <- Response{nil, fmt.Errorf("unknown response for request %v", request)}
		}
	}()

	return returnChan, nil
}

func unmarshalAndSend[P serde.ResponsePayload, T serde.JsonRpcResponse[P]](data []byte, payloadType P, resChan chan<- Response) {
	response := T{}
	err := json.Unmarshal(data, &response)
	if err != nil {
		resChan <- Response{nil, err}
	}

	resChan <- Response{response, nil}
}
