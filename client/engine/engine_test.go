package engine

import (
	"fmt"
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/protocols"
)

func TestRun(t *testing.T) {
	fmt.Println(`constructing engine`)

	e := New(messageservice.TestMessageService{}, chainservice.TestChainService{}, store.TestStore{})

	go e.Run()
	msg := protocols.Message{
		ObjectiveId: `test`,
	}

	for j := 1; j <= 100_000; j++ {
		e.fromMsg <- msg
	}

	// construct an engine with a test store
	// hit the API repeatedly in such a way as to trigger the testobjective's expensive crank
	// try this with and without multithreading

}
