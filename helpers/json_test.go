package helpers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gizo-network/gizo/helpers"
	"github.com/gizo-network/gizo/p2p"
)

func TestDeserialize(t *testing.T) {
	v := p2p.NewVersion(1, 0, []string{})
	_v := p2p.Version{}
	err := helpers.Deserialize(v.Serialize(), &_v)
	assert.NoError(t, err)
	assert.Equal(t, v.GetBlocks(), _v.GetBlocks())
	assert.Equal(t, v.GetVersion(), _v.GetVersion())
	assert.Equal(t, v.GetHeight(), _v.GetHeight())
}

func TestSerialize(t *testing.T) {
	v := p2p.NewVersion(1, 0, []string{})
	b, err := helpers.Serialize(v)
	assert.NoError(t, err)
	assert.NotNil(t, b)
}
