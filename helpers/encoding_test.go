package helpers_test

import (
	"testing"

	"github.com/gizo-network/gizo/helpers"
	"github.com/stretchr/testify/assert"
)

func TestEncode64(t *testing.T) {
	assert.NotNil(t, helpers.Encode64([]byte("testing")))
}

func TestDecode64(t *testing.T) {
	expectedDecoded := []byte("testing")
	encoded := helpers.Encode64(expectedDecoded)
	decoded, err := helpers.Decode64(encoded)
	assert.NoError(t, err)
	assert.Equal(t, expectedDecoded, decoded)
}
