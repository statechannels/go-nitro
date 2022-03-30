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

	someOutcome := makeOutcome(
		Balance{testdata.Actors.Alice.Destination(), *big.NewInt(2)},
		Balance{testdata.Actors.Bob.Destination(), *big.NewInt(7)},
		Guarantee{
			amount: *big.NewInt(1),
			left:   testdata.Actors.Alice.Destination(),
			right:  testdata.Actors.Alice.Destination(),
			target: types.Destination{99},
		})

	t.Run("LedgerOutcome", func(t *testing.T) {
		// got, err := json.MarshalIndent(someOutcome, " ", "    ")
		got, err := someOutcome.MarshalJSON()

		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%+v", someOutcome)
		want := "{something}"

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

		want := "{somethingelse}"

		if string(got) != want {

			t.Fatalf("incorrect json marshalling, expected %v got %v", want, string(got))
		}
	})

}
