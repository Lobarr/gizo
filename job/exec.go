package job

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/gizo-network/gizo/helpers"
)

//Exec - config for job executions
type Exec struct {
	Hash          string
	Timestamp     int64
	Duration      time.Duration //saved in nanoseconds
	Args          []interface{} // parameters
	Err           interface{}
	Priority      int
	Result        interface{}
	Status        string        //job status
	Retries       int           // number of max retries
	RetriesCount  int           //number of retries
	Backoff       time.Duration //backoff time of retries (seconds)
	ExecutionTime int64         // time scheduled to run (unix) - should sleep # of seconds before adding to job queue
	Interval      int           //periodic job exec (seconds)
	By            string        //! ID of the worker node that ran this
	TTL           time.Duration //! time limit of job running
	Pub           string        //! public key for private jobs
	Envs          []byte        //encrypted environment variables
	cancel        chan struct{}
}

//NewExec initializes an exec
func NewExec(args []interface{}, retries, priority int, backoff time.Duration, execTime int64, interval int, ttl time.Duration, pub string, envs EnvironmentVariables, passphrase string) (*Exec, error) {
	if retries > MaxRetries {
		return nil, ErrRetriesOutsideLimit
	}

	if backoff.Seconds() > MaxRetryBackoff {
		return nil, ErrBackoffOutsideLimit
	}

	envsBytes, err := json.Marshal(envs)
	if err != nil {
		return nil, err
	}
	encryptEnvs, err := helpers.Encrypt(envsBytes, passphrase)
	if err != nil {
		return nil, err
	}
	ex := &Exec{
		Args:          args,
		Retries:       retries,
		RetriesCount:  0, //initialized to 0
		Priority:      priority,
		Status:        STARTED,
		Backoff:       backoff,
		Interval:      interval,
		ExecutionTime: execTime,
		TTL:           ttl,
		Envs:          encryptEnvs,
		Pub:           pub,
		cancel:        make(chan struct{}),
	}
	return ex, nil
}

//Cancel cancels job
func (e *Exec) Cancel() {
	e.cancel <- struct{}{}
}

//GetCancelChan returns cancel channel
func (e Exec) GetCancelChan() chan struct{} {
	return e.cancel
}

//GetEnvs returns environment variables
func (e Exec) GetEnvs(passphrase string) (EnvironmentVariables, error) {
	d, err := helpers.Decrypt(e.Envs, passphrase)
	if err != nil {
		return EnvironmentVariables{}, errors.New("Unable to decrypt environment variables")
	}
	var envs EnvironmentVariables
	err = json.Unmarshal(d, &envs)
	if err != nil {
		return EnvironmentVariables{}, err
	}
	return envs, nil
}

//GetEnvsMap returns environment variables as a map
func (e Exec) GetEnvsMap(passphrase string) (map[string]interface{}, error) {
	temp := make(map[string]interface{})
	envs, err := e.GetEnvs(passphrase)
	if err != nil {
		return nil, err
	}
	for _, val := range envs {
		temp[val.GetKey()] = val.GetValue()
	}
	return temp, nil
}

//GetTTL returns ttl
func (e Exec) GetTTL() time.Duration {
	return e.TTL
}

//SetTTL sets ttl
func (e *Exec) SetTTL(ttl time.Duration) {
	e.TTL = ttl
}

//GetInterval returns interval
func (e Exec) GetInterval() int {
	return e.Interval
}

//SetInterval sets interval
func (e *Exec) SetInterval(i int) {
	e.Interval = i
}

//GetPriority returns prioritys
func (e Exec) GetPriority() int {
	return e.Priority
}

//SetPriority sets priority
func (e *Exec) SetPriority(p int) error {
	switch p {
	case HIGH:
	case MEDIUM:
	case LOW:
	case NORMAL:
	default:
		return ErrInvalidPriority
	}
	return nil
}

//GetExecutionTime returne execution time
func (e Exec) GetExecutionTime() int64 {
	return e.ExecutionTime
}

//SetExecutionTime sets execution time - unix
func (e *Exec) SetExecutionTime(t int64) error {
	if time.Now().Unix() > t {
		return ErrExecutionTimeBehind
	}
	e.ExecutionTime = t
	return nil
}

//GetBackoff returns backoff
func (e Exec) GetBackoff() time.Duration {
	return e.Backoff
}

//SetBackoff sets backoff
func (e *Exec) SetBackoff(b time.Duration) error {
	if b > MaxRetryBackoff {
		return ErrRetryDelayOutsideLimit
	}
	e.Backoff = b
	return nil
}

//GetRetriesCount return retries count
func (e Exec) GetRetriesCount() int {
	return e.RetriesCount
}

//IncrRetriesCount increments retries count
func (e *Exec) IncrRetriesCount() {
	e.RetriesCount++
}

//GetRetries return retries
func (e Exec) GetRetries() int {
	return e.Retries
}

//SetRetries sets retries
func (e *Exec) SetRetries(r int) error {
	if r > MaxRetries {
		return ErrRetriesOutsideLimit
	}
	e.Retries = r
	return nil
}

//GetStatus returns status
func (e Exec) GetStatus() string {
	return e.Status
}

//SetStatus sets status
func (e *Exec) SetStatus(s string) {
	e.Status = s
}

//GetArgs return args
func (e Exec) GetArgs() []interface{} {
	return e.Args
}

//SetArgs sets args
func (e *Exec) SetArgs(a []interface{}) {
	e.Args = a
}

//GetHash returns hash
func (e Exec) GetHash() string {
	return e.Hash
}

func (e *Exec) SetHash() error {
	stringified, err := json.Marshal(e.GetErr())
	if err != nil {
		return err
	}
	result, err := json.Marshal(e.GetResult())
	if err != nil {
		return err
	}

	pub, err := hex.DecodeString(e.getPub())
	if err != nil {
		return err
	}

	worker, err := hex.DecodeString(e.GetBy())
	if err != nil {
		return err
	}

	header := bytes.Join(
		[][]byte{
			[]byte(strconv.FormatInt(e.GetTimestamp(), 10)),
			[]byte(strconv.FormatInt(int64(e.GetDuration()), 10)),
			stringified,
			result,
			e.Envs,
			pub,
			worker,
			[]byte(e.GetBy()),
		},
		[]byte{},
	)

	hash := sha256.Sum256(header)
	e.Hash = hex.EncodeToString(hash[:])
	return nil
}

//GetTimestamp return timestamp
func (e Exec) GetTimestamp() int64 {
	return e.Timestamp
}

//SetTimestamp sets timestamp
func (e *Exec) SetTimestamp(t int64) {
	e.Timestamp = t
}

//GetDuration returns duration
func (e Exec) GetDuration() time.Duration {
	return e.Duration
}

//SetDuration sets duration
func (e *Exec) SetDuration(t time.Duration) {
	e.Duration = t
}

//GetErr returns err
func (e Exec) GetErr() interface{} {
	return e.Err
}

//SetErr sets err
func (e *Exec) SetErr(err interface{}) {
	e.Err = err
}

//GetResult returns result
func (e Exec) GetResult() interface{} {
	return e.Result
}

//SetResult sets result
func (e *Exec) SetResult(r interface{}) {
	e.Result = r
}

//GetBy returns by
func (e Exec) GetBy() string {
	return e.By
}

//SetBy sets by
func (e *Exec) SetBy(by string) {
	e.By = by
}

func (e Exec) GetPub() string {
	return e.Pub
}

//UniqExec returns unique values of parameter
func UniqExec(execs []Exec) []Exec {
	temp := []Exec{}
	seen := make(map[string]bool)
	for _, exec := range execs {
		if _, ok := seen[exec.GetHash()]; ok {
			continue
		}
		seen[string(exec.GetHash())] = true
		temp = append(temp, exec)
	}
	return temp
}
