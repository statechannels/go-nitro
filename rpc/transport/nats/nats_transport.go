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

type natsTransport struct {
	nc                *nats.Conn
	natsSubscriptions []*nats.Subscription
}

type natsTransportClient struct {
	natsTransport
	notificationChan chan []byte
}

type natsTransportServer struct {
	natsTransport
	ns *server.Server
}

func newNatsTransport(url string) (*natsTransport, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &natsTransport{
		nc:                nc,
		natsSubscriptions: make([]*nats.Subscription, 0)}, nil

}

func NewNatsTransportAsServer(rpcPort int) (*natsTransportServer, error) {
	opts := &server.Options{Port: rpcPort}
	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, err
	}
	ns.Start()

	natsTransport, err := newNatsTransport(ns.ClientURL())
	if err != nil {
		return nil, err
	}

	con := &natsTransportServer{
		natsTransport: *natsTransport,
		ns:            ns,
	}
	return con, nil
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
	msg, err := c.nc.Request(nitroRequestTopic, data, 10*time.Second)
	if msg == nil {
		return nil, fmt.Errorf("received nill data for request %v with error %w", data, err)
	}
	return msg.Data, err
}

func (c *natsTransportServer) RegisterRequestHandler(handler func([]byte) []byte) error {
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

func (c *natsTransportServer) Notify(data []byte) error {
	return c.nc.Publish(nitroNotificationTopic, data)
}

func (c *natsTransportServer) Url() string {
	return c.ns.ClientURL()
}

func (c *natsTransportServer) Close() {
	c.natsTransport.Close()
	c.ns.Shutdown()
}

func (c *natsTransport) Close() {
	for _, sub := range c.natsSubscriptions {
		err := c.unsubscribeFromTopic(sub, 0)
		if err != nil {
			log.Error().Err(err).Msgf("failed to unsubscribe from a topic: %s", sub.Subject)
		}
	}
	c.nc.Close()
}

func (c *natsTransport) unsubscribeFromTopic(sub *nats.Subscription, try int32) error {
	err := sub.Unsubscribe()
	if err != nil && try < 3 {
		return c.unsubscribeFromTopic(sub, try+1)
	}
	return nil
}
