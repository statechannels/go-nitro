// Package client contains imperative library code for running a go-nitro client inside another application.
package client // import "github.com/statechannels/go-nitro/client"

import (
	"io"
	"math/big"

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

// CreateVirtualChannel creates a virtual channel with the counterParty using ledger channels with the intermediary.
func (c *Client) CreateVirtualChannel(objectiveRequest virtualfund.ObjectiveRequest) virtualfund.ObjectiveResponse {

	// Send the event to the engine
	c.engine.ObjectiveRequestsFromAPI <- objectiveRequest

	return objectiveRequest.Response(*c.Address)
}

// CloseVirtualChannel attempts to close and defund the given virtually funded channel.
func (c *Client) CloseVirtualChannel(channelId types.Destination, paidToBob *big.Int) protocols.ObjectiveId {

	objectiveRequest := virtualdefund.ObjectiveRequest{
		ChannelId: channelId,
		PaidToBob: paidToBob,
	}

	// Send the event to the engine
	c.engine.ObjectiveRequestsFromAPI <- objectiveRequest

	return objectiveRequest.Id(*c.Address)

}

// CreateLedgerChannel creates a directly funded ledger channel with the given counterparty.
// The channel will run under full consensus rules (it is not possible to provide a custom AppDefinition or AppData).
func (c *Client) CreateLedgerChannel(request directfund.ObjectiveRequestForConsensusApp) directfund.ObjectiveResponse {

	objectiveRequest := directfund.ObjectiveRequest{
		ObjectiveRequestForConsensusApp: request,
		AppDefinition:                   c.engine.GetConsensusAppAddress(),
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
