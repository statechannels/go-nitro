package margin

import (
	"math/big"

	"github.com/statechannels/go-nitro/apps"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/internal"
)

// TODO: Extract common errors into a common package
var ErrInvalidRequestType = internal.NewError("invalid request type")

type Balance struct {
	Remaining *big.Int
	Paid      *big.Int
}

type MarginApp struct {
	balances map[string]*Balance

	//
}

var _ apps.App = (*MarginApp)(nil)

func (a *MarginApp) Type() string {
	return "margin"
}

func (a *MarginApp) HandleRequest(ch *channel.Channel, ty string, data interface{}) error {
	switch ty {
	case VoucherRequestType:
		voucher := data.(Voucher)

		// TODO: validate and use voucher to update app balances
		_ = voucher

	default:
		return ErrInvalidRequestType
	}

	return nil
}
