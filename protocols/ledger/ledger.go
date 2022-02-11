package ledger

import (
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// DirectFundObjective is a cache of data computed by reading from the store. It stores (potentially) infinite data
type LedgerRequestHandler interface {
	HandleRequest(request protocols.LedgerRequest, secretKey *[]byte) protocols.SideEffects
}

type LedgerCranker struct {
	ledgers map[types.Destination]*channel.TwoPartyLedger
	nonce   *big.Int
}

func NewLedgerCranker() LedgerCranker {
	return LedgerCranker{
		ledgers: make(map[types.Destination]*channel.TwoPartyLedger),
		nonce:   big.NewInt(0),
	}
}
func (l *LedgerCranker) Update(ledger *channel.TwoPartyLedger) {
	l.ledgers[ledger.Id] = ledger
}
func (l *LedgerCranker) CreateLedger(left outcome.Allocation, right outcome.Allocation, secretKey *[]byte, myIndex uint) *channel.TwoPartyLedger {

	leftAddress, _ := left.Destination.ToAddress()
	rightAddress, _ := right.Destination.ToAddress()
	initialState := state.State{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{leftAddress, rightAddress},
		ChannelNonce:      l.nonce,
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
		AppData:           []byte{},
		Outcome: outcome.Exit{outcome.SingleAssetExit{
			Allocations: outcome.Allocations{left, right},
		}},
		TurnNum: 0, // Start at post fund
		IsFinal: false,
	}

	ledger, lErr := channel.NewTwoPartyLedger(initialState, myIndex)
	if lErr != nil {
		panic(lErr)
	}

	l.ledgers[ledger.Id] = ledger
	return ledger
}

func (l *LedgerCranker) CompleteFunding(ledgerId types.Destination, secretKeys []*[]byte) {
	ledger, ok := l.ledgers[ledgerId]
	if !ok {
		panic(fmt.Sprintf("Ledger %s not found", ledgerId))
	}
	for _, sk := range secretKeys {
		_, _ = ledger.SignAndAddPrefund(sk)

	}
	for _, sk := range secretKeys {
		_, _ = ledger.Channel.SignAndAddPostfund(sk)

	}
	l.ledgers[ledgerId] = ledger

}

func (l *LedgerCranker) Sign(ledgerId types.Destination, turnNum uint64, secretKeys [][]byte) {
	ledger, ok := l.ledgers[ledgerId]
	if !ok {
		panic(fmt.Sprintf("Ledger %s not found", ledgerId))
	}
	toSign := ledger.SignedStateForTurnNum[turnNum]
	for _, secretKey := range secretKeys {
		_ = toSign.SignAndAdd(&secretKey)
	}
	ledger.Channel.AddSignedState(toSign)
}

func (l *LedgerCranker) HandleRequest(request protocols.LedgerRequest, oId protocols.ObjectiveId, secretKey *[]byte) protocols.SideEffects {
	ledger := l.ledgers[request.LedgerId]
	guarantee, _ := outcome.GuaranteeMetadata{
		Left:  request.Left,
		Right: request.Right,
	}.Encode()
	supported, err := ledger.Channel.LatestSupportedState()
	if err != nil {
		panic(err)
	}
	nextState := supported.Clone()
	nextState.Outcome = outcome.Exit{outcome.SingleAssetExit{
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: request.Left,
				Amount:      big.NewInt(0),
			},
			outcome.Allocation{
				Destination: request.Right,
				Amount:      big.NewInt(0),
			},
			outcome.Allocation{
				Destination:    request.Destination,
				Amount:         request.Amount[types.Address{}],
				AllocationType: outcome.GuaranteeAllocationType,
				Metadata:       guarantee,
			},
		},
	}}
	nextState.TurnNum = nextState.TurnNum + 1
	ss := state.NewSignedState(nextState)
	err = ss.SignAndAdd(secretKey)
	if err != nil {
		panic(err)
	}
	if ok := ledger.Channel.AddSignedState(ss); !ok {
		panic("could not add state")
	}

	messages := protocols.CreateSignedStateMessages(oId, ss, ledger.MyIndex)
	return protocols.SideEffects{MessagesToSend: messages}

}
