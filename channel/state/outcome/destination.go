package outcome

import (
	"github.com/statechannels/go-nitro/types"
)

func IsExternalDestination(destination types.Bytes32) bool {

	for i, b := range destination[0:12] {

		if i > 11 {
			break
		}
		if b != 0 {
			return false
		}
	}
	return true

}
