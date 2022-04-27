// todo: #420 delete this file
package virtualfund

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
