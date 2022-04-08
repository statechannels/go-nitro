package virtualfund

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

// Virtual Channel

type actorLedgers struct {
	left  *consensus_channel.ConsensusChannel
	right *consensus_channel.ConsensusChannel
}
type ledgerLookup map[types.Destination]actorLedgers
type testData struct {
	vPreFund  state.State
	vPostFund state.State
	ledgers   ledgerLookup
}

func newTestData() testData {
	var vPreFund = state.State{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{alice.address, p1.address, bob.address}, // A single hop virtual channel
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
		AppData:           []byte{},
		Outcome: outcome.Exit{outcome.SingleAssetExit{
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: alice.destination,
					Amount:      big.NewInt(6),
				},
				outcome.Allocation{
					Destination: bob.destination,
					Amount:      big.NewInt(4),
				},
			},
		}},
		TurnNum: 0,
		IsFinal: false,
	}
	var vPostFund = vPreFund.Clone()
	vPostFund.TurnNum = 1

	ledgers := make(map[types.Destination]actorLedgers)
	ledgers[alice.destination] = actorLedgers{
		right: prepareConsensusChannel(uint(consensus_channel.Leader), alice, p1),
	}
	ledgers[p1.destination] = actorLedgers{
		left:  prepareConsensusChannel(uint(consensus_channel.Follower), alice, p1),
		right: prepareConsensusChannel(uint(consensus_channel.Leader), p1, bob),
	}
	ledgers[bob.destination] = actorLedgers{
		left: prepareConsensusChannel(uint(consensus_channel.Follower), p1, bob),
	}

	return testData{vPreFund, vPostFund, ledgers}
}

func assertNilConnection(t *testing.T, c *Connection) {
	if c != nil {
		t.Fatalf("TestNew: unexpected connection")
	}
}

func assertCorrectConnection(t *testing.T, c *Connection, left, right actor) {
	td := newTestData()
	vPreFund := td.vPreFund

	Id, _ := vPreFund.FixedPart().ChannelId()

	expectedAmount := big.NewInt(0).Set(vPreFund.VariablePart().Outcome[0].TotalAllocated())
	want := consensus_channel.NewGuarantee(expectedAmount, Id, left.destination, right.destination)
	got := c.getExpectedGuarantee()
	if diff := compareGuarantees(want, got); diff != "" {
		t.Fatalf("TestNew: expectedGuarantee mismatch (-want +got):\n%s", diff)
	}
}

func testNew(t *testing.T, a actor) {
	td := newTestData()
	lookup := td.ledgers
	vPreFund := td.vPreFund

	// Assert that a valid set of constructor args does not result in an error
	o, err := constructFromState(
		false,
		vPreFund,
		a.address,
		lookup[a.destination].left,
		lookup[a.destination].right,
	)
	if err != nil {
		t.Fatal(err)
	}

	switch a.role {
	case alice.role:
		assertNilConnection(t, o.ToMyLeft)
		assertCorrectConnection(t, o.ToMyRight, alice, p1)
	case p1.role:
		assertCorrectConnection(t, o.ToMyLeft, alice, p1)
		assertCorrectConnection(t, o.ToMyRight, p1, bob)
	case bob.role:
		assertCorrectConnection(t, o.ToMyLeft, p1, bob)
		assertNilConnection(t, o.ToMyRight)
	}

}

var actors = []actor{alice, p1, bob}

func TestNew(t *testing.T) {
	for _, a := range actors {
		testNew(t, a)
	}
}

func testClone(t *testing.T, my actor) {
	td := newTestData()
	vPreFund := td.vPreFund
	ledgers := td.ledgers

	o, _ := constructFromState(false, vPreFund, my.address, ledgers[my.destination].left, ledgers[my.destination].right)

	clone := o.clone()

	if diff := compareObjectives(o, clone); diff != "" {
		t.Fatalf("Clone: mismatch (-want +got):\n%s", diff)
	}
}

func TestClone(t *testing.T) {
	for _, a := range actors {
		testClone(t, a)
	}
}
