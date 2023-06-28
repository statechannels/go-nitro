package serde

import (
	"github.com/statechannels/go-nitro/node/query"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

type RequestMethod string

const (
	GetAddressMethod                  RequestMethod = "get_address"
	VersionMethod                     RequestMethod = "version"
	CreateLedgerChannelRequestMethod  RequestMethod = "create_ledger_channel"
	CloseLedgerChannelRequestMethod   RequestMethod = "close_ledger_channel"
	CreatePaymentChannelRequestMethod RequestMethod = "create_payment_channel"
	ClosePaymentChannelRequestMethod  RequestMethod = "close_payment_channel"
	PayRequestMethod                  RequestMethod = "pay"
	GetPaymentChannelRequestMethod    RequestMethod = "get_payment_channel"
	GetLedgerChannelRequestMethod     RequestMethod = "get_ledger_channel"
	GetPaymentChannelsByLedgerMethod  RequestMethod = "get_payment_channels_by_ledger"
	GetAllLedgerChannelsMethod        RequestMethod = "get_all_ledger_channels"
	CreateVoucherRequestMethod        RequestMethod = "create_voucher"
	ReceiveVoucherRequestMethod       RequestMethod = "receive_voucher"
)

type NotificationMethod string

const (
	ObjectiveCompleted    NotificationMethod = "objective_completed"
	LedgerChannelUpdated  NotificationMethod = "ledger_channel_updated"
	PaymentChannelUpdated NotificationMethod = "payment_channel_updated"
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
type GetPaymentChannelsByLedgerRequest struct {
	LedgerId types.Destination
}

type (
	NoPayloadRequest = struct{}
)

type RequestPayload interface {
	directfund.ObjectiveRequest |
		directdefund.ObjectiveRequest |
		virtualfund.ObjectiveRequest |
		virtualdefund.ObjectiveRequest |
		PaymentRequest |
		GetLedgerChannelRequest |
		GetPaymentChannelRequest |
		GetPaymentChannelsByLedgerRequest |
		NoPayloadRequest |
		payments.Voucher
}

type NotificationPayload interface {
	protocols.ObjectiveId |
		query.PaymentChannelInfo |
		query.LedgerChannelInfo
}

type JsonRpcRequest[T RequestPayload | NotificationPayload] struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      uint64 `json:"id"`
	Method  string `json:"method"`
	Params  T      `json:"params"`
}

type VersionResponse = string

type (
	GetAllLedgersResponse              = []query.LedgerChannelInfo
	GetPaymentChannelsByLedgerResponse = []query.PaymentChannelInfo
)

type ResponsePayload interface {
	directfund.ObjectiveResponse |
		protocols.ObjectiveId |
		virtualfund.ObjectiveResponse |
		PaymentRequest |
		query.PaymentChannelInfo |
		query.LedgerChannelInfo |
		VersionResponse |
		GetAllLedgersResponse |
		GetPaymentChannelsByLedgerResponse |
		payments.Voucher |
		uint64
}

type JsonRpcResponse[T ResponsePayload] struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      uint64      `json:"id"`
	Result  T           `json:"result"`
	Error   interface{} `json:"error"`
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
