package types

import (
	"errors"
)

// IsExternal returns true if the destination is a blockchain address, and false
// if it is a state channel ID.
func (d Destination) IsExternal() bool {
	for _, b := range d[0:12] {
		if b != 0 {
			return false
		}
	}
	return true
}

// IsZero returns true if the destination is all zeros, and false otherwise.
func (d Destination) IsZero() bool {
	for _, b := range d {
		if b != 0 {
			return false
		}
	}
	return true
}

// ToAddress returns a types.Address encoded external destination, or an error if
// destination is not an external address
func (d Destination) ToAddress() (Address, error) {
	if !d.IsExternal() {
		return Address{}, errors.New("destination is not an external address")
	}

	address := Address{}
	copy(address[:], d[12:])
	return address, nil
}

func (d Destination) String() string {
	return Bytes32(d).String()
}

func (d Destination) Bytes() []byte {
	return Bytes32(d).Bytes()
}

// AddressToDestinaion left-pads the blockchain address with zeros.
func AddressToDestination(a Address) Destination {
	d := Destination{0}
	for i := range a {
		d[i+12] = a[i]
	}
	return d
}
