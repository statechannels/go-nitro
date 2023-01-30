package serde

import (
	"encoding/json"
	"fmt"

	netproto "github.com/statechannels/go-nitro/network/protocol"
)

const JsonRpcVersion = "2.0"

type JsonRpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      uint64        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type JsonRpcResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      uint64      `json:"id"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error"`
}

type JsonRpcRequestResponse struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      uint64        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Result  interface{}   `json:"result"`
	Error   interface{}   `json:"error"`
}

type JsonRpc struct{}

func (j *JsonRpc) Serialize(m *netproto.Message) ([]byte, error) {
	switch m.Type {
	case netproto.TypeRequest:
		return json.Marshal(&JsonRpcRequest{
			Jsonrpc: JsonRpcVersion,
			Id:      m.RequestId,
			Method:  m.Method,
			Params:  m.Args,
		})

	case netproto.TypeError:
		return json.Marshal(&JsonRpcResponse{
			Jsonrpc: JsonRpcVersion,
			Id:      m.RequestId,
			Result:  nil,
			Error:   m.Args,
		})

	case netproto.TypeResponse:
		return json.Marshal(&JsonRpcRequestResponse{
			Jsonrpc: JsonRpcVersion,
			Id:      m.RequestId,
			Result:  m.Args,
			Error:   nil,
			Method:  m.Method,
		})

	}

	return nil, fmt.Errorf("unsupported message type %s", netproto.TypeStr(m.Type))
}

func (j *JsonRpc) Deserialize(data []byte) (*netproto.Message, error) {
	jm := JsonRpcRequestResponse{}
	err := json.Unmarshal(data, &jm)
	if err != nil {
		return nil, err
	}

	if jm.Error != nil {
		m := netproto.Message{
			Type:      netproto.TypeError,
			RequestId: jm.Id,
			Args:      jm.Error.([]interface{}),
		}
		return &m, nil
	}

	if jm.Result != nil {
		m := netproto.Message{
			Type:      netproto.TypeResponse,
			RequestId: jm.Id,
			Method:    jm.Method,
			Args:      []interface{}{jm.Result},
		}
		return &m, nil
	}

	if jm.Method != "" {
		m := netproto.Message{
			Type:      netproto.TypeRequest,
			RequestId: jm.Id,
			Method:    jm.Method,
			Args:      jm.Params,
		}
		return &m, nil
	}

	return nil, fmt.Errorf("unexpected jsonrpc message format: %s", string(data))
}
