package payments

import "github.com/statechannels/go-nitro/types"

const PAYER_INDEX = 0
const PAYEE_INDEX = 2

// GetPayer returns the payer on a payment channel
func GetPayer(participants []types.Address) types.Address {
	return participants[PAYER_INDEX]
}

// GetPayee returns the payee on a payment channel
func GetPayee(participants []types.Address) types.Address {
	return participants[PAYEE_INDEX]
}
