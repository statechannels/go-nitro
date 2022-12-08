package virtualmargin

import "github.com/statechannels/go-nitro/types"

const LEADER_INDEX = 0

func GetLeader(participants []types.Address) types.Address {
	return participants[LEADER_INDEX]
}

// GetPayee returns the payee on a payment channel
func GetFollower(participants []types.Address) types.Address {
	return participants[len(participants)-1]
}
