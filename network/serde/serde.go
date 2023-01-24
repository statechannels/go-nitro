package serde

import netproto "github.com/statechannels/go-nitro/network/protocol"

type Serde interface {
	Serializer
	Deserializer
}

type Serializer interface {
	Serialize(*netproto.Message) ([]byte, error)
}

type Deserializer interface {
	Deserialize([]byte) (*netproto.Message, error)
}
