package client

// Engine implements the business logic of a Nitro wallet and coordinates the other components.
type Engine interface {
	AppChannelManager
	EngineCoordinator
}

type AppChannelManager interface {
	CreateChannels()
	JoinChannels()
	UpdateChannel()
	CloseChannel()
	GetChannels()
	GetState()
	Challenge()
}
type EngineCoordinator interface {
	SyncChannels()
	SyncChannel()
	PushMessage()
	PushUpdate()
	Crank()
}
