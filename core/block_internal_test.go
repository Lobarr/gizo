package core

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/gizo-network/gizo/helpers"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestNewBlock(t *testing.T) {
	os.Setenv("ENV", "dev")
	InitializeDataPath()
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "74657374", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		assert.NoError(t, err)
		nodes = append(nodes, node)
	}
	tree, err := merkletree.NewMerkleTree(nodes)
	assert.NoError(t, err)
	prevHash := "00000000000000000000000000000000000000"
	testBlock, err := NewBlock(*tree, prevHash, 0, 5, "74657374")
	assert.NoError(t, err)

	assert.NotNil(t, testBlock, "returned empty tblock")
	assert.Equal(t, testBlock.Header.PrevBlockHash, prevHash, "prevhashes don't match")
	testBlock.DeleteFile()
}

func TestVerifyBlock(t *testing.T) {
	os.Setenv("ENV", "dev")
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "74657374", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, err := merkletree.NewMerkleTree(nodes)
	assert.NoError(t, err)
	prevHash := "00000000000000000000000000000000000000"
	testBlock, err := NewBlock(*tree, prevHash, 0, 5, "74657374")
	assert.NoError(t, err)

	verify, err := testBlock.VerifyBlock()
	assert.NoError(t, err)
	assert.True(t, verify, "block failed verification")

	testBlock.Header.setNonce(50)
	verify, err = testBlock.VerifyBlock()
	assert.NoError(t, err)
	assert.False(t, verify, "block passed verification")
	testBlock.DeleteFile()
}

func TestIsEmpty(t *testing.T) {
	os.Setenv("ENV", "dev")
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "74657374", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	prevHash := "00000000000000000000000000000000000000"
	testBlock, err := NewBlock(*tree, prevHash, 0, 5, "74657374")
	assert.NoError(t, err)
	b := Block{}
	assert.False(t, testBlock.IsEmpty())
	assert.True(t, b.IsEmpty())
	testBlock.DeleteFile()
}

func TestExport(t *testing.T) {
	os.Setenv("ENV", "dev")
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "74657374", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	prevHash := "00000000000000000000000000000000000000"
	testBlock, err := NewBlock(*tree, prevHash, 0, 5, "74657374")
	assert.NoError(t, err)
	assert.NotNil(t, testBlock.fileStats().Name())
	testBlock.DeleteFile()
}

func TestImport(t *testing.T) {
	os.Setenv("ENV", "dev")
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "74657374", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	prevHash := "00000000000000000000000000000000000000"
	testBlock, err := NewBlock(*tree, prevHash, 0, 5, "74657374")
	assert.NoError(t, err)

	empty := Block{}
	empty.Import(testBlock.Header.GetHash())
	testBlockBytes, err := helpers.Serialize(testBlock)
	assert.NoError(t, err)
	emptyBytes, err := helpers.Serialize(empty)
	assert.NoError(t, err)
	assert.JSONEq(t, string(testBlockBytes), string(emptyBytes))
	testBlock.DeleteFile()
}

func TestFileStats(t *testing.T) {
	os.Setenv("ENV", "dev")
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "74657374", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	prevHash := "00000000000000000000000000000000000000"
	testBlock, err := NewBlock(*tree, prevHash, 0, 5, "74657374")
	assert.NoError(t, err)
	assert.Equal(t, testBlock.fileStats().Name(), fmt.Sprintf(BlockFile, testBlock.Header.GetHash()))
	testBlock.DeleteFile()
}

func TestDeleteFile(t *testing.T) {
	os.Setenv("ENV", "dev")
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "74657374", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			assert.NoError(t, err)
		}
		nodes = append(nodes, node)
	}
	tree, _ := merkletree.NewMerkleTree(nodes)
	prevHash := "00000000000000000000000000000000000000"
	testBlock, err := NewBlock(*tree, prevHash, 0, 5, "74657374")
	assert.NoError(t, err)
	testBlock.DeleteFile()
	_, err = os.Stat(path.Join(BlockPathDev, fmt.Sprintf(BlockFile, testBlock.GetHeader().GetHash())))
	assert.True(t, os.IsNotExist(err))
}
