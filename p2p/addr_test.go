package p2p_test

import (
	"testing"

	"github.com/gizo-network/gizo/p2p"
	"github.com/stretchr/testify/assert"
)

func TestParseAddr(t *testing.T) {
	addr, err := p2p.ParseAddr("gizo://304e301006072a8648ce3d020106052b81040021033a0004f14a7b28af6fdf3136779e0a82e618d5f481ab0377222e71c9473e552785eb4adedfb67030b15ba1d877f9e1a06dd8a58870dd1402da7e6e@99.233.0.99:9995")
	assert.NoError(t, err)
	assert.Equal(t, "304e301006072a8648ce3d020106052b81040021033a0004f14a7b28af6fdf3136779e0a82e618d5f481ab0377222e71c9473e552785eb4adedfb67030b15ba1d877f9e1a06dd8a58870dd1402da7e6e", addr["pub"])
	assert.Equal(t, "99.233.0.99", addr["ip"])
	assert.Equal(t, "9995", addr["port"])
}
