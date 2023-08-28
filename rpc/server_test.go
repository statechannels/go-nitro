package rpc

import (
	"encoding/json"
	"testing"

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

func sendRequestAndExpectError(t *testing.T, request []byte, expectedError serde.JsonRpcError) {
	mockNode := &nitro.Node{}

	mockResponder := &mockResponder{}
	// Since we're using an empty node we want to disable notifications
	// otherwise the server will try to send notifications to the node and fail
	_, err := newRpcServerWithoutNotifications(mockNode, mockResponder)
	if err != nil {
		t.Error(err)
	}

	response := mockResponder.Handler(request)

	jsonResponse := serde.JsonRpcErrorResponse{}
	err = json.Unmarshal(response, &jsonResponse)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expectedError, jsonResponse.Error)
}

func getAuthToken(t *testing.T) string {
	request := serde.JsonRpcSpecificRequest[serde.NoPayloadRequest]{Jsonrpc: "2.0", Id: 1, Method: "get_auth_token"}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}

	mockNode := &nitro.Node{}
	mockResponder := &mockResponder{}
	// Since we're using an empty node we want to disable notifications
	// otherwise the server will try to send notifications to the node and fail
	_, err = newRpcServerWithoutNotifications(mockNode, mockResponder)
	if err != nil {
		t.Error(err)
	}

	response := mockResponder.Handler(jsonRequest)

	jsonResponse := serde.JsonRpcSuccessResponse[string]{}
	err = json.Unmarshal(response, &jsonResponse)
	if err != nil {
		t.Error(err)
	}
	return jsonResponse.Result
}

func TestRpcParseError(t *testing.T) {
	request := []byte{}
	sendRequestAndExpectError(t, request, serde.ParseError)
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
	sendRequestAndExpectError(t, jsonRequest, serde.InvalidRequestError)
}

func TestRpcWrongVersion(t *testing.T) {
	request := serde.JsonRpcSpecificRequest[serde.PaymentRequest]{Jsonrpc: "1.0", Id: 2, Method: "direct_fund"}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	expectedError := serde.InvalidRequestError
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
	sendRequestAndExpectError(t, jsonRequest, serde.InvalidRequestError)
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
	expectedError := serde.InvalidRequestError
	sendRequestAndExpectError(t, jsonRequest, expectedError)
}

func TestRpcMethodNotFound(t *testing.T) {
	request := serde.JsonRpcSpecificRequest[serde.PaymentRequest]{Jsonrpc: "2.0", Id: 2, Method: "fake_method"}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	expectedError := serde.MethodNotFoundError
	sendRequestAndExpectError(t, jsonRequest, expectedError)
}

func TestRpcGetPaymentChannelMissingParam(t *testing.T) {
	authToken := getAuthToken(t)
	request := serde.JsonRpcSpecificRequest[serde.GetPaymentChannelRequest]{
		Jsonrpc: "2.0", Id: 2, Method: "get_payment_channel", Params: serde.Params[serde.GetPaymentChannelRequest]{AuthToken: authToken},
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	expectedError := serde.InvalidParamsError
	sendRequestAndExpectError(t, jsonRequest, expectedError)
}

func TestRpcPayInvalidParam(t *testing.T) {
	authToken := getAuthToken(t)

	paymentRequest := serde.PaymentRequest{
		Amount:  100,
		Channel: types.Destination{},
	}

	request := serde.JsonRpcSpecificRequest[serde.PaymentRequest]{
		Jsonrpc: "2.0",
		Id:      2,
		Method:  "pay",
		Params:  serde.Params[serde.PaymentRequest]{AuthToken: authToken, Payload: paymentRequest},
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	expectedError := serde.InvalidParamsError
	sendRequestAndExpectError(t, jsonRequest, expectedError)
}
