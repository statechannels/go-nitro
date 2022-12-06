package consensus

import (
	"testing"

	"github.com/statechannels/go-nitro/app"
)

func TestConsensusAppType(t *testing.T) {
	var _ app.App = (*ConsensusApp)(nil)
}
