package network

import (
	"sync"
	"testing"

	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/network/serde"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type connMock struct {
	mock.Mock
}

func (c *connMock) Send(s string, bytes []byte) {}

func (c *connMock) Recv() ([]byte, error) {
	return []byte("event"), nil
}

func (c *connMock) Close() {}

func newNetworkServiceMock() *NetworkService {
	service := NetworkService{
		Logger:         zerolog.Logger{},
		Connection:     &connMock{},
		handlerRequest: sync.Map{},
	}

	return &service
}

// maybe logic will get more complicated and gonna keep this as example
func TestRegisterUnregisterRequestHandler(t *testing.T) {
	service := newNetworkServiceMock()

	service.RegisterRequestHandler(serde.DirectFundRequestMethod, func(uint64, []byte) {})
	val, ok := service.handlerRequest.Load(serde.DirectFundRequestMethod)
	assert.NotNil(t, val)
	assert.Equal(t, ok, true)
	service.UnregisterRequestHandler(serde.DirectFundRequestMethod)
	val, ok = service.handlerRequest.Load(serde.DirectFundRequestMethod)
	assert.Nil(t, val)
	assert.Equal(t, ok, false)
}

func TestGetHandler(t *testing.T) {
	service := newNetworkServiceMock()

	service.RegisterRequestHandler(serde.DirectFundRequestMethod, func(uint64, []byte) {})
	service.RegisterResponseHandler(func(uint64, []byte) {})
	service.RegisterErrorHandler(func(uint64, []byte) {})

	val := service.getHandler(serde.DirectFundRequestMethod, serde.TypeRequest)
	assert.NotNil(t, val)

	val = service.getHandler("", serde.TypeResponse)
	assert.NotNil(t, val)

	val = service.getHandler(serde.DirectFundRequestMethod, serde.TypeError)
	assert.NotNil(t, val)

	val = service.getHandler(serde.DirectDefundRequestMethod, serde.TypeRequest)
	assert.Nil(t, val)
}

func TestHandleMessage(t *testing.T) {
	service := newNetworkServiceMock()
	service.RegisterRequestHandler(serde.DirectFundRequestMethod, func(uint64, []byte) {})
	service.handleMessage(0, serde.DirectFundRequestMethod, serde.TypeRequest, []byte{})
}
