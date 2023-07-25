package reverseproxy

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/types"
)

// keyValCollection is an interface that allows us to set, get and delete key value pairs.
// It is used to abstract away the differences between http.Header and url.Values.
type keyValCollection interface {
	Set(key, value string)
	Get(key string) string
	Del(key string)
}

// addVoucher takes in a voucher and adds it to the given keyValCollection.
// It prefixes the keys with the given prefix.
func addVoucher(v payments.Voucher, col keyValCollection, prefix string) {
	col.Set(prefix+CHANNEL_ID_VOUCHER_PARAM, v.ChannelId.String())
	col.Set(prefix+AMOUNT_VOUCHER_PARAM, v.Amount.String())
	col.Set(prefix+SIGNATURE_VOUCHER_PARAM, v.Signature.ToHexString())
}

// parseVoucher takes in an a keyValCollection  parses out a voucher.
func parseVoucher(col keyValCollection, prefix string) (payments.Voucher, error) {
	rawChId := col.Get(prefix + CHANNEL_ID_VOUCHER_PARAM)
	if rawChId == "" {
		return payments.Voucher{}, fmt.Errorf("missing channel ID")
	}
	rawAmt := col.Get(prefix + AMOUNT_VOUCHER_PARAM)
	if rawAmt == "" {
		return payments.Voucher{}, fmt.Errorf("missing amount")
	}
	rawSignature := col.Get(prefix + SIGNATURE_VOUCHER_PARAM)
	if rawSignature == "" {
		return payments.Voucher{}, fmt.Errorf("missing signature")
	}

	amount := big.NewInt(0)
	amount.SetString(rawAmt, 10)

	v := payments.Voucher{
		ChannelId: types.Destination(common.HexToHash(rawChId)),
		Amount:    amount,
		Signature: crypto.SplitSignature(hexutil.MustDecode(rawSignature)),
	}
	return v, nil
}

// removeVoucherParams removes the voucher parameters from the request URL.
func removeVoucher(col keyValCollection, prefix string) {
	col.Del(prefix + CHANNEL_ID_VOUCHER_PARAM)
	col.Del(prefix + AMOUNT_VOUCHER_PARAM)
	col.Del(prefix + SIGNATURE_VOUCHER_PARAM)
}
