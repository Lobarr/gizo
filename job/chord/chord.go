package chord

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

//Chord - jobs executed one after the other and the results passed to a callback
type Chord struct {
	jobs     []job.Request
	bc       *core.BlockChain
	pq       *queue.JobPriorityQueue
	jc       *cache.JobCache
	callback job.Request
	result   job.Request
	logger   *glg.Glg
	length   int
	status   string
	cancel   chan struct{}
}

//NewChord returns chord
func NewChord(j []job.Request, callback job.Request, bc *core.BlockChain, pq *queue.JobPriorityQueue, jc *cache.JobCache) (*Chord, error) {
	//FIXME: count callback execs too
	length := len(callback.GetExec())
	for _, jr := range j {
		length += len(jr.GetExec())
	}

	if length > job.MaxExecs {
		return nil, job.ErrJobsLenRange
	}
	c := &Chord{
		jobs:     j,
		bc:       bc,
		pq:       pq,
		jc:       jc,
		callback: callback,
		length:   length,
		logger:   helpers.Logger(),
		cancel:   make(chan struct{}),
	}
	return c, nil
}

//Cancel terminates chain
func (c *Chord) Cancel() {
	c.cancel <- struct{}{}
}

//GetCancelChan returns cancel channel
func (c Chord) GetCancelChan() chan struct{} {
	return c.cancel
}

//GetCallback returns callback exec
func (c Chord) GetCallback() job.Request {
	return c.callback
}

func (c *Chord) setCallback(j job.Request) {
	c.callback = j
}

//GetJobs returns jobs
func (c Chord) GetJobs() []job.Request {
	return c.jobs
}

func (c *Chord) setJobs(j []job.Request) {
	c.jobs = j
}

//GetStatus returns status
func (c Chord) GetStatus() string {
	return c.status
}

func (c *Chord) setStatus(s string) {
	c.status = s
}

func (c Chord) getLength() int {
	return c.length
}

func (c *Chord) setBC(bc *core.BlockChain) {
	c.bc = bc
}

func (c Chord) getPQ() *queue.JobPriorityQueue {
	return c.pq
}

func (c Chord) getBC() *core.BlockChain {
	return c.bc
}

func (c Chord) getJC() *cache.JobCache {
	return c.jc
}

func (c *Chord) setResults(res job.Request) {
	c.result = res
}

//Result returns result
func (c Chord) Result() job.Request {
	return c.result
}

//Dispatch executes the chord
func (c *Chord) Dispatch() {
	c.setStatus(job.RUNNING)
	var items []qItem.Item // used to hold results
	resChan := make(chan qItem.Item)
	cancelled := false
	closeCancel := make(chan struct{})
	var wg sync.WaitGroup
	//! watch cancel channel
	wg.Add(1)
	go func() {
		select {
		case <-c.cancel:
			cancelled = true
			//! for jobs
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

			//! for callback
			for _, exec := range c.GetCallback().GetExec() {
				if exec.GetStatus() == job.RUNNING || exec.GetStatus() == job.RETRYING {
					exec.Cancel()
				}
				if exec.GetResult() != nil {
					exec.SetStatus(job.FINISHED)
				} else {
					exec.SetStatus(job.CANCELLED)
				}
			}
			break
		case <-closeCancel:
			break
		}
		wg.Done()
	}()
	for _, jr := range c.GetJobs() {
		c.setStatus("Queueing execs of job - " + jr.GetID())
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
				if cancelled == true {
					task, err := j.GetTask()
					if err != nil {
						jr.GetExec()[i].SetErr(err)
					}
					items = append(items, qItem.NewItem(job.Job{
						ID:             j.GetID(),
						Hash:           j.GetHash(),
						Name:           j.GetName(),
						Task:           task,
						Signature:      j.GetSignature(),
						SubmissionTime: j.GetSubmissionTime(),
						Private:        j.GetPrivate(),
					}, jr.GetExec()[i], resChan, c.GetCancelChan()))
				} else {
					if jr.GetExec()[i].GetExecutionTime() != 0 {
						c.logger.Info("Chord: Queuing in " + strconv.FormatFloat(time.Unix(jr.GetExec()[i].GetExecutionTime(), 0).Sub(time.Now()).Seconds(), 'f', -1, 64) + " nanoseconds")
						time.Sleep(time.Nanosecond * time.Duration(time.Unix(jr.GetExec()[i].GetExecutionTime(), 0).Sub(time.Now()).Nanoseconds()))

					}
					c.getPQ().Push(*j, jr.GetExec()[i], resChan, c.GetCancelChan()) //? queues first job
					items = append(items, <-resChan)
				}
			}
		}
	}
	close(resChan)

	var callbackResults []qItem.Item
	var callbackArgs []interface{}        //holds result of execs
	callbackChan := make(chan qItem.Item) //causes program to pause
	for _, item := range items {
		callbackArgs = append(callbackArgs, item.GetExec().GetResult())
	}

	//! sets args as results or jobs
	for _, exec := range c.GetCallback().GetExec() {
		exec.SetArgs(callbackArgs)
	}

	cj, err := c.getBC().FindJob(c.GetCallback().GetID())
	if err != nil {
		for _, exec := range c.GetCallback().GetExec() {
			exec.SetErr("Unable to find job - " + c.GetCallback().GetID())
		}
	} else {
		for _, exec := range c.GetCallback().GetExec() {
			if cancelled == true {
				task, err := cj.GetTask()
				if err != nil {
					exec.SetErr(err)
				}
				callbackResults = append(callbackResults, qItem.NewItem(job.Job{
					ID:             cj.GetID(),
					Hash:           cj.GetHash(),
					Name:           cj.GetName(),
					Task:           task,
					Signature:      cj.GetSignature(),
					SubmissionTime: cj.GetSubmissionTime(),
					Private:        cj.GetPrivate(),
				}, exec, callbackChan, c.GetCancelChan()))
			} else {
				if exec.GetExecutionTime() != 0 {
					c.logger.Info("Chord: Queuing in " + strconv.FormatFloat(time.Unix(exec.GetExecutionTime(), 0).Sub(time.Now()).Seconds(), 'f', -1, 64) + " nanoseconds")
					time.Sleep(time.Nanosecond * time.Duration(time.Unix(exec.GetExecutionTime(), 0).Sub(time.Now()).Nanoseconds()))
				}
				c.getPQ().Push(*cj, exec, callbackChan, c.GetCancelChan())
				callbackResults = append(callbackResults, <-callbackChan)
			}
		}
	}

	close(callbackChan)

	var callback job.Request
	callback.SetID(c.GetCallback().GetID())
	for _, item := range callbackResults {
		callback.AppendExec(item.GetExec())
	}

	if cancelled == false {
		closeCancel <- struct{}{}
		c.setStatus(job.FINISHED)
	} else {
		c.setStatus(job.CANCELLED)
	}
	wg.Wait()
	c.setResults(callback)
}
