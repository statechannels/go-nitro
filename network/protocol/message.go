package netproto

import "strconv"

type MessageType int8

const (
	TypeRequest      MessageType = 1
	TypeResponse     MessageType = 2
	TypePublicEvent  MessageType = 3
	TypePrivateEvent MessageType = 4
	TypeError        MessageType = 5
)

type Message struct {
	Type      MessageType
	RequestId uint64
	Method    string
	Args      []any
}

func NewMessage(msgType MessageType, requestId uint64, method string, args []any) *Message {
	return &Message{
		Type:      msgType,
		RequestId: requestId,
		Method:    method,
		Args:      args,
	}
}

func TypeStr(t MessageType) string {
	switch t {
	case 1:
		return "TypeRequest"
	case 2:
		return "TypeResponse"
	case 3:
		return "TypePublicEvent"
	case 4:
		return "TypePrivateEvent"
	case 5:
		return "TypeError"
	}
	return strconv.Itoa(int(t))
}
