package utils

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	p2pms "github.com/statechannels/go-nitro/node/engine/messageservice/p2p-message-service"
)

// WaitForPeerInfoExchange waits for all the P2PMessageServices to receive peer info from each other
func WaitForPeerInfoExchange(services ...*p2pms.P2PMessageService) {
	for sNum, s := range services {
		for i := 0; i < len(services)-1; i++ {
			peerInfo := <-s.PeerInfoReceived()
			fmt.Printf("Service num: %d, peer num: %d, peerInfo: %v\n", sNum, i, peerInfo)
		}
	}
}

// StopCommands stops the given executing commands
func StopCommands(cmds ...*exec.Cmd) {
	for _, cmd := range cmds {
		fmt.Printf("Stopping process %v\n", cmd.Args)
		err := cmd.Process.Signal(syscall.SIGINT)
		if err != nil {
			panic(err)
		}
		err = cmd.Process.Kill()
		if err != nil {
			panic(err)
		}
	}
}

// GenerateTempStoreFolder generates a temporary folder for storing store data and a cleanup function to clean up the folder
func GenerateTempStoreFolder() (dataFolder string, cleanup func()) {
	var err error

	dataFolder, err = os.MkdirTemp("", "nitro-store-*")
	if err != nil {
		panic(err)
	}

	cleanup = func() {
		err := os.RemoveAll(dataFolder)
		if err != nil {
			panic(err)
		}
	}

	return
}
