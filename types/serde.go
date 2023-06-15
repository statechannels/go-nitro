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

type JsonRpcError struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Id      uint64      `json:"id"`
}

func (e JsonRpcError) Error() string {
	return e.Message
}

var (
	ParseError                       = JsonRpcError{Code: -32700, Message: "Parse error"}
	InvalidRequestError              = JsonRpcError{Code: -32600, Message: "Invalid Request"}
	MethodNotFoundError              = JsonRpcError{Code: -32601, Message: "Method not found"}
	InvalidParamsError               = JsonRpcError{Code: -32602, Message: "Invalid params"}
	InternalServerError              = JsonRpcError{Code: -32603, Message: "Internal error"}
	UnexpectedRequestUnmarshalError  = JsonRpcError{Code: -32010, Message: "Unexpected unmarshal error"}
	UnexpectedRequestUnmarshalError2 = JsonRpcError{Code: -32009, Message: "Unexpected unmarshal error"}
)
