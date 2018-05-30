package difficulty

import (
	"math/rand"
	"time"

	"github.com/gizo-network/gizo/helpers"

	"github.com/gizo-network/gizo/benchmark"
	"github.com/gizo-network/gizo/core"
)

//Blockrate block per minute
const Blockrate = 15

//Difficulty returns a difficulty based on the blockrate and the number of blocks in the last minute
func Difficulty(benchmarks []benchmark.Benchmark, bc core.BlockChain) int {
	blocks, err := bc.GetBlocksWithinMinute()
	if err != nil {
		helpers.Logger().Fatal(err)
	}
	latest := len(blocks)
	for _, val := range benchmarks {
		rate := float64(time.Minute) / val.GetAvgTime() // how many can be run in minute
		if int(rate)+latest <= Blockrate {
			return int(val.GetDifficulty())
		}
	}
	return rand.Intn(int(benchmarks[len(benchmarks)-1].GetDifficulty())) // returns random difficulty between benchmark limit
}
