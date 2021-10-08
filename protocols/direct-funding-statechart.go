package protocols

// A linear state machine
// PreFundIncomplete => NotYetMyTurnToFund => FundingIncomplete => PostFundIncomplete
type DirectFundingProtocolState int

const (
	PreFundIncomplete DirectFundingProtocolState = iota // 0
	NotYetMyTurnToFund
	FundingIncomplete
	PostFundIncomplete
)

// The events for the state machine
type DirectFundingProtocolEvent int // TODO these need to be richer than ints
const (
	PrefundReceived DirectFundingProtocolEvent = iota // 0
	FundingUpdated
	PostFundReceived
)

// TODO we need to also have a context to turn this into a state chart

func (s DirectFundingProtocolState) NextState(e DirectFundingProtocolEvent) (DirectFundingProtocolState, error) {
	// it is better to switch on the state than on the event
	// https://dev.to/davidkpiano/you-don-t-need-a-library-for-state-machines-k7h
	switch s {
	case PreFundIncomplete:
		return s.nextStateFromPrefundIncomplete(e)
	case NotYetMyTurnToFund:
		fallthrough // TODO
	case FundingIncomplete:
		fallthrough // TODO
	case PostFundIncomplete:
		fallthrough // TODO
	default:
		return s, nil
	}
}

func (s DirectFundingProtocolState) nextStateFromPrefundIncomplete(e DirectFundingProtocolEvent) (DirectFundingProtocolState, error) {
	// TODO checks on the context
	// here we need the context plus the event to be a full postfund setup?
	if e == PostFundReceived {
		return NotYetMyTurnToFund, nil
	} else {
		return s, nil
	}
}
