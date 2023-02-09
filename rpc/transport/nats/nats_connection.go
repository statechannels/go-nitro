package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"github.com/statechannels/go-nitro/rpc/serde"
)

type natsConnection struct {
	nc                *nats.Conn
	natsSubscriptions []*nats.Subscription
}

// NewNatsConnection creates an instance of a Connection interface that uses the nats transport
func NewNatsConnection(nc *nats.Conn) *natsConnection {
	natsConnection := &natsConnection{
		nc:                nc,
		natsSubscriptions: make([]*nats.Subscription, 0),
	}

	return natsConnection
}

// Request sends a blocking request for a topic with the given data
// It returns the response data and an error
func (c *natsConnection) Request(topic serde.RequestMethod, data []byte) ([]byte, error) {
	msg, err := c.nc.Request(methodToTopic(topic), data, 10*time.Second)
	if msg == nil {
		return nil, fmt.Errorf("received nill data for request %v with error %w", topic, err)
	}
	return msg.Data, err
}

// Respond subscribes to a topic and calls the handler function when a message is received
// It returns an error if the subscription fails
// The handler processes the incoming data and returns the response data
func (c *natsConnection) Respond(topic serde.RequestMethod, handler func([]byte) []byte) error {
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

func (c *natsConnection) Notify(topic serde.NotificationMethod, data []byte) error {
	return c.nc.Publish(methodToTopic(topic), data)
}

func (c *natsConnection) Subscribe(topic serde.NotificationMethod, handler func([]byte)) error {
	sub, err := c.nc.Subscribe(methodToTopic(topic), func(msg *nats.Msg) {
		handler(msg.Data)
	})
	c.natsSubscriptions = append(c.natsSubscriptions, sub)

	return err
}

// Close shuts down the connection
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

func methodToTopic[T serde.RequestMethod | serde.NotificationMethod](method T) string {
	return fmt.Sprintf("nitro.%s", method)
}
