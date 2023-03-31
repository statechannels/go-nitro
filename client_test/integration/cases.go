package integration_test

import "github.com/statechannels/go-nitro/internal/testactors"

var simpleCase = TestCase{
	Description:    "Simple test",
	Chain:          MockChain,
	MessageService: TestMessageService,
	NumOfChannels:  1,
	MessageDelay:   0,
	LogName:        "simple_integration_run.log",
	NumOfHops:      1,
	NumOfPayments:  1,
	Participants: []TestParticipant{
		{StoreType: MemStore, Name: testactors.AliceName},
		{StoreType: MemStore, Name: testactors.BobName},
		{StoreType: MemStore, Name: testactors.IreneName},
	},
}

var complexCase = TestCase{
	Description:    "Complex test",
	Chain:          SimulatedChain,
	MessageService: P2PMessageService,
	NumOfChannels:  5,
	MessageDelay:   0,
	LogName:        "complex_integration_run.log",
	NumOfHops:      1,
	NumOfPayments:  5,
	Participants: []TestParticipant{
		{StoreType: DurableStore, Name: testactors.AliceName},
		{StoreType: DurableStore, Name: testactors.BobName},
		{StoreType: DurableStore, Name: testactors.IreneName},
		{StoreType: DurableStore, Name: testactors.BrianName},
	},
}

var cases = []TestCase{simpleCase, complexCase}
