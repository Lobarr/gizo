package queue

import (
	lane "github.com/Lobarr/lane"
	"github.com/gizo-network/gizo/helpers"
	"github.com/gizo-network/gizo/job"
	"github.com/gizo-network/gizo/job/queue/qitem"
	"github.com/kpango/glg"
)

//JobPriorityQueue priority queue of jobs
type JobPriorityQueue struct {
	pq     *lane.PQueue
	logger *glg.Glg
}

//Push adds item config to priority queue
func (pq JobPriorityQueue) Push(j job.Job, exec *job.Exec, results chan<- qitem.Item, cancel chan struct{}) error {
	task, err := j.GetTask()
	if err != nil {
		return err
	}
	temp := job.Job{
		ID:             j.GetID(),
		Hash:           j.GetHash(),
		Name:           j.GetName(),
		Task:           task,
		Signature:      j.GetSignature(),
		SubmissionTime: j.GetSubmissionTime(),
		Private:        j.GetPrivate(),
	}
	pq.GetPQ().Push(qitem.NewItem(temp, exec, results, cancel), exec.GetPriority())
	pq.logger.Log("JobPriotityQueue: received job - " + j.GetID())
	return nil
}

//PushItem adds item to priority queue
func (pq JobPriorityQueue) PushItem(i qitem.Item, piority int) {
	pq.GetPQ().Push(i, piority)
	pq.logger.Log("JobPriotityQueue: received job - " + i.GetID())

}

//Pop returns next item in the queue
func (pq JobPriorityQueue) Pop() qitem.Item {
	i, _ := pq.GetPQ().Pop()
	return i.(qitem.Item)
}

//Remove removes item from the priority queue
func (pq JobPriorityQueue) Remove(hash string) {
	pq.pq.RemoveHash(hash)
}

//GetPQ returns priority queue
func (pq JobPriorityQueue) GetPQ() *lane.PQueue {
	return pq.pq

}

//Len returns the size of the priority queue
func (pq JobPriorityQueue) Len() int {
	return pq.GetPQ().Size()
}

//NewJobPriorityQueue initializes a job priority queue
func NewJobPriorityQueue() *JobPriorityQueue {
	pq := lane.NewPQueue(lane.MAXPQ)
	q := &JobPriorityQueue{
		pq:     pq,
		logger: helpers.Logger(),
	}
	return q
}
