package p2p

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gizo-network/gizo/helpers"

	"github.com/gizo-network/gizo/job/batch"

	"github.com/gizo-network/gizo/job/chain"

	"github.com/gizo-network/gizo/job/chord"

	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job/solo"

	"github.com/gizo-network/gizo/job"
)

//RPC exposes rpc functions
func (d Dispatcher) RPC() {
	d.GetRPC().AddFunction("Version", d.Version)
	d.GetRPC().AddFunction("PeerCount", d.PeerCount)
	d.GetRPC().AddFunction("BlockByHash", d.BlockByHash)
	d.GetRPC().AddFunction("BlockByHeight", d.BlockByHeight)
	d.GetRPC().AddFunction("Latest15Blocks", d.Latest15Blocks)
	d.GetRPC().AddFunction("LatestBlock", d.LatestBlock)
	d.GetRPC().AddFunction("PendingCount", d.PendingCount)
	d.GetRPC().AddFunction("Score", d.Score)
	d.GetRPC().AddFunction("Peers", d.Peers)
	d.GetRPC().AddFunction("PublicKey", d.PublicKey)
	d.GetRPC().AddFunction("NewJob", d.NewJob)
	d.GetRPC().AddFunction("NewExec", d.NewExec)
	d.GetRPC().AddFunction("WorkersCount", d.WorkersCount)
	d.GetRPC().AddFunction("WorkersCountBusy", d.WorkersCountBusy)
	d.GetRPC().AddFunction("WorkersCountNotBusy", d.WorkersCountNotBusy)
	d.GetRPC().AddFunction("ExecStatus", d.ExecStatus)
	d.GetRPC().AddFunction("CancelExec", d.CancelExec)
	d.GetRPC().AddFunction("ExecTimestamp", d.ExecTimestamp)
	d.GetRPC().AddFunction("ExecTimestampString", d.ExecTimestampString)
	d.GetRPC().AddFunction("ExecDurationNanoseconds", d.ExecDurationNanoseconds)
	d.GetRPC().AddFunction("ExecDurationSeconds", d.ExecDurationSeconds)
	d.GetRPC().AddFunction("ExecDurationMinutes", d.ExecDurationMinutes)
	d.GetRPC().AddFunction("ExecDurationString", d.ExecDurationString)
	d.GetRPC().AddFunction("ExecArgs", d.ExecArgs)
	d.GetRPC().AddFunction("ExecErr", d.ExecErr)
	d.GetRPC().AddFunction("ExecPriority", d.ExecPriority)
	d.GetRPC().AddFunction("ExecResult", d.ExecResult)
	d.GetRPC().AddFunction("ExecRetries", d.ExecRetries)
	d.GetRPC().AddFunction("ExecBackoff", d.ExecBackoff)
	d.GetRPC().AddFunction("ExecExecutionTime", d.ExecExecutionTime)
	d.GetRPC().AddFunction("ExecExecutionTimeString", d.ExecExecutionTimeString)
	d.GetRPC().AddFunction("ExecInterval", d.ExecInterval)
	d.GetRPC().AddFunction("ExecBy", d.ExecBy)
	d.GetRPC().AddFunction("ExecTTLNanoseconds", d.ExecTTLNanoseconds)
	d.GetRPC().AddFunction("ExecTTLSeconds", d.ExecTTLSeconds)
	d.GetRPC().AddFunction("ExecTTLMinutes", d.ExecTTLMinutes)
	d.GetRPC().AddFunction("ExecTTLHours", d.ExecTTLHours)
	d.GetRPC().AddFunction("ExecTTLString", d.ExecTTLString)
	d.GetRPC().AddFunction("JobQueueCount", d.JobQueueCount)
	d.GetRPC().AddFunction("LatestBlockHeight", d.LatestBlockHeight)
	d.GetRPC().AddFunction("Job", d.Job)
	d.GetRPC().AddFunction("JobSubmissionTimeUnix", d.JobSubmisstionTimeUnix)
	d.GetRPC().AddFunction("JobSubmissionTimeString", d.JobSubmisstionTimeString)
	d.GetRPC().AddFunction("IsJobPrivate", d.IsJobPrivate)
	d.GetRPC().AddFunction("JobName", d.JobName)
	d.GetRPC().AddFunction("JobLatestExec", d.JobLatestExec)
	d.GetRPC().AddFunction("JobExecs", d.JobExecs)
	d.GetRPC().AddFunction("BlockHashesHex", d.BlockHashesHex)
	d.GetRPC().AddFunction("KeyPair", d.KeyPair)
	d.GetRPC().AddFunction("Solo", d.Solo)
	d.GetRPC().AddFunction("Chord", d.Chord)
	d.GetRPC().AddFunction("Chain", d.Chain)
	d.GetRPC().AddFunction("Batch", d.Batch)
}

//Version returns dispatcher node's version information
func (d Dispatcher) Version() (string, error) {
	height, err := d.GetBC().GetLatestHeight()
	if err != nil {
		return "", err
	}
	hashes, err := d.GetBC().GetBlockHashesHex()
	if err != nil {
		return "", err
	}
	versionBytes, err := helpers.Serialize(NewVersion(GizoVersion, height, hashes))
	return string(versionBytes), nil
}

//PeerCount returns the number of peers a node has
func (d Dispatcher) PeerCount() int {
	return len(d.GetPeers())
}

//BlockByHash returns block of specified hash
func (d Dispatcher) BlockByHash(hash string) (string, error) {
	b, err := d.GetBC().GetBlockInfo(hash)
	if err != nil {
		return "", err
	}
	blockBytes, err := helpers.Serialize(b.GetBlock())
	if err != nil {
		return "", err
	}
	return string(blockBytes), nil
}

//BlockByHeight returns block at specified height
func (d Dispatcher) BlockByHeight(height int) (string, error) {
	b, err := d.GetBC().GetBlockByHeight(height)
	if err != nil {
		return "", err
	}
	blockBytes, err := helpers.Serialize(b)
	if err != nil {
		return "", err
	}
	return string(blockBytes), nil
}

//Latest15Blocks returns array of most recent 15 blocks
func (d Dispatcher) Latest15Blocks() (string, error) {
	blocks, err := d.GetBC().GetLatest15()
	if err != nil {
		return "", err
	}
	blocksBytes, err := helpers.Serialize(blocks)
	if err != nil {
		return "", err
	}
	return string(blocksBytes), nil
}

//LatestBlock returns latest block in the blockchain
func (d Dispatcher) LatestBlock() (string, error) {
	block, err := d.GetBC().GetLatestBlock()
	if err != nil {
		return "", err
	}
	blockBytes, err := helpers.Serialize(block)
	if err != nil {
		return "", err
	}
	return string(blockBytes), nil
}

//PendingCount returns number of jobs waiting to be written to the blockchain
func (d Dispatcher) PendingCount() int {
	return d.GetWriteQ().Size()
}

//Score returns benchmark score of node
func (d Dispatcher) Score() float64 {
	return d.GetBench().GetScore()
}

//Peers returns the public keys of its peers
func (d Dispatcher) Peers() []string {
	return d.GetPeersPubs()
}

//PublicKey returns public key of node
func (d Dispatcher) PublicKey() string {
	return d.GetPubString()
}

//NewJob deploys job to the blockchain and returns ID
func (d Dispatcher) NewJob(task string, name string, priv bool, privKey string) (string, error) {
	if privKey == "" {
		return "", errors.New("Empty private key")
	}
	j, err := job.NewJob(task, name, priv, privKey)
	if err != nil {
		return "", err
	}
	d.AddJob(*j)
	return j.GetID(), nil
}

//NewExec returns exec with specified config
func (d Dispatcher) NewExec(args []interface{}, retries, priority int, backoff int64, execTime int64, interval int, ttl int64, pub string, envs string) (string, error) {
	var e job.EnvironmentVariables
	err := helpers.Deserialize([]byte(envs), &e)
	if err != nil {
		return "", err
	}
	_backoff := time.Second.Seconds() * float64(backoff)
	_ttl := time.Minute.Minutes() * float64(ttl)

	exec, err := job.NewExec(args, retries, priority, time.Duration(_backoff), execTime, interval, time.Duration(_ttl), pub, e, d.GetPubString())
	if err != nil {
		return "", err
	}
	execBytes, err := helpers.Serialize(exec)
	if err != nil {
		return "", err
	}
	return string(execBytes), nil
}

//WorkersCount returns number of workers in a dispatchers standard area
func (d Dispatcher) WorkersCount() int {
	return len(d.GetWorkers())
}

//WorkersCountBusy returns number of workers in a dispatchers standard area that are busy
func (d Dispatcher) WorkersCountBusy() int {
	temp := 0
	d.mu.Lock()
	for _, info := range d.GetWorkers() {
		if info.GetJob() != nil {
			temp++
		}
	}
	d.mu.Unlock()
	return temp
}

//WorkersCountNotBusy returns number of workers in a dispatchers standard area that are not busy
func (d Dispatcher) WorkersCountNotBusy() int {
	temp := 0
	d.mu.Lock()
	for _, info := range d.GetWorkers() {
		if info.GetJob() == nil {
			temp++
		}
	}
	d.mu.Unlock()
	return temp
}

//ExecStatus returns status of exec
func (d Dispatcher) ExecStatus(id string, hash string) (string, error) {
	if d.GetJobPQ().GetPQ().InQueueHash(hash) {
		return job.QUEUED, nil
	} else if worker := d.GetAssignedWorker(hash); worker != nil {
		return job.RUNNING, nil
	}
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return "", err
	}
	return e.GetStatus(), nil
}

//CancelExec cancels exc
func (d Dispatcher) CancelExec(hash string) error {
	if worker := d.GetAssignedWorker(hash); worker != nil {
		cm, err := CancelMessage(d.GetPrivByte())
		if err != nil {
			return err
		}
		worker.Write(cm)
		return nil
	}
	return errors.New("Exec not running")
}

//ExecTimestamp returns timestamp of exec - when the job started running (unix)
func (d Dispatcher) ExecTimestamp(id string, hash string) (int64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetTimestamp(), nil
}

//ExecTimestampString returns timestamp of exec - when the job started running (string)
func (d Dispatcher) ExecTimestampString(id string, hash string) (string, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return "", err
	}
	return time.Unix(e.GetTimestamp(), 0).String(), nil
}

//ExecDurationNanoseconds returns duration of an exec in nanoseconds
func (d Dispatcher) ExecDurationNanoseconds(id string, hash string) (int64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetDuration().Nanoseconds(), nil
}

//ExecDurationSeconds returns duration of an exec in seconds
func (d Dispatcher) ExecDurationSeconds(id string, hash string) (float64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetDuration().Seconds(), nil
}

//ExecDurationMinutes returns duration of an exec in minutes
func (d Dispatcher) ExecDurationMinutes(id string, hash string) (float64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetDuration().Minutes(), nil
}

//ExecDurationString returns duration of an exec as string
func (d Dispatcher) ExecDurationString(id string, hash string) (string, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return "", err
	}
	return e.GetDuration().String(), nil
}

//ExecArgs returns arguments of an exec
func (d Dispatcher) ExecArgs(id string, hash string) ([]interface{}, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return make([]interface{}, 1), err
	}
	return e.GetArgs(), nil
}

//ExecErr returns error of an exec - None if no error occured
func (d Dispatcher) ExecErr(id string, hash string) (interface{}, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return "", err
	}
	return e.GetErr(), nil
}

//ExecPriority returns priority of an exec
func (d Dispatcher) ExecPriority(id string, hash string) (int, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetPriority(), nil
}

//ExecResult result of an exec - None if error occurs
func (d Dispatcher) ExecResult(id string, hash string) (interface{}, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return "", err
	}
	return e.GetResult(), nil
}

//ExecRetries returns number of retries attempted by the worker
func (d Dispatcher) ExecRetries(id string, hash string) (int, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetRetries(), nil
}

//ExecBackoff returns time between retries of an exec(seconds)
func (d Dispatcher) ExecBackoff(id string, hash string) (float64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetBackoff().Seconds(), nil
}

//ExecExecutionTime returns scheduled time of exec (unix)
func (d Dispatcher) ExecExecutionTime(id string, hash string) (int64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetExecutionTime(), nil
}

//ExecExecutionTimeString returns scheduled time of exec (string)
func (d Dispatcher) ExecExecutionTimeString(id string, hash string) (string, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return "", err
	}
	return time.Unix(e.GetExecutionTime(), 0).String(), nil
}

//ExecInterval returns time between retries of an exec(seconds)
func (d Dispatcher) ExecInterval(id string, hash string) (int, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetInterval(), nil
}

//ExecBy returns public key of worker that executed the job
func (d Dispatcher) ExecBy(id string, hash string) (string, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return "", err
	}
	return e.GetBy(), nil
}

//ExecTTLNanoseconds returns tll of exec (nanoseconds)
func (d Dispatcher) ExecTTLNanoseconds(id string, hash string) (int64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetTTL().Nanoseconds(), nil
}

//ExecTTLSeconds returns ttl of exec (seconds)
func (d Dispatcher) ExecTTLSeconds(id string, hash string) (float64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetTTL().Seconds(), nil
}

//ExecTTLMinutes ttl of exec (minutes)
func (d Dispatcher) ExecTTLMinutes(id string, hash string) (float64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetTTL().Minutes(), nil
}

//ExecTTLHours returns ttl of exec (hours)
func (d Dispatcher) ExecTTLHours(id string, hash string) (float64, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return 0, err
	}
	return e.GetTTL().Hours(), nil
}

//ExecTTLString returns ttl of exec (string)
func (d Dispatcher) ExecTTLString(id string, hash string) (string, error) {
	e, err := d.GetBC().FindExec(id, hash)
	if err != nil {
		return "", err
	}
	return e.GetTTL().String(), nil
}

//JobQueueCount returns nubmer of jobs waiting to be executed
func (d Dispatcher) JobQueueCount() int {
	return d.GetJobPQ().Len()
}

//LatestBlockHeight height of latest block in the blockchain
func (d Dispatcher) LatestBlockHeight() (int, error) {
	block, err := d.GetBC().GetLatestBlock()
	if err != nil {
		return 0, err
	}
	return int(block.GetHeight()), nil
}

//Job returns a job
func (d Dispatcher) Job(id string) (string, error) {
	j, err := d.GetBC().FindJob(id)
	if err != nil {
		return "", err
	}
	jBytes, err := helpers.Serialize(j)
	if err != nil {
		return "", err
	}
	return string(jBytes), nil
}

//JobSubmisstionTimeUnix returns submission time of job (unix)
func (d Dispatcher) JobSubmisstionTimeUnix(id string) (int64, error) {
	j, err := d.GetBC().FindJob(id)
	if err != nil {
		return 0, err
	}
	return j.GetSubmissionTime(), nil
}

//JobSubmisstionTimeString returns submission time of job (string)
func (d Dispatcher) JobSubmisstionTimeString(id string) (string, error) {
	j, err := d.GetBC().FindJob(id)
	if err != nil {
		return "", err
	}
	return time.Unix(j.GetSubmissionTime(), 0).String(), nil
}

//IsJobPrivate returns if job is  private (true) / public (false)
func (d Dispatcher) IsJobPrivate(id string) (bool, error) {
	j, err := d.GetBC().FindJob(id)
	if err != nil {
		return false, err
	}
	return j.GetPrivate(), nil
}

//JobName returns name of job
func (d Dispatcher) JobName(id string) (string, error) {
	j, err := d.GetBC().FindJob(id)
	if err != nil {
		return "", err
	}
	return j.GetName(), nil
}

//JobLatestExec returns latest exec of job
func (d Dispatcher) JobLatestExec(id string) (string, error) {
	j, err := d.GetBC().FindJob(id)
	if err != nil {
		return "", err
	}
	jBytes, err := helpers.Serialize(j)
	if err != nil {
		return "", err
	}
	return string(jBytes), nil
}

//JobExecs returns all execs of a job
func (d Dispatcher) JobExecs(id string) (string, error) {
	execs, err := d.GetBC().GetJobExecs(id)
	if err != nil {
		return "", err
	}
	execsBytes, err := helpers.Serialize(execs)
	if err != nil {
		return "", err
	}
	return string(execsBytes), nil
}

//BlockHashesHex returns hashes of all blocks in the blockchain
func (d Dispatcher) BlockHashesHex() ([]string, error) {
	hashes, err := d.GetBC().GetBlockHashesHex()
	if err != nil {
		return []string{}, err
	}
	return hashes, nil
}

//KeyPair returns new pub and priv keypair
func (d Dispatcher) KeyPair() (string, error) {
	priv, pub := crypt.GenKeys()
	temp := make(map[string]string)
	temp["priv"] = priv
	temp["pub"] = pub
	keysBytes, err := helpers.Serialize(temp)
	if err != nil {
		return "", err
	}
	return string(keysBytes), nil
}

//Solo executes a single exec
func (d Dispatcher) Solo(jr string) (string, error) {
	//TODO: send result to message broker
	var request job.Request
	err := helpers.Deserialize([]byte(jr), &request)
	if err != nil {
		return "", err
	}
	solo := solo.NewSolo(request, d.GetBC(), d.GetJobPQ(), d.GetJC())
	solo.Dispatch()
	rBytes, err := helpers.Serialize(solo.Result())
	if err != nil {
		return "", err
	}
	return string(rBytes), nil
}

//Chord executes execs one after the other then passes results into callback exec as a list (allows multiple jobs and multiple execs)
func (d Dispatcher) Chord(jrs []string, callbackJr string) (string, error) {
	//TODO: send result to message broker
	var requests []job.Request
	for _, jr := range jrs {
		var request job.Request
		err := helpers.Deserialize([]byte(jr), &request)
		if err != nil {
			return "", err
		}
		requests = append(requests, request)
	}
	var callbackRequest job.Request
	err := helpers.Deserialize([]byte(callbackJr), &callbackRequest)
	if err != nil {
		return "", err
	}
	c, err := chord.NewChord(requests, callbackRequest, d.GetBC(), d.GetJobPQ(), d.GetJC())
	if err != nil {
		return "", err
	}
	c.Dispatch()
	result, err := helpers.Serialize(c.Result())
	if err != nil {
		return "", err
	}
	return string(result), nil
}

//Chain executes execs one after the other (allows multiple jobs and multiple execs)
func (d Dispatcher) Chain(jrs []string) (string, error) {
	//TODO: send result to message broker
	var requests []job.Request
	for _, jr := range jrs {
		var request job.Request
		err := helpers.Deserialize([]byte(jr), &request)
		if err != nil {
			return "", err
		}
		requests = append(requests, request)
	}

	c, err := chain.NewChain(requests, d.GetBC(), d.GetJobPQ(), d.GetJC())
	if err != nil {
		return "", err
	}
	c.Dispatch()
	result, err := json.Marshal(c.Result())
	if err != nil {
		return "", err
	}
	return string(result), nil
}

//Batch executes execs in parallel
func (d Dispatcher) Batch(jrs []string) (string, error) {
	//TODO: send result to message broker
	var requests []job.Request
	for _, jr := range jrs {
		var request job.Request
		err := helpers.Deserialize([]byte(jr), &request)
		if err != nil {
			return "", err
		}
		requests = append(requests, request)
	}

	b, err := batch.NewBatch(requests, d.GetBC(), d.GetJobPQ(), d.GetJC())
	if err != nil {
		return "", err
	}
	b.Dispatch()
	result, err := json.Marshal(b.Result())
	if err != nil {
		return "", err
	}
	return string(result), nil
}
