// Package messageservice is a messaging service responsible for routing messages to peers and relaying messages recieved from peers
package messageservice // import "github.com/statechannels/go-nitro/client/messageservice"

import "github.com/statechannels/go-nitro/protocols"

type MessageService interface {
	GetRecieveChan() chan protocols.Message
	GetSendChan() chan protocols.Message
	Send(message protocols.Message)
}

var recieveChan chan protocols.Message = make(chan protocols.Message)
var sendChan chan protocols.Message = make(chan protocols.Message)

type TestMessageService struct{}

func (TestMessageService) GetRecieveChan() chan protocols.Message { return recieveChan }
func (TestMessageService) GetSendChan() chan protocols.Message    { return sendChan }
func (TestMessageService) Send(message protocols.Message)         {}
