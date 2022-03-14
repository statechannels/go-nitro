// Package client contains imperative library code for running a go-nitro client inside another application.
package client // import "github.com/statechannels/go-nitro/client"

import (
	"io"
	"math/rand"

	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

// Client provides the interface for the consuming application
type Client struct {
	engine              engine.Engine // The core business logic of the client
	Address             *types.Address
	completedObjectives chan protocols.ObjectiveId
}

// New is the constructor for a Client. It accepts a messaging service, a chain service, and a store as injected dependencies.
func New(messageService messageservice.MessageService, chainservice chainservice.ChainService, store store.Store, logDestination io.Writer) Client {
	c := Client{}
	c.Address = store.GetAddress()
	c.engine = engine.New(messageService, chainservice, store, logDestination)
	c.completedObjectives = make(chan protocols.ObjectiveId, 100)

	// Start the engine in a go routine
	go c.engine.Run()

	// Start the event handler in a go routine
	// It will listen for events from the engine and dispatch events to client channels
	go c.handleEngineEvents()

	return c
}

// handleEngineEvents is responsible for monitoring the ToApi channel on the engine.
// It parses events from the ToApi chan and then dispatches events to the necessary client chan.
func (c *Client) handleEngineEvents() {
	for update := range c.engine.ToApi() {

		for _, completed := range update.CompletedObjectives {

			c.completedObjectives <- completed.Id()

		}

	}
}

// Begin API

// CompletedObjectives returns a chan that receives a objective id whenever that objective is completed
func (c *Client) CompletedObjectives() <-chan protocols.ObjectiveId {
	return c.completedObjectives
}

// CreateVirtualChannel creates a virtual channel with the counterParty using ledger channels with the intermediary.
func (c *Client) CreateVirtualChannel(counterParty types.Address, intermediary types.Address, appDefinition types.Address, appData types.Bytes, outcome outcome.Exit, challengeDuration *types.Uint256) protocols.ObjectiveId {
	objectiveRequest := virtualfund.ObjectiveRequest{
		MyAddress:         *c.Address,
		CounterParty:      counterParty,
		Intermediary:      intermediary,
		AppDefinition:     appDefinition,
		AppData:           appData,
		Outcome:           outcome,
		ChallengeDuration: challengeDuration,
		Nonce:             rand.Int63(), // TODO: Proper nonce management!
	}
	apiEvent := engine.APIEvent{
		ObjectiveToSpawn: objectiveRequest,
	}
	// Send the event to the engine
	c.engine.FromAPI <- apiEvent

	return objectiveRequest.Id()
}

// CreateDirectChannel creates a directly funded channel with the given counterparty
func (c *Client) CreateDirectChannel(counterparty types.Address, appDefinition types.Address, appData types.Bytes, outcome outcome.Exit, challengeDuration *types.Uint256) protocols.ObjectiveId {
	// Convert the API call into an internal event.
	objectiveRequest := directfund.ObjectiveRequest{
		MyAddress:         *c.Address,
		CounterParty:      counterparty,
		AppDefinition:     appDefinition,
		AppData:           appData,
		Outcome:           outcome,
		ChallengeDuration: challengeDuration,
		Nonce:             rand.Int63(), // TODO: Proper nonce management!
	}
	apiEvent := engine.APIEvent{
		ObjectiveToSpawn: objectiveRequest,
	}
	// Send the event to the engine
	c.engine.FromAPI <- apiEvent

	return objectiveRequest.Id()

}
