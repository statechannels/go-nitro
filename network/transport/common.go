package transport

import "github.com/nats-io/nats.go"

func InitNats(connectionUrl string) (*nats.Conn, error) {
	nc, err := nats.Connect(connectionUrl)
	return nc, err
}
