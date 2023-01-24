package parser

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestI2Uint8(t *testing.T) {
	assert.Equal(t, uint8(42), i2Uint8(int64(42)))
	assert.Equal(t, uint8(42), i2Uint8(float64(42)))
	assert.Equal(t, uint8(42), i2Uint8(uint8(42)))
	assert.Equal(t, uint8(42), i2Uint8(uint32(42)))
	assert.Panics(t, func() {
		i2Uint8(int8(42))
	})
	assert.Panics(t, func() {
		i2Uint8(int32(42))
	})
	assert.Panics(t, func() {
		i2Uint8(uint64(42))
	})
}

func TestI2Uint32(t *testing.T) {
	assert.Equal(t, uint32(42), i2Uint32(uint32(42)))
	assert.Equal(t, uint32(42), i2Uint32(float64(42)))
	assert.Equal(t, uint32(42), i2Uint32(int64(42)))
	assert.Panics(t, func() {
		i2Uint32(int32(42))
	})
	assert.Panics(t, func() {
		i2Uint32(uint64(42))
	})
}

func TestI2Uint64(t *testing.T) {
	assert.Equal(t, uint64(42), i2Uint64(int64(42)))
	assert.Equal(t, uint64(42), i2Uint64(float64(42)))
	assert.Equal(t, uint64(42), i2Uint64(uint64(42)))
	assert.Panics(t, func() {
		i2Uint64(int32(42))
	})
	assert.Panics(t, func() {
		i2Uint64(uint32(42))
	})
}

func TestI2Uint256(t *testing.T) {
	assert.Equal(t, big.NewInt(42), i2Uint256("42"))
	assert.Panics(t, func() {
		i2Uint256(int32(42))
	})
	assert.Equal(t, big.NewInt(42), i2Uint256(float64(42)))
}

func TestToByteArray(t *testing.T) {
	assert.Equal(t, []byte(nil), toByteArray(nil))
	assert.Equal(t, []byte("test"), toByteArray([]byte("test")))
}
