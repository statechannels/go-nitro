package serde

import (
	"github.com/statechannels/go-nitro/types"
)

func ValidatePaymentRequest(req PaymentRequest) error {
	if req.Amount == 0 {
		return types.InvalidParamsError
	}
	if (req.Channel == types.Destination{}) {
		return types.InvalidParamsError
	}
	return nil
}

func ValidateGetPaymentChannelRequest(req GetPaymentChannelRequest) error {
	if (req.Id == types.Destination{}) {
		return types.InvalidParamsError
	}
	return nil
}

func ValidateGetPaymentChannelsByLedgerRequest(req GetPaymentChannelsByLedgerRequest) error {
	if (req.LedgerId == types.Destination{}) {
		return types.InvalidParamsError
	}
	return nil
}
