package p2p

import (
	"github.com/Lobarr/lane"
	"gopkg.in/olahol/melody.v1"
)

//WorkerPriorityQueue priotity queue of workers
type WorkerPriorityQueue struct {
	pq *lane.PQueue
}

//Push adds worker to the priority queue
func (pq WorkerPriorityQueue) Push(s *melody.Session, priority int) {
	pq.getPQ().Push(s, priority)
}

//Pop returns next worker in the priority queue
func (pq WorkerPriorityQueue) Pop() *melody.Session {
	i, _ := pq.getPQ().Pop()
	return i.(*melody.Session)
}

//Remove removes worker from the priority queue
func (pq WorkerPriorityQueue) Remove(s *melody.Session) {
	pq.getPQ().RemoveSession(s)
}

func (pq WorkerPriorityQueue) getPQ() *lane.PQueue {
	return pq.pq
}

//Len returns priority queues length
func (pq WorkerPriorityQueue) Len() int {
	return pq.getPQ().Size()
}

//NewWorkerPriorityQueue initializes worker priority queue
func NewWorkerPriorityQueue() *WorkerPriorityQueue {
	pq := lane.NewPQueue(lane.MINPQ)
	q := &WorkerPriorityQueue{
		pq: pq,
	}
	return q
}
