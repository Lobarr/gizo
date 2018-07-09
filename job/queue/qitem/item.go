package qitem

import (
	"github.com/gizo-network/gizo/job"
)

//Item used for keeping track of execs
type Item struct {
	Job     job.Job
	Exec    *job.Exec
	results chan<- Item
	cancel  chan<- struct{}
}

//NewItem initializes a new item
func NewItem(j job.Job, exec *job.Exec, results chan<- Item, cancel chan<- struct{}) Item {
	return Item{
		Job:     j,
		Exec:    exec,
		results: results,
		cancel:  cancel,
	}
}

//GetCancel returns cancel channel
func (i Item) GetCancel() chan<- struct{} {
	return i.cancel
}

//SetExec sets exec
func (i *Item) SetExec(ex *job.Exec) {
	i.Exec = ex
}

//GetExec returns exec
func (i Item) GetExec() *job.Exec {
	return i.Exec
}

//GetID returns job id
func (i Item) GetID() string {
	return i.Job.GetID()
}

//GetJob returns job
func (i Item) GetJob() job.Job {
	return i.Job
}

//ResultsChan return result chan
func (i Item) ResultsChan() chan<- Item {
	return i.results
}
