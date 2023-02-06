package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type natsConnection struct {
	nc *nats.Conn

	subTopicNames     []string
	natsSubscriptions []*nats.Subscription
}

func NewNatsConnection(nc *nats.Conn, subTopicNames []string) *natsConnection {
	natsConnection := &natsConnection{
		nc:                nc,
		subTopicNames:     subTopicNames,
		natsSubscriptions: make([]*nats.Subscription, len(subTopicNames)),
	}

	return natsConnection
}

func (c *natsConnection) Request(t string, data []byte) ([]byte, error) {
	msg, err := c.nc.Request(t, data, 10*time.Second)
	if msg == nil {
		return nil, fmt.Errorf("Received nill data for request %v with error %w", t, err)
	}
	return msg.Data, err
}

func (c *natsConnection) Subscribe(t string, handler func([]byte) []byte) error {
	sub, err := c.nc.Subscribe(t, func(msg *nats.Msg) {
		responseData := handler(msg.Data)
		err := c.nc.Publish(msg.Reply, responseData)
		if err != nil {
			panic(err)
		}
	})
	c.natsSubscriptions = append(c.natsSubscriptions, sub)

	return err
}

func (c *natsConnection) unsubscribeFromTopics() {
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

func (c *natsConnection) Close() {
	c.unsubscribeFromTopics()
}
