package network

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/client/engine/store/safesync"
	"github.com/statechannels/go-nitro/network/serde"
	"github.com/statechannels/go-nitro/network/transport"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
)

type ClientConnection struct {
	Connection transport.Connection
}

type Response struct {
	Data  any
	Error error
}

func Request[T serde.RequestPayload](cc *ClientConnection, request T, logger zerolog.Logger, idsToMethods *safesync.Map[serde.RequestMethod]) (<-chan Response, error) {
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
		return nil, fmt.Errorf("Unknown request type %v", request)
	}
	requestId := rand.Uint64()
	message := serde.NewJsonRpcRequest(requestId, method, request)
	data, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	idsToMethods.Store(string(fmt.Sprintf("%d", requestId)), method)

	topic := fmt.Sprintf("nitro.%s", method)

	logger.Trace().
		Str("method", string(method)).
		Msg("sent message")

	go func() {
		responseData, err := cc.Connection.Request(topic, data)
		if err != nil {
			returnChan <- Response{nil, err}
		}

		logger.Trace().Msgf("Rpc client received response: %+v", responseData)
		switch any(request).(type) {
		case directfund.ObjectiveRequest:
			unmarshalAndSend(responseData, directfund.ObjectiveResponse{}, returnChan)
		case directdefund.ObjectiveRequest:
		case virtualdefund.ObjectiveRequest:
			unmarshalAndSend(responseData, protocols.ObjectiveId(""), returnChan)
		case virtualfund.ObjectiveRequest:
			unmarshalAndSend(responseData, virtualfund.ObjectiveResponse{}, returnChan)
		case serde.PaymentRequest:
			unmarshalAndSend(responseData, serde.PaymentRequest{}, returnChan)
		default:
			returnChan <- Response{nil, fmt.Errorf("Unknown response for request %v", request)}
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
