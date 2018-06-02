package merkletree

import (
	"testing"

	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	node1, _ := NewNode(*j, nil, nil)
	node2, _ := NewNode(*j, nil, nil)

	parent, _ := merge(*node1, *node2)
	assert.NotNil(t, parent)
	assert.NotNil(t, parent.GetHash())
	assert.Equal(t, node1.GetHash(), parent.GetLeftNode().GetHash())
	assert.Equal(t, node2.GetHash(), parent.GetRightNode().GetHash())
}
