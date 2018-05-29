package merkletree_test

import (
	"encoding/hex"
	"testing"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	assert.NotNil(t, merkletree.ErrTooMuchLeafNodes)
	assert.NotNil(t, merkletree.ErrTreeRebuildAttempt)
	assert.NotNil(t, merkletree.ErrTreeNotBuilt)
	assert.NotNil(t, merkletree.ErrLeafNodesEmpty)
}

func TestBuild(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, hex.EncodeToString(priv))
	nodes_8 := []*merkletree.MerkleNode{}
	for i := 0; i < 8; i++ {
		nodes_8 = append(nodes_8, merkletree.NewNode(*j, nil, nil))
	}

	tree := merkletree.MerkleTree{LeafNodes: nodes_8}
	assert.Equal(t, "", tree.GetRoot())
	tree.Build()
	assert.NotNil(t, tree.GetRoot())
	assert.Error(t, tree.Build())

	nodes_max := []*merkletree.MerkleNode{}
	for i := 0; i < merkletree.MaxTreeJobs+1; i++ {
		nodes_max = append(nodes_8, merkletree.NewNode(*j, nil, nil))
	}
	tree2 := merkletree.MerkleTree{LeafNodes: nodes_max}
	// assert.Error(t, tree2.Build())
	t.Fatal(len(tree2.GetLeafNodes()))
	t.Fatal(tree2.Build())
}

func TestNewMerkleTree(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, hex.EncodeToString(priv))
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}

	tree, _ := merkletree.NewMerkleTree(nodes)
	assert.NotNil(t, tree.GetRoot())
	assert.NotNil(t, tree.GetLeafNodes())
}

func TestVerifyTree(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, hex.EncodeToString(priv))
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}

	tree, _ := merkletree.NewMerkleTree(nodes)
	assert.True(t, tree.VerifyTree())

	tree.SetLeafNodes(tree.GetLeafNodes()[4:])
	assert.False(t, tree.VerifyTree())
}

func TestSearchNode(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, hex.EncodeToString(priv))
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}

	tree, _ := merkletree.NewMerkleTree(nodes)
	f, err := tree.SearchNode(node5.GetHash())
	assert.NoError(t, err)
	assert.NotNil(t, f)
}

func TestSearchJob(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, hex.EncodeToString(priv))
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}

	tree, _ := merkletree.NewMerkleTree(nodes)
	f, err := tree.SearchJob(j.GetID())
	assert.NoError(t, err)
	assert.NotNil(t, f)
}
