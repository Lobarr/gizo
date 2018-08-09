package helpers_test

import (
	"testing"

	"github.com/gizo-network/gizo/helpers"
	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	encrypt, err := helpers.Encrypt([]byte("test"), "test")
	assert.NoError(t, err)
	assert.NotNil(t, encrypt)
}

func TestDecrypt(t *testing.T) {
	test := []byte("test")
	passphrase := "test"
	encrypted, err := helpers.Encrypt(test, passphrase)
	assert.NoError(t, err)
	decrypted, err := helpers.Decrypt(encrypted, passphrase)
	assert.NoError(t, err)
	assert.Equal(t, decrypted, test)
}
