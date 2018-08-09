package job

import (
	"errors"
	"time"
)

var (
	//ErrExecNotFound occurs when exec not found
	ErrExecNotFound = errors.New("Exec Not Found")
	//ErrInvalidPriority occurs when invalid priority number
	ErrInvalidPriority = errors.New("Invalid priority number")
	//ErrRetriesOutsideLimit occurs when retries outside limit
	ErrRetriesOutsideLimit = errors.New("Retries outside limit")
	//ErrRetryDelayOutsideLimit occurs when retry delay is outside limit
	ErrRetryDelayOutsideLimit = errors.New("Retry Delay outside limit")
	//ErrBackoffOutsideLimit occurs when backoff duration is outside limit
	ErrBackoffOutsideLimit = errors.New("Backoff duration outside limit")
	//ErrExecutionTimeBehind occurs when execution time is past
	ErrExecutionTimeBehind = errors.New("Execution time is past")
	//ErrJobsLenRange occurs when number of jobs is more than allowed
	ErrJobsLenRange = errors.New("Number of jobs is more than allowed")
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

//Gizo Priorities
const (
	BOOST  = 4
	HIGH   = 3
	MEDIUM = 2
	LOW    = 1
	NORMAL = 0 //! default
)

//Gizo Statuses
const (
	CANCELLED   = "CANCELLED"  //cancelled by sender
	QUEUED      = "QUEUED"     //job added to job queue
	TIMEOUT     = "TIMEOUT"    //job timed out
	RUNNING     = "RUNNING"    //job being executed
	FINISHED    = "FINISHED"   //job done
	RETRYING    = "RETRYING"   //job retrying
	DISPATHCHED = "DISPATCHED" //job dispatched to worker
	STARTED     = "STARTED"    //job received by dispatcher (prior to dispatch)
)
