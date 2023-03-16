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
	mockLogger := &zerolog.Logger{}
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
	assert.Equal(t, expectedError, jsonError)
}

func TestRpcParseError(t *testing.T) {
	request := []byte{}
	sendRequestAndExpectError(t, request, parseError)
}

func TestRpcMissingRequiredFields(t *testing.T) {
	type InvalidRequest struct {
		Message string `json:"message"`
	}

	request := InvalidRequest{Message: "direct_fund"}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	sendRequestAndExpectError(t, jsonRequest, invalidRequestError)
}

func TestRpcWrongVersion(t *testing.T) {
	request := serde.JsonRpcRequest[serde.PaymentRequest]{Jsonrpc: "1.0", Id: 2, Method: "direct_fund"}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	expectedError := invalidRequestError
	expectedError.Id = 2
	sendRequestAndExpectError(t, jsonRequest, expectedError)
}

func TestRpcIncorrectId(t *testing.T) {
	type InvalidRequest struct {
		Jsonrpc string  `json:"jsonrpc"`
		Id      float64 `json:"id"`
		Method  string  `json:"method"`
	}
	request := InvalidRequest{Jsonrpc: "1.0", Id: 2.2, Method: "direct_fund"}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	sendRequestAndExpectError(t, jsonRequest, invalidRequestError)
}

func TestRpcMissingMethod(t *testing.T) {
	type InvalidRequest struct {
		Jsonrpc string `json:"jsonrpc"`
		Id      uint64 `json:"id"`
	}
	request := InvalidRequest{Jsonrpc: "1.0", Id: 2}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	expectedError := invalidRequestError
	expectedError.Id = 2
	sendRequestAndExpectError(t, jsonRequest, expectedError)
}

func TestRpcMethodNotFound(t *testing.T) {
	request := serde.JsonRpcRequest[serde.PaymentRequest]{Jsonrpc: "2.0", Id: 2, Method: "direct_funds"}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	expectedError := methodNotFoundError
	expectedError.Id = 2
	sendRequestAndExpectError(t, jsonRequest, expectedError)
}
