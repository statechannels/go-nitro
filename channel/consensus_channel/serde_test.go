package consensus_channel

import (
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

func TestSerde(t *testing.T) {
	someGuarantee := Guarantee{
		amount: big.NewInt(1),
		left:   alice.Destination(),
		right:  alice.Destination(),
		target: types.Destination{99},
	}
	someGuaranteeJSON := `{"Amount":1,"Target":"0x6300000000000000000000000000000000000000000000000000000000000000","Left":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce","Right":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce"}`

	someAdd := Add{
		Guarantee:   someGuarantee,
		LeftDeposit: big.NewInt(77),
	}
	someAddJSON := `{"Guarantee":{"Amount":1,"Target":"0x6300000000000000000000000000000000000000000000000000000000000000","Left":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce","Right":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce"},"LeftDeposit":77}`

	someOutcome := makeOutcome(
		Balance{alice.Destination(), big.NewInt(2)},
		Balance{bob.Destination(), big.NewInt(7)},
		someGuarantee)
	someOutcomeJSON := `{"AssetAddress":"0x0000000000000000000000000000000000000000","Leader":{"Destination":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce","Amount":2},"Follower":{"Destination":"0x000000000000000000000000bbb676f9cff8d242e9eac39d063848807d3d1d94","Amount":7},"Guarantees":{"0x6300000000000000000000000000000000000000000000000000000000000000":{"Amount":1,"Target":"0x6300000000000000000000000000000000000000000000000000000000000000","Left":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce","Right":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce"}}}`

	someConsensusChannel := ConsensusChannel{
		MyIndex: Leader,
		fp:      fp(),
		Id:      types.Destination{1},
		OnChainFunding: types.Funds{
			common.HexToAddress("0x00"): big.NewInt(9),
		},
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
		proposalQueue: []SignedProposal{
			{
				Signature: crypto.Signature{
					S: common.Hex2Bytes(`704b3afcc6e702102ca1af3f73cf3b37f3007f368c40e8b81ca823a65740a053`),
					R: common.Hex2Bytes(`14040ad4c598dbb055a50430142a13518e1330b79d24eed86fcbdff1a7a95589`),
					V: byte(0),
				},
				Proposal: Proposal{ToAdd: add(1, types.Destination{3}, alice, bob)},
			},
			{
				Signature: crypto.Signature{
					S: common.Hex2Bytes(`704b3afcc6e702102ca1af3f73cf3b37f3007f368c40e8b81ca823a65740a053`),
					R: common.Hex2Bytes(`14040ad4c598dbb055a50430142a13518e1330b79d24eed86fcbdff1a7a95589`),
					V: byte(0),
				},
				Proposal: Proposal{ToRemove: remove(types.Destination{3}, 1)},
			},
		},
	}
	someConsensusChannelJSON := `{"Id":"0x0100000000000000000000000000000000000000000000000000000000000000","OnChainFunding":{"0x0000000000000000000000000000000000000000":9},"MyIndex":0,"FP":{"Participants":["0xaaa6628ec44a8a742987ef3a114ddfe2d4f7adce","0xbbb676f9cff8d242e9eac39d063848807d3d1d94"],"ChannelNonce":9001,"AppDefinition":"0x0000000000000000000000000000000000000000","ChallengeDuration":100},"Current":{"TurnNum":0,"Outcome":{"AssetAddress":"0x0000000000000000000000000000000000000000","Leader":{"Destination":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce","Amount":2},"Follower":{"Destination":"0x000000000000000000000000bbb676f9cff8d242e9eac39d063848807d3d1d94","Amount":7},"Guarantees":{"0x6300000000000000000000000000000000000000000000000000000000000000":{"Amount":1,"Target":"0x6300000000000000000000000000000000000000000000000000000000000000","Left":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce","Right":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce"}}},"Signatures":["0x704b3afcc6e702102ca1af3f73cf3b37f3007f368c40e8b81ca823a65740a05314040ad4c598dbb055a50430142a13518e1330b79d24eed86fcbdff1a7a9558900","0x14040ad4c598dbb055a50430142a13518e1330b79d24eed86fcbdff1a7a95589704b3afcc6e702102ca1af3f73cf3b37f3007f368c40e8b81ca823a65740a05300"]},"ProposalQueue":[{"Signature":"0x14040ad4c598dbb055a50430142a13518e1330b79d24eed86fcbdff1a7a95589704b3afcc6e702102ca1af3f73cf3b37f3007f368c40e8b81ca823a65740a05300","Proposal":{"LedgerID":"0x0000000000000000000000000000000000000000000000000000000000000000","ToAdd":{"Guarantee":{"Amount":1,"Target":"0x0300000000000000000000000000000000000000000000000000000000000000","Left":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce","Right":"0x000000000000000000000000bbb676f9cff8d242e9eac39d063848807d3d1d94"},"LeftDeposit":1},"ToRemove":{"Target":"0x0000000000000000000000000000000000000000000000000000000000000000","LeftAmount":null}},"TurnNum":0},{"Signature":"0x14040ad4c598dbb055a50430142a13518e1330b79d24eed86fcbdff1a7a95589704b3afcc6e702102ca1af3f73cf3b37f3007f368c40e8b81ca823a65740a05300","Proposal":{"LedgerID":"0x0000000000000000000000000000000000000000000000000000000000000000","ToAdd":{"Guarantee":{"Amount":null,"Target":"0x0000000000000000000000000000000000000000000000000000000000000000","Left":"0x0000000000000000000000000000000000000000000000000000000000000000","Right":"0x0000000000000000000000000000000000000000000000000000000000000000"},"LeftDeposit":null},"ToRemove":{"Target":"0x0300000000000000000000000000000000000000000000000000000000000000","LeftAmount":1}},"TurnNum":0}]}`

	type testCase struct {
		name string
		rich interface{}
		json string
	}

	testCases := []testCase{
		{
			"Guarantee",
			someGuarantee,
			someGuaranteeJSON,
		},
		{
			"Add",
			someAdd,
			someAddJSON,
		},
		{
			"LedgerOutcome",
			someOutcome,
			someOutcomeJSON,
		},
		{
			"ConsensusChannel",
			someConsensusChannel,
			someConsensusChannelJSON,
		},
	}

	for _, c := range testCases {

		t.Run("Marshaling "+c.name, func(t *testing.T) {
			got, err := json.Marshal(c.rich)
			if err != nil {
				t.Fatal(err)
			}
			want := c.json
			if string(got) != want {
				t.Fatalf("incorrect json marshaling, expected %v got \n%v", want, string(got))
			}
		})

		t.Run("Unmarshaling "+c.name, func(t *testing.T) {
			want := c.rich
			var got interface{}
			var err error
			switch c.rich.(type) {
			case Guarantee:
				g := Guarantee{}
				err = json.Unmarshal([]byte(c.json), &g)
				got = g
			case LedgerOutcome:
				lo := LedgerOutcome{}
				err = json.Unmarshal([]byte(c.json), &lo)
				got = lo
			case ConsensusChannel:
				cc := ConsensusChannel{}
				err = json.Unmarshal([]byte(c.json), &cc)
				got = cc
			case Add:
				a := Add{}
				err = json.Unmarshal([]byte(c.json), &a)
				got = a
			default:
				panic("unimplemented")
			}
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("incorrect json unmarshaling, expected \n%+v got \n%+v", want, got)
			}
		})
	}
}
