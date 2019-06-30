package core

import (
	"testing"

	"github.com/gizo-network/gizo/helpers"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestGetBlock(t *testing.T) {
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
	block, err := NewBlock(*tree, "00000000000000000000000000000000000000", 1, 10, "74657374")
	assert.NoError(t, err)
	blockinfo := BlockInfo{
		Header:    block.GetHeader(),
		Height:    block.GetHeight(),
		TotalJobs: uint(len(block.GetNodes())),
		FileName:  block.fileStats().Name(),
		FileSize:  block.fileStats().Size(),
	}
	bBytes, err := json.Marshal(block)
	assert.NoError(t, err)
	biBytes, err := json.Marshal(blockinfo.GetBlock())
	assert.Equal(t, bBytes, biBytes)
	block.DeleteFile()
}
