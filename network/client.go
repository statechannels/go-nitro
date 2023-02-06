package network

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/client/engine/store/safesync"
	"github.com/statechannels/go-nitro/network/serde"
	"github.com/statechannels/go-nitro/network/transport"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
)

type ClientConnection struct {
	Connection transport.Connection
}

func Request[T serde.RequestPayload](cc *ClientConnection, request T, logger zerolog.Logger, idsToMethods *safesync.Map[serde.RequestMethod]) error {
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
		return fmt.Errorf("Unknown request type %v", request)
	}
	requestId := rand.Uint64()
	message := serde.NewJsonRpcRequest(requestId, method, request)
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	idsToMethods.Store(string(fmt.Sprintf("%d", requestId)), method)

	topic := fmt.Sprintf("nitro.%s", method)
	cc.Connection.Send(topic, data)

	logger.Trace().
		Str("method", string(method)).
		Msg("sent message")

	return nil
}
