package chain_test

// import (
// 	"encoding/hex"
// 	"testing"

// 	"github.com/stretchr/testify/assert"

// 	"github.com/gizo-network/gizo/cache"
// 	"github.com/gizo-network/gizo/core"
// 	"github.com/gizo-network/gizo/core/merkletree"
// 	"github.com/gizo-network/gizo/crypt"
// 	"github.com/gizo-network/gizo/job"
// 	"github.com/gizo-network/gizo/job/chain"
// 	"github.com/gizo-network/gizo/job/queue"
// )

// func TestChain(t *testing.T) {
// 	core.RemoveDataPath()
// 	priv, _ := crypt.GenKeys()
// 	pq := queue.NewJobPriorityQueue()
// 	j, _ := job.NewJob(`
// 	func Factorial(n){
// 	 if(n > 0){
// 	  result = n * Factorial(n-1)
// 	  return result
// 	 }
// 	 return 1
// 	}`, "Factorial", false, hex.EncodeToString(priv))
// 	j2, _ := job.NewJob(`
// 		func Test(){
// 			return "Testing"
// 		}
// 		`, "Test", false, hex.EncodeToString(priv))
// 	envs := job.NewEnvVariables(*job.NewEnv("Env", "Anko"), *job.NewEnv("By", "Lobarr"))
// 	exec1, err := job.NewExec([]interface{}{10}, 5, job.NORMAL, 0, 0, 0, 0, "", envs, "test")
// 	assert.NoError(t, err)
// 	exec2, err := job.NewExec([]interface{}{11}, 5, job.NORMAL, 0, 0, 0, 0, "", envs, "test")
// 	assert.NoError(t, err)
// 	exec3, err := job.NewExec([]interface{}{12}, 5, job.NORMAL, 0, 0, 0, 0, "", envs, "test")
// 	assert.NoError(t, err)
// 	exec4, err := job.NewExec([]interface{}{}, 5, job.NORMAL, 0, 0, 0, 0, "", envs, "test")
// 	assert.NoError(t, err)
// 	node1 := merkletree.NewNode(*j, nil, nil)
// 	node2 := merkletree.NewNode(*j2, nil, nil)
// 	nodes := []*merkletree.MerkleNode{node1, node2}
// 	tree := merkletree.NewMerkleTree(nodes)
// 	bc := core.CreateBlockChain("74657374")
// 	block := core.NewBlock(*tree, bc.GetLatestBlock().GetHeader().GetHash(), bc.GetLatestHeight()+1, 10, "74657374")
// 		err = bc.AddBlock(block) 	assert.NoError(t, err)(block)
// 	jr := job.NewJobRequest(j.GetID(), exec1, exec2, exec3)
// 	jr2 := job.NewJobRequest(j2.GetID(), exec4, exec4, exec4, exec4, exec4)
// 	chain, err := chain.NewChain([]job.JobRequest{*jr, *jr2}, bc, pq, cache.NewJobCacheNoWatch(bc))
// 	assert.NoError(t, err)
// 	chain.Dispatch()
// 	assert.NotNil(t, chain.Result())
// }
