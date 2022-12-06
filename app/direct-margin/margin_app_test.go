package directmargin

import (
	"testing"

	"github.com/statechannels/go-nitro/app"
)

func TestMarginAppType(t *testing.T) {
	var _ app.App = (*MarginApp)(nil)
}
