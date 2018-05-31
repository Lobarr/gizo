package cache_test

import (
	"os"
	"testing"

	"github.com/gizo-network/gizo/crypt"

	"github.com/stretchr/testify/assert"

	"github.com/gizo-network/gizo/cache"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/job"
)

func TestJobCache(t *testing.T) {
	os.Setenv("ENV", "dev")
	core.RemoveDataPath()
	priv, _ := crypt.GenKeys()
	j1, _ := job.NewJob(`
		func Factorial(n){
		 if(n > 0){
		  result = n * Factorial(n-1)
		  return result
		 }
		 return 1
		}`, "Factorial", false, priv)
	j1.AddExec(job.Exec{})
	j1.AddExec(job.Exec{})
	j1.AddExec(job.Exec{})
	j1.AddExec(job.Exec{})
	j2, _ := job.NewJob(`
			func Factorial(n){
			 if(n > 0){
			  result = n * Factorial(n-1)
			  return result
			 }
			 return 1
			}`, "Factorial", false, priv)
	j2.AddExec(job.Exec{})
	j2.AddExec(job.Exec{})
	j2.AddExec(job.Exec{})
	j2.AddExec(job.Exec{})
	j3, _ := job.NewJob(`
				func Factorial(n){
				 if(n > 0){
				  result = n * Factorial(n-1)
				  return result
				 }
				 return 1
				}`, "Factorial", false, priv)
	j3.AddExec(job.Exec{})
	j3.AddExec(job.Exec{})
	j3.AddExec(job.Exec{})
	j3.AddExec(job.Exec{})
	j3.AddExec(job.Exec{})
	j3.AddExec(job.Exec{})

	node1, _ := merkletree.NewNode(*j1, nil, nil)
	node2, _ := merkletree.NewNode(*j2, nil, nil)
	node3, _ := merkletree.NewNode(*j3, nil, nil)
	tree1, _ := merkletree.NewMerkleTree([]*merkletree.MerkleNode{node1, node3})
	tree2, _ := merkletree.NewMerkleTree([]*merkletree.MerkleNode{node2, node1})
	tree3, _ := merkletree.NewMerkleTree([]*merkletree.MerkleNode{node3, node2})
	bc := core.CreateBlockChain("74657374")
	latest_block, err := bc.GetLatestBlock()
	assert.NoError(t, err)
	next_height, err := bc.GetNextHeight()
	assert.NoError(t, err)
	blk1, err := core.NewBlock(*tree1, latest_block.GetHeader().GetHash(), next_height, 10, "74657374")
	assert.NoError(t, err)
	err = bc.AddBlock(blk1)
	assert.NoError(t, err)
	latest_block, err = bc.GetLatestBlock()
	assert.NoError(t, err)
	next_height, err = bc.GetNextHeight()
	assert.NoError(t, err)
	blk2, err := core.NewBlock(*tree2, latest_block.GetHeader().GetHash(), next_height, 10, "74657374")
	assert.NoError(t, err)
	err = bc.AddBlock(blk2)
	assert.NoError(t, err)
	latest_block, err = bc.GetLatestBlock()
	assert.NoError(t, err)
	next_height, err = bc.GetNextHeight()
	assert.NoError(t, err)
	blk3, err := core.NewBlock(*tree3, latest_block.GetHeader().GetHash(), next_height, 10, "74657374")
	assert.NoError(t, err)
	err = bc.AddBlock(blk3)
	assert.NoError(t, err)
	c := cache.NewJobCacheNoWatch(bc)
	cj1, err := c.Get(j1.GetID())
	t.Log(cj1)
	assert.NoError(t, err)
	assert.NotNil(t, cj1)
	cj2, err := c.Get(j2.GetID())
	t.Log(cj2)
	assert.NoError(t, err)
	assert.NotNil(t, cj2)
	cj3, err := c.Get(j3.GetID())
	t.Log(cj3)
	assert.NoError(t, err)
	assert.NotNil(t, cj3)
	assert.False(t, c.IsFull())
}
