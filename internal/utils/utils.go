package utils

import (
	"fmt"
	"os/exec"
	"syscall"

	p2pms "github.com/statechannels/go-nitro/node/engine/messageservice/p2p-message-service"
)

// WaitForPeerInfoExchange waits for all the P2PMessageServices to receive peer info from each other
func WaitForPeerInfoExchange(services ...*p2pms.P2PMessageService) {
	for _, s := range services {
		for i := 0; i < len(services)-1; i++ {
			<-s.PeerInfoReceived()
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
