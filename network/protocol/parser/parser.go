package parser

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

func createAllocations(allocationInterfaces []any) []outcome.Allocation {
	allocationsArray := make([]outcome.Allocation, len(allocationInterfaces))
	for i := 0; i < len(allocationInterfaces); i++ {
		alloc := allocationInterfaces[i].(map[string]any)
		allocationsArray[i] = outcome.Allocation{
			Destination:    types.AddressToDestination(common.HexToAddress(alloc["Destination"].(string))),
			Amount:         i2Uint256(alloc["Amount"]),
			AllocationType: outcome.AllocationType(i2Uint8(alloc["AllocationType"])),
			Metadata:       toByteArray(alloc["Metadata"]),
		}
	}
	return allocationsArray
}

func createExit(outcomesInterfaces []any) outcome.Exit {
	var e = outcome.Exit{}
	for i := 0; i < len(outcomesInterfaces); i++ {
		out := outcomesInterfaces[i].(map[string]any)
		allocations := out["Allocations"].([]any)
		allocationsArray := createAllocations(allocations)

		e = append(e, outcome.SingleAssetExit{
			Asset:       common.HexToAddress(out["Asset"].(string)),
			Metadata:    toByteArray(out["Metadata"]),
			Allocations: allocationsArray,
		})
	}

	return e
}

func ParseDirectFundRequest(data map[string]any) directfund.ObjectiveRequest {
	outcomes := data["Outcome"].([]any)

	counterParty := common.HexToAddress(data["CounterParty"].(string))
	challengeDuration := i2Uint32(data["ChallengeDuration"])
	outcome := createExit(outcomes)
	appDefinition := common.HexToAddress(data["AppDefinition"].(string))
	appData := toByteArray(data["AppData"])
	nonce := i2Uint64(data["Nonce"])

	r := directfund.NewObjectiveRequest(counterParty, challengeDuration, outcome, nonce, appDefinition)
	r.AppData = appData

	return r
}

func ParseDirectFundResponse(data map[string]any) directfund.ObjectiveResponse {
	r := directfund.ObjectiveResponse{
		Id:        protocols.ObjectiveId(data["Id"].(string)),
		ChannelId: types.Destination(common.HexToHash(data["ChannelId"].(string))),
	}

	return r
}

func ParseDirectDefundRequest(data map[string]any) directdefund.ObjectiveRequest {
	channelId := types.Destination(common.HexToHash(data["ChannelId"].(string)))
	r := directdefund.NewObjectiveRequest(channelId)

	return r
}

func ParseVirtualFundRequest(data map[string]any) virtualfund.ObjectiveRequest {
	outcomes := data["outcome"].([]any)

	intermediaries := hexesToAddresses(data["Intermediaries"].([]string))
	counterParty := common.HexToAddress(data["CounterParty"].(string))
	challengeDuration := i2Uint32(data["ChallengeDuration"])
	outcome := createExit(outcomes)
	nonce := i2Uint64(data["Nonce"])
	appDefinition := common.HexToAddress(data["AppDefinition"].(string))

	r := virtualfund.NewObjectiveRequest(intermediaries, counterParty, challengeDuration, outcome, nonce, appDefinition)

	return r
}

func ParseVirtualFundResponse(data map[string]any) virtualfund.ObjectiveResponse {
	r := virtualfund.ObjectiveResponse{
		Id:        protocols.ObjectiveId(data["Id"].(string)),
		ChannelId: types.Destination(common.HexToHash(data["ChannelId"].(string))),
	}

	return r
}

func ParseVirtualDefundRequest(data map[string]any) virtualdefund.ObjectiveRequest {
	channelId := types.Destination(common.HexToHash(data["ChannelId"].(string)))
	r := virtualdefund.NewObjectiveRequest(channelId)

	return r
}
