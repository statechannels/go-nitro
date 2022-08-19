package messageservice

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func TestSummarizeMessage(t *testing.T) {
	msg1 := protocols.CreateSignedProposalMessage(testactors.Alice.Address(), testactors.Alice.Address(), consensus_channel.SignedProposal{
		Proposal: consensus_channel.Proposal{LedgerID: types.Destination{3}, ToAdd: consensus_channel.Add{
			Guarantee:   consensus_channel.NewGuarantee(big.NewInt(4), types.Destination{9}, types.Destination{8}, types.Destination{7}),
			LeftDeposit: big.NewInt(3),
		}},
		TurnNum: 2,
	})

	got1 := summarizeMessageSend(msg1)
	want1 := "752B:propose 0x0300000000000000000000000000000000000000000000000000000000000000 funds 0x0900000000000000000000000000000000000000000000000000000000000000"

	if got1 != want1 {
		t.Fatalf("wrong message summary: got %s, wanted %s", got1, want1)
	}

	msg2 := protocols.CreateSignedProposalMessage(testactors.Alice.Address(), testactors.Alice.Address(), consensus_channel.SignedProposal{
		Proposal: consensus_channel.Proposal{LedgerID: types.Destination{3}, ToRemove: consensus_channel.Remove{
			Target:     types.Destination{7},
			LeftAmount: big.NewInt(2),
		}},
		TurnNum: 2,
	})

	got2 := summarizeMessageSend(msg2)
	want2 := "757B:propose 0x0300000000000000000000000000000000000000000000000000000000000000 defunds 0x0700000000000000000000000000000000000000000000000000000000000000"

	if got2 != want2 {
		t.Fatalf("wrong message summary: got %s, wanted %s", got2, want2)
	}

}
