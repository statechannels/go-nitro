package types

import "math/big"

//////////
// bigutils.go is a parking spot for small, reusable utility functions extending the big package
//////////

// Gt returns true if a > b, false otherwise
func Gt(a *big.Int, b *big.Int) bool {
	return a.Cmp(b) > 0
}

// Lt returns true if a < b, false otherwise
func Lt(a *big.Int, b *big.Int) bool {
	return a.Cmp(b) < 0
}

// Equal returns true if a == b or if both of
// a and b are nil, and false otherwise.
func Equal(a *big.Int, b *big.Int) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return a == nil
	}
	return a.Cmp(b) == 0
}

// IsZero returns true if a is zero, false otherwise.
func IsZero(a *big.Int) bool {
	return a.Cmp(big.NewInt(0)) == 0
}

// Max returns a if a > b, b otherwise.
func Max(a *big.Int, b *big.Int) *big.Int {
	if Gt(a, b) {
		return a
	}
	return b
}

// Max returns a if a > b, b otherwise.
func Min(a *big.Int, b *big.Int) *big.Int {
	if Lt(a, b) {
		return a
	}
	return b
}
