package core

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestNewBlock(t *testing.T) {
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
	prevHash := "00000000000000000000000000000000000000"
	testBlock := NewBlock(*tree, prevHash, 0, 5, "test")

	assert.NotNil(t, testBlock, "returned empty tblock")
	assert.Equal(t, testBlock.Header.PrevBlockHash, prevHash, "prevhashes don't match")
	testBlock.DeleteFile()
}

func TestVerifyBlock(t *testing.T) {
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
	prevHash := []byte("00000000000000000000000000000000000000")
	testBlock := NewBlock(*tree, prevHash, 0, 5, "test")

	assert.True(t, testBlock.VerifyBlock(), "block failed verification")

	testBlock.Header.setNonce(50)
	assert.False(t, testBlock.VerifyBlock(), "block passed verification")
	testBlock.DeleteFile()
}

func TestSerialize(t *testing.T) {
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
	prevHash := []byte("00000000000000000000000000000000000000")
	testBlock := NewBlock(*tree, prevHash, 0, 5, "test")
	stringified := testBlock.Serialize()
	assert.NotEmpty(t, stringified)
	testBlock.DeleteFile()
}

func TestDeserializeBlock(t *testing.T) {
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
	prevHash := []byte("00000000000000000000000000000000000000")
	testBlock := NewBlock(*tree, prevHash, 0, 5, "test")
	stringified := testBlock.Serialize()
	unmarshaled, err := DeserializeBlock(stringified)
	confirm := unmarshaled.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, stringified, confirm)
	testBlock.DeleteFile()
}

func TestIsEmpty(t *testing.T) {
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
	prevHash := []byte("00000000000000000000000000000000000000")
	testBlock := NewBlock(*tree, prevHash, 0, 5, "test")
	b := Block{}
	assert.False(t, testBlock.IsEmpty())
	assert.True(t, b.IsEmpty())
	testBlock.DeleteFile()
}

func TestExport(t *testing.T) {
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
	prevHash := []byte("00000000000000000000000000000000000000")
	testBlock := NewBlock(*tree, prevHash, 0, 5, "test")
	assert.NotNil(t, testBlock.fileStats().Name())
	testBlock.DeleteFile()
}

func TestImport(t *testing.T) {
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
	prevHash := []byte("00000000000000000000000000000000000000")
	testBlock := NewBlock(*tree, prevHash, 0, 5, "test")

	empty := Block{}
	empty.Import(testBlock.Header.GetHash())
	testBlockBytes := testBlock.Serialize()
	emptyBytes := empty.Serialize()
	assert.JSONEq(t, string(testBlockBytes), string(emptyBytes))
	testBlock.DeleteFile()
}

func TestFileStats(t *testing.T) {
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
	prevHash := []byte("00000000000000000000000000000000000000")
	testBlock := NewBlock(*tree, prevHash, 0, 5, "test")
	assert.Equal(t, testBlock.fileStats().Name(), fmt.Sprintf(BlockFile, hex.EncodeToString(testBlock.Header.GetHash())))
	testBlock.DeleteFile()
}
