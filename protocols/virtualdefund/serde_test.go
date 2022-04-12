package virtualdefund

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/internal/testhelpers"
)

func TestSerde(t *testing.T) {
	data := generateTestData()

	objective := newObjective(true, data.vFixed, data.initialOutcome, big.NewInt(int64(data.paid)), 0)
	objectiveJson := `{\"Status\":1,\"VFixed\":{\"ChainId\":9001,\"Participants\":[\"0xaaa6628ec44a8a742987ef3a114ddfe2d4f7adce\",\"0x111a00868581f73ab42feef67d235ca09ca1e8db\",\"0xbbb676f9cff8d242e9eac39d063848807d3d1d94\"],\"ChannelNonce\":0,\"AppDefinition\":\"0x0000000000000000000000000000000000000000\",\"ChallengeDuration\":45},\"InitialOutcome\":{\"Asset\":\"0x0000000000000000000000000000000000000000\",\"Metadata\":null,\"Allocations\":[{\"Destination\":\"0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce\",\"Amount\":7,\"AllocationType\":0,\"Metadata\":null},{\"Destination\":\"0x000000000000000000000000bbb676f9cff8d242e9eac39d063848807d3d1d94\",\"Amount\":2,\"AllocationType\":0,\"Metadata\":null}]},\"Signatures\":[{\"R\":null,\"S\":null,\"V\":0},{\"R\":null,\"S\":null,\"V\":0},{\"R\":null,\"S\":null,\"V\":0}],\"PaidToBob\":1,\"ToMyLeft\":\"bnVsbA==\",\"ToMyRight\":\"bnVsbA==\",\"MyRole\":0}`

	gotJson, err := json.Marshal(objective)
	if err != nil {
		t.Fatal(err)
	}

	testhelpers.Equals(t, objectiveJson, string(gotJson))

	var gotObjective Objective
	err = json.Unmarshal([]byte(objectiveJson), &gotObjective)
	if err != nil {
		t.Fatal(err)
	}

	testhelpers.Equals(t, objective, gotObjective)
}
