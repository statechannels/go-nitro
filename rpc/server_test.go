package rpc

import (
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	nitro "github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/rpc/serde"
	"github.com/stretchr/testify/assert"
)

type mockResponder struct {
	Handler func([]byte) []byte
}

func (*mockResponder) Close() {}
func (*mockResponder) Url() string {
	return ""
}
func (m *mockResponder) RegisterRequestHandler(handler func([]byte) []byte) error {
	m.Handler = handler
	return nil
}
func (*mockResponder) Notify([]byte) error {
	return nil
}

func sendRequestAndExpectError(t *testing.T, request []byte, expectedError serde.JsonRpcError) {
	mockClient := &nitro.Client{}
	mockLogger := zerolog.Logger{}
	mockResponder := &mockResponder{}
	_, err := NewRpcServer(mockClient, mockLogger, mockResponder)
	if err != nil {
		t.Error(err)
	}

	response := mockResponder.Handler(request)

	jsonError := serde.JsonRpcError{}
	err = json.Unmarshal(response, &jsonError)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, jsonError, expectedError)
}

func TestRpcParseError(t *testing.T) {
	request := []byte{}
	expectedError := serde.JsonRpcError{Code: -32700, Message: "Parse error"}
	sendRequestAndExpectError(t, request, expectedError)
}

// TestRpcInvalidRequest exists to point out that the server needs request validation improvement.
// The server receives a valid json object that does not contain any of the required feilds of a request
// like Id or Method. The server should return an error like {-32600:	Invalid Request} or
// {-32602:	Invalid params}.
func TestRpcInvalidRequest(t *testing.T) {
	type InvalidRequest struct {
		Message string
	}

	request := InvalidRequest{Message: "hello"}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	expectedError := serde.JsonRpcError{Code: -32601, Message: "Method not found", Id: 0}
	sendRequestAndExpectError(t, jsonRequest, expectedError)
}

func TestRpcMethodNotFound(t *testing.T) {
	request := serde.JsonRpcMessage{Method: "invalid", Id: 1}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	expectedError := serde.JsonRpcError{Code: -32601, Message: "Method not found", Id: 1}
	sendRequestAndExpectError(t, jsonRequest, expectedError)
}
