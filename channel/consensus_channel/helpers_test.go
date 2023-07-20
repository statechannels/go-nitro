package consensus_channel

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/types"
)

const TEST_CHAIN_ID = 1337

var (
	channel1Id    = types.Destination{1}
	targetChannel = types.Destination{2}
)

const (
	aBal    = uint64(200)
	bBal    = uint64(300)
	vAmount = uint64(5)
)

var alice, bob, ivan testactors.Actor = testactors.Alice, testactors.Bob, testactors.Ivan

func fp() state.FixedPart {
	participants := [2]types.Address{
		alice.Address(), bob.Address(),
	}
	return state.FixedPart{
		Participants:      participants[:],
		ChannelNonce:      9001,
		ChallengeDuration: 100,
	}
}

func allocation(d testactors.Actor, a uint64) Balance {
	return Balance{destination: d.Destination(), amount: big.NewInt(int64(a))}
}

func guarantee(amount uint64, target types.Destination, left, right testactors.Actor) Guarantee {
	return Guarantee{
		target: target,
		amount: big.NewInt(int64(amount)),
		left:   left.Destination(),
		right:  right.Destination(),
	}
}

func makeOutcome(leader, follower Balance, guarantees ...Guarantee) LedgerOutcome {
	mappedGuarantees := make(map[types.Destination]Guarantee)
	for _, g := range guarantees {
		mappedGuarantees[g.target] = g
	}
	return LedgerOutcome{leader: leader, follower: follower, guarantees: mappedGuarantees}
}

// ledgerOutcome constructs the LedgerOutcome with items
//   - alice: 200,
//   - bob: 300,
//   - guarantee(target: 1, left: alice, right: bob, amount: 5)
func ledgerOutcome() LedgerOutcome {
	return makeOutcome(
		allocation(alice, aBal),
		allocation(bob, bBal),
		guarantee(vAmount, channel1Id, alice, bob),
	)
}

func add(amount uint64, vId types.Destination, left, right testactors.Actor) Add {
	bigAmount := big.NewInt(int64(amount))
	return Add{
		Guarantee: Guarantee{
			amount: bigAmount,
			target: vId,
			left:   left.Destination(),
			right:  right.Destination(),
		},
		LeftDeposit: bigAmount,
	}
}

func remove(vId types.Destination, leftAmount uint64) Remove {
	return Remove{
		Target:     vId,
		LeftAmount: big.NewInt(int64(leftAmount)),
	}
}

// createSignedProposal generates a signed proposal given the vars, proposal fixed parts and private key
// The vars passed in are NOT mutated!
func createSignedProposal(vars Vars, proposal Proposal, fp state.FixedPart, pk []byte) SignedProposal {
	proposalVars := Vars{TurnNum: vars.TurnNum, Outcome: vars.Outcome.clone()}
	_ = proposalVars.HandleProposal(proposal)

	state := proposalVars.AsState(fp)
	sig, _ := state.Sign(pk)

	signedProposal := SignedProposal{
		Proposal:  proposal,
		Signature: sig,
		TurnNum:   state.TurnNum,
	}

	return signedProposal
}

// fingerprint computes a fingerprint for vars by encoding and returning the hash when provided
// with a consisted FixedPart
func fingerprint(v Vars) string {
	h, err := v.AsState(state.TestState.FixedPart()).Hash()
	if err != nil {
		panic(err)
	}

	return h.String()
}

// equals checks that v is other by comparing fingerprints
func (v *Vars) equals(other Vars) bool {
	return fingerprint(*v) == fingerprint(other)
}
