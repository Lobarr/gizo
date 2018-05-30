package solo_test

import (
	"encoding/hex"
	"testing"

	"github.com/gizo-network/gizo/cache"

	"github.com/stretchr/testify/assert"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job"
	"github.com/gizo-network/gizo/job/queue"
	"github.com/gizo-network/gizo/job/solo"
)

func TestSolo(t *testing.T) {
	core.RemoveDataPath()
	priv, pub := crypt.GenKeys()
	pq := queue.NewJobPriorityQueue()
	j, _ := job.NewJob(`	
	func Factorial(n){
	 if(n > 0){
	  result = n * Factorial(n-1)
	  return result
	 }
	 return 1
	}`, "Factorial", false, hex.EncodeToString(priv))
	envs := job.NewEnvVariables(*job.NewEnv("Env", "Anko"), *job.NewEnv("By", "Lobarr"))
	exec1, err := job.NewExec([]interface{}{2}, 5, job.NORMAL, 0, 0, 0, 0, hex.EncodeToString(pub), envs, "test")
	assert.NoError(t, err)
	node1 := merkletree.NewNode(*j, nil, nil)
	nodes := []*merkletree.MerkleNode{node1}
	tree := merkletree.NewMerkleTree(nodes)
	bc := core.CreateBlockChain("74657374")
	block, err := core.NewBlock(*tree, bc.GetLatestBlock().GetHeader().GetHash(), bc.GetNextHeight(), 10, "74657374")
	assert.NoError(t, err)
	err = bc.AddBlock(block)
	assert.NoError(t, err)
	s := solo.NewSolo(*job.NewJobRequestSingle(j.GetID(), exec1), bc, pq, cache.NewJobCacheNoWatch(bc))
	s.Dispatch()
	assert.NotNil(t, s.Result())
}
