// Package client WILL contain imperative library code for running a go-nitro client inside another application.
// CURRENTLY it contains demonstration code (TODO)
package client // import "github.com/statechannels/go-nitro/client"

import (
	"io"
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/protocols"
	directfund "github.com/statechannels/go-nitro/protocols/direct-fund"
	"github.com/statechannels/go-nitro/types"
)

type ChannelReadyEvent struct {
	ChannelId   types.Destination
	ObjectiveId protocols.ObjectiveId
}

// Client provides the interface for the consuming application
type Client struct {
	engine       engine.Engine // The core business logic of the client
	Address      *types.Address
	ChannelReady chan ChannelReadyEvent // All Objective updates from the engine
	listeners    map[protocols.ObjectiveId][]chan<- ChannelReadyEvent
}

// New is the constructor for a Client. It accepts a messaging service, a chain service, and a store as injected dependencies.
func New(messageService messageservice.MessageService, chainservice chainservice.ChainService, store store.Store, logDestination io.Writer) Client {
	c := Client{}
	c.Address = store.GetAddress()
	c.engine = engine.New(messageService, chainservice, store, logDestination)
	c.ChannelReady = make(chan ChannelReadyEvent, 100)
	c.listeners = make(map[protocols.ObjectiveId][]chan<- ChannelReadyEvent)
	// Start the engine in a go routine
	go c.engine.Run()

	go func() {
		for update := range c.engine.ToApi() {

			for _, completed := range update.CompletedObjectives {

				channelId := completed.Channels()[0]

				event := ChannelReadyEvent{ObjectiveId: completed.Id(), ChannelId: channelId}
				// We dispatch an event to the channel that handles **all** objective updates.
				// This provides a central place to monitor for objective updates.
				c.ChannelReady <- event

				// Dispatch an event to any listeners that have been registered
				if listeners, ok := c.listeners[completed.Id()]; ok {
					for _, l := range listeners {
						l <- event
					}
				}
			}
		}
	}()

	return c
}

// Begin API

// CreateDirectChannel creates a directly funded channel with the given counterparty
func (c *Client) CreateDirectChannel(counterparty types.Address, appDefinition types.Address, appData types.Bytes, outcome outcome.Exit, challengeDuration *types.Uint256) (protocols.ObjectiveId, chan ChannelReadyEvent) {
	// Convert the API call into an internal event.
	objective, _ := directfund.New(true,
		state.State{
			ChainId:           big.NewInt(0), // TODO
			Participants:      []types.Address{*c.Address, counterparty},
			ChannelNonce:      big.NewInt(0), // TODO -- how do we get a fresh nonce safely without race conditions? Could we conisder a random nonce?
			AppDefinition:     appDefinition,
			ChallengeDuration: challengeDuration,
			AppData:           appData,
			Outcome:           outcome,
			TurnNum:           0,
			IsFinal:           false,
		},
		*c.Address,
	)

	// Pass in a fresh, dedicated go channel to communicate the response:
	apiEvent := engine.APIEvent{
		ObjectiveToSpawn: objective,
	}
	// Send the event to the engine
	c.engine.FromAPI <- apiEvent

	clientResponse := make(chan ChannelReadyEvent)
	// We register a listener for the objective id
	c.listeners[objective.Id()] = append(c.listeners[objective.Id()], clientResponse)

	return objective.Id(), clientResponse
}
