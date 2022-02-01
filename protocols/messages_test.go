package protocols

import (
	"testing"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

func TestMessage(t *testing.T) {
	msg := Message{
		To:          types.Address{'a'},
		ObjectiveId: `say-hello-to-my-little-friend`,
		SignedStates: []state.SignedState{
			state.NewSignedState(state.TestState),
		},
	}

	// TODO this doesn't contain the signed states!
	msgString := `{"To":"0x6100000000000000000000000000000000000000","ObjectiveId":"say-hello-to-my-little-friend","SignedStates":[{}],"Proposal":null}`

	t.Run(`serialize`, func(t *testing.T) {
		got, err := msg.Serialize()
		if err != nil {
			t.Error(err)
		}
		want := msgString
		if got != want {
			t.Errorf(`incorrect serialization: got %v wanted %v`, got, want)
		}
	})

	t.Run(`deserialize`, func(t *testing.T) {
		got, err := DeserialiseMessage(msgString)
		want := msg
		if err != nil {
			t.Error(err)
		}
		if !got.Equal(want) {
			t.Errorf(`incorrect deserialization: got %v wanted %v`, got, want)
		}
	})

}
