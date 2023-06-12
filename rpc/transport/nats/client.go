package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type natsTransportClient struct {
	natsTransport
	notificationChan chan []byte
}

func NewNatsTransportAsClient(url string) (*natsTransportClient, error) {
	natsTransport, err := newNatsTransport(url)
	if err != nil {
		return nil, err
	}
	return &natsTransportClient{
		natsTransport: *natsTransport,
	}, nil
}

func (c *natsTransportClient) Request(data []byte) ([]byte, error) {
	msg, err := c.nc.Request(nitroRequestTopic+apiVersionPath, data, 10*time.Second)
	if msg == nil {
		return nil, fmt.Errorf("received nill data for request %v with error %w", data, err)
	}
	return msg.Data, err
}

func (c *natsTransportClient) Subscribe() (<-chan []byte, error) {
	if c.notificationChan != nil {
		return c.notificationChan, nil
	}
	c.notificationChan = make(chan []byte)
	subscription, err := c.nc.Subscribe(nitroNotificationTopic, func(msg *nats.Msg) {
		c.notificationChan <- msg.Data
	})
	c.natsSubscriptions = append(c.natsSubscriptions, subscription)

	return c.notificationChan, err
}

func (c *natsTransportClient) Close() {
	c.natsTransport.Close()
	close(c.notificationChan)
}
