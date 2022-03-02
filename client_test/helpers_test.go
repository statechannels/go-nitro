package client_test

import (
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/protocols"
)

// waitForCompletedObjectiveId waits for completed objectives and returns when the completed objective id matchs the id waitForCompletedObjectiveId has been given
func waitForCompletedObjectiveId(id protocols.ObjectiveId, client *client.Client) {
	got := <-client.CompletedObjectives()
	for got != id {
		got = <-client.CompletedObjectives()
	}
}
