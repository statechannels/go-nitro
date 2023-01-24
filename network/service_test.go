package network

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/rs/zerolog"
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	regReqHandler = "regReqHandler"
	regErrHandler = "regErrHandler"
	regResHandler = "regResHandler"
)

type connMock struct {
	mock.Mock
}

func (c *connMock) Send(s string, bytes []byte) {}

func (c *connMock) Recv() ([]byte, error) {
	return []byte("event"), nil
}

func (c *connMock) Close() {}

type serializerMock struct {
	mock.Mock
}

func (s *serializerMock) Serialize(message *netproto.Message) ([]byte, error) {
	return json.Marshal(message)
}

func (s *serializerMock) Deserialize(bytes []byte) (*netproto.Message, error) {
	msg := new(netproto.Message)
	err := json.Unmarshal(bytes, msg)
	return msg, err
}

func newNetworkServiceMock() *NetworkService {
	service := NetworkService{
		Logger:              zerolog.Logger{},
		Connection:          &connMock{},
		Serde:               &serializerMock{},
		handlerRequest:      sync.Map{},
		handlerResponse:     sync.Map{},
		handlerError:        sync.Map{},
		handlerPublicEvent:  sync.Map{},
		handlerPrivateEvent: sync.Map{},
	}

	return &service
}

// maybe logic will get more complicated and gonna keep this as example
func TestRegisterUnregisterRequestHandler(t *testing.T) {
	service := newNetworkServiceMock()

	service.RegisterRequestHandler(regReqHandler, func(message *netproto.Message) {})
	val, ok := service.handlerRequest.Load(regReqHandler)
	assert.NotNil(t, val)
	assert.Equal(t, ok, true)
	service.UnregisterRequestHandler(regReqHandler)
	val, ok = service.handlerRequest.Load(regReqHandler)
	assert.Nil(t, val)
	assert.Equal(t, ok, false)
}

func TestGetHandler(t *testing.T) {
	service := newNetworkServiceMock()
	msg1 := &netproto.Message{
		Type:      netproto.TypeRequest,
		RequestId: 0,
		Method:    regReqHandler,
		Args:      nil,
	}
	msg2 := &netproto.Message{
		Type:      netproto.TypeResponse,
		RequestId: 1,
		Method:    regResHandler,
		Args:      nil,
	}
	msg3 := &netproto.Message{
		Type:      netproto.TypeError,
		RequestId: 2,
		Method:    regErrHandler,
		Args:      nil,
	}
	msg4 := &netproto.Message{
		Type:      netproto.TypePrivateEvent,
		RequestId: 5,
		Method:    "",
		Args:      nil,
	}

	service.RegisterRequestHandler(regReqHandler, func(message *netproto.Message) {})
	service.RegisterResponseHandler(regResHandler, func(message *netproto.Message) {})
	service.RegisterErrorHandler(regErrHandler, func(message *netproto.Message) {})

	val := service.getHandler(msg1)
	assert.NotNil(t, val)

	val = service.getHandler(msg2)
	assert.NotNil(t, val)

	val = service.getHandler(msg3)
	assert.NotNil(t, val)

	val = service.getHandler(msg4)
	assert.Nil(t, val)
}

func TestHandleMessage(t *testing.T) {
	service := newNetworkServiceMock()
	msg1 := &netproto.Message{
		Type:      netproto.TypeRequest,
		RequestId: 0,
		Method:    regReqHandler,
		Args:      nil,
	}
	//msg2 := &netproto.Message{
	//	Type:      netproto.TypeResponse,
	//	RequestId: 1,
	//	Method:    regResHandler,
	//	Args:      nil,
	//}

	service.RegisterRequestHandler(regReqHandler, func(message *netproto.Message) {})

	service.handleMessage(msg1)
}
