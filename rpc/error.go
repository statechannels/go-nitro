package rpc

import "github.com/statechannels/go-nitro/rpc/serde"

var (
	parseError                       = serde.JsonRpcError{Code: -32700, Message: "Parse error"}
	invalidRequestError              = serde.JsonRpcError{Code: -32600, Message: "Invalid Request"}
	methodNotFoundError              = serde.JsonRpcError{Code: -32601, Message: "Method not found"}
	unexpectedRequestUnmarshalError  = serde.JsonRpcError{Code: -32010, Message: "Unexpected unmarshal error"}
	unexpectedRequestUnmarshalError2 = serde.JsonRpcError{Code: -32009, Message: "Unexpected unmarshal error"}
)
