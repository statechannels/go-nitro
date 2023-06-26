package NitroAdjudicator

import (
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

// computeReclaimEffects mirrors on chain code.
// It computes side effects for the reclaim function. Returns updated allocations for the source, computed by finding the guarantee in the source for the target, and moving money out of the guarantee and back into the ledger channel as regular allocations for the participants.
func computeReclaimEffects(sourceAllocations []outcome.Allocation, targetAllocations []outcome.Allocation, indexOfTargetInSource uint) ([]outcome.Allocation, error) {
	newSourceAllocations := make([]outcome.Allocation, len(sourceAllocations)-1)
	guarantee := sourceAllocations[indexOfTargetInSource]
	guaranteeData, err := outcome.DecodeIntoGuaranteeMetadata(guarantee.Metadata)
	if err != nil {
		return []outcome.Allocation{}, err
	}

	var foundTarget, foundLeft, foundRight bool

	totalReclaimed := big.NewInt(0)

	k := 0
	for i := 0; i < len(sourceAllocations); i++ {
		if i == int(indexOfTargetInSource) {
			foundTarget = true
			continue
		}
		newSourceAllocations[k] = outcome.Allocation{
			Destination:    sourceAllocations[i].Destination,
			Amount:         sourceAllocations[i].Amount,
			AllocationType: sourceAllocations[i].AllocationType,
			Metadata:       sourceAllocations[i].Metadata,
		}
		if !foundLeft && sourceAllocations[i].Destination == guaranteeData.Left {
			newSourceAllocations[k].Amount.Add(newSourceAllocations[k].Amount, targetAllocations[0].Amount)
			totalReclaimed = totalReclaimed.Add(totalReclaimed, targetAllocations[0].Amount)
			foundLeft = true
		}
		if !foundRight && sourceAllocations[i].Destination == guaranteeData.Right {
			newSourceAllocations[k].Amount.Add(newSourceAllocations[k].Amount, targetAllocations[1].Amount)
			totalReclaimed = totalReclaimed.Add(totalReclaimed, targetAllocations[1].Amount)
			foundRight = true
		}
		k++
	}

	if !foundTarget {
		return []outcome.Allocation{}, fmt.Errorf("could not find target")
	}
	if !foundLeft {
		return []outcome.Allocation{}, fmt.Errorf("could not find left")
	}
	if !foundRight {
		return []outcome.Allocation{}, fmt.Errorf("could not find right")
	}
	if !types.Equal(totalReclaimed, guarantee.Amount) {
		return []outcome.Allocation{}, fmt.Errorf("totalReclaimed!=guarantee.amount")
	}
	return newSourceAllocations, nil
}
