package chainservice

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	"github.com/statechannels/go-nitro/client/engine/store/safesync"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var depositedTopic = crypto.Keccak256Hash([]byte("Deposited(bytes32,address,uint256,uint256)"))

type eventSource interface {
	SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- ethTypes.Log) (ethereum.Subscription, error)
}

type EthChainService struct {
	ChainServiceBase
	na       *NitroAdjudicator.NitroAdjudicator
	txSigner *bind.TransactOpts
}

// NewEthChainService constructs a chain service that submits transactions to a NitroAdjudicator
// and listens to events from an eventSource
func NewEthChainService(na *NitroAdjudicator.NitroAdjudicator, naAddress common.Address, txSigner *bind.TransactOpts, es eventSource) *EthChainService {
	ecs := EthChainService{ChainServiceBase: newChainServiceBase()}
	ecs.out = safesync.Map[chan Event]{}
	ecs.na = na
	ecs.txSigner = txSigner

	go ecs.listenForLogEvents(na, naAddress, es)

	return &ecs
}

// SendTransaction sends the transaction and blocks until it has been submitted.
func (ecs *EthChainService) SendTransaction(tx protocols.ChainTransaction) {
	txOpt := bind.TransactOpts{
		From:     ecs.txSigner.From,
		Nonce:    ecs.txSigner.Nonce,
		Signer:   ecs.txSigner.Signer,
		GasPrice: big.NewInt(10000000000),
	}
	switch tx.(type) {
	case protocols.DepositTransaction:
		depositTx := tx.(protocols.DepositTransaction)
		for tokenAddress, amount := range depositTx.Deposit {
			txOpt.Value = amount
			holdings, err := ecs.na.Holdings(&bind.CallOpts{}, tokenAddress, depositTx.ChannelId())
			if err != nil {
				panic(err)
			}

			_, err = ecs.na.Deposit(&txOpt, tokenAddress, depositTx.ChannelId(), holdings, amount)

			if err != nil {
				panic(err)
			}
		}
	// TODO handle other transaction types
	default:
		panic("unexpected chain transaction")
	}
}

func (ecs *EthChainService) listenForLogEvents(na *NitroAdjudicator.NitroAdjudicator, naAddress common.Address, es eventSource) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{naAddress},
	}
	logs := make(chan ethTypes.Log)
	sub, err := es.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case chainEvent := <-logs:
			switch chainEvent.Topics[0] {
			case depositedTopic:
				nad, err := na.ParseDeposited(chainEvent)
				if err != nil {
					log.Fatal(err)
				}

				holdings := types.Funds{}
				holdings[nad.Asset] = nad.DestinationHoldings
				// TODO fill out other event fields once the event data structure is settled.
				event := DepositedEvent{
					CommonEvent: CommonEvent{
						channelID: nad.Destination,
						BlockNum:  chainEvent.BlockNumber,
					},
					Holdings: holdings,
				}
				ecs.broadcast(event)
			// TODO introduce the remaining events
			default:
				panic("Unknown chain event")
			}
		}
	}
}
