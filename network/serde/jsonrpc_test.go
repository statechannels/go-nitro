package serde

import (
	"testing"

	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonRpcSerializeRequest(t *testing.T) {
	rpc := JsonRpc{}
	data, err := rpc.Serialize(&netproto.Message{
		Type:      netproto.TypeRequest,
		RequestId: 4242,
		Method:    "test",
		Args:      []interface{}{"foo"},
	})
	require.NoError(t, err)
	assert.Equal(t, `{"jsonrpc":"2.0","id":4242,"method":"test","params":["foo"]}`, string(data))
}

func TestJsonRpcSerializeResponse(t *testing.T) {
	rpc := JsonRpc{}
	data, err := rpc.Serialize(&netproto.Message{
		Type:      netproto.TypeResponse,
		Method:    "bar",
		RequestId: 4242,
		Args:      []interface{}{"foo"},
	})
	require.NoError(t, err)
	assert.Equal(t, `{"jsonrpc":"2.0","id":4242,"method":"bar","params":null,"result":"foo","error":null}`, string(data))
}

func TestJsonRpcSerializeError(t *testing.T) {
	rpc := JsonRpc{}
	data, err := rpc.Serialize(&netproto.Message{
		Type:      netproto.TypeError,
		RequestId: 123,
		Method:    "test",
		Args:      []interface{}{-32601, "Method not found"},
	})
	require.NoError(t, err)
	assert.Equal(t, `{"jsonrpc":"2.0","id":123,"result":null,"error":[-32601,"Method not found"]}`, string(data))
}

func TestJsonRpcDeserializeRequest(t *testing.T) {
	rpc := JsonRpc{}
	m, err := rpc.Deserialize([]byte(`{"method":"test","params":["foo"],"jsonrpc":"2.0","id":4242}`))
	require.NoError(t, err)
	assert.Equal(t, &netproto.Message{
		Type:      netproto.TypeRequest,
		RequestId: 4242,
		Method:    "test",
		Args:      []interface{}{"foo"},
	}, m)
}

func TestJsonRpcDeserializeResponse(t *testing.T) {
	rpc := JsonRpc{}
	m, err := rpc.Deserialize([]byte(`{"result":"foo","error":null,"id":4242}`))
	require.NoError(t, err)
	assert.Equal(t, &netproto.Message{
		Type:      netproto.TypeResponse,
		RequestId: 4242,
		Args:      []interface{}{"foo"},
	}, m)
}

func TestJsonRpcDeserializeError(t *testing.T) {
	rpc := JsonRpc{}
	m, err := rpc.Deserialize([]byte(`{"jsonrpc":"2.0","error":[-32601,"Method not found"],"id":456}`))
	require.NoError(t, err)
	assert.Equal(t, &netproto.Message{
		Type:      netproto.TypeError,
		RequestId: 456,
		Args:      []interface{}{float64(-32601), "Method not found"},
	}, m)
}

func BenchmarkJsonRpcSerializeRequest(b *testing.B) {
	rpc := JsonRpc{}
	m := &netproto.Message{
		Type:      netproto.TypeRequest,
		RequestId: 4242,
		Method:    "test",
		Args:      []interface{}{"foo"},
	}
	for i := 0; i < b.N; i++ {
		_, err := rpc.Serialize(m)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkJsonRpcDeserializeRequest(b *testing.B) {
	rpc := JsonRpc{}
	data := []byte(`{"method":"test","params":["foo"],"jsonrpc":"2.0","id":4242}`)
	for i := 0; i < b.N; i++ {
		_, err := rpc.Deserialize(data)
		if err != nil {
			b.Error(err)
		}
	}
}
