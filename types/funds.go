package types

import "math/big"

// IsNonZero returns true if the Holdings structure has any non-zero asset
func (h Funds) IsNonZero() bool {
	for _, v := range h {
		if v.Cmp(big.NewInt(0)) == 1 {
			return true
		}
	}
	return false
}

// String returns a bracket-separaged list of assets: {[0x0a,0x01][0x0b,0x01]}
func (h Funds) String() string {
	if len(h) == 0 {
		return "{}"
	}
	var s string = "{"
	for asset, amount := range h {
		s += "[" + asset.Hex() + "," + amount.Text(64) + "]"
	}
	s = s + "}"
	return s
}

// todo:
// ToFunds returns a Funds map from its string representation
// func ToFunds(s string) Funds {}

// Add sums all assets and returns the result. Does not modify the calling Holdings object
func (h Funds) Add(a ...Funds) Funds {
	a = append(a, h)
	return Sum(a...)
}

// Sum returns the sum of all input Funds maps
func Sum(a ...Funds) Funds {
	sum := Funds{}

	for _, holdings := range a {
		for asset, amount := range holdings {
			sum[asset] = sum[asset].Add(sum[asset], amount)
		}
	}

	return sum
}
