package parser

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/stretchr/testify/assert"
)

func TestCreateAllocations(t *testing.T) {
	allocations := make([]any, 2)
	allocation1 := map[string]any{
		"Destination":    testactors.Bob.Destination().String(),
		"Amount":         float64(123),
		"AllocationType": float64(0),
		"Metadata":       nil,
	}
	allocation2 := map[string]any{
		"Destination":    testactors.Alice.Destination().String(),
		"Amount":         float64(456),
		"AllocationType": float64(1),
		"Metadata":       []byte("test"),
	}
	allocations[0] = allocation1
	allocations[1] = allocation2

	result := createAllocations(allocations)

	amount1, _ := math.ParseBig256("123")
	amount2, _ := math.ParseBig256("456")
	assert.Equal(t, result[0], outcome.Allocation{
		Destination:    testactors.Bob.Destination(),
		Amount:         amount1,
		AllocationType: 0,
		Metadata:       nil,
	})
	assert.Equal(t, result[1], outcome.Allocation{
		Destination:    testactors.Alice.Destination(),
		Amount:         amount2,
		AllocationType: 1,
		Metadata:       []byte("test"),
	})
}

func TestCreateExit(t *testing.T) {

}
