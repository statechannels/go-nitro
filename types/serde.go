package types

import "github.com/ethereum/go-ethereum/common"

// MarshalText encodes the receiver into UTF-8-encoded text and returns the result.
//
// This makes Destination an encoding.TextMarshaler, meaning it can be safely used as the key in a map which will be encoded into json.
func (d Destination) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

// UnmarshalText unmarshals the supplied text (assumed to be a valid marshaling) into the receiver.
//
// This makes Destination an encoding.TextUnmarshaler.
func (d *Destination) UnmarshalText(text []byte) error {
	*d = Destination(common.HexToHash(string(text)))
	return nil
}

type JsonRpcRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      uint64      `json:"id"`
	Method  string      `json:"code"`
	Params  interface{} `json:"params"`
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
		Jsonrpc: "2.0",
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
)
