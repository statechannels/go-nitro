package protocols

import (
	"encoding/json"

	"github.com/statechannels/go-nitro/types"
)

type jsonChainTransactionBase struct {
	ChannelId types.Destination
}

// MarshalJSON returns a JSON representation of the ChainTransactionBase
func (tx ChainTransactionBase) MarshalJSON() ([]byte, error) {
	jsonTx := jsonChainTransactionBase{tx.channelId}
	bytes, err := json.Marshal(jsonTx)

	if err != nil {
		return []byte{}, err
	}

	return bytes, err
}

// UnmarshalJSON unmarshals the passed JSON into a ChainTransactionBase, implementing the Unmarshaler interface.
func (tx *ChainTransactionBase) UnmarshalJSON(data []byte) error {
	var jsonTx jsonChainTransactionBase
	err := json.Unmarshal(data, &jsonTx)

	if err != nil {
		return err
	}

	tx.channelId = jsonTx.ChannelId
	return nil

}
