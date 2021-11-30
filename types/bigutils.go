package types

import "math/big"

//////////
// bigutils.go is a parking spot for small, reusable utility functions extending the big package
//////////

// Gt return true if a > b, false otherwise
func Gt(a *big.Int, b *big.Int) bool {
	return a.Cmp(b) > 0
}
