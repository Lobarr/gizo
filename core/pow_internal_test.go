package core

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestPrepareData(t *testing.T) {
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

	block := &Block{
		Header: BlockHeader{
			Timestamp:     time.Now().Unix(),
			PrevBlockHash: []byte("00000000000000000000000000000000000000"),
			MerkleRoot:    tree.GetRoot(),
			Difficulty:    big.NewInt(int64(10)),
		},
		Jobs:   tree.GetLeafNodes(),
		Height: 0,
	}
	pow := NewPOW(block)
	assert.NotNil(t, pow)
	assert.NotNil(t, pow.prepareData(5))
}

func TestNewPOW(t *testing.T) {
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

	block := &Block{
		Header: BlockHeader{
			Timestamp:     time.Now().Unix(),
			PrevBlockHash: []byte("00000000000000000000000000000000000000"),
			MerkleRoot:    tree.GetRoot(),
			Difficulty:    big.NewInt(int64(10)),
		},
		Jobs:   tree.GetLeafNodes(),
		Height: 0,
	}
	pow := NewPOW(block)
	assert.NotNil(t, pow)
}

func TestRun(t *testing.T) {
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

	block := &Block{
		Header: BlockHeader{
			Timestamp:     time.Now().Unix(),
			PrevBlockHash: []byte("00000000000000000000000000000000000000"),
			MerkleRoot:    tree.GetRoot(),
			Difficulty:    big.NewInt(int64(10)),
		},
		Jobs:   tree.GetLeafNodes(),
		Height: 0,
	}
	pow := NewPOW(block)
	assert.NotNil(t, pow)
	pow.run()
	assert.NotNil(t, block.GetHeader().GetHash())
	assert.NotNil(t, block.GetHeader().GetNonce())
}

func TestValidate(t *testing.T) {
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

	block := &Block{
		Header: BlockHeader{
			Timestamp:     time.Now().Unix(),
			PrevBlockHash: []byte("00000000000000000000000000000000000000"),
			MerkleRoot:    tree.GetRoot(),
			Difficulty:    big.NewInt(int64(10)),
		},
		Jobs:   tree.GetLeafNodes(),
		Height: 1,
	}
	pow := NewPOW(block)
	pow.run()
	assert.NotNil(t, pow)
	assert.True(t, pow.Validate())
}
