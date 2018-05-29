package core

import (
	"encoding/hex"
	"testing"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestGetBlock(t *testing.T) {
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
	block := NewBlock(*tree, []byte("00000000000000000000000000000000000000"), 1, 10, "test")
	blockinfo := BlockInfo{
		Header:    block.GetHeader(),
		Height:    block.GetHeight(),
		TotalJobs: uint(len(block.GetNodes())),
		FileName:  block.fileStats().Name(),
		FileSize:  block.fileStats().Size(),
	}
	assert.Equal(t, block.Serialize(), blockinfo.GetBlock().Serialize())
	block.DeleteFile()
}
