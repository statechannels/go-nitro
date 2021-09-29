package outcome

import (
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
