package cache

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gizo-network/gizo/helpers"

	"github.com/allegro/bigcache"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/job"
)

//MaxCacheLen number of jobs held in cache
const MaxCacheLen = 128

//ErrCacheFull occurs when cache is filled up
var ErrCacheFull = errors.New("Cache: cache filled up")

//JobCache holds most likely to be executed jobs
type JobCache struct {
	cache       IBigCache
	blockchain  core.IBlockChain
	logger      helpers.ILogger
	watchMode   bool
	quitChannel chan (struct{})
}

func (jobCache JobCache) getCache() IBigCache {
	return jobCache.cache
}

func (jobCache JobCache) getBC() core.IBlockChain {
	return jobCache.blockchain
}

func (jobCache JobCache) getLogger() helpers.ILogger {
	return jobCache.logger
}

func (jobCache JobCache) getQuitChnnel() chan (struct{}) {
	return jobCache.quitChannel
}

//StopWatchMode exits watch mode
func (jobCache JobCache) StopWatchMode() {
	jobCache.getQuitChnnel() <- struct{}{}
}

//IsFull returns true if cache is full
func (jobCache JobCache) IsFull() bool {
	if jobCache.getCache().Len() >= MaxCacheLen {
		return true
	}
	return false
}

//updates cache every minute
func (jobCache JobCache) watch() error {
	ticker := time.NewTicker(time.Minute)
	quitChannel := jobCache.getQuitChnnel()
	for {
		select {
		case <-ticker.C:
			if err := jobCache.fill(); err != nil {
				return err
			}
		case <-quitChannel:
			ticker.Stop()
			return nil
		}
	}
}

//Set adds key and value to cache
func (jobCache JobCache) Set(key string, val []byte) error {
	jobCache.getLogger().Debugf("Cache: setting job %s", key)
	if jobCache.IsFull() {
		return ErrCacheFull
	}
	if err := jobCache.getCache().Set(key, val); err != nil {
		return err
	}
	return nil
}

//Get returns job from cache
func (jobCache JobCache) Get(key string) (*job.Job, error) {
	jobCache.getLogger().Debugf("Cache: getting job %s", key)
	jobBytes, err := jobCache.getCache().Get(key)
	if err != nil {
		return nil, err
	}
	var job *job.Job
	if err = json.Unmarshal(jobBytes, &job); err != nil {
		return nil, err
	}
	return job, nil
}

//fills up the cache with jobs with most execs in the last 15 blocks
func (jobCache JobCache) fill() error {
	jobCache.getLogger().Debug("Cache: filling up cache")
	var jobs []job.Job
	blocks, err := jobCache.getBC().GetLatest15()
	if err != nil {
		return err
	}
	if len(blocks) != 0 {
		for _, block := range blocks {
			var nodes = block.GetNodes()
			for _, job := range nodes {
				jobs = append(jobs, job.GetJob())
			}
		}
		var endIndex int
		sorted := job.UniqJob(mergeSort(jobs))
		if len(sorted) > MaxCacheLen {
			endIndex = MaxCacheLen + 1
		} else {
			endIndex = len(sorted)
		}
		for i := 0; i < endIndex; i++ {
			jobBytes, err := json.Marshal(sorted[i])
			if err != nil {
				return err
			}
			jobCache.Set(sorted[i].GetID(), jobBytes)
		}
	}
	return nil
}

//NewJobCache return s initialized job cache
func NewJobCache(blockchain core.IBlockChain, cacheArg IBigCache, loggerArg helpers.ILogger, watch bool) (IJobCache, error) {
	var cache IBigCache
	var logger helpers.ILogger
	quitChannel := make(chan (struct{}))

	if cacheArg == nil {
		bigCache, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Minute))
		if err != nil {
			return nil, err
		}
		cache = bigCache
	} else {
		cache = cacheArg
	}

	if loggerArg == nil {
		logger = helpers.Logger()
	} else {
		logger = loggerArg
	}

	jobCache := JobCache{cache, blockchain, logger, watch, quitChannel}
	if err := jobCache.fill(); err != nil {
		return nil, err
	}

	if watch {
		go jobCache.watch()
	}

	return &jobCache, nil
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
