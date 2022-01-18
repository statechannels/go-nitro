package channel

const (
	PREFUNDTURNUM uint64 = iota
	POSTFUNDTURNNUM
	MAXTURNNUM = ^uint64(0) // MAXTURNNUM is a reserved value which is taken to mean "there is not yet a supported state"
)
