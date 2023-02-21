package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

const nitroRequestTopic = "nitro-request"
const nitroNotificationTopic = "nitro-notify"

type natsConnection struct {
	nc                *nats.Conn
	natsSubscriptions []*nats.Subscription
}

type natsConnectionClient struct {
	natsConnection
}

type natsConnectionServer struct {
	natsConnection
	ns *server.Server
}

func newNatsConnection(url string) (*natsConnection, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &natsConnection{
		nc:                nc,
		natsSubscriptions: make([]*nats.Subscription, 0)}, nil

}

func NewNatsConnectionAsServer(rpcPort int) (*natsConnectionServer, error) {
	opts := &server.Options{Port: rpcPort}
	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, err
	}
	ns.Start()

	natsConnection, err := newNatsConnection(ns.ClientURL())
	if err != nil {
		return nil, err
	}

	con := &natsConnectionServer{
		natsConnection: *natsConnection,
		ns:             ns,
	}
	return con, nil
}

func NewNatsConnectionAsClient(url string) (*natsConnectionClient, error) {
	natsConnection, err := newNatsConnection(url)
	if err != nil {
		return nil, err
	}
	return &natsConnectionClient{
		natsConnection: *natsConnection,
	}, nil
}

// Request sends a blocking request for a topic with the given data
// It returns the response data and an error
func (c *natsConnectionClient) Request(data []byte) ([]byte, error) {
	msg, err := c.nc.Request(nitroRequestTopic, data, 10*time.Second)
	if msg == nil {
		return nil, fmt.Errorf("received nill data for request %v with error %w", data, err)
	}
	return msg.Data, err
}

// Respond subscribes to a topic and calls the handler function when a message is received
// It returns an error if the subscription fails
// The handler processes the incoming data and returns the response data
func (c *natsConnectionServer) Respond(handler func([]byte) []byte) error {
	sub, err := c.nc.Subscribe(nitroRequestTopic, func(msg *nats.Msg) {
		responseData := handler(msg.Data)
		err := c.nc.Publish(msg.Reply, responseData)
		if err != nil {
			panic(err)
		}
	})
	c.natsSubscriptions = append(c.natsSubscriptions, sub)

	return err
}

func (c *natsConnectionClient) Subscribe() (<-chan []byte, error) {
	notificationChan := make(chan []byte)
	subscription, err := c.nc.Subscribe(nitroNotificationTopic, func(msg *nats.Msg) {
		notificationChan <- msg.Data
	})
	c.natsSubscriptions = append(c.natsSubscriptions, subscription)

	return notificationChan, err
}

func (c *natsConnectionServer) Notify(data []byte) error {
	return c.nc.Publish(nitroNotificationTopic, data)
}

func (c *natsConnectionServer) Url() string {
	return c.ns.ClientURL()
}

func (c *natsConnectionServer) Close() {
	c.natsConnection.Close()
	c.ns.Shutdown()

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
