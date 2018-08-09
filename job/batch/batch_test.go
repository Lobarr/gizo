package batch_test

// import (
// 	"testing"

// 	"github.com/gizo-network/gizo/cache"
// 	"github.com/gizo-network/gizo/core"
// 	"github.com/gizo-network/gizo/core/merkletree"
// 	"github.com/gizo-network/gizo/crypt"
// 	"github.com/gizo-network/gizo/job"
// 	"github.com/gizo-network/gizo/job/batch"
// 	"github.com/gizo-network/gizo/job/queue"
// 	"github.com/stretchr/testify/assert"
// )

// func TestBatch(t *testing.T) {
// 	core.RemoveDataPath()
// 	pq := queue.NewJobPriorityQueue()
// 	priv, _ := crypt.GenKeys()
// 	j, _ := job.NewJob(`
// 	func Factorial(n){
// 	 if(n > 0){
// 	  result = n * Factorial(n-1)
// 	  return result
// 	 }
// 	 return 1
// 	}`, "Factorial", false, priv)
// 	j2, _ := job.NewJob(`
// 		func Test(){
// 			return "Testing"
// 		}
// 		`, "Test", false, priv)
// 	envs := job.NewEnvVariables(*job.NewEnv("Env", "Anko"), *job.NewEnv("By", "Lobarr"))
// 	exec1, err := job.NewExec([]interface{}{10}, 5, job.NORMAL, 0, 0, 0, 0, "", envs, "test")
// 	assert.NoError(t, err)
// 	exec2, err := job.NewExec([]interface{}{11}, 5, job.NORMAL, 0, 0, 0, 0, "", envs, "test")
// 	assert.NoError(t, err)
// 	exec3, err := job.NewExec([]interface{}{12}, 5, job.NORMAL, 0, 0, 0, 0, "", envs, "test")
// 	assert.NoError(t, err)
// 	exec4, err := job.NewExec([]interface{}{}, 5, job.NORMAL, 0, 0, 0, 0, "", envs, "test")
// 	assert.NoError(t, err)
// 	node1, err := merkletree.NewNode(*j, nil, nil)
// 	assert.NoError(t, err)
// 	node2, err := merkletree.NewNode(*j2, nil, nil)
// 	assert.NoError(t, err)
// 	nodes := []*merkletree.MerkleNode{node1, node2}
// 	tree, err := merkletree.NewMerkleTree(nodes)
// 	assert.NoError(t, err)
// 	bc := core.CreateBlockChain("74657374")
// 	prevHash, err := bc.GetPrevHash()
// 	assert.NoError(t, err)
// 	nextHeight, err := bc.GetNextHeight()
// 	assert.NoError(t, err)
// 	block, err := core.NewBlock(*tree, prevHash, nextHeight, 10, "74657374")
// 	assert.NoError(t, err)
// 	err = bc.AddBlock(block)
// 	assert.NoError(t, err)
// 	jr := job.NewRequest(j.GetID(), exec1, exec2, exec3)
// 	jr2 := job.NewRequest(j2.GetID(), exec4, exec4, exec4)
// 	batch, err := batch.NewBatch([]job.Request{*jr, *jr2}, bc, pq, cache.NewJobCacheNoWatch(bc))
// 	assert.NoError(t, err)
// 	go func() {
// 		for {
// 			if pq.GetPQ().Empty() == false {
// 				item := pq.Pop()
// 				exec := item.Job.Execute(item.GetExec(), "test")
// 				item.SetExec(exec)
// 				item.ResultsChan() <- item
// 			}
// 		}
// 	}()
// 	batch.Dispatch()
// 	assert.NotNil(t, batch.Result())
// }
