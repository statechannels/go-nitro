package outcome

import (
	"errors"

	"github.com/statechannels/go-nitro/types"
)

// IsExternalDestination returns true if the destination has the 12 leading bytes as zero, false otherwise
func IsExternalDestination(destination types.Bytes32) bool {
	for _, b := range destination[0:12] {
		if b != 0 {
			return false
		}
	}
	return true
}

// ToExternalDestination returns a types.Address encoded external destination, or an error if
// destination is not external
func ToExternalDestination(destination types.Bytes32) (types.Address, error) {
	if IsExternalDestination(destination) {
		address := types.Address{}
		for i, b := range destination[12:] {
			address[i] = b
		}
		return address, nil
	}

	return types.Address{}, errors.New("destination is not external")
}
