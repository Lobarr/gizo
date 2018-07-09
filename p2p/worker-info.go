package p2p

import (
	"github.com/gizo-network/gizo/job/queue/qItem"
)

//WorkerInfo information about worker
type WorkerInfo struct {
	pub  string
	job  *qItem.Item
	shut bool
}

//NewWorkerInfo initalizes worker info
func NewWorkerInfo(pub string) *WorkerInfo {
	return &WorkerInfo{pub: pub}
}

//GetPub returns public key
func (w WorkerInfo) GetPub() string {
	return w.pub
}

//SetPub sets public key
func (w *WorkerInfo) SetPub(pub string) {
	w.pub = pub
}

//GetJob returns job
func (w WorkerInfo) GetJob() *qItem.Item {
	return w.job
}

//SetJob sets job
func (w *WorkerInfo) SetJob(j *qItem.Item) {
	w.job = j
}

//Assign assigns job to worker
func (w *WorkerInfo) Assign(j *qItem.Item) {
	w.job = j
}

//GetShut returns shut
func (w WorkerInfo) GetShut() bool {
	return w.shut
}

//SetShut sets shut
func (w *WorkerInfo) SetShut(s bool) {
	w.shut = s
}

//Busy checks if worker is busy
func (w *WorkerInfo) Busy() bool {
	return w.GetJob() == nil
}
