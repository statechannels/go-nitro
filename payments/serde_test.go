package payments

import (
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

func TestSerde(t *testing.T) {
	someVoucher := Voucher{types.Destination{1}, big.NewInt(2), crypto.Signature{
		R: common.Hex2Bytes(`704b3afcc6e702102ca1af3f73cf3b37f3007f368c40e8b81ca823a65740a053`),
		S: common.Hex2Bytes(`14040ad4c598dbb055a50430142a13518e1330b79d24eed86fcbdff1a7a95589`),
		V: byte(0),
	}}

	someVoucherJson := `{"ChannelId":"0x0100000000000000000000000000000000000000000000000000000000000000","Amount":2,"Signature":{"R":"cEs6/MbnAhAsoa8/c887N/MAfzaMQOi4HKgjpldAoFM=","S":"FAQK1MWY27BVpQQwFCoTUY4TMLedJO7Yb8vf8aepVYk=","V":0}}`

	t.Run("Marshalling", func(t *testing.T) {
		got, err := json.Marshal(someVoucher)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != someVoucherJson {
			t.Fatalf("incorrect json marshaling, expected %v got \n%v", someVoucherJson, string(got))
		}
	})

	t.Run("Unmarshalling", func(t *testing.T) {
		got := Voucher{}
		err := json.Unmarshal([]byte(someVoucherJson), &got)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, someVoucher) {
			t.Fatalf("incorrect json unmarshaling, expected \n%+v got \n%+v", someVoucher, got)
		}
	})
}
