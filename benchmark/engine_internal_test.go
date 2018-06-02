package benchmark

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	os.Setenv("ENV", "dev")
	b := Engine{}
	assert.NotNil(t, b.block(10))
}
