package serde

import (
	"github.com/ethereum/go-ethereum/common"

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
	GetAuthTokenMethod                RequestMethod = "get_auth_token"
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

type AuthRequest struct {
	Id string
}
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
		AuthRequest |
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

type Params[T RequestPayload | NotificationPayload] struct {
	AuthToken string `json:"authtoken"`
	Payload   T      `json:"payload"`
}

type JsonRpcSpecificRequest[T RequestPayload | NotificationPayload] struct {
	Jsonrpc string    `json:"jsonrpc"`
	Id      uint64    `json:"id"`
	Method  string    `json:"method"`
	Params  Params[T] `json:"params"`
}

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
		GetAllLedgersResponse |
		GetPaymentChannelsByLedgerResponse |
		payments.Voucher |
		common.Address |
		string |
		payments.ReceiveVoucherSummary
}

type JsonRpcSuccessResponse[T ResponsePayload] struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      uint64 `json:"id"`
	Result  T      `json:"result"`
}

func NewJsonRpcSpecificRequest[T RequestPayload | NotificationPayload, U RequestMethod | NotificationMethod](requestId uint64, method U, objectiveRequest T, authToken string) *JsonRpcSpecificRequest[T] {
	return &JsonRpcSpecificRequest[T]{
		Jsonrpc: JsonRpcVersion,
		Id:      requestId,
		Method:  string(method),
		Params:  Params[T]{AuthToken: authToken, Payload: objectiveRequest},
	}
}

func NewJsonRpcResponse[T ResponsePayload](requestId uint64, objectiveResponse T) *JsonRpcSuccessResponse[T] {
	return &JsonRpcSuccessResponse[T]{
		Jsonrpc: JsonRpcVersion,
		Id:      requestId,
		Result:  objectiveResponse,
	}
}

type JsonRpcGeneralRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      uint64      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"payload"`
}

type JsonRpcGeneralResponse struct {
	Jsonrpc string       `json:"jsonrpc"`
	Id      uint64       `json:"id"`
	Error   JsonRpcError `json:"error"`
	Result  interface{}  `json:"result"`
}

type JsonRpcErrorResponse struct {
	Jsonrpc string       `json:"jsonrpc"`
	Id      uint64       `json:"id"`
	Error   JsonRpcError `json:"error"`
}

type JsonRpcError struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (e JsonRpcError) Error() string {
	return e.Message
}

func NewJsonRpcErrorResponse(requestId uint64, error JsonRpcError) *JsonRpcErrorResponse {
	return &JsonRpcErrorResponse{
		Jsonrpc: JsonRpcVersion,
		Id:      requestId,
		Error:   error,
	}
}

var (
	ParseError            = JsonRpcError{Code: -32700, Message: "Parse error"}
	InvalidRequestError   = JsonRpcError{Code: -32600, Message: "Invalid Request"}
	MethodNotFoundError   = JsonRpcError{Code: -32601, Message: "Method not found"}
	InvalidParamsError    = JsonRpcError{Code: -32602, Message: "Invalid params"}
	InternalServerError   = JsonRpcError{Code: -32603, Message: "Internal error"}
	RequestUnmarshalError = JsonRpcError{Code: -32010, Message: "Could not unmarshal request object"}
	ParamsUnmarshalError  = JsonRpcError{Code: -32009, Message: "Could not unmarshal params object"}
	InvalidAuthTokenError = JsonRpcError{Code: -32008, Message: "Invalid auth token"}
)
