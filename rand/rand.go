// Package rand is a convenience wrapper aroung golang rand
// go math/rand is deterministic unless a random seed is provided
// see https://gobyexample.com/random-numbers
package rand

import (
	"math/rand"
	"time"
)

// getRandGenerator seeds a random number generator based on current time
func getRandGenerator() *rand.Rand {
	source := rand.NewSource(time.Now().UnixNano())
	return rand.New(source)
}

func Uint64() uint64 {
	return getRandGenerator().Uint64()
}

func Int63n(i int64) int64 {
	return getRandGenerator().Int63n(i)
}

func Int63() int64 {
	return getRandGenerator().Int63()
}
