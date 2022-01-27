package example

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	typ "github.com/statechannels/go-nitro/types"
)

type StateChannel interface {
	Transfer(opts *bind.TransactOpts, assetIndex *big.Int, fromChannelId [32]byte, outcomeBytes []byte, stateHash [32]byte, indices []*big.Int) (*types.Transaction, error)
	TransferAllAssets(opts *bind.TransactOpts, channelId [32]byte, outcomeBytes []byte, stateHash [32]byte) (*types.Transaction, error)
	Deposit(opts *bind.TransactOpts, asset common.Address, channelId [32]byte, expectedHeld *big.Int, amount *big.Int) (*types.Transaction, error)
	ValidTransition(opts *bind.CallOpts, nParticipants *big.Int, isFinalAB [2]bool, ab [2]IForceMoveAppVariablePart, turnNumB *big.Int, appDefinition common.Address) (bool, error)
	GetChainID(opts *bind.CallOpts) (*big.Int, error)
}

type Client struct {
	// TODO
	// Add more fields
	// What else do we need here?
	contract StateChannel
}

func NewClient(contractAddr, rpcUrl string) (*Client, error) {
	contractAddress := common.HexToAddress(contractAddr)
	ethClient, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}

	adjucator, err := NewNitroAdjucator(contractAddress, ethClient)
	if err != nil {
		return nil, err
	}

	return &Client{
		contract: adjucator,
	}, nil

}

// To define participants
type actor struct {
	address     typ.Address
	destination typ.Destination
	privateKey  []byte
	role        uint
}

func main() {
	// ganache-configurations
	// ganache-cli -m "pistol kiwi shrug future ozone ostrich match remove crucial oblige cream critic" --block-time 5 -e 1000
	// STEP 1 - deploy smart contract (nitro adjucator)
	// Ganache || Hardhat

	// STEP 2 - initialize client
	// contract address from already deployed contract
	// rpc url - ganache or hardhat node url
	client, err := NewClient("", "127.0.0.1:8545")
	if err != nil {
		panic(err)
	}

	// STEP 3 - Create brokers
	// There are 2 participants (broker 1 and broker 2) in the system
	// TODO change their addresses and private keys
	var broker1 = actor{
		address:     common.HexToAddress(`0x2EE1ac154435f542ECEc55C5b0367650d8A5343B`),
		destination: typ.AddressToDestination(common.HexToAddress(`0x2EE1ac154435f542ECEc55C5b0367650d8A5343B`)),
		privateKey:  common.Hex2Bytes(`0xb691bc22c5a30f64876c6136553023d522dcdf0744306dccf4f034a465532e27`),
		role:        0,
	}

	var broker2 = actor{
		address:     common.HexToAddress(`0x70765701b79a4e973dAbb4b30A72f5a845f22F9E`),
		destination: typ.AddressToDestination(common.HexToAddress(`0x70765701b79a4e973dAbb4b30A72f5a845f22F9E`)),
		privateKey:  common.Hex2Bytes(`0xb5dc82fc5f4d82b59a38ac963a15eaaedf414f496a037bb4a52310915ac84097`),
		role:        1,
	}

	// STEP 4 - Open a channel
	chainId, err := client.contract.GetChainID(nil)
	if err != nil {
		panic(err)
	}

	var preFundState = state.State{
		ChainId:           chainId,
		Participants:      []typ.Address{broker1.address, broker2.address},
		ChannelNonce:      big.NewInt(0),                                                     // TODO
		AppDefinition:     common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`), // TODO What we should put here?
		ChallengeDuration: big.NewInt(60),
		AppData:           []byte{},
		Outcome: outcome.Exit{
			outcome.SingleAssetExit{
				Asset: typ.Address{}, //TODO
				Allocations: outcome.Allocations{
					outcome.Allocation{
						Destination: broker1.destination,
						Amount:      big.NewInt(3),
					},
					outcome.Allocation{
						Destination: broker2.destination,
						Amount:      big.NewInt(2),
					},
				},
			},
		},
		TurnNum: 0,
		IsFinal: false,
	}

	c, err := channel.New(preFundState, broker1.role)
	if err != nil {
		panic(err)
	}

	// STEP 5 - Sign prefund state
	_, err = c.PreFundState().Sign(broker1.privateKey)
	if err != nil {
		panic(err)
	}

	_, err = c.PreFundState().Sign(broker2.privateKey)
	if err != nil {
		panic(err)
	}

	// STEP 6 - Deposit process
	// client.contract.Deposit(opts *bind.TransactOpts, asset common.Address, channelId [32]byte, expectedHeld *big.Int, amount *big.Int)
	// STEP 7 - Sign post fund state

	_, err = c.PostFundState().Sign(broker1.privateKey)
	if err != nil {
		panic(err)
	}

	_, err = c.PostFundState().Sign(broker2.privateKey)
	if err != nil {
		panic(err)
	}

	// STEP 8 - withdraw

	// try to withdraw from chanel
	// client.contract.Withdraw(opts *bind.TransactOpts)
}
