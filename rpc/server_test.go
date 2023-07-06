package rpc

import (
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	nitro "github.com/statechannels/go-nitro/node"
	"github.com/statechannels/go-nitro/rpc/serde"
	"github.com/statechannels/go-nitro/types"
	"github.com/stretchr/testify/assert"
)

type mockResponder struct {
	Handler func([]byte) []byte
}

func (*mockResponder) Close() error {
	return nil
}

func (*mockResponder) Url() string {
	return ""
}

func (m *mockResponder) RegisterRequestHandler(apiVersion string, handler func([]byte) []byte) error {
	m.Handler = handler
	return nil
}

func (*mockResponder) Notify([]byte) error {
	return nil
}

func sendRequestAndExpectError(t *testing.T, request []byte, expectedError types.JsonRpcError) {
	mockNode := &nitro.Node{}
	mockLogger := &zerolog.Logger{}
	mockResponder := &mockResponder{}
	// Since we're using an empty node we want to disable notifications
	// otherwise the server will try to send notifications to the node and fail
	_, err := newRpcServerWithoutNotifications(mockNode, mockLogger, mockResponder)
	if err != nil {
		t.Error(err)
	}

	response := mockResponder.Handler(request)

	jsonResponse := types.JsonRpcResponse{}
	err = json.Unmarshal(response, &jsonResponse)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expectedError, jsonResponse.ErrorObj)
}

func TestRpcParseError(t *testing.T) {
	request := []byte{}
	sendRequestAndExpectError(t, request, types.ParseError)
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
	sendRequestAndExpectError(t, jsonRequest, types.InvalidRequestError)
}

func TestRpcWrongVersion(t *testing.T) {
	request := serde.JsonRpcRequest[serde.PaymentRequest]{Jsonrpc: "1.0", Id: 2, Method: "direct_fund"}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	expectedError := types.InvalidRequestError
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
	sendRequestAndExpectError(t, jsonRequest, types.InvalidRequestError)
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
	expectedError := types.InvalidRequestError
	sendRequestAndExpectError(t, jsonRequest, expectedError)
}

func TestRpcMethodNotFound(t *testing.T) {
	request := serde.JsonRpcRequest[serde.PaymentRequest]{Jsonrpc: "2.0", Id: 2, Method: "fake_method"}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	expectedError := types.MethodNotFoundError
	sendRequestAndExpectError(t, jsonRequest, expectedError)
}

func TestRpcGetPaymentChannelMissingParam(t *testing.T) {
	request := serde.JsonRpcRequest[serde.GetPaymentChannelRequest]{Jsonrpc: "2.0", Id: 2, Method: "get_payment_channel"}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	expectedError := types.InvalidParamsError
	sendRequestAndExpectError(t, jsonRequest, expectedError)
}

func TestRpcPayInvalidParam(t *testing.T) {
	paymentRequest := serde.PaymentRequest{
		Amount:  100,
		Channel: types.Destination{},
	}

	request := serde.JsonRpcRequest[serde.PaymentRequest]{
		Jsonrpc: "2.0",
		Id:      2,
		Method:  "pay",
		Params:  paymentRequest,
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	expectedError := types.InvalidParamsError
	sendRequestAndExpectError(t, jsonRequest, expectedError)
}
