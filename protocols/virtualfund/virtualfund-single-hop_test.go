// todo: #420 delete this file
package virtualfund

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/testactors"

	"github.com/statechannels/go-nitro/types"
)

// signPreAndPostFundingStates is a test utility function which applies signatures from
// multiple participants to pre and post fund states
// func signPreAndPostFundingStates(ledger *channel.TwoPartyLedger, secretKeys []*[]byte) {
// 	for _, sk := range secretKeys {
// 		_, _ = ledger.SignAndAddPrefund(sk)
// 		_, _ = ledger.SignAndAddPostfund(sk)
// 	}
// }

// signLatest is a test utility function which applies signatures from
// multiple participants to the latest recorded state
// func signLatest(ledger *consensus_channel.ConsensusChannel, secretKeys [][]byte) {

// Find the largest turn num and therefore the latest state
// turnNum := uint64(0)
// for t := range ledger.SignedStateForTurnNum {
// 	if t > turnNum {
// 		turnNum = t
// 	}
// }
// // Sign it
// toSign := ledger.SignedStateForTurnNum[turnNum]
// for _, secretKey := range secretKeys {
// 	_ = toSign.Sign(&secretKey)
// }
// ledger.Channel.AddSignedState(toSign)
// }

// addLedgerProposal calculates the ledger proposal state, signs it and adds it to the ledger.
// func addLedgerProposal(
// 	ledger *channel.TwoPartyLedger,
// 	left types.Destination,
// 	right types.Destination,
// 	guaranteeDestination types.Destination,
// 	secretKey *[]byte,
// ) {

// 	supported, _ := ledger.LatestSupportedState()
// 	nextState := constructLedgerProposal(supported, left, right, guaranteeDestination)
// 	_, _ = ledger.SignAndAddState(nextState, secretKey)
// }

// constructLedgerProposal returns a new ledger state with an updated outcome that includes the proposal
// func constructLedgerProposal(
// 	supported state.State,
// 	left types.Destination,
// 	right types.Destination,
// 	guaranteeDestination types.Destination,
// ) state.State {
// 	leftAmount := types.Funds{types.Address{}: big.NewInt(6)}
// 	rightAmount := types.Funds{types.Address{}: big.NewInt(4)}
// 	nextState := supported.Clone()

// 	nextState.TurnNum = nextState.TurnNum + 1
// 	nextState.Outcome, _ = nextState.Outcome.DivertToGuarantee(left, right, leftAmount, rightAmount, guaranteeDestination)
// 	return nextState
// }

func TestSingleHopVirtualFund(t *testing.T) {
	t.Skip()

	/////////////////////
	// VIRTUAL CHANNEL //
	/////////////////////

	// Virtual Channel
	vPreFund := state.State{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{alice.Address, p1.Address, bob.Address}, // A single hop virtual channel
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
		AppData:           []byte{},
		Outcome: outcome.Exit{outcome.SingleAssetExit{
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: alice.Destination(),
					Amount:      big.NewInt(6),
				},
				outcome.Allocation{
					Destination: bob.Destination(),
					Amount:      big.NewInt(4),
				},
			},
		}},
		TurnNum: 0,
		IsFinal: false,
	}
	vPostFund := vPreFund.Clone()
	vPostFund.TurnNum = 1

	TestAs := func(my testactors.Actor, t *testing.T) {

	}

	t.Run(`AsAlice`, func(t *testing.T) { TestAs(alice, t) })
	t.Run(`AsBob`, func(t *testing.T) { TestAs(bob, t) })
	t.Run(`AsP1`, func(t *testing.T) { TestAs(p1, t) })
}
