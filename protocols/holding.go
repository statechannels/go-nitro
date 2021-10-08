package protocols

import (
	"errors"
	"math/big"
)

var ErrIncorrectChannelId = errors.New("incorrect channel id")

// Mirrors the on-chain holdings for each channel and each asset
type Holding interface {
	ChannelId() string
	Asset() string
	Amount() big.Int
}
