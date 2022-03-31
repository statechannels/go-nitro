package consensus_channel

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

type actor struct {
	Address    types.Address
	PrivateKey []byte
}

func (a actor) Destination() types.Destination {
	return types.AddressToDestination(a.Address)
}

// actors namespaces the actors exported for test consumption
type actors struct {
	Alice actor
	Bob   actor
	Brian actor
	Irene actor
}

// Actors is the endpoint for tests to consume constructed statechannel
// network participants (public-key secret-key pairs)
var Actors actors = actors{
	Alice: actor{
		common.HexToAddress(`0xAAA6628Ec44A8a742987EF3A114dDFE2D4F7aDCE`),
		common.Hex2Bytes(`2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d`),
	},
	Bob: actor{
		common.HexToAddress(`0xBBB676f9cFF8D242e9eaC39D063848807d3D1D94`),
		common.Hex2Bytes(`0279651921cd800ac560c21ceea27aab0107b67daf436cdd25ce84cad30159b4`),
	},
	Brian: actor{
		common.HexToAddress("0xB2B22ec3889d11f2ddb1A1Db11e80D20EF367c01"),
		common.Hex2Bytes("0aca28ba64679f63d71e671ab4dbb32aaa212d4789988e6ca47da47601c18fe2"),
	},
	Irene: actor{
		common.HexToAddress(`0x111A00868581f73AB42FEEF67D235Ca09ca1E8db`),
		common.Hex2Bytes(`febb3b74b0b52d0976f6571d555f4ac8b91c308dfa25c7b58d1e6a7c3f50c781`),
	},
}

func TestConsensusChannel(t *testing.T) {
	var alice = types.Destination(common.HexToHash("0x0a"))
	var bob = types.Destination(common.HexToHash("0x0b"))

	allocation := func(d types.Destination, a uint64) Balance {
		return Balance{destination: d, amount: big.NewInt(int64(a))}
	}

	guarantee := func(amount uint64, target, left, right types.Destination) Guarantee {
		return Guarantee{
			target: target,
			amount: big.NewInt(int64(amount)),
			left:   left,
			right:  right,
		}
	}

	makeOutcome := func(left, right Balance, guarantees ...Guarantee) LedgerOutcome {
		mappedGuarantees := make(map[types.Destination]Guarantee)
		for _, g := range guarantees {
			mappedGuarantees[g.target] = g
		}
		return LedgerOutcome{left: left, right: right, guarantees: mappedGuarantees}
	}

	add := func(turnNum, amount uint64, vId, left, right types.Destination) Add {
		bigAmount := big.NewInt(int64(amount))
		return Add{
			turnNum: turnNum,
			Guarantee: Guarantee{
				amount: bigAmount,
				target: vId,
				left:   left,
				right:  right,
			},
			LeftDeposit: bigAmount,
		}
	}

	existingChannel := types.Destination{1}
	targetChannel := types.Destination{2}
	aBal := uint64(200)
	bBal := uint64(300)
	vAmount := uint64(5)

	proposal := add(10, vAmount, targetChannel, alice, bob)

	outcome := func() LedgerOutcome {
		return makeOutcome(
			allocation(alice, aBal),
			allocation(bob, bBal),
			guarantee(vAmount, existingChannel, alice, bob),
		)

	}

	fingerprint := func(v Vars) string {
		h, err := v.AsState(state.TestState.FixedPart()).Hash()

		if err != nil {
			panic(err)
		}
		return h.String()
	}

	vars := Vars{TurnNum: 9, Outcome: outcome()}

	f1 := fingerprint(vars)
	clone1 := vars.Outcome.clone()

	if fingerprint(Vars{TurnNum: vars.TurnNum, Outcome: clone1}) != f1 {
		t.Fatal("vars incorrectly cloned")
	}

	mutatedG := clone1.guarantees[existingChannel]
	mutatedG.amount.SetInt64(111)
	if f1 != fingerprint(vars) {
		t.Fatal("vars shares data with clone")
	}

	clone2 := vars.Outcome.clone()
	clone2.left.amount.SetInt64(111)
	if f1 != fingerprint(vars) {
		t.Fatal("vars shares data with clone")
	}

	clone3 := vars.Outcome.clone()
	clone3.right.amount.SetInt64(111)
	if f1 != fingerprint(vars) {
		t.Fatal("vars shares data with clone")
	}

	testApplyingAddProposalToVars := func(t *testing.T) {
		startingTurnNum := uint64(9)
		vars := Vars{TurnNum: startingTurnNum, Outcome: outcome()}

		err := vars.Add(proposal)

		if err != nil {
			t.Fatalf("unable to compute next state: %v", err)
		}

		if vars.TurnNum != startingTurnNum+1 {
			t.Fatalf("incorrect state calculation: %v", err)
		}

		expected := makeOutcome(
			allocation(alice, aBal-vAmount),
			allocation(bob, bBal),
			guarantee(vAmount, existingChannel, alice, bob),
			guarantee(vAmount, targetChannel, alice, bob),
		)

		if diff := cmp.Diff(vars.Outcome, expected, cmp.AllowUnexported(expected, Balance{}, big.Int{}, Guarantee{})); diff != "" {
			t.Fatalf("incorrect outcome: %v", diff)
		}

		// Proposing the same change again should fail
		duplicateProposal := proposal
		duplicateProposal.turnNum += 1
		err = vars.Add(duplicateProposal)

		if !errors.Is(err, ErrDuplicateGuarantee) {
			t.Fatalf("expected error when adding duplicate guarantee: %v", err)
		}

		// Proposing a change that depletes a balance should fail
		vars = Vars{TurnNum: startingTurnNum, Outcome: outcome()}
		largeProposal := proposal
		leftAmount := big.NewInt(0).Set(vars.Outcome.left.amount)
		largeProposal.amount = leftAmount.Add(leftAmount, big.NewInt(1))
		err = vars.Add(largeProposal)
		if !errors.Is(err, ErrInsufficientFunds) {
			t.Fatalf("expected error when adding too large a guarantee: %v", err)
		}
	}

	fp := func() state.FixedPart {
		participants := [2]types.Address{
			Actors.Alice.Address, Actors.Bob.Address,
		}
		return state.FixedPart{
			Participants:      participants[:],
			ChainId:           big.NewInt(0),
			ChannelNonce:      big.NewInt(9001),
			ChallengeDuration: big.NewInt(100),
		}
	}

	initialVars := Vars{Outcome: outcome(), TurnNum: 0}
	aliceSig, _ := initialVars.AsState(fp()).Sign(Actors.Alice.PrivateKey)
	bobsSig, _ := initialVars.AsState(fp()).Sign(Actors.Bob.PrivateKey)
	sigs := [2]state.Signature{aliceSig, bobsSig}

	testConsensusChannelFunctionality := func(t *testing.T) {
		channel, err := newConsensusChannel(fp(), leader, 0, outcome(), sigs)

		if err != nil {
			t.Fatalf("unable to construct a new consensus channel: %v", err)
		}

		_, err = channel.sign(initialVars, Actors.Bob.PrivateKey)
		if err == nil {
			t.Fatalf("channel should check that signer is participant")
		}

		f := fingerprint(channel.current.Vars)

		latest, err := channel.latestProposedVars()
		if err != nil {
			t.Fatalf("latest proposed vars returned err: %v", err)
		}

		latest.Outcome.guarantees[targetChannel] = guarantee(10, targetChannel, alice, bob)
		if f != fingerprint(channel.current.Vars) {
			t.Fatalf("latestProposedVars did not return a copy")
		}

		briansSig, _ := initialVars.AsState(fp()).Sign(Actors.Brian.PrivateKey)
		wrongSigs := [2]state.Signature{sigs[1], briansSig}
		_, err = newConsensusChannel(fp(), leader, 0, outcome(), wrongSigs)
		if err == nil {
			t.Fatalf("channel should check that signers are participants")
		}
	}

	t.Run(`TestApplyingAddProposalToVars`, testApplyingAddProposalToVars)
	t.Run(`TestConsensusChannelFunctionality`, testConsensusChannelFunctionality)
}
