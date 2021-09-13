package outcome

import (
	"math/big"
)

// types
type Allocation struct {
	Destination    string
	Amount         big.Int
	AllocationType uint
	Metadata       string
}

type Allocations []Allocation

type SingleAssetExit struct {
	Asset, Metadata string
	Allocations     Allocations
}

type Exit []SingleAssetExit

// methods
func (a Allocation) Equals(b Allocation) bool {
	return a.Destination == b.Destination && a.AllocationType == b.AllocationType && a.Metadata == b.Metadata && a.Amount.Cmp(&b.Amount) == 0
}

func (a Allocations) Equals(b Allocations) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !a[i].Equals(b[i]) {
			return false
		}
	}
	return true
}