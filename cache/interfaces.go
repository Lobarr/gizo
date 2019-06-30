package cache

import "github.com/gizo-network/gizo/job"

//IBigCache interface for big cache
type IBigCache interface {
	Get(string) ([]byte, error)
	Len() int
	Set(string, []byte) error
}

//IJobCache interface for job cache
type IJobCache interface {
	Get(string) (*job.Job, error)
	IsFull() bool
	Set(string, []byte) error
	StopWatchMode()
}
