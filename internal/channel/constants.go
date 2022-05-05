package channel

const (
	PreFundTurnNum uint64 = iota
	PostFundTurnNum
	MaxTurnNum = ^uint64(0) // MaxTurnNum is a reserved value which is taken to mean "there is not yet a supported state"
)
