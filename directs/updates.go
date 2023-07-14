package directs

import (
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

// A ChannelUpdate is a message sent by a participant to the client to either validate
// data and countersign it or to notify the client, that other party agreed on the earlier
// appdata change proposal.
type ChannelUpdates struct {
	ChannelId  types.Destination
	AppData    types.Bytes
	Signatures []crypto.Signature
}
