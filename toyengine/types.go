package main

type LedgerStore struct {
	ledgerChannels map[uint]LedgerChannelState
}
type LedgerChannelState struct {
	hubId, hubBal, leafId, leafBal, turnNum uint            // hubs and leaves have integer ids
	virtualChannelBal                       map[string]uint // channels have 2d uint ids [joiner][proposer]
	signedByHub                             bool
	signedByLeaf                            bool
}

type VirtualChannelState struct {
	proposerId, joinerId, proposerBal, joinerBal, turnNum uint
}

type LedgerRequest struct {
	virtualChannelProposer, virtualChannelJoiner uint
	amount                                       int // note this is a *signed* quantity
	sucess                                       chan bool
}
