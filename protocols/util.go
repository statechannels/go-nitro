package protocols

import "math/big"

//////////
// util.go is a parking spot for small, reusable utility functions
//////////

// gte returns true if a >= b
func gte(a *big.Int, b *big.Int) bool {
	return a.Cmp(b) > -1
}

// gt return true if a > b
func gt(a *big.Int, b *big.Int) bool {
	return a.Cmp(b) > 0
}
