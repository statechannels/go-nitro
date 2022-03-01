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

func NewLedgerManager() *LedgerManager {
	return &LedgerManager{}
}

// HandleRequest accepts a ledger request and updates the ledger channel based on the request.
// It returns a signed state message that can be sent to other participants.
func (l *LedgerManager) HandleRequest(ledger *channel.TwoPartyLedger, request protocols.GuaranteeRequest, secretKey *[]byte) (protocols.SideEffects, error) {

	supported, err := ledger.Channel.LatestSupportedState()
	if err != nil {
		return protocols.SideEffects{}, fmt.Errorf("error finding a supported state: %w", err)
	}
	nextState := supported.Clone()
	nextState.Outcome, err = nextState.Outcome.DivertToGuarantee(request.Left, request.Right, request.LeftAmount, request.RightAmount, request.Destination)
	if err != nil {
		return protocols.SideEffects{}, fmt.Errorf("error diverting to guarantee: %w", err)
	}
	nextState.TurnNum = nextState.TurnNum + 1

	ss := state.NewSignedState(nextState)
	err = ss.Sign(secretKey)
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
		_ = toSign.Sign(&secretKey)
	}
	ledger.Channel.AddSignedState(toSign)
}

// NewTestTwoPartyLedger creates a new two party ledger channel based on the provided allocations. The channel will appear to be fully-funded on chain.
// ONLY FOR TESTING PURPOSES
func NewTestTwoPartyLedger(allocations []outcome.Allocation, myAddress types.Address, nonce *big.Int) (*channel.TwoPartyLedger, error) {

	initialState := state.State{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{},
		ChannelNonce:      nonce,
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
		AppData:           []byte{},
		Outcome: outcome.Exit{outcome.SingleAssetExit{
			Allocations: outcome.Allocations{},
		}},
		TurnNum: 0,
		IsFinal: false,
	}
	foundMyAddress := false
	var myIndex uint
	for i, alloc := range allocations {
		a, err := alloc.Destination.ToAddress()
		if err != nil {
			ntpl := channel.TwoPartyLedger{}
			return &ntpl, fmt.Errorf("could not extract address: %w", err)
		}
		initialState.Participants = append(initialState.Participants, a)
		initialState.Outcome[0].Allocations = append(initialState.Outcome[0].Allocations, alloc)
		if a == myAddress {
			foundMyAddress = true
			myIndex = uint(i)
		}
	}

	if !foundMyAddress {
		panic("Destination corresponding to myAddress not found in outcome")
	}

	ledger, lErr := channel.NewTwoPartyLedger(initialState, myIndex)
	if lErr != nil {
		return ledger, fmt.Errorf("error creating ledger: %w", lErr)
	}
	ledger.OnChainFunding = ledger.PreFundState().Outcome.TotalAllocated()

	return ledger, nil
}
