package types

import (
	"errors"
)

// IsExternal returns true if the destination has the 12 leading bytes as zero, false otherwise
func (d Destination) IsExternal() bool {
	for _, b := range d[0:12] {
		if b != 0 {
			return false
		}
	}
	return true
}

// ToAddress returns a types.Address encoded external destination, or an error if
// destination is not an external address
func (d Destination) ToAddress() (Address, error) {
	if d.IsExternal() {
		address := Address{}
		for i, b := range d[12:] {
			address[i] = b
		}
		return address, nil
	}

	return Address{}, errors.New("destination is not external")
}

func (d Destination) String() string {
	return Bytes32(d).String()
}

func (d Destination) Bytes() []byte {
	return Bytes32(d).Bytes()
}
