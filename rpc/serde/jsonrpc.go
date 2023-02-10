package serde

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

type RequestMethod string

const (
	DirectFundRequestMethod    RequestMethod = "direct_fund"
	DirectDefundRequestMethod  RequestMethod = "direct_defund"
	VirtualFundRequestMethod   RequestMethod = "virtual_fund"
	VirtualDefundRequestMethod RequestMethod = "virtual_defund"
	PayRequestMethod           RequestMethod = "pay"
)

type NotificationMethod string

const (
	ObjectiveCompleted NotificationMethod = "objective_completed"
)

const JsonRpcVersion = "2.0"

type PaymentRequest struct {
	Amount  uint64
	Channel types.Destination
}
type RequestPayload interface {
	directfund.ObjectiveRequest |
		directdefund.ObjectiveRequest |
		virtualfund.ObjectiveRequest |
		virtualdefund.ObjectiveRequest |
		PaymentRequest
}

type NotificationPayload interface {
	protocols.ObjectiveId
}

type JsonRpcRequest[T RequestPayload | NotificationPayload] struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      uint64 `json:"id"`
	Method  string `json:"method"`
	Params  T      `json:"params"`
}
type ResponsePayload interface {
	directfund.ObjectiveResponse | protocols.ObjectiveId | virtualfund.ObjectiveResponse | PaymentRequest
}
type JsonRpcResponse[T ResponsePayload] struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      uint64      `json:"id"`
	Result  T           `json:"result"`
	Error   interface{} `json:"error"`
}

type JsonRpcError struct {
	Code    uint64      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewJsonRpcRequest[T RequestPayload | NotificationPayload, U RequestMethod | NotificationMethod](requestId uint64, method U, objectiveRequest T) *JsonRpcRequest[T] {

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
