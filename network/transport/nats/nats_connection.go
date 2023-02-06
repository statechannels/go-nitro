package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"github.com/statechannels/go-nitro/network/serde"
)

type natsConnection struct {
	nc                *nats.Conn
	natsSubscriptions []*nats.Subscription
}

func NewNatsConnection(nc *nats.Conn) *natsConnection {
	natsConnection := &natsConnection{
		nc:                nc,
		natsSubscriptions: make([]*nats.Subscription, 0),
	}

	return natsConnection
}

func (c *natsConnection) Request(topic serde.RequestMethod, data []byte) ([]byte, error) {
	msg, err := c.nc.Request(methodToTopic(topic), data, 10*time.Second)
	if msg == nil {
		return nil, fmt.Errorf("received nill data for request %v with error %w", topic, err)
	}
	return msg.Data, err
}

func (c *natsConnection) Subscribe(topic serde.RequestMethod, handler func([]byte) []byte) error {
	sub, err := c.nc.Subscribe(methodToTopic(topic), func(msg *nats.Msg) {
		responseData := handler(msg.Data)
		err := c.nc.Publish(msg.Reply, responseData)
		if err != nil {
			panic(err)
		}
	})
	c.natsSubscriptions = append(c.natsSubscriptions, sub)

	return err
}

func (c *natsConnection) Close() {
	for _, sub := range c.natsSubscriptions {
		err := c.unsubscribeFromTopic(sub, 0)
		if err != nil {
			log.Error().Err(err).Msgf("failed to unsubscribe from a topic: %s", sub.Subject)
		}
	}
}

func (c *natsConnection) unsubscribeFromTopic(sub *nats.Subscription, try int32) error {
	err := sub.Unsubscribe()
	if err != nil && try < 3 {
		return c.unsubscribeFromTopic(sub, try+1)
	}
	return nil
}

func methodToTopic(method serde.RequestMethod) string {
	return fmt.Sprintf("nitro.%s", method)
}
