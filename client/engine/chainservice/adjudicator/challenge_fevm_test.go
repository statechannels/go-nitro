package NitroAdjudicator

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"

	"github.com/statechannels/go-nitro/types"
)

// This is the mnemonic for the prefunded accounts on wallaby.
// The first 25 accounts will be prefunded.
const WALLABY_MNEMONIC = "army forest resource shop tray cluster teach cause spice judge link oppose"

// This is the HD path to use when deriving accounts from the mnemonic
const WALLABY_HD_PATH = "m/44'/1'/0'/0"

func TestChallengeFEVM(t *testing.T) {
	t.Skip() // We only want to run this test manually for now.
	hyperspaceChainId := big.NewInt(3141)

	// Setup transacting EOA
	wallet, err := hdwallet.NewFromMnemonic(WALLABY_MNEMONIC)
	if err != nil {
		panic(err)
	}

	// The 0th account is usually used for deployment so we grab the 1st account
	a, err := wallet.Derive(hdwallet.MustParseDerivationPath(fmt.Sprintf("%s/%d", WALLABY_HD_PATH, 1)), false)
	if err != nil {
		panic(err)
	}

	//PK: 0x1688820ffc6a811e09ff17eccec23d8dec4850c3098ffc03ac4aa38dd8f3a994
	// corresponding ETH address is 0x280c53E2C574418D8d6d8d651d4c3323F4b194Be
	// corresponding f4 address (delegated) is t410ffagfhywforay3dlnrvsr2tbtep2ldff6xuxkrjq.
	pk, err := wallet.PrivateKey(a)

	if err != nil {
		panic(err)
	}
	client, err := ethclient.Dial("https://api.hyperspace.node.glif.io/rpc/v0")
	if err != nil {
		t.Fatal(err)
	}

	// When submitting a transaction it's signed against a specific chain id
	// To get the correct signature we need to use the correct chain id that wallaby is expecting
	txSubmitter, err := bind.NewKeyedTransactorWithChainID(pk, hyperspaceChainId)
	if err != nil {
		log.Fatal(err)
	}
	// By setting the GasTipCap we signal this is a type 2 transaction
	// FEVM does NOT support type 1 transactions
	txSubmitter.GasTipCap = big.NewInt(300000)

	// This is the deployed contract on wallaby
	// If wallaby gets reset this will need to be redeployed by running:
	// WALLABY_DEPLOYER_PK="f4d69c36885541f56f4728ddc002a6fa2fcb26c9f608910310a776c83b7fde47" npx hardhat deploy --network hyperspace --deploy-scripts ./hardhat-deploy-fvm --reset
	// The PK corresponds to account 0xE39dce95b1A924E2472E24C20C55eA3559a09251.
	// It should be prefunded after every wallaby reset.
	naAddress := common.HexToAddress("0x4fBeCDA4735eaF21C8ba5BD40Ab97dFa2Ed88E80")
	caAddress := common.HexToAddress("0xC57875E317f67F2bE5D62f5c7C696D2eb7Fe79FE")

	na, err := NewNitroAdjudicator(naAddress, client)
	if err != nil {
		t.Fatal(err)
	}

	balance := new(big.Int)
	balance.SetString("10000000000000000000", 10) // 10 eth in wei

	TestChallengeWithTurnNum := func(t *testing.T, turnNum uint64) {

		var s = state.State{
			Participants: []types.Address{
				Actors.Alice.Address,
				Actors.Bob.Address,
			},
			ChannelNonce:      rand.Uint64(),
			AppDefinition:     caAddress,
			ChallengeDuration: 60,
			AppData:           []byte{},
			Outcome:           outcome.Exit{},
			TurnNum:           turnNum,
			IsFinal:           false,
		}

		// Generate Signatures
		aSig, _ := s.Sign(Actors.Alice.PrivateKey)
		bSig, _ := s.Sign(Actors.Bob.PrivateKey)
		challengerSig, err := SignChallengeMessage(s, Actors.Alice.PrivateKey)

		if err != nil {
			t.Fatal(err)
		}

		// Construct support proof
		candidate := INitroTypesSignedVariablePart{
			ConvertVariablePart(s.VariablePart()),
			[]INitroTypesSignature{ConvertSignature(aSig), ConvertSignature(bSig)},
		}
		proof := make([]INitroTypesSignedVariablePart, 0)

		// Fire off a Challenge tx
		tx, err := na.Challenge(
			txSubmitter,
			INitroTypesFixedPart(ConvertFixedPart(s.FixedPart())),
			proof,
			candidate,
			ConvertSignature(challengerSig),
		)
		if err != nil {
			t.Log(tx)
			t.Fatal(err)
		}

		_, _ = bind.WaitMined(context.Background(), client, tx)

		// Compute challenge time
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			t.Fatal(err)
		}
		header, err := client.HeaderByNumber(context.Background(), receipt.BlockNumber)
		if err != nil {
			t.Fatal(err)
		}

		// Generate expectation
		expectedFinalizesAt := header.Time + uint64(s.ChallengeDuration)
		cId := s.ChannelId()
		expectedOnChainStatus, err := generateStatus(s, expectedFinalizesAt)
		if err != nil {
			t.Fatal(err)
		}

		// Inspect state of chain (call StatusOf)
		statusOnChain, err := na.StatusOf(&bind.CallOpts{}, cId)
		if err != nil {
			t.Fatal(err)
		}

		// Make assertion
		if !bytes.Equal(statusOnChain[:], expectedOnChainStatus) {
			t.Fatalf("Adjudicator not updated as expected, got %v wanted %v", common.Bytes2Hex(statusOnChain[:]), common.Bytes2Hex(expectedOnChainStatus[:]))
		}

	}

	for _, turnNum := range []uint64{0} {
		t.Run("turnNum = "+fmt.Sprint(turnNum), func(t *testing.T) { TestChallengeWithTurnNum(t, turnNum) })
	}
}
