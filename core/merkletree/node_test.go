package merkletree_test

import (
	"testing"

	"github.com/gizo-network/gizo/crypt"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	n, _ := merkletree.NewNode(*j, nil, nil)
	assert.NotNil(t, n.GetHash(), "empty hash value")
	assert.NotNil(t, n, "returned empty node")
}

func TestIsLeaf(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	n, _ := merkletree.NewNode(*j, nil, nil)
	assert.True(t, n.IsLeaf())
}

func TestIsEmpty(t *testing.T) {
	n := merkletree.MerkleNode{}
	assert.True(t, n.IsEmpty())
}

func TestIsEqual(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	n, _ := merkletree.NewNode(*j, nil, nil)
	equal, err := n.IsEqual(*n)
	assert.NoError(t, err)
	assert.True(t, equal)
}

func TestMergeJobs(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	n, _ := merkletree.NewNode(*j, nil, nil)
	merged := merkletree.MergeJobs(*n, *n)
	assert.Equal(t, j.GetID()+j.GetID(), merged.GetID())
}
