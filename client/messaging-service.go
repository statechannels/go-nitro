package client

// - A `MessagingService` will be implemented to enable the peer-to-peer messaging required to execute Nitro protocols.
type MessagingService interface {

	// Send sendss out the messages to the specified by participants.
	Send()

	// RegisterPeer registers a peer so messages can be sent to that peer
	RegisterPeer()

	// Destory is called when the message service is no longer needed
	Destroy()
}
