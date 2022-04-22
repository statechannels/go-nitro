// Package store contains the interface for a go-nitro store.
package store // import "github.com/statechannels/go-nitro/client/engine/store"

import (
	"errors"
)

var ErrNoSuchObjective error = errors.New("store: no such objective")
var ErrNoSuchChannel error = errors.New("store: failed to find required channel data")
