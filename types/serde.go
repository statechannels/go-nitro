package types

import "github.com/ethereum/go-ethereum/common"

// MarshalText encodes the receiver into UTF-8-encoded text and returns the result.
//
// This makes Destination an encoding.TextMarshaler, meaning it can be safely used as the key in a map which will be encoded into json.
func (d Destination) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

// UnmarshalText unmarshals the supplied text (assumed to be a valid marshaling) into the receiver.
//
// This makes Destination an encoding.TextUnmarshaler.
func (d *Destination) UnmarshalText(text []byte) error {
	*d = Destination(common.HexToHash(string(text)))
	return nil
}
