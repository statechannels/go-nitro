package main

import (
	"crypto/ecdsa"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/websocket"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/types"
)

var cs chainservice.ChainService
var pk *ecdsa.PrivateKey
var upgrader = websocket.Upgrader{}

type asyncLogger struct {
	r *io.PipeReader
	w *io.PipeWriter
}

func newAsyncLogger() asyncLogger {
	r, w := io.Pipe()
	return asyncLogger{r, w}
}
func (a asyncLogger) Write(p []byte) (n int, err error) {
	go a.w.Write(p)
	return 0, nil
}
func (a asyncLogger) Read(p []byte) (n int, err error) {
	return a.r.Read(p)
}

func init() {

}

func main() {
	rpcPort := os.Args[1]
	msgPort, _ := strconv.Atoi(rpcPort)
	msgPort++

	fileServer := http.FileServer(http.Dir("./rpc-server/static"))
	http.Handle("/", fileServer)

	setupChainService()

	a := newAsyncLogger()

	c := client.New(
		p2pms.NewMessageService("127.0.0.1", msgPort, crypto.FromECDSA(pk)),
		cs,
		store.NewMemStore(crypto.FromECDSA(pk)),
		a,
		&engine.PermissivePolicy{},
		nil,
	)

	GetAddressHandler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, c.GetMyAddress().String())
	}

	http.HandleFunc("/address", GetAddressHandler)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade failed: ", err)
			return
		}

		defer conn.Close()

		go func() {
			for {
				data := []byte{}
				a.Read(data)
				err = conn.WriteMessage(1, data)
				if err != nil {
					log.Println("write failed:", err)
					break
				}
			}
		}()

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

		c.CreateLedgerChannel(types.Address{}, 4, outcome.Exit{})

	})
	fmt.Println("rpc server listening on port " + rpcPort)
	_ = http.ListenAndServe(":"+rpcPort, nil)

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
	chain, err := ethclient.Dial("http://127.0.0.1:8545/")

	if err != nil {
		panic(err)
	}
	naAddress := common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")
	caAddress := common.HexToAddress("0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512")
	vpaAddress := common.HexToAddress("0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0")
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
