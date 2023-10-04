package p2pms

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	DHT_RECORD_PREFIX      = "/" + DHT_NAMESPACE + "/"
	DHT_NAMESPACE          = "scaddr"
	DHT_RECORD_MAX_AGE     = 24 * time.Hour
	DHT_REPUBLSIH_INTERVAL = 4 * time.Hour
)

type stateChannelAddrToPeerIDValidator struct{}

// dhtRecord represents the data stored in the DHT record
type dhtRecord struct {
	Data      dhtData
	Signature []byte
}

type dhtData struct {
	SCAddr    string // state channel address
	PeerID    string
	Timestamp int64 // Unix timestamp (seconds since January 1, 1970)
}

func (v stateChannelAddrToPeerIDValidator) Validate(key string, value []byte) error {
	// Trim the DHT_RECORD_PREFIX from the key to get the state channel address
	signingAddrStr := strings.TrimPrefix(key, DHT_RECORD_PREFIX)

	// Check if it's a valid state channel address
	if !common.IsHexAddress(signingAddrStr) {
		return errors.New("invalid state channel address used for key")
	}

	var dhtRecord dhtRecord
	if err := json.Unmarshal(value, &dhtRecord); err != nil {
		return errors.New("malformed record value")
	}

	if common.HexToAddress(dhtRecord.Data.SCAddr) != common.HexToAddress(signingAddrStr) {
		return errors.New("record key does not match state channel address")
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

// Choose the most recent record if we receive multiple records for the same key
func (v stateChannelAddrToPeerIDValidator) Select(key string, values [][]byte) (int, error) {
	var mostRecentIndex int
	var mostRecentTimestamp int64

	for i, value := range values {
		var record dhtRecord
		err := json.Unmarshal(value, &record)
		if err != nil {
			return -1, fmt.Errorf("error unmarshalling record: %w", err)
		}

		if record.Data.Timestamp > mostRecentTimestamp {
			mostRecentIndex = i
			mostRecentTimestamp = record.Data.Timestamp
		}
	}

	return mostRecentIndex, nil
}
