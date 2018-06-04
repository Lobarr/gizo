package p2p_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/p2p"
)

func TestPeerMessage(t *testing.T) {
	priv, pub := crypt.GenKeys()
	_, pub2 := crypt.GenKeys()
	privBytes, err := hex.DecodeString(priv)
	assert.NoError(t, err)
	m := p2p.NewPeerMessage("TEST", []byte("test"), privBytes)
	assert.NotNil(t, m.GetSignature())
	verify, err := m.VerifySignature(pub)
	assert.NoError(t, err)
	assert.True(t, verify)
	verify, err = m.VerifySignature(pub2)
	assert.NoError(t, err)
	assert.False(t, verify)
}
