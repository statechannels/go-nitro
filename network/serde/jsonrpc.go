package serde

import (
	"encoding/json"
	"fmt"

	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
)

const JsonRpcVersion = "2.0"

type MessageType int8

const (
	TypeRequest  MessageType = 1
	TypeResponse MessageType = 2
	TypeError    MessageType = 3
)

type MethodType string

const (
	DirectFund    MethodType = "direct_fund"
	DirectDefund  MethodType = "direct_defund"
	VirtualFund   MethodType = "virtual_fund"
	VirtualDefund MethodType = "virtual_defund"
)

type JsonRpcRequestResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      uint64      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error"`
}
type RequestPayload interface {
	directfund.ObjectiveRequest | directdefund.ObjectiveRequest | virtualfund.ObjectiveRequest | virtualdefund.ObjectiveRequest
}
type JsonRpcRequest[T RequestPayload] struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      uint64 `json:"id"`
	Method  string `json:"method"`
	Params  T      `json:"params"`
}
type ResponsePayload interface {
	directfund.ObjectiveResponse | protocols.ObjectiveId | virtualfund.ObjectiveResponse
}
type JsonRpcResponse[T ResponsePayload] struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      uint64      `json:"id"`
	Result  T           `json:"result"`
	Error   interface{} `json:"error"`
}

func NewJsonRpcRequest[T RequestPayload](requestId uint64, method MethodType, objectiveRequest T) *JsonRpcRequest[T] {

	return &JsonRpcRequest[T]{
		Jsonrpc: JsonRpcVersion,
		Id:      requestId,
		Method:  string(method),
		Params:  objectiveRequest,
	}
}

func NewJsonRpcResponse[T ResponsePayload](requestId uint64, objectiveResponse T) *JsonRpcResponse[T] {
	return &JsonRpcResponse[T]{
		Jsonrpc: JsonRpcVersion,
		Id:      requestId,
		Result:  objectiveResponse,
		Error:   nil,
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
