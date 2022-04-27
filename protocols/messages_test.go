package protocols

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/statechannels/go-nitro/channel/ledger"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

func addProposal() ledger.SignedProposal {
	amount := big.NewInt(1)
	add := ledger.NewAddProposal(types.Destination{'l'}, 1, ledger.NewGuarantee(
		amount,
		types.Destination{'a'},
		types.Destination{'b'},
		types.Destination{'c'},
	),
		amount,
	)

	return ledger.SignedProposal{Proposal: add, Signature: state.Signature{}}
}

func TestMessage(t *testing.T) {

	msgs := []Message{{
		To: types.Address{'a'},

		Payloads: []MessagePayload{{
			ObjectiveId: `say-hello-to-my-little-friend`,
			SignedState: state.NewSignedState(state.TestState),
		}}}, {
		To: types.Address{'a'},

		Payloads: []MessagePayload{{
			ObjectiveId:    `say-hello-to-my-little-friend2`,
			SignedProposal: addProposal(),
		}}}}

	msgStrings := []string{
		`{"To":"0x6100000000000000000000000000000000000000","Payloads":[{"ObjectiveId":"say-hello-to-my-little-friend","SignedState":{"State":{"ChainId":9001,"Participants":["0xf5a1bb5607c9d079e46d1b3dc33f257d937b43bd","0x760bf27cd45036a6c486802d30b5d90cffbe31fe"],"ChannelNonce":37140676580,"AppDefinition":"0x5e29e5ab8ef33f050c7cc10b5a0456d975c5f88d","ChallengeDuration":60,"AppData":"","Outcome":[{"Asset":"0x0000000000000000000000000000000000000000","Metadata":null,"Allocations":[{"Destination":"0x000000000000000000000000f5a1bb5607c9d079e46d1b3dc33f257d937b43bd","Amount":5,"AllocationType":0,"Metadata":null},{"Destination":"0x000000000000000000000000ee18ff1575055691009aa246ae608132c57a422c","Amount":5,"AllocationType":0,"Metadata":null}]}],"TurnNum":5,"IsFinal":false},"Sigs":{}},"SignedProposal":{"R":null,"S":null,"V":0,"Proposal":{"ChannelID":"0x0000000000000000000000000000000000000000000000000000000000000000","ToAdd":{"TurnNum":0,"Guarantee":{"Amount":null,"Target":"0x0000000000000000000000000000000000000000000000000000000000000000","Left":"0x0000000000000000000000000000000000000000000000000000000000000000","Right":"0x0000000000000000000000000000000000000000000000000000000000000000"},"LeftDeposit":null},"ToRemove":{"Target":"0x0000000000000000000000000000000000000000000000000000000000000000","LeftAmount":null,"RightAmount":null}}}}]}`,
		`{"To":"0x6100000000000000000000000000000000000000","Payloads":[{"ObjectiveId":"say-hello-to-my-little-friend2","SignedState":{"State":{"ChainId":null,"Participants":null,"ChannelNonce":null,"AppDefinition":"0x0000000000000000000000000000000000000000","ChallengeDuration":null,"AppData":null,"Outcome":null,"TurnNum":0,"IsFinal":false},"Sigs":null},"SignedProposal":{"R":null,"S":null,"V":0,"Proposal":{"ChannelID":"0x6c00000000000000000000000000000000000000000000000000000000000000","ToAdd":{"TurnNum":1,"Guarantee":{"Amount":1,"Target":"0x6100000000000000000000000000000000000000000000000000000000000000","Left":"0x6200000000000000000000000000000000000000000000000000000000000000","Right":"0x6300000000000000000000000000000000000000000000000000000000000000"},"LeftDeposit":1},"ToRemove":{"Target":"0x0000000000000000000000000000000000000000000000000000000000000000","LeftAmount":null,"RightAmount":null}}}}]}`,
	}

	for i, msg := range msgs {
		t.Run(`serialize`, func(t *testing.T) {
			got, err := msg.Serialize()
			if err != nil {
				t.Error(err)
			}
			want := msgStrings[i]
			if got != want {
				t.Fatalf("incorrect serialization: got:\n%v\nwanted:\n%v", got, want)
			}
		})

		t.Run(`deserialize`, func(t *testing.T) {
			got, err := DeserializeMessage(msgStrings[i])
			want := msg
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("incorrect deserialization: got:\n%v\nwanted:\n%v", got, want)
			}
		})
	}
}
