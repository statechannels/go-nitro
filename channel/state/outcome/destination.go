package outcome

import (
	"errors"

	"github.com/statechannels/go-nitro/types"
)

// IsExternal returns true if the destination has the 12 leading bytes as zero, false otherwise
func IsExternal(destination types.Destination) bool {
	for _, b := range destination[0:12] {
		if b != 0 {
			return false
		}
	}
	return true
}

// ToAddress returns a types.Address encoded external destination, or an error if
// destination is not an external address
func ToAddress(destination types.Destination) (types.Address, error) {
	if IsExternal(destination) {
		address := types.Address{}
		for i, b := range destination[12:] {
			address[i] = b
		}
		return address, nil
	}

	return types.Address{}, errors.New("destination is not external")
}
