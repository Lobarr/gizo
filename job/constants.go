package job

import (
	"errors"
	"time"
)

var (
	ErrExecNotFound           = errors.New("Exec Not Found")
	ErrInvalidPriority        = errors.New("Invalid priority number")
	ErrRetriesOutsideLimit    = errors.New("Retries outside limit")
	ErrRetryDelayOutsideLimit = errors.New("Retry Delay outside limit")
	ErrExecutionTimeBehind    = errors.New("Execution time is past")
	ErrJobsLenRange           = errors.New("Number of jobs is more than allowed")
)

const (
	//MaxExecs max number of jobs allowed in the chain
	MaxExecs = 10
	//MaxRetries max number of retries allowed
	MaxRetries = 5
	//MaxRetryBackoff time limit between retries
	MaxRetryBackoff = 120 //! 2 minutes
	//DefaultMaxTTL time limit of job
	DefaultMaxTTL = time.Minute * 10
	//DefaultRetries default number of retries
	DefaultRetries = 0
	//DefaultPriority default priority
	DefaultPriority = NORMAL
)

//! priorities
const (
	BOOST  = 4
	HIGH   = 3
	MEDIUM = 2
	LOW    = 1
	NORMAL = 0 //! default
)

//! statuses
const (
	CANCELLED   = "CANCELLED"  // cancelled by sender
	QUEUED      = "QUEUED"     // job added to job queue
	TIMEOUT     = "TIMEOUT"    // job timed out
	RUNNING     = "RUNNING"    //job being executed
	FINISHED    = "FINISHED"   //job done
	RETRYING    = "RETRYING"   //job retrying
	DISPATHCHED = "DISPATCHED" //job dispatched to worker
	STARTED     = "STARTED"    //job received by dispatcher (prior to dispatch)
)
