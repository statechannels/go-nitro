// Package client contains imperative library code for running a go-nitro client inside another application.
package client // import "github.com/statechannels/go-nitro/client"

import (
	"io"
	"math/big"
	"math/rand"

	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

// Client provides the interface for the consuming application
type Client struct {
	engine              engine.Engine // The core business logic of the client
	Address             *types.Address
	completedObjectives chan protocols.ObjectiveId
	failedObjectives    chan protocols.ObjectiveId
	receivedVouchers    chan payments.Voucher
}

// New is the constructor for a Client. It accepts a messaging service, a chain service, and a store as injected dependencies.
func New(messageService messageservice.MessageService, chainservice chainservice.ChainService, store store.Store, logDestination io.Writer, policymaker engine.PolicyMaker, metricsApi engine.MetricsApi) Client {
	c := Client{}
	c.Address = store.GetAddress()
	// If a metrics API is not provided we used the no-op version which does nothing.
	if metricsApi == nil {
		metricsApi = &engine.NoOpMetrics{}
	}

	c.engine = engine.New(messageService, chainservice, store, logDestination, policymaker, metricsApi)
	c.completedObjectives = make(chan protocols.ObjectiveId, 100)
	c.failedObjectives = make(chan protocols.ObjectiveId, 100)
	// Using a larger buffer since payments can be sent frequently.
	c.receivedVouchers = make(chan payments.Voucher, 1000)
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

		for _, erred := range update.FailedObjectives {
			c.failedObjectives <- erred
		}

		for _, payment := range update.ReceivedVouchers {

			c.receivedVouchers <- payment
		}

	}
}

// Begin API

// CompletedObjectives returns a chan that receives a objective id whenever that objective is completed
func (c *Client) CompletedObjectives() <-chan protocols.ObjectiveId {
	return c.completedObjectives
}

// FailedObjectives returns a chan that receives an objective id whenever that objective has failed
func (c *Client) FailedObjectives() <-chan protocols.ObjectiveId {
	return c.failedObjectives
}

// ReceivedVouchers returns a chan that receives a voucher every time we receive a payment voucher
func (c *Client) ReceivedVouchers() <-chan payments.Voucher {
	return c.receivedVouchers
}

// CreateVirtualChannel creates a virtual channel with the counterParty using ledger channels
// with the supplied intermediaries.
func (c *Client) CreateVirtualPaymentChannel(Intermediaries []types.Address, CounterParty types.Address, ChallengeDuration uint32, Outcome outcome.Exit) virtualfund.ObjectiveResponse {

	objectiveRequest := virtualfund.ObjectiveRequest{
		Intermediaries:    Intermediaries,
		CounterParty:      CounterParty,
		ChallengeDuration: ChallengeDuration,
		Outcome:           Outcome,
		Nonce:             rand.Uint64(),
		AppDefinition:     c.engine.GetVirtualPaymentAppAddress(),
	}

	// Send the event to the engine
	c.engine.ObjectiveRequestsFromAPI <- objectiveRequest

	return objectiveRequest.Response(*c.Address)
}

// CreateVirtualChannel creates a virtual channel with the counterParty using ledger channels
// with the supplied intermediaries.
func (c *Client) CreateVirtualChannel(Intermediaries []types.Address, CounterParty types.Address, ChallengeDuration uint32, Outcome outcome.Exit) virtualfund.ObjectiveResponse {
	objectiveRequest := virtualfund.ObjectiveRequest{
		Intermediaries:    Intermediaries,
		CounterParty:      CounterParty,
		ChallengeDuration: ChallengeDuration,
		Outcome:           Outcome,
		Nonce:             rand.Uint64(),
		// TODO: This should be the address of the virtual channel app
		// (adding more rules, based on AppData, and being backed up by ledger channel consensus app)
		AppDefinition: types.Address{},
	}

	// Send the event to the engine
	c.engine.ObjectiveRequestsFromAPI <- objectiveRequest

	return objectiveRequest.Response(*c.Address)
}

// ProgressChannel will send a signed state to the counterparty that they can use to progress the channel.
func (c *Client) ProgressChannel(channelId types.Destination, outcome outcome.Exit) protocols.ObjectiveId {
	// TODO: introduce new protocol, that will move a state from turnNum N to turnNum N+1
	// This should work with a direct or virtual channel

	// TODO: introduce AppData later on, but for now lets only focus on the outcome

	// NOTE: The protocol will internaly take charge of incrementing the state turn number

	// RFC: name of the protocol: progress, stateprogress, ...?
	// This name should reflect that it is an objective protocol that simply ask for each participant to agree and sign a state protopoal
	// A bit like how LedgerProposals message works for direct channels

	objectiveRequest := stateprogress.ObjectiveRequest{
		ChannelId: channelId,
		Outcome:   outcome,
	}

	// Send the event to the engine
	c.engine.ObjectiveRequestsFromAPI <- objectiveRequest

	return objectiveRequest.Id(*c.Address)
}

// Pay will send a payment to the channel counterparty.
func (c *Client) _Pay(channelId types.Destination, asset types.Address, amount *big.Int) protocols.ObjectiveId {
	// NOTE: Here we reuse the "progress" channel protocol directily
	// This method is simply an helper for outcome manipulation
	// This kind of method is usefull for both direct and virtual channels

	return c.ProgressChannel(channelId, outcome.Exit{
		outcome.SingleAssetExit{
			Asset:       asset,
			Allocations: []outcome.Allocation{
				// TODO: modify allocations according to the payment
			},
		},
	})
}

// CloseVirtualChannel attempts to close and defund the given virtually funded channel.
func (c *Client) CloseVirtualChannel(channelId types.Destination) protocols.ObjectiveId {

	objectiveRequest := virtualdefund.ObjectiveRequest{
		ChannelId: channelId,
	}

	// Send the event to the engine
	c.engine.ObjectiveRequestsFromAPI <- objectiveRequest

	return objectiveRequest.Id(*c.Address)

}

// CreateLedgerChannel creates a directly funded ledger channel with the given counterparty.
// The channel will run under full consensus rules (it is not possible to provide a custom AppDefinition or AppData).
func (c *Client) CreateLedgerChannel(Counterparty types.Address, ChallengeDuration uint32, outcome outcome.Exit) directfund.ObjectiveResponse {
	// TODO: introduce a CreateChannel method, that will simply use objectives to create a channel (direct or virtual)
	// This method will take all channel params (or a struct), allowing to create very customized channels (outcomes and appData)

	// NOTE: CreateLedgerChannel will then simply be a specialaized version of CreateChannel, that will use the consensus app
	// And CreateVirtualChannel will be a specialized version of CreateChannel, that will update outcome of ledger channel with guarentee allocations

	objectiveRequest := directfund.ObjectiveRequest{
		CounterParty:      Counterparty,
		ChallengeDuration: ChallengeDuration,
		Outcome:           outcome,
		AppDefinition:     c.engine.GetConsensusAppAddress(),
		Nonce:             rand.Uint64(),
		// Appdata implicitly zero
	}

	// Send the event to the engine
	c.engine.ObjectiveRequestsFromAPI <- objectiveRequest

	return objectiveRequest.Response(*c.Address)

}

// CloseLedgerChannel attempts to close and defund the given directly funded channel.
func (c *Client) CloseLedgerChannel(channelId types.Destination) protocols.ObjectiveId {

	objectiveRequest := directdefund.ObjectiveRequest{
		ChannelId: channelId,
	}

	// Send the event to the engine
	c.engine.ObjectiveRequestsFromAPI <- objectiveRequest

	return objectiveRequest.Id(*c.Address)

}

// Pay will send a signed voucher to the payee that they can redeem for the given amount.
func (c *Client) Pay(channelId types.Destination, amount *big.Int) {
	// Send the event to the engine
	c.engine.PaymentRequestsFromAPI <- engine.PaymentRequest{ChannelId: channelId, Amount: amount}
}
