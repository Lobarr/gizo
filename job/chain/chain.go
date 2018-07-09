package chain

import (
	"strconv"
	"sync"
	"time"

	"github.com/gizo-network/gizo/helpers"

	"github.com/gizo-network/gizo/cache"

	"github.com/kpango/glg"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/job"
	"github.com/gizo-network/gizo/job/queue"
)

//Chain - jobs executed one after the other
type Chain struct {
	jobs   []job.Request
	bc     *core.BlockChain
	pq     *queue.JobPriorityQueue
	jc     *cache.JobCache
	result []job.Request
	logger *glg.Glg
	length int
	status string
	cancel chan struct{}
}

//NewChain returns chain
func NewChain(j []job.Request, bc *core.BlockChain, pq *queue.JobPriorityQueue, jc *cache.JobCache) (*Chain, error) {
	length := 0
	for _, jr := range j {
		length += len(jr.GetExec())
	}
	if length > job.MaxExecs {
		return nil, job.ErrJobsLenRange
	}
	c := &Chain{
		jobs:   j,
		bc:     bc,
		pq:     pq,
		jc:     jc,
		length: length,
		cancel: make(chan struct{}),
		logger: helpers.Logger(),
	}
	return c, nil
}

//Cancel terminates exec
func (c *Chain) Cancel() {
	c.cancel <- struct{}{}
}

//GetCancelChan returns cancel channel
func (c Chain) GetCancelChan() chan struct{} {
	return c.cancel
}

//GetJobs returns jobs
func (c Chain) GetJobs() []job.Request {
	return c.jobs
}

func (c *Chain) setJobs(j []job.Request) {
	c.jobs = j
}

//GetStatus returns status
func (c Chain) GetStatus() string {
	return c.status
}

func (c *Chain) setStatus(s string) {
	c.status = s
}

func (c Chain) getLength() int {
	return c.length
}

func (c *Chain) setBC(bc *core.BlockChain) {
	c.bc = bc
}

func (c Chain) getPQ() *queue.JobPriorityQueue {
	return c.pq
}

func (c Chain) getBC() *core.BlockChain {
	return c.bc
}

func (c Chain) getJC() *cache.JobCache {
	return c.jc
}

func (c *Chain) setResults(res []job.Request) {
	c.result = res
}

//Result returns result
func (c Chain) Result() []job.Request {
	return c.result
}

//Dispatch executes the chain
func (c *Chain) Dispatch() {
	c.setStatus(job.RUNNING)
	var results []qItem.Item // used to hold results
	res := make(chan qItem.Item)
	cancelled := false
	closeCancel := make(chan struct{})
	var wg sync.WaitGroup
	//! watch cancel channel
	wg.Add(1)
	go func() {
		select {
		case <-c.cancel:
			cancelled = true
			for _, jr := range c.GetJobs() {
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
	var jobIDs []string
	for _, jr := range c.GetJobs() {
		c.setStatus("Queueing execs of job - " + jr.GetID())
		jobIDs = append(jobIDs, jr.GetID())
		var j *job.Job
		var err error
		j, err = c.getJC().Get(jr.GetID())
		if j == nil {
			j, err = c.getBC().FindJob(jr.GetID())
		}
		if err != nil {
			for _, exec := range jr.GetExec() {
				exec.SetErr("Unable to find job - " + jr.GetID())
			}
		} else {
			for i := 0; i < len(jr.GetExec()); i++ {
				task, err := j.GetTask()
				if err != nil {
					jr.GetExec()[i].SetErr(err)
				}
				if cancelled == true {
					results = append(results, qItem.NewItem(job.Job{
						ID:             j.GetID(),
						Hash:           j.GetHash(),
						Name:           j.GetName(),
						Task:           task,
						Signature:      j.GetSignature(),
						SubmissionTime: j.GetSubmissionTime(),
						Private:        j.GetPrivate(),
					}, jr.GetExec()[i], res, c.GetCancelChan()))
				} else {
					if jr.GetExec()[i].GetExecutionTime() != 0 {
						c.logger.Info("Chain: Queuing in " + strconv.FormatFloat(time.Unix(jr.GetExec()[i].GetExecutionTime(), 0).Sub(time.Now()).Seconds(), 'f', -1, 64) + " nanoseconds")
						time.Sleep(time.Nanosecond * time.Duration(time.Unix(jr.GetExec()[i].GetExecutionTime(), 0).Sub(time.Now()).Nanoseconds()))
					}
					c.getPQ().Push(*j, jr.GetExec()[i], res, c.GetCancelChan()) //? queues first job
					results = append(results, <-res)
				}
			}
		}
	}
	close(res)

	var grouped []job.Request
	for _, jID := range jobIDs {
		var req job.Request
		req.SetID(jID)
		for _, item := range results {
			if item.GetID() == jID {
				req.AppendExec(item.GetExec())
			}
		}
		grouped = append(grouped, req)
	}

	if cancelled == false {
		closeCancel <- struct{}{}
		c.setStatus(job.FINISHED)
	} else {
		c.setStatus(job.CANCELLED)
	}
	wg.Wait()
	c.setResults(grouped)
}
