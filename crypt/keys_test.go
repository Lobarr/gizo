package crypt_test

import (
	"testing"

	"github.com/gizo-network/gizo/crypt"
	"github.com/stretchr/testify/assert"
)

func TestGenKeys(t *testing.T) {
	priv, pub := crypt.GenKeys()
	assert.NotNil(t, priv)
	assert.NotNil(t, pub)
}
