package core

import (
	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/helpers"
	"github.com/gizo-network/gizo/job"
)

//! modify on job engine creation

//GenesisBlock returns genesis block
func GenesisBlock(by string) *Block {
	logger := helpers.Logger()
	logger.Info("Core: creating Genesis Block")
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func Genesis(){return 1+1}", "Genesis", false, priv)
	node, err := merkletree.NewNode(*j, nil, nil)
	if err != nil {
		logger.Fatal(err)
	}
	tree := merkletree.MerkleTree{
		Root:      node.GetHash(),
		LeafNodes: []*merkletree.MerkleNode{node},
	}
	prevHash := "47656e65736973"
	block, err := NewBlock(tree, prevHash, 0, 10, by)
	if err != nil {
		logger.Fatal(err)
	}
	return block
}
