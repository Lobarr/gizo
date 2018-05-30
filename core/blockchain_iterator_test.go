package core

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	bc := CreateBlockChain("74657374")
	bci := bc.iterator()
	next, err := bci.Next()
	assert.NoError(t, err)
	assert.NotNil(t, next)
}

func TestNextBlockinf0(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	bc := CreateBlockChain("74657374")
	bci := bc.iterator()
	next, err := bci.NextBlockinfo()
	assert.NoError(t, err)
	assert.NotNil(t, next)
}
