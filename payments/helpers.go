package payments

import "github.com/statechannels/go-nitro/types"

const PAYMENT_SENDER_INDEX = 0
const PAYMENT_RECEIVER_INDEX = 2

// GetPaymentSender returns the sender on a payment channel
func GetPaymentSender(participants []types.Address) types.Address {
	return participants[PAYMENT_SENDER_INDEX]
}

// GetPaymentReceiver returns the receiver on a payment channel
func GetPaymentReceiver(participants []types.Address) types.Address {
	return participants[PAYMENT_RECEIVER_INDEX]
}
