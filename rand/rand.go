// Package rand is a convenience wrapper aroung golang rand
package rand

import (
	"math/rand"
	"time"
)

var r *rand.Rand

// GetRandGenerator seeds a random number generator based on current time
func GetRandGenerator() *rand.Rand {
	if r != nil {
		return r
	}
	source := rand.NewSource(time.Now().UnixNano())
	r = rand.New(source)
	return r
}
