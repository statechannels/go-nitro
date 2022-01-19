// Package client WILL contain imperative library code for running a go-nitro client inside another application.
// CURRENTLY it contains demonstration code (TODO)
package client // import "github.com/statechannels/go-nitro/client"

import (
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
)

// Client provides the interface for the consuming application
type Client struct {
	engine engine.Engine // The core business logic of the client
}

// New is the constructor for a Client. It accepts a messaging service and a chain service as injected dependencies.
func New(messageService messageservice.MessageService, chainservice chainservice.ChainService, store store.Store) Client {
	c := Client{}
	c.engine = engine.New(messageService, chainservice, store)

	// Start the engine in a go routine
	go c.engine.Run()

	return c
}

// Begin API

// CreateChannel creates a channel
func (c *Client) CreateChannel() chan engine.Response {
	// Convert the API call into an internal event.
	// Pass in a fresh, dedicated go channel to communicate the response:
	apiEvent := engine.APIEvent{Response: make(chan engine.Response)}
	// Send the event to the engine
	c.engine.FromAPI <- apiEvent
	// Return the go channel where the response will be sent.
	return apiEvent.Response
}
