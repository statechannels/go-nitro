package consensus_channel

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/types"
)

func TestSerde(t *testing.T) {
	// TODO unskip this test when we have solved the issue for persisting big.Ints
	// https://github.com/statechannels/go-nitro/issues/439
	t.Skip()
	someGuarantee := Guarantee{
		amount: *big.NewInt(1),
		left:   testdata.Actors.Alice.Destination(),
		right:  testdata.Actors.Alice.Destination(),
		target: types.Destination{99},
	}
	someOutcome := makeOutcome(
		Balance{testdata.Actors.Alice.Destination(), *big.NewInt(2)},
		Balance{testdata.Actors.Bob.Destination(), *big.NewInt(7)},
		someGuarantee)

	t.Run("Guarantee", func(t *testing.T) {
		got, err := json.Marshal(someGuarantee)
		if err != nil {
			t.Fatal(err)
		}
		// TODO: this expectation is not right and currently has some placeholders  / gaps
		want := `{"Amount":{SOMETHING!},"Target":"0x6300000000000000000000000000000000000000000000000000000000000000","Left":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce","Right":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce"}`
		if string(got) != want {
			t.Fatalf("incorrect json marshalling, expected %v got %v", want, string(got))
		}
	})

	t.Run("LedgerOutcome", func(t *testing.T) {
		got, err := json.Marshal(someOutcome)
		if err != nil {
			t.Fatal(err)
		}
		// TODO: this expectation is not right and currently has some placeholders  / gaps
		want := `{"AssetAddress":"0x0000000000000000000000000000000000000000","Left":{"Destination":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce","Amount":{SOMETHING!}},"Right":{"Destination":"0x000000000000000000000000bbb676f9cff8d242e9eac39d063848807d3d1d94","Amount":{}},"Guarantees":{"0x6300000000000000000000000000000000000000000000000000000000000000":{"Amount":{},"Target":"0x6300000000000000000000000000000000000000000000000000000000000000","Left":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce","Right":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce"}}}`
		if string(got) != want {
			t.Fatalf("incorrect json marshalling, expected %v got %v", want, string(got))
		}
	})

	t.Run("ConsensusChannel", func(t *testing.T) {
		cc := ConsensusChannel{
			myIndex: leader,
			fp:      fp(),
			Id:      types.Destination{1},
			current: SignedVars{
				Vars: Vars{
					TurnNum: 0,
					Outcome: someOutcome,
				},
				Signatures: [2]crypto.Signature{{
					R: common.Hex2Bytes(`704b3afcc6e702102ca1af3f73cf3b37f3007f368c40e8b81ca823a65740a053`),
					S: common.Hex2Bytes(`14040ad4c598dbb055a50430142a13518e1330b79d24eed86fcbdff1a7a95589`),
					V: byte(0),
				}, {
					S: common.Hex2Bytes(`704b3afcc6e702102ca1af3f73cf3b37f3007f368c40e8b81ca823a65740a053`),
					R: common.Hex2Bytes(`14040ad4c598dbb055a50430142a13518e1330b79d24eed86fcbdff1a7a95589`),
					V: byte(0),
				}},
			},
			proposalQueue: []SignedProposal{{
				Signature: crypto.Signature{
					S: common.Hex2Bytes(`704b3afcc6e702102ca1af3f73cf3b37f3007f368c40e8b81ca823a65740a053`),
					R: common.Hex2Bytes(`14040ad4c598dbb055a50430142a13518e1330b79d24eed86fcbdff1a7a95589`),
					V: byte(0),
				},
				Proposal: add(1, 2, types.Destination{3}, types.Destination{4}, types.Destination{5}),
			}},
		}

		got, err := json.Marshal(cc)

		if err != nil {
			t.Fatal(err)
		}

		// TODO: this expectation is not right and currently has some placeholders  / gaps
		want := `{"Id":"0x0100000000000000000000000000000000000000000000000000000000000000","MyIndex":0,"FP":{"ChainId":0,"Participants":["0xaaa6628ec44a8a742987ef3a114ddfe2d4f7adce","0xbbb676f9cff8d242e9eac39d063848807d3d1d94"],"ChannelNonce":9001,"AppDefinition":"0x0000000000000000000000000000000000000000","ChallengeDuration":100},"Current":{"TurnNum":0,"Outcome":{"AssetAddress":"0x0000000000000000000000000000000000000000","Left":{"Destination":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce","Amount":{}},"Right":{"Destination":"0x000000000000000000000000bbb676f9cff8d242e9eac39d063848807d3d1d94","Amount":{}},"Guarantees":{"0x6300000000000000000000000000000000000000000000000000000000000000":{"Amount":{},"Target":"0x6300000000000000000000000000000000000000000000000000000000000000","Left":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce","Right":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce"}}},"Signatures":[{"R":"cEs6/MbnAhAsoa8/c887N/MAfzaMQOi4HKgjpldAoFM=","S":"FAQK1MWY27BVpQQwFCoTUY4TMLedJO7Yb8vf8aepVYk=","V":0},{"R":"FAQK1MWY27BVpQQwFCoTUY4TMLedJO7Yb8vf8aepVYk=","S":"cEs6/MbnAhAsoa8/c887N/MAfzaMQOi4HKgjpldAoFM=","V":0}]},"ProposalQueue":[{"R":"FAQK1MWY27BVpQQwFCoTUY4TMLedJO7Yb8vf8aepVYk=","S":"cEs6/MbnAhAsoa8/c887N/MAfzaMQOi4HKgjpldAoFM=","V":0,"Proposal":{"Amount":{SOMETHING},"Target":"0x0300000000000000000000000000000000000000000000000000000000000000","Left":"0x0400000000000000000000000000000000000000000000000000000000000000","Right":"0x0500000000000000000000000000000000000000000000000000000000000000"}}]}`

		if string(got) != want {

			t.Fatalf("incorrect json marshalling, expected %v got %v", want, string(got))
		}
	})

}
