package batch

import (
	"strconv"
	"sync"
	"time"

	"github.com/gizo-network/gizo/helpers"

	"github.com/gizo-network/gizo/cache"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/job"
	"github.com/gizo-network/gizo/job/queue"
	"github.com/gizo-network/gizo/job/queue/qItem"
	"github.com/kpango/glg"
)

//Batch - jobs executed in parralel
type Batch struct {
	jobs   []job.Request
	bc     *core.BlockChain
	pq     *queue.JobPriorityQueue
	jc     *cache.JobCache
	logger *glg.Glg
	result []job.Request
	length int
	status string
	cancel chan struct{}
}

//NewBatch returns batch
func NewBatch(j []job.Request, bc *core.BlockChain, pq *queue.JobPriorityQueue, jc *cache.JobCache) (*Batch, error) {
	length := 0
	for _, jr := range j {
		length += len(jr.GetExec())
	}
	if length > job.MaxExecs {
		return nil, job.ErrJobsLenRange
	}
	b := &Batch{
		jobs:   j,
		bc:     bc,
		pq:     pq,
		jc:     jc,
		length: length,
		logger: helpers.Logger(),
		cancel: make(chan struct{}),
	}

	return b, nil
}

//Cancel terminates batch
func (b *Batch) Cancel() {
	b.cancel <- struct{}{}
}

//GetCancelChan returns cancel channel
func (b Batch) GetCancelChan() chan struct{} {
	return b.cancel
}

//GetJobs return jobs
func (b Batch) GetJobs() []job.Request {
	return b.jobs
}

func (b *Batch) setJobs(j []job.Request) {
	b.jobs = j
}

//GetStatus returns status
func (b Batch) GetStatus() string {
	return b.status
}

func (b *Batch) setStatus(s string) {
	b.status = s
}

func (b *Batch) setBC(bc *core.BlockChain) {
	b.bc = bc
}

func (b Batch) getBC() *core.BlockChain {
	return b.bc
}

func (b Batch) getJC() *cache.JobCache {
	return b.jc
}

func (b Batch) getPQ() *queue.JobPriorityQueue {
	return b.pq
}

func (b Batch) getLength() int {
	return b.length
}

func (b *Batch) setResults(res []job.Request) {
	b.result = res
}

//Result returns result
func (b Batch) Result() []job.Request {
	return b.result
}

//Dispatch executes the batch
func (b *Batch) Dispatch() {
	//! should be run in a go routine because it blocks till all jobs are complete
	var items []qItem.Item
	b.setStatus(job.RUNNING)
	cancelled := false
	closeCancel := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case <-b.cancel:
			cancelled = true
			b.logger.Info("Batch: Cancelling jobs")
			for _, jr := range b.GetJobs() {
				for _, exec := range jr.GetExec() {
					if exec.GetStatus() == job.RUNNING || exec.GetStatus() == job.RETRYING {
						exec.Cancel()
					}
					if exec.GetResult() == nil {
						exec.SetStatus(job.CANCELLED)
					}
				}
			}
			break
		case <-closeCancel:
			break
		}
		wg.Done()
	}()

	results := make(chan qItem.Item, b.getLength())
	var jobIDs []string
	var sleepWG sync.WaitGroup
	for _, jr := range b.GetJobs() {
		b.setStatus("Queueing execs of job - " + jr.GetID())
		jobIDs = append(jobIDs, jr.GetID())
		var j *job.Job
		var err error
		j, err = b.getJC().Get(jr.GetID())
		if j == nil {
			j, err = b.getBC().FindJob(jr.GetID())
		}
		if err != nil {
			for _, exec := range jr.GetExec() {
				exec.SetErr("Unable to find job - " + jr.GetID())
			}
		} else {
			for _, exec := range jr.GetExec() {
				sleepWG.Add(1)
				go func(ex *job.Exec) {
					if ex.GetExecutionTime() != 0 {
						b.logger.Info("Batch: Queuing in " + strconv.FormatFloat(time.Unix(ex.GetExecutionTime(), 0).Sub(time.Now()).Seconds(), 'f', -1, 64) + " nanoseconds")
						time.Sleep(time.Nanosecond * time.Duration(time.Unix(ex.GetExecutionTime(), 0).Sub(time.Now()).Nanoseconds()))
					}
					b.getPQ().Push(*j, ex, results, b.GetCancelChan())
					sleepWG.Done()
				}(exec)
			}
		}
	}
	sleepWG.Wait()

	//! wait for all jobs to be done
	for {
		if len(results) == cap(results) {
			close(results)
			break
		}
	}

	for item := range results {
		items = append(items, item)
	}

	var grouped []job.Request
	for _, jID := range jobIDs {
		var req job.Request
		req.SetID(jID)
		for _, item := range items {
			if item.GetID() == jID {
				req.AppendExec(item.GetExec())
			}
		}
		grouped = append(grouped, req)
	}
	if cancelled == false {
		closeCancel <- struct{}{}
		b.setStatus(job.FINISHED)
	} else {
		b.setStatus(job.CANCELLED)
	}
	wg.Wait()
	b.setResults(grouped)
}
