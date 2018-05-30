package core

import (
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
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("test")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, 1, 10, "test")
	err = bc.AddBlock(block)
	assert.NoError(t, err)
	latest_height, err := bc.GetLatestHeight()
	assert.NoError(t, err)
	assert.Equal(t, 1, latest_height)
}

func TestVerify(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("test")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "test")
	bc.AddBlock(block)
	verify, err := bc.Verify()
	assert.NoError(t, err)
	assert.True(t, verify)
}

func TestGetBlockInfo(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("test")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "test")
	bc.AddBlock(block)
	blockinfo, err := bc.GetBlockInfo(block.GetHeader().GetHash())
	assert.NoError(t, err)
	assert.NotNil(t, blockinfo)
}

func TestGetBlocksWithinMinute(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("test")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "test")
	bc.AddBlock(block)
	lastMinute, err := bc.GetBlocksWithinMinute()
	assert.NoError(t, err)
	assert.NotNil(t, lastMinute)
}

func TestGetLatestHeight(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("test")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "test")
	bc.AddBlock(block)
	latest_height, err := bc.GetLatestHeight()
	assert.NoError(t, err)
	assert.NotNil(t, latest_height)
}

func TestFindJob(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("test")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "test")
	bc.AddBlock(block)
	f, err := bc.FindJob(nodes[4].GetJob().GetID())
	assert.NoError(t, err)
	assert.NotNil(t, f)
}

func TestGetBlockHashes(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	bc := CreateBlockChain("test")
	prevHash, err := bc.GetPrevHash()
	assert.NoError(t, err)
	nextHeight, err := bc.GetNextHeight()
	assert.NoError(t, err)
	block := NewBlock(*tree, prevHash, nextHeight, 10, "test")
	err = bc.AddBlock(block)
	assert.NoError(t, err)
	hashes, err := bc.GetBlockHashes()
	assert.NoError(t, err)
	assert.NotNil(t, hashes)
}

func TestCreateBlockChain(t *testing.T) {
	os.Setenv("ENV", "dev")
	RemoveDataPath()
	bc := CreateBlockChain("test")
	assert.NotNil(t, bc)
}
