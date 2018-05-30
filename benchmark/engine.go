package benchmark

import (
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
	j, _ := job.NewJob("func test(){return 1+1}", "test", false, priv)
	nodes := []*merkletree.MerkleNode{}
	for i := 0; i < 16; i++ {
		node, err := merkletree.NewNode(*j, nil, nil)
		if err != nil {
			glg.Fatal(err)
		}
		nodes = append(nodes, node)
	}
	tree, err := merkletree.NewMerkleTree(nodes)
	if err != nil {
		glg.Fatal(err)
	}
	return core.NewBlock(*tree, "47656e65736973", uint64(rand.Int()), difficulty, "benchmark-engine")
}

// Run executes the benchmark engine
func (b *Engine) run() {
	// mines blocks from difficulty to when it takes longer than 60 seconds to mine block
	b.logger.Log("Benchmark: starting benchmark")
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
	b.logger.Info("Benchmark: benchmark done")
}

//NewEngine returns a Engine with benchmarks run
func NewEngine() Engine {
	b := Engine{logger: helpers.Logger()}
	b.run()
	return b
}
