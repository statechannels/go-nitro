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

// Client provides the interface for the consuming application
type Client struct {
	engine              engine.Engine // The core business logic of the client
	Address             *types.Address
	CompletedObjectives chan<- engine.CompletedObjectiveEvent // All Objective updates from the engine
}

// New is the constructor for a Client. It accepts a messaging service, a chain service, and a store as injected dependencies.
func New(messageService messageservice.MessageService, chainservice chainservice.ChainService, store store.Store, logDestination io.Writer) Client {
	c := Client{}
	c.Address = store.GetAddress()
	c.engine = engine.New(messageService, chainservice, store, logDestination)
	c.CompletedObjectives = make(chan<- engine.CompletedObjectiveEvent, 100)
	// Start the engine in a go routine
	go c.engine.Run()

	return c
}

// Begin API

// CreateDirectChannel creates a directly funded channel with the given counterparty
func (c *Client) CreateDirectChannel(counterparty types.Address, appDefinition types.Address, appData types.Bytes, outcome outcome.Exit, challengeDuration *types.Uint256) (protocols.ObjectiveId, chan engine.CompletedObjectiveEvent) {
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

	clientResponse := make(chan engine.CompletedObjectiveEvent)

	// Starts a go function that listens for our objective to be completed and then dispatchs events to different client channels
	go func() {
		for update := range c.engine.ToApi() {
			for _, completed := range update.CompletedObjectives {

				if completed == objective.Id() {

					// We dispatch an event to the channel that handles **all** objective updates.
					// This provides a central place to monitor for objective updates.
					c.CompletedObjectives <- engine.CompletedObjectiveEvent{Id: completed}
					// We dispatch an event to the channel returned by this function. This channel is only used for the one objective.
					// This allows for promise like behaviour where we can wait for the objective to be completed before continuing.
					clientResponse <- engine.CompletedObjectiveEvent{Id: completed}
					// We can close the channel as the objective is completed so there will be no further events.
					close(clientResponse)

					// Exit the go function now that our objective is completed.
					return
				}
			}
		}
	}()

	return objective.Id(), clientResponse
}
