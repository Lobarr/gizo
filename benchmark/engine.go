package benchmark

import (
	"encoding/hex"
	"math/rand"
	"sync"
	"time"

	"github.com/gizo-network/gizo/helpers"

	"github.com/gizo-network/gizo/crypt"

	"github.com/kpango/glg"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/job"
)

//Engine hold's an array of benchmarks and a score of the node
type Engine struct {
	Data   []Benchmark
	Score  float64
	logger *glg.Glg
}

func (b *Engine) setScore(s float64) {
	b.Score = s
}

//GetScore returns the score
func (b Engine) GetScore() float64 {
	return b.Score
}

func (b *Engine) addBenchmark(benchmark Benchmark) {
	b.Data = append(b.Data, benchmark)
}

//GetData returns an array of benchmarks
func (b Engine) GetData() []Benchmark {
	return b.Data
}

//returns a block with mock data
func (b Engine) block(difficulty uint8) *core.Block {
	//random data
	priv, _ := crypt.GenKeys()
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, hex.EncodeToString(priv))
	node1 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node2 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node3 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node4 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node5 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node6 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node7 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node8 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node9 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node10 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node11 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node12 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node13 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node14 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node15 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node16 := merkletree.NewNode(*j, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})

	tree := merkletree.NewMerkleTree([]*merkletree.MerkleNode{node1, node2, node3, node4, node5, node6, node7, node8, node9, node10, node11, node12, node13, node14, node15, node16})
	return core.NewBlock(*tree, []byte("TestingPreviousHash"), uint64(rand.Int()), difficulty, "benchmark-engine")
}

// Run executes the benchmark engine
func (b *Engine) run() {
	// mines blocks from difficulty to when it takes longer than 60 seconds to mine block
	b.logger.Log("Benchmarking node")
	done := false
	close := make(chan struct{})
	difficulty := 10 //! difficulty starts at 10
	go func() {
		var wg sync.WaitGroup
		for done == false {
			for i := 0; i < 3; i++ { //? runs 3 difficulties in parralel to increase benchmark speed
				wg.Add(1)
				go func(myDifficulty int) { //? mines block of myDifficulty 5 times and adds average time and difficulty
					var avg []float64
					var mu sync.Mutex
					var mineWG sync.WaitGroup
					for j := 0; j < 5; j++ {
						mineWG.Add(1)
						go func() {
							start := time.Now()
							block := b.block(uint8(myDifficulty))
							end := time.Now()
							block.DeleteFile()
							diff := end.Sub(start).Seconds()
							mu.Lock()
							avg = append(avg, diff)
							mu.Unlock()
							mineWG.Done()
						}()
					}
					mineWG.Wait()
					var avgSum float64
					for _, val := range avg {
						avgSum += val
					}
					average := avgSum / float64(len(avg))
					if average > float64(60) { // a minute
						done = true
						close <- struct{}{}
					} else {
						if done == true {
							wg.Done()
							return
						}
						b.addBenchmark(NewBenchmark(average, uint8(myDifficulty)))
					}
					wg.Done()
				}(difficulty)
				difficulty++
			}
			wg.Wait()
		}
	}()

	<-close //wait for close signal

	score := float64(b.GetData()[len(b.GetData())-1].GetDifficulty()) - 10 //! 10 is subtracted to allow the score start from 1 since difficulty starts at 10
	scoreDecimal := 1 - b.GetData()[len(b.GetData())-1].GetAvgTime()/100   // determine decimal part of score
	b.setScore(score + scoreDecimal)
	b.logger.Info("Benchmark: Benchmark done")
}

//NewEngine returns a Engine with benchmarks run
func NewEngine() Engine {
	b := Engine{logger: helpers.Logger()}
	b.run()
	return b
}
