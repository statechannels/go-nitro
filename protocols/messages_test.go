package protocols

import (
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/types"
)

func removeProposal() consensus_channel.SignedProposal {
	remove := consensus_channel.NewRemoveProposal(types.Destination{'l'}, types.Destination{'a'}, big.NewInt(1))
	return consensus_channel.SignedProposal{Proposal: remove, Signature: state.Signature{}}
}

func addProposal() consensus_channel.SignedProposal {
	amount := big.NewInt(1)
	add := consensus_channel.NewAddProposal(types.Destination{'l'}, consensus_channel.NewGuarantee(
		amount,
		types.Destination{'a'},
		types.Destination{'b'},
		types.Destination{'c'},
	),
		amount,
	)

	return consensus_channel.SignedProposal{Proposal: add, Signature: state.Signature{}}
}

func TestMessage(t *testing.T) {
	msg := Message{
		To:   types.Address{'a'},
		From: types.Address{'b'},
		payloads: []messagePayload{{
			ObjectiveId: `say-hello-to-my-little-friend`,
			SignedState: state.NewSignedState(state.TestState),
		}, {
			ObjectiveId:    `say-hello-to-my-little-friend2`,
			SignedProposal: addProposal(),
		},
			{
				ObjectiveId:    `say-hello-to-my-little-friend3`,
				SignedProposal: removeProposal(),
			},

			{

				Voucher: payments.Voucher{ChannelId: types.Destination{'d'}, Amount: big.NewInt(123), Signature: state.Signature{}},
			},
		},
	}

	msgString :=
		`{"To":"0x6100000000000000000000000000000000000000","From":"0x6200000000000000000000000000000000000000","Payloads":[{"ObjectiveId":"say-hello-to-my-little-friend","SignedState":{"State":{"ChainId":9001,"Participants":["0xf5a1bb5607c9d079e46d1b3dc33f257d937b43bd","0x760bf27cd45036a6c486802d30b5d90cffbe31fe"],"ChannelNonce":37140676580,"AppDefinition":"0x5e29e5ab8ef33f050c7cc10b5a0456d975c5f88d","ChallengeDuration":60,"AppData":"","Outcome":[{"Asset":"0x0000000000000000000000000000000000000000","Metadata":null,"Allocations":[{"Destination":"0x000000000000000000000000f5a1bb5607c9d079e46d1b3dc33f257d937b43bd","Amount":5,"AllocationType":0,"Metadata":null},{"Destination":"0x000000000000000000000000ee18ff1575055691009aa246ae608132c57a422c","Amount":5,"AllocationType":0,"Metadata":null}]}],"TurnNum":5,"IsFinal":false},"Sigs":{}}},{"ObjectiveId":"say-hello-to-my-little-friend2","SignedProposal":{"R":null,"S":null,"V":0,"Proposal":{"LedgerID":"0x6c00000000000000000000000000000000000000000000000000000000000000","ToAdd":{"Guarantee":{"Amount":1,"Target":"0x6100000000000000000000000000000000000000000000000000000000000000","Left":"0x6200000000000000000000000000000000000000000000000000000000000000","Right":"0x6300000000000000000000000000000000000000000000000000000000000000"},"LeftDeposit":1},"ToRemove":{"Target":"0x0000000000000000000000000000000000000000000000000000000000000000","LeftAmount":null}},"TurnNum":0}},{"ObjectiveId":"say-hello-to-my-little-friend3","SignedProposal":{"R":null,"S":null,"V":0,"Proposal":{"LedgerID":"0x6c00000000000000000000000000000000000000000000000000000000000000","ToAdd":{"Guarantee":{"Amount":null,"Target":"0x0000000000000000000000000000000000000000000000000000000000000000","Left":"0x0000000000000000000000000000000000000000000000000000000000000000","Right":"0x0000000000000000000000000000000000000000000000000000000000000000"},"LeftDeposit":null},"ToRemove":{"Target":"0x6100000000000000000000000000000000000000000000000000000000000000","LeftAmount":1}},"TurnNum":0}},{"ObjectiveId":"","Voucher":{"ChannelId":"0x6400000000000000000000000000000000000000000000000000000000000000","Amount":123,"Signature":{"R":null,"S":null,"V":0}}}]}`

	t.Run(`serialize`, func(t *testing.T) {
		got, err := msg.Serialize()
		if err != nil {
			t.Error(err)
		}
		want := msgString
		if got != want {
			t.Fatalf("incorrect serialization: got:\n%v\nwanted:\n%v", got, want)
		}
	})

	t.Run(`deserialize`, func(t *testing.T) {
		got, err := DeserializeMessage(msgString)
		want := msg
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("incorrect deserialization: got:\n%v\nwanted:\n%v", got, want)
		}
	})

	t.Run(`validation`, func(t *testing.T) {

		invalidMsg := `{"To":"0x6100000000000000000000000000000000000000","Payloads":[{"ObjectiveId":"say-hello-to-my-little-friend","SignedState":{"State":{"ChainId":9001,"Participants":["0xf5a1bb5607c9d079e46d1b3dc33f257d937b43bd","0x760bf27cd45036a6c486802d30b5d90cffbe31fe"],"ChannelNonce":37140676580,"AppDefinition":"0x5e29e5ab8ef33f050c7cc10b5a0456d975c5f88d","ChallengeDuration":60,"AppData":"","Outcome":[{"Asset":"0x0000000000000000000000000000000000000000","Metadata":null,"Allocations":[{"Destination":"0x000000000000000000000000f5a1bb5607c9d079e46d1b3dc33f257d937b43bd","Amount":5,"AllocationType":0,"Metadata":null},{"Destination":"0x000000000000000000000000ee18ff1575055691009aa246ae608132c57a422c","Amount":5,"AllocationType":0,"Metadata":null}]}],"TurnNum":5,"IsFinal":false},"Sigs":{}},"SignedProposal":{"R":null,"S":null,"V":0,"Proposal":{"LedgerID":"0x6c00000000000000000000000000000000000000000000000000000000000000","ToAdd":{"Guarantee":{"Amount":1,"Target":"0x6100000000000000000000000000000000000000000000000000000000000000","Left":"0x6200000000000000000000000000000000000000000000000000000000000000","Right":"0x6300000000000000000000000000000000000000000000000000000000000000"},"LeftDeposit":1},"ToRemove":{"Target":"0x0000000000000000000000000000000000000000000000000000000000000000","LeftAmount":null}},"TurnNum":0}}]}`

		_, err := DeserializeMessage(invalidMsg)
		if !errors.Is(err, ErrInvalidPayload) {
			t.Fatalf("expected error deserializing invalid payload")
		}
	})
}
