package helpers_test

import (
	"testing"

	"github.com/gizo-network/gizo/helpers"
	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	assert.NotNil(t, helpers.Encrypt([]byte("test"), "test"))
}

func TestDecrypt(t *testing.T) {
	test := []byte("test")
	passphrase := "test"
	encrypted := helpers.Encrypt(test, passphrase)
	decrypted, err := helpers.Decrypt(encrypted, passphrase)
	assert.NoError(t, err)
	assert.Equal(t, decrypted, test)
}
