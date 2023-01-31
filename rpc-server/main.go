package main

import (
	"crypto/ecdsa"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/websocket"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
)

var cs chainservice.ChainService
var pk *ecdsa.PrivateKey
var upgrader = websocket.Upgrader{}

func init() {

}

func main() {
	port := os.Args[1]

	fileServer := http.FileServer(http.Dir("./rpc-server/static"))
	http.Handle("/", fileServer)

	// setupChainService()
	// c := client.New(
	// 	p2pms.NewMessageService("127.0.0.1", 2828, crypto.FromECDSA(pk)),
	// 	cs,
	// 	store.NewMemStore(crypto.FromECDSA(pk)),
	// 	io.Discard,
	// 	&engine.PermissivePolicy{},
	// 	nil,
	// )

	// GetAddressHandler := func(w http.ResponseWriter, req *http.Request) {
	// 	fmt.Fprint(w, c.GetMyAddress().String())
	// }

	// http.HandleFunc("/address", GetAddressHandler)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade failed: ", err)
			return
		}
		defer conn.Close()

		// Continuosly read and write message
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read failed:", err)
				break
			}
			fmt.Print(string(message))
			err = conn.WriteMessage(mt, message)
			if err != nil {
				log.Println("write failed:", err)
				break
			}
		}
	})
	_ = http.ListenAndServe(":"+port, nil)
}

func setupChainService() {
	// This is the mnemonic for the prefunded accounts on wallaby.
	// The first 25 accounts will be prefunded.
	const WALLABY_MNEMONIC = "army forest resource shop tray cluster teach cause spice judge link oppose"

	// This is the HD path to use when deriving accounts from the mnemonic
	const WALLABY_HD_PATH = "m/44'/1'/0'/0"

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
	pk, err = wallet.PrivateKey(a)
	if err != nil {
		panic(err)
	}
	chain, err := ethclient.Dial("https://api.hyperspace.node.glif.io/rpc/v0")

	if err != nil {
		panic(err)
	}
	naAddress := common.HexToAddress("0x4fBeCDA4735eaF21C8ba5BD40Ab97dFa2Ed88E80")
	caAddress := common.HexToAddress("0xC57875E317f67F2bE5D62f5c7C696D2eb7Fe79FE")
	vpaAddress := common.HexToAddress("0xc1AcE8075ee548AA2284b61C5eD8f1a69c4cE756")
	na, err := NitroAdjudicator.NewNitroAdjudicator(naAddress, chain)
	if err != nil {
		panic(err)
	}
	hyperspaceChainId := big.NewInt(3141)
	txSubmitter, err := bind.NewKeyedTransactorWithChainID(pk, hyperspaceChainId)
	if err != nil {
		log.Fatal(err)
	}

	cs, err = chainservice.NewEthChainService(
		chain, na, naAddress, caAddress, vpaAddress, txSubmitter, io.Discard,
	)

	if err != nil {
		panic(err)
	}
}
