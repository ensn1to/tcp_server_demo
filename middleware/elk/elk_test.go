package elk

import (
	"testing"

	"gotest.tools/assert"
)

func TestNewElkLogger(t *testing.T) {
	Logger = NewElkLogger("101.43.84.106", 4560, 5)
	Logger.Infof("connect success")

	assert.Equal(t, false, true)
}
