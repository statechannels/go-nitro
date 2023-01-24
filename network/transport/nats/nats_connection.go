package nats

import (
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type natsConnection struct {
	nc *nats.Conn

	subTopicNames     []string
	msgChannel        chan *nats.Msg
	mutex             sync.Mutex
	natsSubscriptions []*nats.Subscription
}

func NewNatsConnection(nc *nats.Conn, subTopicNames []string) *natsConnection {
	natsConnection := &natsConnection{
		nc:                nc,
		subTopicNames:     subTopicNames,
		msgChannel:        make(chan *nats.Msg, 128),
		natsSubscriptions: make([]*nats.Subscription, len(subTopicNames)),
		mutex:             sync.Mutex{},
	}
	go natsConnection.subscribeToTopics()

	return natsConnection
}

func (c *natsConnection) subscribeToTopics() {
	for _, topic := range c.subTopicNames {
		sub, err := c.subscribeToTopic(topic, 0)
		if err != nil {
			log.Error().Err(err).Msgf("failed to connect on topic: %s", topic)
		} else {
			c.mutex.Lock()
			c.natsSubscriptions = append(c.natsSubscriptions, sub)
		}
	}
}

func (c *natsConnection) subscribeToTopic(topic string, try int32) (*nats.Subscription, error) {
	sub, err := c.nc.ChanSubscribe(topic, c.msgChannel)

	if err != nil && try < 3 {
		time.Sleep(time.Millisecond * 100)
		return c.subscribeToTopic(topic, try+1)
	}

	return sub, nil
}

func (c *natsConnection) Send(t string, data []byte) {
	err := c.nc.Publish(t, data)
	if err != nil {
		log.Error().Err(err).Msgf("failed to send message on topic: %s. msg: %s", t, string(data))
	}
	log.Trace().Msgf("published message on %s.\ndata: %v", t, string(data))
}

func (c *natsConnection) Recv() ([]byte, error) {
	msg := <-c.msgChannel

	return msg.Data, nil
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
	close(c.msgChannel)
}
