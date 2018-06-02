package core_test

import (
	"os"
	"testing"

	"github.com/gizo-network/gizo/core"
	"github.com/stretchr/testify/assert"
)

func TestGenesis(t *testing.T) {
	os.Setenv("ENV", "dev")
	assert.NotNil(t, core.GenesisBlock("74657374"))
}
