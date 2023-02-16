package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type natsConnection struct {
	nc                *nats.Conn
	natsSubscriptions []*nats.Subscription
	ns                *server.Server
}

func NewNatsConnectionAsServer(rpcPort int) (*natsConnection, error) {
	opts := &server.Options{Port: rpcPort}
	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, err
	}
	ns.Start()

	con, err := NewNatsConnectionAsClient(ns.ClientURL())
	if err != nil {
		return nil, err
	}
	con.ns = ns
	return con, nil
}

func NewNatsConnectionAsClient(url string) (*natsConnection, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	con := &natsConnection{
		nc:                nc,
		natsSubscriptions: make([]*nats.Subscription, 0),
		ns:                nil,
	}
	return con, nil
}

// Request sends a blocking request for a topic with the given data
// It returns the response data and an error
func (c *natsConnection) Request(id uint64, data []byte) ([]byte, error) {
	msg, err := c.nc.Request("nitro-request", data, 10*time.Second)
	if msg == nil {
		return nil, fmt.Errorf("received nill data for request %v with error %w", data, err)
	}
	return msg.Data, err
}

// Respond subscribes to a topic and calls the handler function when a message is received
// It returns an error if the subscription fails
// The handler processes the incoming data and returns the response data
func (c *natsConnection) Respond(handler func([]byte) []byte) error {
	sub, err := c.nc.Subscribe("nitro-request", func(msg *nats.Msg) {
		responseData := handler(msg.Data)
		err := c.nc.Publish(msg.Reply, responseData)
		if err != nil {
			panic(err)
		}
	})
	c.natsSubscriptions = append(c.natsSubscriptions, sub)

	return err
}

func (c *natsConnection) Notify(data []byte) error {
	return c.nc.Publish("nitro-notify", data)
}

func (c *natsConnection) Subscribe() (<-chan []byte, error) {
	notificationChan := make(chan []byte)
	subscription, err := c.nc.Subscribe("nitro-notify", func(msg *nats.Msg) {
		notificationChan <- msg.Data
	})
	c.natsSubscriptions = append(c.natsSubscriptions, subscription)

	return notificationChan, err
}

// Close shuts down the connection
func (c *natsConnection) Close() {
	for _, sub := range c.natsSubscriptions {
		err := c.unsubscribeFromTopic(sub, 0)
		if err != nil {
			log.Error().Err(err).Msgf("failed to unsubscribe from a topic: %s", sub.Subject)
		}
	}
	if c.ns != nil {
		c.ns.Shutdown()
	}
}

func (c *natsConnection) Url() string {
	return c.ns.ClientURL()
}

func (c *natsConnection) unsubscribeFromTopic(sub *nats.Subscription, try int32) error {
	err := sub.Unsubscribe()
	if err != nil && try < 3 {
		return c.unsubscribeFromTopic(sub, try+1)
	}
	return nil
}
