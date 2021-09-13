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
