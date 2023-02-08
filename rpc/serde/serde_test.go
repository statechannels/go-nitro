package serde

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols/directfund"
)

func TestMarshalJSON(t *testing.T) {

	req := JsonRpcRequest[directfund.ObjectiveRequest]{
		Jsonrpc: JsonRpcVersion,
		Id:      123,
		Method:  "CreateLedgerChannel",
		Params: directfund.NewObjectiveRequest(
			testactors.Alice.Address(),
			345,
			testdata.Outcomes.Create(
				testactors.Alice.Address(),
				testactors.Bob.Address(),
				9,
				7,
				common.Address{},
			),
			18446744073709551615, // max uint64
			common.Address{123},
		),
	}

	want := `{"jsonrpc":"2.0","id":123,"method":"CreateLedgerChannel","params":{"CounterParty":"0xaaa6628ec44a8a742987ef3a114ddfe2d4f7adce","ChallengeDuration":345,"Outcome":[{"Asset":"0x0000000000000000000000000000000000000000","AssetMetadata":{"AssetType":0,"Metadata":null},"Allocations":[{"Destination":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce","Amount":9,"AllocationType":0,"Metadata":null},{"Destination":"0x000000000000000000000000bbb676f9cff8d242e9eac39d063848807d3d1d94","Amount":7,"AllocationType":0,"Metadata":null}]}],"AppDefinition":"0x7b00000000000000000000000000000000000000","AppData":null,"Nonce":18446744073709551615}}`
	enc, err := json.Marshal(req)

	if err != nil {
		t.Fatal(err)
	}
	if string(enc) != want {
		t.Fatalf("incorrect json marshaling, expected %v got \n%v", want, string(enc))
	}

}

func TestUnmarshalJSON(t *testing.T) {
	t.Skip()
}
