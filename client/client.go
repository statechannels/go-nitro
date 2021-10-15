// Package client contains imperative library code for running a go-nitro client inside another application
package client // import "github.com/statechannels/go-nitro/protocols"

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/protocols"

	"github.com/statechannels/go-nitro/types"
)

type Client struct {
	engine Engine
	api    chan APIEvent
}

func NewClient() Client {
	c := Client{}

	// create the engine's inbound channels
	c.engine.api = make(chan APIEvent)
	c.engine.chain = make(chan ChainEvent)
	c.engine.inbox = make(chan Message)

	// create the engine's outbound channel
	c.engine.client = make(chan Response)

	go c.engine.Run()

	return c
}

// CreateChannel creates a channel
func (c *Client) CreateChannel() chan Response {
	apiEvent := APIEvent{Response: make(chan Response)} // Make a dedicated go channel to communicate the response
	c.engine.api <- apiEvent                            // The API call is "converted" into an internal event sent to the engine
	return apiEvent.Response

}

type APIEvent struct {
	ObjectiveToSpawn   protocols.Objective   // try this first
	ObjectiveToReject  protocols.ObjectiveId // then this
	ObjectiveToApprove protocols.ObjectiveId // then this

	Response chan Response
}
type ChainEvent struct {
	ChannelId          types.Bytes32
	Holdings           big.Int
	AdjudicationStatus protocols.Status
}
type Message struct {
	ObjectiveId protocols.ObjectiveId
	Sigs        map[*state.State]state.Signature // mapping from state to signature TODO consider using a hash of the state
}

type Response struct{}
