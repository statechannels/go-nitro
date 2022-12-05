package margin

import (
	"github.com/statechannels/go-nitro/apps"
	"github.com/statechannels/go-nitro/channel"
)

type PaymentApp struct {
	//
}

var _ apps.App = (*PaymentApp)(nil)

func (a *PaymentApp) Type() string {
	return "payment"
}

func (a *PaymentApp) HandleRequest(ch *channel.Channel, ty string, data interface{}) error {
	//

	return nil
}
