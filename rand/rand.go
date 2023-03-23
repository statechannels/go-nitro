// Package rand is a convenience wrapper aroung golang rand
// go math/rand is deterministic unless a random seed is provided
// see https://gobyexample.com/random-numbers
package rand

import (
	"math/rand"
	"time"
)

var r *rand.Rand

// getRandGenerator seeds a random number generator based on current time
func getRandGenerator() *rand.Rand {
	if r != nil {
		return r
	}
	source := rand.NewSource(time.Now().UnixNano())
	r = rand.New(source)
	return r
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
