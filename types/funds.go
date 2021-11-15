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

// Add returns the sum of the receiver and the input Funds objects
func (h Funds) Add(a ...Funds) Funds {
	a = append(a, h)
	return Sum(a...)
}

// Sum returns a new Funds object with all of the asset keys from the supplied Funds objects,
// each having an amount summed across that asset's amount in each input object.
//
// e.g. {[0x0a,0x01][0x0b,0x01]} + {[0x0a,0x02]} = {[0x0a,0x03][0x0b,0x01]}
func Sum(a ...Funds) Funds {
	sum := Funds{}

	for _, holdings := range a {
		for asset, amount := range holdings {
			sum[asset] = sum[asset].Add(sum[asset], amount)
		}
	}

	return sum
}

// Equal returns true if reciever `f` and input `g` are identical.
//
// Note that a zero-balance does NOT equal a non-balance: {[0x0a,0x00]} != {}
func (f Funds) Equal(g Funds) bool {
	if len(f) != len(g) {
		return false
	}

	for asset, amount := range f {
		if g[asset].Cmp(amount) != 0 {
			return false
		}
	}

	return true
}
