package job

import "time"

// IJobType provides for multiple types of jobs
type IJobType interface {
	Cancel()
	GetCancelChan() chan struct{}
	GetJob() Request
	GetStatus() string
	Result() Request
	Dispatch()
}

//IJob interface for jobs
type IJob interface {
	Sign(string) error
	VerifySignature(string) (bool, error)
	GetSubmissionTime() int64
	IsEmpty() (bool, error)
	GetPrivate() bool
	GetName() string
	GetID() string
	GetHash() string
	Verify() (bool, error)
	GetExec(string) (IExec, error)
	GetLatestExec() IExec
	GetExecs() []IExec
	GetSignature() string
	AddExec(IExec)
	GetTask() (string, error)
	Execute(exec IExec, passphrase string) IExec
}

//IExec interface for exec
type IExec interface {
	Cancel()
	GetCancelChan() chan struct{}
	GetEnvs(string) (EnvironmentVariables, error)
	GetEnvsMap(string) (map[string]interface{}, error)
	GetTTL() time.Duration
	SetTTL(time.Duration)
	GetInterval() int
	SetInterval(int)
	GetPriority() int
	SetPriority(int) error
	GetExecutionTime() int64
	SetExecutionTime(int64) error
	GetBackoff() time.Duration
	SetBackoff(time.Duration) error
	GetRetriesCount() int
	IncrRetriesCount()
	GetRetries() int
	SetRetries(int) error
	GetStatus() string
	SetStatus(string)
	GetArgs() []interface{}
	SetArgs([]interface{})
	GetHash() string
	GetTimestamp() int64
	SetTimestamp(int64)
	GetDuration() time.Duration
	SetDuration(time.Duration)
	GetErr() interface{}
	SetErr(interface{})
	GetResult() interface{}
	SetResult(interface{})
	GetBy() string
	SetBy(string)
	GetPub() string
	SetHash() error
}
