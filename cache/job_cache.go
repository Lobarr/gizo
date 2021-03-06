package cache

import (
	"errors"
	"time"

	"github.com/kpango/glg"

	"github.com/gizo-network/gizo/helpers"

	"github.com/allegro/bigcache"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/job"
)

const (
	//MaxCacheLen number of jobs held in cache
	MaxCacheLen = 128
)

var (
	//ErrCacheFull occurs when cache is filled up
	ErrCacheFull = errors.New("Cache: Cache filled up")
)

//JobCache holds most likely to be executed jobs
type JobCache struct {
	cache  *bigcache.BigCache
	bc     *core.BlockChain
	logger *glg.Glg
}

func (c JobCache) getCache() *bigcache.BigCache {
	return c.cache
}

func (c JobCache) getBC() *core.BlockChain {
	return c.bc
}

//IsFull returns true if cache is full
func (c JobCache) IsFull() bool {
	if c.getCache().Len() >= MaxCacheLen {
		return true
	}
	return false
}

//updates cache every minute
func (c JobCache) watch() {
	ticker := time.NewTicker(time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				c.fill()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

//Set adds key and value to cache
func (c JobCache) Set(key string, val []byte) error {
	if c.getCache().Len() >= MaxCacheLen {
		return ErrCacheFull
	}
	c.getCache().Set(key, val)
	return nil
}

//Get returns job from cache
func (c JobCache) Get(key string) (*job.Job, error) {
	jBytes, err := c.getCache().Get(key)
	if err != nil {
		return nil, err
	}
	var j *job.Job
	err = helpers.Deserialize(jBytes, &j)
	if err != nil {
		return nil, err
	}
	return j, nil
}

//fills up the cache with jobs with most execs in the last 15 blocks
func (c JobCache) fill() {
	var jobs []job.Job
	blks, err := c.getBC().GetLatest15()
	if err != nil {
		c.logger.Fatal(err)
	}
	if len(blks) != 0 {
		for _, blk := range blks {
			for _, job := range blk.GetNodes() {
				jobs = append(jobs, job.GetJob())
			}
		}
		sorted := job.UniqJob(mergeSort(jobs))
		if len(sorted) > MaxCacheLen {
			for i := 0; i <= MaxCacheLen; i++ {
				jobBytes, err := helpers.Serialize(sorted[i])
				if err != nil {
					c.logger.Fatal(err)
				}
				c.Set(sorted[i].GetID(), jobBytes)
			}
		} else {
			for _, job := range sorted {
				jobBytes, err := helpers.Serialize(job)
				if err != nil {
					c.logger.Fatal(err)
				}
				c.Set(job.GetID(), jobBytes)
			}
		}
	}
}

//NewJobCache return s initialized job cache
func NewJobCache(bc *core.BlockChain) *JobCache {
	c, _ := bigcache.NewBigCache(bigcache.DefaultConfig(time.Minute))
	jc := JobCache{c, bc, helpers.Logger()}
	jc.fill()
	go jc.watch()
	return &jc
}

//NewJobCacheNoWatch creates a new jobcache without updating every minute
func NewJobCacheNoWatch(bc *core.BlockChain) *JobCache {
	c, _ := bigcache.NewBigCache(bigcache.DefaultConfig(time.Minute))
	jc := JobCache{c, bc, helpers.Logger()}
	jc.fill()
	return &jc
}

// merge returns an array of job in order of number of execs in the job from max to min
func merge(left, right []job.Job) []job.Job {
	size, i, j := len(left)+len(right), 0, 0
	slice := make([]job.Job, size, size)
	for k := 0; k < size; k++ {
		if i > len(left)-1 && j <= len(right)-1 {
			slice[k] = right[j]
			j++
		} else if j > len(right)-1 && i <= len(left)-1 {
			slice[k] = left[i]
			i++
		} else if len(left[i].GetExecs()) > len(right[j].GetExecs()) {
			slice[k] = left[i]
			i++
		} else {
			slice[k] = right[j]
			j++
		}
	}
	return slice
}

//used to quickly sort the jobs from max execs to min execs
func mergeSort(slice []job.Job) []job.Job {
	if len(slice) < 2 {
		return slice
	}
	mid := (len(slice)) / 2
	return merge(mergeSort(slice[:mid]), mergeSort(slice[mid:]))
}
