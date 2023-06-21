package nats

import (
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

const (
	nitroRequestTopic      = "nitro-request"
	nitroNotificationTopic = "nitro-notify"
	apiVersionPath         = "/api/v1"
)

type natsTransport struct {
	nc                *nats.Conn
	natsSubscriptions []*nats.Subscription
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
		natsSubscriptions: make([]*nats.Subscription, 0),
	}, nil
}

func (c *natsTransport) Close() error {
	for _, sub := range c.natsSubscriptions {
		err := c.unsubscribeFromTopic(sub, 3)
		if err != nil {
			log.Error().Err(err).Msgf("failed to unsubscribe from a topic: %s", sub.Subject)
			return err
		}
	}
	c.nc.Close()
	return nil
}

// unsubscribeFromTopic will attempt to unsubscribe the supplied subscription. On error, it will retry up to retries times.
func (c *natsTransport) unsubscribeFromTopic(sub *nats.Subscription, retries uint32) error {
	var try uint32
	err := sub.Unsubscribe()
	for err != nil && try < retries {
		err = sub.Unsubscribe()
		try++
	}
	return err
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

func (c *natsTransportServer) RegisterRequestHandler(apiVersion string, handler func([]byte) []byte) error {
	sub, err := c.nc.Subscribe(nitroRequestTopic+"/api/"+apiVersion, func(msg *nats.Msg) {
		responseData := handler(msg.Data)
		err := c.nc.Publish(msg.Reply, responseData)
		if err != nil {
			panic(err)
		}
	})
	c.natsSubscriptions = append(c.natsSubscriptions, sub)

	return err
}

func (c *natsTransportServer) Notify(data []byte) error {
	return c.nc.Publish(nitroNotificationTopic, data)
}

func (c *natsTransportServer) Url() string {
	return c.ns.ClientURL()
}

func (c *natsTransportServer) Close() error {
	err := c.natsTransport.Close()
	if err != nil {
		return err
	}
	c.ns.Shutdown()
	return nil
}
