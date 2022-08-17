package payments

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

type jsonVoucher struct {
	ChannelId types.Destination
	Amount    *big.Int
	Signature state.Signature
}

// MarshalJSON returns a JSON representation of the Voucher
func (v Voucher) MarshalJSON() ([]byte, error) {
	jsonV := jsonVoucher{
		v.channelId, v.amount, v.signature,
	}
	return json.Marshal(jsonV)
}

// UnmarshalJSON populates the receiver with the
// json-encoded data
func (v *Voucher) UnmarshalJSON(data []byte) error {
	var jsonV jsonVoucher
	err := json.Unmarshal(data, &jsonV)
	if err != nil {
		return fmt.Errorf("error unmarshaling voucher data: %w", err)
	}

	v.channelId = jsonV.ChannelId
	v.amount = jsonV.Amount
	v.signature = jsonV.Signature

	return nil
}
