package p2pms

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	DHT_RECORD_PREFIX = "/scaddr/"
)

type stateChannelAddrToPeerIDValidator struct{}

// DhtRecord represents the data stored in the DHT record
type DhtRecord struct {
	Data      DhtData `json:"data"`
	Signature []byte  `json:"signature"`
}

type DhtData struct {
	SCAddr    string `json:"scaddr"` // state channel address
	PeerID    string `json:"peerid"`
	Timestamp int64  `json:"timestamp"` // Unix timestamp (seconds since January 1, 1970)
}

func (v stateChannelAddrToPeerIDValidator) Validate(key string, value []byte) error {
	// Trim the DHT_RECORD_PREFIX from the key to get the state channel address
	signingAddrStr := strings.TrimPrefix(key, DHT_RECORD_PREFIX)

	// Check if it's a valid state channel address
	if !common.IsHexAddress(signingAddrStr) {
		return errors.New("invalid state channel address used for key")
	}

	// Parse the value into a RecordData object
	var dhtRecord DhtRecord
	if err := json.Unmarshal(value, &dhtRecord); err != nil {
		return errors.New("malformed record value")
	}

	// Make sure the timestamp is not in the future or negative number
	if dhtRecord.Data.Timestamp > time.Time.Unix(time.Now()) || dhtRecord.Data.Timestamp < 0 {
		return errors.New("invalid timestamp")
	}

	// Check if the value can be parsed into a valid libp2p peer.ID
	peerId, err := peer.Decode(dhtRecord.Data.PeerID)
	if err != nil {
		return errors.New("invalid libp2p peer ID")
	}

	pubKey, err := peerId.ExtractPublicKey()
	if err != nil {
		return err
	}

	dataBytes, err := json.Marshal(dhtRecord.Data)
	if err != nil {
		return err
	}

	// Check the signature to ensure it is the signed hash of dataBytes
	valid, err := pubKey.Verify(dataBytes, dhtRecord.Signature)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("invalid signature")
	}

	return nil
}

// Simply return the first record as the best record.
// In a more complex scenario, we could add logic to select the best record
// based on some criteria.
func (v stateChannelAddrToPeerIDValidator) Select(key string, values [][]byte) (int, error) {
	return 0, nil
}
