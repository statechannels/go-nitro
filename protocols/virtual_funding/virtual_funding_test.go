package virtual_funding

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

type Player struct {
	address     types.Address
	destination types.Bytes32
}

var (
	alice  Player = Player{common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`), common.HexToHash("0x00000000000000000000000096f7123E3A80C9813eF50213ADEd0e4511CB820f")}
	bob           = Player{common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`), common.HexToHash("0x00000000000000000000000096f7123E3A80C9813eF50213ADEd0e4511CB820f")}
	ingrid        = Player{common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`), common.HexToHash("0x00000000000000000000000096f7123E3A80C9813eF50213ADEd0e4511CB820f")}
)

func TestNew(t *testing.T) {

	params := ConstructorParams{
		chainId:              big.NewInt(0),
		appDefinition:        common.Address{0},
		appData:              []byte{},
		outcome:              outcome.Exit{},
		myExitDestination:    alice.destination,
		theirExitDestination: bob.destination,
		mySigningAddress:     alice.address,
		theirSigningAddress:  bob.address,
		hubSigningAddress:    ingrid.address,
		challengeDuration:    big.NewInt(5),
	}

	got, err := New(params)
	if err != nil {
		t.Error(err)
	}

	want := Objective{
		status: `Approved`,
	}

	if got.status != want.status {
		t.Error(`Objective initialized with incorrect status`)
	}
	// TODO check all of the other fields
}

// TODO testUpdate
// TODO testCrank
