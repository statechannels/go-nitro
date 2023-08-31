package NitroAdjudicator

import (
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
)

// ConvertBindingsExitToExit converts the exit type returned from abigen bindings to an outcome.Exit
func ConvertBindingsExitToExit(e []ExitFormatSingleAssetExit) outcome.Exit {
	exit := make([]outcome.SingleAssetExit, 0, len(e))
	for _, sae := range e {
		exit = append(exit, convertBindingsSingleAssetExitToSingleAssetExit(sae))
	}
	return exit
}

func convertBindingsSingleAssetExitToSingleAssetExit(e ExitFormatSingleAssetExit) outcome.SingleAssetExit {
	return outcome.SingleAssetExit{
		Asset: e.Asset,
		AssetMetadata: outcome.AssetMetadata{
			AssetType: outcome.AssetType(e.AssetMetadata.AssetType),
			Metadata:  e.AssetMetadata.Metadata,
		},
		Allocations: convertBindingsAllocationsToAllocations(e.Allocations),
	}
}

func convertBindingsAllocationsToAllocations(as []ExitFormatAllocation) outcome.Allocations {
	allocations := make([]outcome.Allocation, 0, len(as))
	for _, a := range as {
		allocations = append(allocations, outcome.Allocation{
			Destination:    a.Destination,
			Amount:         a.Amount,
			Metadata:       a.Metadata,
			AllocationType: outcome.AllocationType(a.AllocationType),
		})
	}
	return allocations
}

// ConvertBindingsSignatureToSignature converts the signature type returned from abigien bindings to a state.Signature
func ConvertBindingsSignatureToSignature(s INitroTypesSignature) state.Signature {
	return state.Signature{
		R: s.R[:],
		S: s.S[:],
		V: s.V,
	}
}

// ConvertBindingsSignatureToSignature converts a slice of the signature type returned from abigien bindings to a []state.Signature
func ConvertBindingsSignaturesToSignatures(ss []INitroTypesSignature) []state.Signature {
	sigs := make([]state.Signature, 0, len(ss))
	for _, s := range ss {
		sigs = append(sigs, ConvertBindingsSignatureToSignature(s))
	}
	return sigs
}
