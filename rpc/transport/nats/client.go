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
	requestFn := func(data []byte) (*nats.Msg, error) {
		return c.nc.Request(nitroRequestTopic+apiVersionPath, data, 10*time.Second)
	}

	numTries := 2
	var err error
	var msg *nats.Msg
	for i := 0; i < numTries; i++ {
		msg, err = requestFn(data)
		if msg != nil && err == nil {
			return msg.Data, nil
		}

		// Skip sleep after the last try
		if lastTry := i == numTries-1; lastTry {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	return nil, fmt.Errorf("received nill data for request %v with error %w", string(data), err)
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

const unsubRetries = 3

func (c *natsTransportClient) Close() error {
	// TODO: This is a workaround for https://github.com/nats-io/nats.go/issues/1396
	// See https://github.com/nats-io/nats.go/issues/1396#issuecomment-1714261643
	for _, sub := range c.natsSubscriptions {
		err := c.unsubscribeFromTopic(sub, unsubRetries)
		if err != nil {
			return err
		}
	}
	err := c.natsTransport.Close()
	if err != nil {
		return err
	}
	close(c.notificationChan)
	return nil
}
