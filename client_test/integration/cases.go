package integration_test

import "github.com/statechannels/go-nitro/internal/testactors"

var simpleCase = TestCase{
	Description:    "Simple test: 1 channel, 1 hop, MockChain, MockMessageService",
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

var cases = []TestCase{simpleCase}
