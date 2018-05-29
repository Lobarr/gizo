package core

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestNewBlockChain(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	bc := CreateBlockChain("test")
	assert.NotNil(t, bc)
}
func TestAddBlock(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
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
	tree := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("test")
	block := NewBlock(*tree, bc.GetPrevHash(), 1, 10, "test")
	err := bc.AddBlock(block)
	assert.NoError(t, err)
	assert.Equal(t, 1, int(bc.GetLatestHeight()))
}

func TestVerify(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
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
	tree := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("test")
	block := NewBlock(*tree, bc.GetPrevHash(), bc.GetNextHeight(), 10, "test")
	bc.AddBlock(block)
	assert.True(t, bc.Verify())
}

func TestGetBlockInfo(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
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
	tree := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("test")
	block := NewBlock(*tree, bc.GetPrevHash(), bc.GetNextHeight(), 10, "test")
	bc.AddBlock(block)
	blockinfo, err := bc.GetBlockInfo(block.GetHeader().GetHash())
	assert.NoError(t, err)
	assert.NotNil(t, blockinfo)
}

func TestGetBlocksWithinMinute(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
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
	tree := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("test")
	block := NewBlock(*tree, bc.GetPrevHash(), bc.GetNextHeight(), 10, "test")
	bc.AddBlock(block)
	assert.NotNil(t, bc.GetBlocksWithinMinute())
}

func TestGetLatestHeight(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
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
	tree := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("test")
	block := NewBlock(*tree, bc.GetPrevHash(), bc.GetNextHeight(), 10, "test")
	bc.AddBlock(block)
	assert.NotNil(t, bc.GetLatestHeight())
}

func TestFindJob(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
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
	tree := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("test")
	block := NewBlock(*tree, bc.GetPrevHash(), bc.GetNextHeight(), 10, "test")
	bc.AddBlock(block)
	f, err := bc.FindJob(node5.GetJob().GetID())
	assert.NoError(t, err)
	assert.NotNil(t, f)
}

func TestGetBlockHashes(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
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
	tree := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("test")
	block := NewBlock(*tree, bc.GetPrevHash(), bc.GetNextHeight(), 10, "test")
	bc.AddBlock(block)
	assert.NotNil(t, bc.GetBlockHashes())
}

func TestCreateBlockChain(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	bc := CreateBlockChain("test")
	assert.NotNil(t, bc)
}
