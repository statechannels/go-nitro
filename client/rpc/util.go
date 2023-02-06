package rpc

import (
	"fmt"

	"github.com/statechannels/go-nitro/network/serde"
)

func methodToTopic(method serde.RequestMethod) string {
	return fmt.Sprintf("nitro.%s", method)
}
