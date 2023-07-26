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

// RecordData represents the data stored in the DHT record
type RecordData struct {
	PeerID    string `json:"peerid"`
	Signature string `json:"signature"`
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
	var recordData RecordData
	if err := json.Unmarshal(value, &recordData); err != nil {
		return errors.New("malformed record value")
	}

	// Check if the value can be parsed into a valid libp2p peer.ID
	_, err := peer.Decode(recordData.PeerID)
	if err != nil {
		return errors.New("invalid libp2p peer ID")
	}

	// Make sure the timestamp is not in the future or negative number
	if recordData.Timestamp > time.Time.Unix(time.Now()) || recordData.Timestamp < 0 {
		return errors.New("invalid timestamp")
	}

	//// Check the signature
	//sigBytes, err := hex.DecodeString(recordData.Signature)
	//if err != nil {
		//return errors.New("signature malformed")
	//}

	//addrBytes, err := hex.DecodeString(signingAddrStr[2:]) // remove "0x" prefix
	//if err != nil {
		//return errors.New("signature malformed")
	//}

	//sigPubKey, err := crypto.SigToPub(crypto.Keccak256([]byte(recordData.PeerID)), sigBytes)
	//if err != nil {
		//return errors.New("failed to extract public key from signature")
	//}

	//signatureAddr := crypto.PubkeyToAddress(*sigPubKey)
	//expectedAddr := common.BytesToAddress(addrBytes)

	//if signatureAddr != expectedAddr {
		//return errors.New("signature does not match address")
	//}

	// If no errors, the record is valid
	return nil
}

// Simply return the first record as the best record.
// In a more complex scenario, we could add logic to select the best record
// based on some criteria.
func (v stateChannelAddrToPeerIDValidator) Select(key string, values [][]byte) (int, error) {
	return 0, nil
}
