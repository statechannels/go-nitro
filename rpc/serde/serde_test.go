package serde

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols/directfund"
)

var someRequest JsonRpcRequest[directfund.ObjectiveRequest] = JsonRpcRequest[directfund.ObjectiveRequest]{
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

var someRequestJSONString = `{"jsonrpc":"2.0","id":123,"method":"CreateLedgerChannel","params":{"CounterParty":"0xaaa6628ec44a8a742987ef3a114ddfe2d4f7adce","ChallengeDuration":345,"Outcome":[{"Asset":"0x0000000000000000000000000000000000000000","AssetMetadata":{"AssetType":0,"Metadata":null},"Allocations":[{"Destination":"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce","Amount":9,"AllocationType":0,"Metadata":null},{"Destination":"0x000000000000000000000000bbb676f9cff8d242e9eac39d063848807d3d1d94","Amount":7,"AllocationType":0,"Metadata":null}]}],"AppDefinition":"0x7b00000000000000000000000000000000000000","AppData":null,"Nonce":18446744073709551615}}`

func TestMarshalJSON(t *testing.T) {
	enc, err := json.Marshal(someRequest)
	if err != nil {
		t.Fatal(err)
	}
	if string(enc) != someRequestJSONString {
		t.Fatalf("incorrect json marshaling, expected %v got \n%v", someRequestJSONString, string(enc))
	}
}

func TestUnmarshalJSON(t *testing.T) {
	got := JsonRpcRequest[directfund.ObjectiveRequest]{} // This test assumes we have a way to know the "type" (method) of the request.
	err := json.Unmarshal([]byte(someRequestJSONString), &got)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(someRequest, got, cmpopts.IgnoreUnexported(directfund.ObjectiveRequest{})); diff != "" {
		t.Fatalf("TestUnmarshalJSON: mismatch (-want +got):\n%s", diff)
	}
}
