package Create2Deployer

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
)

func TestCreate2(t *testing.T) {
	sb, _, txOpts, err := chainservice.SetupSimulatedBackend(2)
	if err != nil {
		t.Fatal(err)
	}

	// Deploy Create2Deployer contract
	_, _, deployer, err := DeployCreate2Deployer(txOpts[0], sb)
	if err != nil {
		t.Fatal(err)
	}
	sb.Commit()

	hexBytecode, err := hex.DecodeString(NitroAdjudicator.NitroAdjudicatorMetaData.Bin[2:])
	if err != nil {
		t.Fatal(err)
	}
	_, err = deployer.Deploy(txOpts[0], big.NewInt(0), [32]byte{}, hexBytecode)
	if err != nil {
		t.Fatal(err)
	}
	sb.Commit()

	_, err = deployer.ComputeAddress(&bind.CallOpts{}, [32]byte{}, crypto.Keccak256Hash(hexBytecode))
	if err != nil {
		t.Fatal(err)
	}
}
