package protocols

import (
	"testing"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

func TestEquals(t *testing.T) {
	stateOne := state.TestState.Clone()
	stateTwo := state.TestState.Clone()
	stateTwo.TurnNum = 1
	msg1 := Message{
		To:          types.Address{'a'},
		ObjectiveId: `say-hello-to-my-little-friend`,
		SignedStates: []state.SignedState{
			state.NewSignedState(stateOne),
		},
	}

	msg2 := Message{
		To:          types.Address{'a'},
		ObjectiveId: `say-hello-to-my-little-friend`,
		SignedStates: []state.SignedState{
			state.NewSignedState(stateTwo),
		},
	}
	if (msg1).Equal(msg2) {
		t.Error("Equal returned true for two different messages")
	}
}

func TestMessage(t *testing.T) {

	msg := Message{
		To:          types.Address{'a'},
		ObjectiveId: `say-hello-to-my-little-friend`,
		SignedStates: []state.SignedState{
			state.NewSignedState(state.TestState),
		},
		Proposal: true,
	}

	msgString := `{"To":"0x6100000000000000000000000000000000000000","ObjectiveId":"say-hello-to-my-little-friend","SignedStates":[{"State":{"ChainId":9001,"Participants":["0xf5a1bb5607c9d079e46d1b3dc33f257d937b43bd","0x760bf27cd45036a6c486802d30b5d90cffbe31fe"],"ChannelNonce":37140676580,"AppDefinition":"0x5e29e5ab8ef33f050c7cc10b5a0456d975c5f88d","ChallengeDuration":60,"AppData":"","Outcome":[{"Asset":"0x0000000000000000000000000000000000000000","Metadata":null,"Allocations":[{"Destination":[0,0,0,0,0,0,0,0,0,0,0,0,245,161,187,86,7,201,208,121,228,109,27,61,195,63,37,125,147,123,67,189],"Amount":5,"AllocationType":0,"Metadata":null},{"Destination":[0,0,0,0,0,0,0,0,0,0,0,0,238,24,255,21,117,5,86,145,0,154,162,70,174,96,129,50,197,122,66,44],"Amount":5,"AllocationType":0,"Metadata":null}]}],"TurnNum":5,"IsFinal":false},"Sigs":{}}],"Proposal":null}`

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
		got, err := DeserializeMessage(msgString)
		want := msg
		if err != nil {
			t.Error(err)
		}
		if !got.Equal(want) {
			t.Errorf(`incorrect deserialization: got %v wanted %v`, got, want)
		}
	})

}
