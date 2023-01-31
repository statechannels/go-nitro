package serde

import (
	"encoding/json"
	"fmt"

	"github.com/statechannels/go-nitro/protocols/directfund"
)

const JsonRpcVersion = "2.0"

type MessageType int8

const (
	TypeRequest  MessageType = 1
	TypeResponse MessageType = 2
	TypeError    MessageType = 3
)

type JsonRpcRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      uint64      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type JsonRpcResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      uint64      `json:"id"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error"`
}

type JsonRpcRequestResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      uint64      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error"`
}

type JsonRpcDirectFundRequest struct {
	Jsonrpc          string                      `json:"jsonrpc"`
	Id               uint64                      `json:"id"`
	Method           string                      `json:"method"`
	ObjectiveRequest directfund.ObjectiveRequest `json:"params"`
}

type JsonRpcDirectFundResponse struct {
	Jsonrpc           string                       `json:"jsonrpc"`
	Id                uint64                       `json:"id"`
	ObjectiveResponse directfund.ObjectiveResponse `json:"result"`
	Error             interface{}                  `json:"error"`
}

type JsonRpc struct{}

func NewDirectFundRequestMessage(requestId uint64, objectiveRequest directfund.ObjectiveRequest) *JsonRpcDirectFundRequest {
	return &JsonRpcDirectFundRequest{
		Jsonrpc:          JsonRpcVersion,
		Id:               requestId,
		Method:           "direct_fund",
		ObjectiveRequest: objectiveRequest,
	}
}

func NewDirectFundResponseMessage(requestId uint64, objectiveResponse directfund.ObjectiveResponse) *JsonRpcDirectFundResponse {
	return &JsonRpcDirectFundResponse{
		Jsonrpc:           JsonRpcVersion,
		Id:                requestId,
		ObjectiveResponse: objectiveResponse,
		Error:             nil,
	}
}

func Deserialize(data []byte) (*JsonRpcRequestResponse, MessageType, error) {
	jm := JsonRpcRequestResponse{}
	err := json.Unmarshal(data, &jm)
	if jm.Error != nil {
		return &jm, TypeError, err
	}
	if jm.Result != nil {
		return &jm, TypeResponse, err
	}
	if jm.Method != "" {
		return &jm, TypeRequest, err
	}

	return nil, TypeError, fmt.Errorf("unexpected jsonrpc message format: %s", string(data))
}
