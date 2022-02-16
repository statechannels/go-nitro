package ledger

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type LedgerManager struct {
}

func NewLedgerManager() LedgerManager {
	return LedgerManager{}
}

// HandleRequest accepts a ledger request and updates the ledger channel based on the request.
// It returns a signed state message that can be sent to other participants.
func (l *LedgerManager) HandleRequest(ledger *channel.TwoPartyLedger, request protocols.LedgerRequest, secretKey *[]byte) (protocols.SideEffects, error) {

	guarantee, _ := outcome.GuaranteeMetadata{
		Left:  request.Left,
		Right: request.Right,
	}.Encode()

	supported, err := ledger.Channel.LatestSupportedState()
	if err != nil {
		return protocols.SideEffects{}, fmt.Errorf("error finding a supported state: %w", err)
	}

	asset := types.Address{} // todo: loop over request.amount's assets
	nextState := supported.Clone()

	// Calculate the amounts
	leftAmount := big.NewInt(0).Sub(nextState.Outcome.TotalAllocatedFor(request.Left)[asset], request.LeftAmount[asset])
	rightAmount := big.NewInt(0).Sub(nextState.Outcome.TotalAllocatedFor(request.Right)[asset], request.RightAmount[asset])
	total := big.NewInt(0).Add(request.LeftAmount[asset], request.RightAmount[asset])
	if types.Lt(leftAmount, big.NewInt(0)) {
		return protocols.SideEffects{}, fmt.Errorf("Allocation for %x cannot afford the amount %d", request.Left, request.LeftAmount[asset])
	}
	if types.Lt(rightAmount, big.NewInt(0)) {
		return protocols.SideEffects{}, fmt.Errorf("Allocation for %x cannot afford the amount %d", request.Right, request.RightAmount[asset])
	}

	nextState.Outcome = outcome.Exit{outcome.SingleAssetExit{
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: request.Left,
				Amount:      leftAmount,
			},
			outcome.Allocation{
				Destination: request.Right,
				Amount:      rightAmount,
			},
			outcome.Allocation{
				Destination:    request.Destination,
				Amount:         total,
				AllocationType: outcome.GuaranteeAllocationType,
				Metadata:       guarantee,
			},
		},
	}}

	nextState.TurnNum = nextState.TurnNum + 1

	ss := state.NewSignedState(nextState)
	err = ss.SignAndAdd(secretKey)
	if err != nil {
		return protocols.SideEffects{}, fmt.Errorf("Could not sign state: %w", err)
	}
	if ok := ledger.Channel.AddSignedState(ss); !ok {
		return protocols.SideEffects{}, errors.New("Could not add signed state to channel")
	}

	messages := protocols.CreateSignedStateMessages(request.ObjectiveId, ss, ledger.MyIndex)
	return protocols.SideEffects{MessagesToSend: messages}, nil

}

// SignPreAndPostFundingStates is a test utility function which applies signatures from
// multiple participants to pre and post fund states
func SignPreAndPostFundingStates(ledger *channel.TwoPartyLedger, secretKeys []*[]byte) {
	for _, sk := range secretKeys {
		_, _ = ledger.SignAndAddPrefund(sk)
		_, _ = ledger.SignAndAddPostfund(sk)
	}
}

// Signlatest is a test utility function which applies signatures from
// multiple participants to the latest recorded state
func SignLatest(ledger *channel.TwoPartyLedger, secretKeys [][]byte) {

	// Find the largest turn num and therefore the latest state
	turnNum := uint64(0)
	for t := range ledger.SignedStateForTurnNum {
		if t > turnNum {
			turnNum = t
		}
	}
	// Sign it
	toSign := ledger.SignedStateForTurnNum[turnNum]
	for _, secretKey := range secretKeys {
		_ = toSign.SignAndAdd(&secretKey)
	}
	ledger.Channel.AddSignedState(toSign)
}

// CreateTestLedger creates a new  two party ledger channel based on the provided left and right outcomes.
func CreateTestLedger(left outcome.Allocation, right outcome.Allocation, secretKey *[]byte, myIndex uint, nonce *big.Int) (*channel.TwoPartyLedger, error) {

	leftAddress, _ := left.Destination.ToAddress()
	rightAddress, _ := right.Destination.ToAddress()
	initialState := state.State{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{leftAddress, rightAddress},
		ChannelNonce:      nonce,
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
		AppData:           []byte{},
		Outcome: outcome.Exit{outcome.SingleAssetExit{
			Allocations: outcome.Allocations{left, right},
		}},
		TurnNum: 0,
		IsFinal: false,
	}

	ledger, lErr := channel.NewTwoPartyLedger(initialState, myIndex)
	if lErr != nil {
		return ledger, fmt.Errorf("error creating ledger: %w", lErr)
	}

	return ledger, nil
}
