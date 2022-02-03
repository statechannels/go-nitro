// Package client contains imperative library code for running a go-nitro client inside another application.
package client // import "github.com/statechannels/go-nitro/client"

import (
	"io"
	"log"
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/types"
)

// ChannelFundedEvents is returned in a channel when a state channel has been successfly funded.
type ChannelFundedEvent struct {
	ChannelId   types.Destination
	ObjectiveId protocols.ObjectiveId
}

// Client provides the interface for the consuming application
type Client struct {
	engine        engine.Engine // The core business logic of the client
	Address       *types.Address
	ChannelFunded chan ChannelFundedEvent
	listeners     map[protocols.ObjectiveId]chan<- ChannelFundedEvent
}

// New is the constructor for a Client. It accepts a messaging service, a chain service, and a store as injected dependencies.
func New(messageService messageservice.MessageService, chainservice chainservice.ChainService, store store.Store, logDestination io.Writer) Client {
	c := Client{}
	c.Address = store.GetAddress()
	c.engine = engine.New(messageService, chainservice, store, logDestination)
	c.ChannelFunded = make(chan ChannelFundedEvent, 100)
	c.listeners = make(map[protocols.ObjectiveId]chan<- ChannelFundedEvent)
	// Start the engine in a go routine
	go c.engine.Run()

	// Start the event handler in a go routine
	// It will listen for events from the engine and dispatch events to client channels
	go c.handleEngineEvents()

	return c
}

func (c *Client) registerListener(objectiveId protocols.ObjectiveId, listener chan<- ChannelFundedEvent) {
	c.listeners[objectiveId] = listener
}
func (c *Client) removeListener(objectiveId protocols.ObjectiveId) {
	l, ok := c.listeners[objectiveId]
	if !ok {
		log.Fatalf("Could not find listener for objective id %s", objectiveId)
	}
	close(l)
	delete(c.listeners, objectiveId)
}

// handleEngineEvents is responsible for monitoring the ToApi channel on the engine.
// It parses events from the ToApi channel and then dispatches events to the necessary client channels
func (c *Client) handleEngineEvents() {
	for update := range c.engine.ToApi() {

		for _, completed := range update.CompletedObjectives {
			// TODO: We're assuming the first channel id is the one we're interested in.
			channelId := completed.Channels()[0]
			event := ChannelFundedEvent{ObjectiveId: completed.Id(), ChannelId: channelId}

			// We dispatch an event to the channel that handles **all** objective updates.
			// This provides a central place to monitor for objective updates.
			c.ChannelFunded <- event

			// Dispatch an event to any listeners that have been registered by calls to CreateDirectChannel
			if l, ok := c.listeners[completed.Id()]; ok {
				l <- event
				// Since the objective is completed we no longer need the listener
				c.removeListener(completed.Id())
			}

		}

	}
}

// Begin API

// CreateDirectChannel creates a directly funded channel with the given counterparty
func (c *Client) CreateDirectChannel(counterparty types.Address, appDefinition types.Address, appData types.Bytes, outcome outcome.Exit, challengeDuration *types.Uint256) (protocols.ObjectiveId, chan ChannelFundedEvent) {
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

	clientResponse := make(chan ChannelFundedEvent)
	// We register a listener so a ChannelFundedEvent will get sent to clientResponse when the objective completes
	c.registerListener(objective.Id(), clientResponse)

	return objective.Id(), clientResponse
}
