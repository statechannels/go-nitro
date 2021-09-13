package outcome

import (
	"math/big"
)


func min(a big.Int, b big.Int) big.Int {
	switch a.Cmp(&b) {
	case -1:
		return *big.NewInt(0).Set(&a)
	default:
		return *big.NewInt(0).Set(&b)
	}
}

func ComputeTransferEffectsAndInteractions(initialHoldings big.Int, allocations Allocations, indices []uint) (newAllocations Allocations, exitAllocations Allocations) {
	// TODO here we assume indices = [], so pay out all allocations
	surplus := big.NewInt(0).Set(&initialHoldings)
	newAllocations = make([]Allocation, len(allocations))
	exitAllocations = make([]Allocation, len(allocations))

	for i := 0; i < len(allocations); i++ {
		// copy allocation
		newAllocations[i] = Allocation{Destination: allocations[i].Destination, Amount: *big.NewInt(0).Set(&allocations[i].Amount), AllocationType: allocations[i].AllocationType, Metadata: allocations[i].Metadata}
		// compute payout amount
		affordsForDestination := min(allocations[i].Amount, *surplus)
		// decrease allocation amount
		newAllocations[i].Amount.Sub(&newAllocations[i].Amount, &affordsForDestination)
		// increase exit allocation amount
		exitAllocations[i] = Allocation{Destination: allocations[i].Destination, Amount: *big.NewInt(0).Set(&affordsForDestination), AllocationType: allocations[i].AllocationType, Metadata: allocations[i].Metadata}
		// decrease surplus
		surplus.Sub(surplus, &affordsForDestination)
	}

	return

}
