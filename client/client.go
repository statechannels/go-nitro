// Package client WILL contain imperative library code for running a go-nitro client inside another application.
// CURRENTLY it contains demonstration code (TODO)
package client // import "github.com/statechannels/go-nitro/client"

import (
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	directfund "github.com/statechannels/go-nitro/protocols/direct-fund"
	"github.com/statechannels/go-nitro/types"
)

// Client provides the interface for the consuming application
type Client struct {
	engine  engine.Engine // The core business logic of the client
	Address types.Address // Identifier for this client
}

func GeneratePrivateKeyAndAddress() (types.Bytes, types.Address) {
	channelSecretKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	channelSecretKeyBytes := crypto.FromECDSA(channelSecretKey)

	publicKey := channelSecretKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return channelSecretKeyBytes, address
}

// New is the constructor for a Client. It accepts a messaging service, achain service and a store as injected dependencies.
func New(
	messageService messageservice.MessageService,
	chainservice chainservice.ChainService,
	store store.Store,
	channelSecretKeyBytes types.Bytes,
	address types.Address) Client {

	c := Client{}

	// Store channel secret key and associated address
	store.SetChannelSecretKey(channelSecretKeyBytes)
	c.Address = address

	// Construct a new Engine
	c.engine = engine.New(messageService, chainservice, store)

	// Start the engine in a go routine
	go c.engine.Run()

	return c
}

// Begin API

// CreateDirectChannel creates a directly funded channel with the given counterparty
func (c *Client) CreateDirectChannel(counterparty types.Address, appDefinition types.Address, appData types.Bytes, outcome outcome.Exit, challengeDuration *types.Uint256) chan engine.Response {
	// Convert the API call into an internal event.
	objective, _ := directfund.New(
		state.State{
			ChainId:           big.NewInt(0), // TODO
			Participants:      []types.Address{c.Address, counterparty},
			ChannelNonce:      big.NewInt(0), // TODO -- how do we get a fresh nonce safely without race conditions? Could we conisder a random nonce?
			AppDefinition:     appDefinition,
			ChallengeDuration: challengeDuration,
			AppData:           appData,
			Outcome:           outcome,
			TurnNum:           0,
			IsFinal:           false,
		},
		c.Address,
	)

	// Pass in a fresh, dedicated go channel to communicate the response:
	apiEvent := engine.APIEvent{
		ObjectiveToSpawn: objective,
		Response:         make(chan engine.Response)}

	// Send the event to the engine
	c.engine.FromAPI <- apiEvent
	// Return the go channel where the response will be sent.
	return apiEvent.Response
}
