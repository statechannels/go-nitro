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
	engine            engine.Engine // The core business logic of the client
	Address           *types.Address
	ObjectiveStatuses chan<- engine.ObjectiveStatus
	ReceivedVouchers  chan<- payments.Voucher
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

	c.ObjectiveStatuses = c.engine.ObjectiveStatuses
	c.ReceivedVouchers = c.engine.ReceivedVouchers

	// Start the engine in a go routine
	go c.engine.Run()

	return c
}

// Begin API

// CreateVirtualChannel creates a virtual channel with the counterParty using ledger channels with the intermediary.
func (c *Client) CreateVirtualChannel(objectiveRequest virtualfund.ObjectiveRequest) virtualfund.ObjectiveResponse {

	apiEvent := engine.APIEvent{
		ObjectiveToSpawn: objectiveRequest,
	}
	// Send the event to the engine
	c.engine.FromAPI <- apiEvent

	return objectiveRequest.Response(*c.Address)
}

// CloseVirtualChannel attempts to close and defund the given virtually funded channel.
func (c *Client) CloseVirtualChannel(channelId types.Destination, paidToBob *big.Int) protocols.ObjectiveId {

	objectiveRequest := virtualdefund.ObjectiveRequest{
		ChannelId: channelId,
		PaidToBob: paidToBob,
	}
	apiEvent := engine.APIEvent{
		ObjectiveToSpawn: objectiveRequest,
	}
	// Send the event to the engine
	c.engine.FromAPI <- apiEvent

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

	apiEvent := engine.APIEvent{
		ObjectiveToSpawn: objectiveRequest,
	}
	// Send the event to the engine
	c.engine.FromAPI <- apiEvent

	return objectiveRequest.Response(*c.Address)

}

// CloseLedgerChannel attempts to close and defund the given directly funded channel.
func (c *Client) CloseLedgerChannel(channelId types.Destination) protocols.ObjectiveId {

	objectiveRequest := directdefund.ObjectiveRequest{
		ChannelId: channelId,
	}
	apiEvent := engine.APIEvent{
		ObjectiveToSpawn: objectiveRequest,
	}
	// Send the event to the engine
	c.engine.FromAPI <- apiEvent

	return objectiveRequest.Id(*c.Address)

}

// Pay will send a signed voucher to the payee that they can redeem for the given amount.
func (c *Client) Pay(channelId types.Destination, amount *big.Int) {

	apiEvent := engine.APIEvent{
		PaymentToMake: engine.PaymentRequest{ChannelId: channelId, Amount: amount},
	}
	// Send the event to the engine
	c.engine.FromAPI <- apiEvent
}
