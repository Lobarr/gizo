package difficulty_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gizo-network/gizo/crypt"

	"github.com/stretchr/testify/assert"

	"github.com/gizo-network/gizo/benchmark"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/core/difficulty"
	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/job"
)

func TestDifficulty(t *testing.T) {
	os.Setenv("ENV", "dev")
	core.RemoveDataPath()
	bc := core.CreateBlockChain("test")
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		assert.NoError(t, err)
		nodes = append(nodes, node)
	}

	tree, err := merkletree.NewMerkleTree(nodes)
	assert.NoError(t, err)
	fmt.Println("passed here")
	latest_block, err := bc.GetLatestBlock()
	assert.NoError(t, err)
	latest_height, err := bc.GetLatestHeight()
	assert.NoError(t, err)
	block := core.NewBlock(*tree, latest_block.GetHeader().GetHash(), uint64(latest_height), 10, "test")
	bc.AddBlock(block)
	d10 := benchmark.NewBenchmark(0.0115764096, 10)
	d11 := benchmark.NewBenchmark(0.13054728, 11)
	d12 := benchmark.NewBenchmark(0.0740971, 12)
	d13 := benchmark.NewBenchmark(0.28987127999999995, 13)
	d14 := benchmark.NewBenchmark(1.36593388, 14)
	d15 := benchmark.NewBenchmark(1.8645611, 15)
	d16 := benchmark.NewBenchmark(3.82076494, 16)
	d17 := benchmark.NewBenchmark(7.12966816, 17)
	d18 := benchmark.NewBenchmark(28.470944839999998, 18)
	d19 := benchmark.NewBenchmark(42.251310620000005, 19)
	benchmarks := []benchmark.Benchmark{d10, d11, d12, d13, d14, d15, d16, d17, d18, d19}
	assert.NotNil(t, difficulty.Difficulty(benchmarks, *bc))
}
