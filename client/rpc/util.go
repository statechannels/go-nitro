package rpc

import (
	"fmt"

	"github.com/statechannels/go-nitro/network/serde"
)

// getTopics returns a list of topics that the client/server should subscribe to.
func getTopics() []string {

	return []string{
		fmt.Sprintf("nitro.%s",
			serde.DirectDefundRequestMethod),
		fmt.Sprintf("nitro.%s",
			serde.VirtualFundRequestMethod),
		fmt.Sprintf("nitro.%s",
			serde.VirtualDefundRequestMethod),
		fmt.Sprintf("nitro.%s",
			serde.PayRequestMethod),
	}
}
