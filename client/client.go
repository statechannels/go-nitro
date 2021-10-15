// Package client contains imperative library code for running a go-nitro client inside another application
package client // import "github.com/statechannels/go-nitro/protocols"

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/protocols"

	"github.com/statechannels/go-nitro/types"
)

// APIEvent is an internal representation of an API call
type APIEvent struct {
	ObjectiveToSpawn   protocols.Objective   // try this first
	ObjectiveToReject  protocols.ObjectiveId // then this
	ObjectiveToApprove protocols.ObjectiveId // then this

	Response chan Response
}

// ChainEvent is an internal representation of a blockchain event
type ChainEvent struct {
	ChannelId          types.Bytes32
	Holdings           map[types.Address]big.Int // indexed by asset
	AdjudicationStatus protocols.AdjudicationStatus
}

// Message is an internal representation of a message from another client
type Message struct {
	ObjectiveId protocols.ObjectiveId
	Sigs        map[types.Bytes32]state.Signature // mapping from state hash to signature
}

// Response is the return type that asynchronous API calls "resolve to". Such a call returns a go channel of type Response.
type Response struct{}

// Client provides the interface for the consuming application
type Client struct {
	engine Engine // The core business logic of the client
}

// NewClient is the constructor for a Client
func NewClient() Client {
	c := Client{}

	// create the engine's inbound channels
	c.engine.api = make(chan APIEvent)
	c.engine.chain = make(chan ChainEvent)
	c.engine.inbox = make(chan Message)

	// create the engine's outbound channel
	c.engine.client = make(chan Response)

	// Start the engine in a go routine
	go c.engine.Run()

	return c
}

// Begin API

// CreateChannel creates a channel
func (c *Client) CreateChannel() chan Response {
	// Convert the API call into an internal event.
	// Pass in a fresh, dedicated go channel to communicate the response:
	apiEvent := APIEvent{Response: make(chan Response)}
	// Send the event to the engine
	c.engine.api <- apiEvent
	// Return the go channel where the response will be sent.
	return apiEvent.Response
}
