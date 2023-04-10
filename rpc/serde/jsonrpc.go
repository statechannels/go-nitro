package serde

import (
	"github.com/statechannels/go-nitro/client/query"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

type RequestMethod string

const (
	DirectFundRequestMethod        RequestMethod = "direct_fund"
	DirectDefundRequestMethod      RequestMethod = "direct_defund"
	VirtualFundRequestMethod       RequestMethod = "virtual_fund"
	VirtualDefundRequestMethod     RequestMethod = "virtual_defund"
	PayRequestMethod               RequestMethod = "pay"
	GetPaymentChannelRequestMethod RequestMethod = "get_payment_channel"
	GetLedgerChannelRequestMethod  RequestMethod = "get_ledger_channel"
)

type NotificationMethod string

const (
	ObjectiveCompleted NotificationMethod = "objective_completed"
)

type NotificationOrRequest interface {
	RequestMethod | NotificationMethod
}

const JsonRpcVersion = "2.0"

type PaymentRequest struct {
	Amount  uint64
	Channel types.Destination
}
type GetPaymentChannelRequest struct {
	Id types.Destination
}
type GetLedgerChannelRequest struct {
	Id types.Destination
}

type RequestPayload interface {
	directfund.ObjectiveRequest |
		directdefund.ObjectiveRequest |
		virtualfund.ObjectiveRequest |
		virtualdefund.ObjectiveRequest |
		PaymentRequest |
		GetLedgerChannelRequest |
		GetPaymentChannelRequest
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
	directfund.ObjectiveResponse |
		protocols.ObjectiveId |
		virtualfund.ObjectiveResponse |
		PaymentRequest |
		query.PaymentChannelInfo |
		query.LedgerChannelInfo
}

type JsonRpcResponse[T ResponsePayload] struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      uint64      `json:"id"`
	Result  T           `json:"result"`
	Error   interface{} `json:"error"`
}

type JsonRpcError struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Id      uint64      `json:"id"`
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
