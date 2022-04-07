package virtualfund

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

func TestSingleHopVirtualFundNew(t *testing.T) {
	/////////////////////
	// VIRTUAL CHANNEL //
	/////////////////////

	// Virtual Channel
	vPreFund := state.State{
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
	vPostFund := vPreFund.Clone()
	vPostFund.TurnNum = 1

	TestAs := func(my actor, t *testing.T) {
		prepareConsensusChannels := func(role uint) (*consensus_channel.ConsensusChannel, *consensus_channel.ConsensusChannel) {
			var left *consensus_channel.ConsensusChannel
			var right *consensus_channel.ConsensusChannel

			switch role {
			case 0:
				right = prepareConsensusChannel(uint(consensus_channel.Leader), alice, p1)
			case 1:
				left = prepareConsensusChannel(uint(consensus_channel.Leader), alice, p1)
				right = prepareConsensusChannel(uint(consensus_channel.Follower), p1, bob)
			case 2:
				left = prepareConsensusChannel(uint(consensus_channel.Leader), p1, bob)
			}

			return left, right
		}

		testNew := func(t *testing.T) {
			ledgerChannelToMyLeft, ledgerChannelToMyRight := prepareConsensusChannels(my.role)

			// Assert that a valid set of constructor args does not result in an error
			o, err := constructFromState(false, vPreFund, my.address, ledgerChannelToMyLeft, ledgerChannelToMyRight) // todo: #420 deprecate TwoPartyLedgers
			if err != nil {
				t.Fatal(err)
			}

			var expectedGuaranteeMetadataLeft outcome.GuaranteeMetadata
			var expectedGuaranteeMetadataRight outcome.GuaranteeMetadata
			switch my.role {
			case alice.role:
				{
					expectedGuaranteeMetadataRight = outcome.GuaranteeMetadata{Left: alice.destination, Right: p1.destination}
				}
			case p1.role:
				{
					expectedGuaranteeMetadataLeft = outcome.GuaranteeMetadata{Left: alice.destination, Right: p1.destination}
					expectedGuaranteeMetadataRight = outcome.GuaranteeMetadata{Left: p1.destination, Right: bob.destination}
				}
			case bob.role:
				{
					expectedGuaranteeMetadataLeft = outcome.GuaranteeMetadata{Left: p1.destination, Right: bob.destination}
				}
			}
			amount := big.NewInt(0).Set(vPreFund.VariablePart().Outcome[0].TotalAllocated())
			if (expectedGuaranteeMetadataLeft != outcome.GuaranteeMetadata{}) {
				gotLeft := o.ToMyLeft.getExpectedGuarantee()

				left := expectedGuaranteeMetadataLeft.Left
				right := expectedGuaranteeMetadataLeft.Left

				wantLeft := consensus_channel.NewGuarantee(amount, o.V.Id, left, right)
				if diff := compareGuarantees(wantLeft, gotLeft); diff != "" {
					t.Fatalf("TestNew: expectedGuarantee mismatch (-want +got):\n%s", diff)
				}
			}
			if (expectedGuaranteeMetadataRight != outcome.GuaranteeMetadata{}) {
				gotRight := o.ToMyRight.getExpectedGuarantee()
				left := expectedGuaranteeMetadataRight.Left
				right := expectedGuaranteeMetadataRight.Left

				wantRight := consensus_channel.NewGuarantee(amount, o.V.Id, left, right)
				if diff := compareGuarantees(wantRight, gotRight); diff != "" {
					t.Fatalf("TestNew: expectedGuarantee mismatch (-want +got):\n%s", diff)
				}
			}
		}

		t.Run(`New`, testNew)
	}

	t.Run(`AsAlice`, func(t *testing.T) { TestAs(alice, t) })
	t.Run(`AsBob`, func(t *testing.T) { TestAs(bob, t) })
	t.Run(`AsP1`, func(t *testing.T) { TestAs(p1, t) })
}
