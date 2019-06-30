package job

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"reflect"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/gizo-network/gizo/helpers"

	"github.com/kpango/glg"
	anko_core "github.com/mattn/anko/builtins"
	anko_vm "github.com/mattn/anko/vm"
)

var (
	//ErrUnverifiedSignature when unauthorized person tries to execute private job
	ErrUnverifiedSignature = errors.New("Signature not verified")
	//ErrUnableToConvert when args is not able to be converted to string
	ErrUnableToConvert = errors.New("Unable to convert to string")
	//ErrUnableToSign unable to sign job
	ErrUnableToSign = errors.New("Unable to sign job")
)

//Job holds deployed job
type Job struct {
	ID             string
	Hash           string
	Execs          []IExec
	Name           string
	Task           string
	Signature      string // signature of owner
	SubmissionTime int64
	Private        bool //private job flag (default to false - public)
}

//UniqJob returns unique values of parameter
func UniqJob(jobs []Job) []Job {
	tempJob := []Job{}
	seen := make(map[string]bool)
	for _, job := range jobs {
		if _, ok := seen[job.GetID()]; ok {
			continue
		}
		seen[job.GetID()] = true
		tempJob = append(tempJob, job)
	}
	return tempJob
}

//Sign signature of owner of job
func (j *Job) Sign(priv string) error {
	task, err := j.GetTask()
	if err != nil {
		return err
	}
	hash := sha256.Sum256([]byte(task))
	privBytes, err := hex.DecodeString(priv)
	if err != nil {
		return err
	}
	privateKey, _ := x509.ParseECPrivateKey(privBytes)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return err
	}
	j.setSignature(hex.EncodeToString(r.Bytes()) + hex.EncodeToString(s.Bytes()))
	return nil
}

//VerifySignature verifies if pub is for the owner
func (j Job) VerifySignature(pub string) (bool, error) {
	pubBytes, err := hex.DecodeString(pub)
	if err != nil {
		return false, err
	}

	rBytes, err := hex.DecodeString(j.GetSignature()[:len(j.GetSignature())/2])
	if err != nil {
		return false, err
	}
	sBytes, err := hex.DecodeString(j.GetSignature()[len(j.GetSignature())/2:])
	if err != nil {
		return false, err
	}
	var r big.Int
	var s big.Int
	r.SetBytes(rBytes)
	s.SetBytes(sBytes)

	publicKey, _ := x509.ParsePKIXPublicKey(pubBytes)
	task, err := j.GetTask()
	if err != nil {
		return false, err
	}
	hash := sha256.Sum256([]byte(task))
	switch pubConv := publicKey.(type) {
	case *ecdsa.PublicKey:
		return ecdsa.Verify(pubConv, hash[:], &r, &s), nil
	default:
		return false, nil
	}
}

//GetSubmissionTime returns submission time
func (j Job) GetSubmissionTime() int64 {
	return j.SubmissionTime
}

func (j *Job) setSubmissionTime(t time.Time) {
	j.SubmissionTime = t.Unix()
}

//IsEmpty check if job is empty
func (j Job) IsEmpty() (bool, error) {
	task, err := j.GetTask()
	if err != nil {
		return false, err
	}
	return j.GetID() == "" && j.GetHash() == "" && reflect.ValueOf(j.GetExecs()).IsNil() && task == "" && j.GetSignature() == "" && j.GetName() == "", nil
}

//NewJob initializes a new job
func NewJob(task string, name string, priv bool, privKey string) (*Job, error) {
	j := &Job{
		SubmissionTime: time.Now().Unix(),
		ID:             uuid.NewV4().String(),
		Execs:          []IExec{},
		Name:           name,
		Task:           helpers.Encode64([]byte(task)),
		Private:        priv,
	}
	err := j.Sign(privKey)
	if err != nil {
		return nil, err
	}
	j.setHash()
	return j, nil
}

//GetPrivate check if job is private
func (j Job) GetPrivate() bool {
	return j.Private
}

func (j *Job) setPrivate(p bool) {
	j.Private = p
}

//GetName returns name of job
func (j Job) GetName() string {
	return j.Name
}

func (j *Job) setName(n string) {
	j.Name = n
}

//GetID returns job id
func (j Job) GetID() string {
	return j.ID
}

//GetHash return job hash
func (j Job) GetHash() string {
	return j.Hash
}

func (j *Job) setHash() error {
	rBytes, err := hex.DecodeString(j.GetSignature()[:len(j.GetSignature())/2])
	if err != nil {
		return err
	}
	sBytes, err := hex.DecodeString(j.GetSignature()[len(j.GetSignature())/2:])
	if err != nil {
		return err
	}
	task, err := j.GetTask()
	if err != nil {
		return err
	}
	headers := bytes.Join(
		[][]byte{
			[]byte(j.GetID()),
			[]byte(task),
			[]byte(j.GetName()),
			rBytes,
			sBytes,
			[]byte(string(j.GetSubmissionTime())),
			[]byte(strconv.FormatBool(j.GetPrivate())),
		},
		[]byte{},
	)
	hash := sha256.Sum256(headers)
	j.Hash = hex.EncodeToString(hash[:])
	return nil
}

//Verify checks if the job has been modified
func (j Job) Verify() (bool, error) {
	rBytes, err := hex.DecodeString(j.GetSignature()[:len(j.GetSignature())/2])
	if err != nil {
		return false, err
	}
	sBytes, err := hex.DecodeString(j.GetSignature()[len(j.GetSignature())/2:])
	if err != nil {
		return false, err
	}
	task, err := j.GetTask()
	if err != nil {
		return false, err
	}
	headers := bytes.Join(
		[][]byte{
			[]byte(j.GetID()),
			[]byte(task),
			[]byte(j.GetName()),
			rBytes,
			sBytes,
			[]byte(string(j.GetSubmissionTime())),
			[]byte(strconv.FormatBool(j.GetPrivate())),
		},
		[]byte{},
	)
	hash := sha256.Sum256(headers)
	return j.GetHash() == hex.EncodeToString(hash[:]), nil
}

//return json bytes of the execs
func (j Job) serializeExecs() []byte {
	temp, err := json.Marshal(j.GetExecs())
	if err != nil {
		glg.Error(err)
	}
	return temp
}

//GetExec return exec of specified hash
func (j Job) GetExec(hash string) (IExec, error) {
	for _, exec := range j.GetExecs() {
		if exec.GetHash() == hash {
			return exec, nil
		}
	}
	return nil, ErrExecNotFound
}

//GetLatestExec return latest exec of job
func (j Job) GetLatestExec() IExec {
	return j.Execs[len(j.GetExecs())-1]
}

//GetExecs returns exec of job
func (j Job) GetExecs() []IExec {
	return j.Execs
}

//GetSignature returns job signature
func (j Job) GetSignature() string {
	return j.Signature
}

func (j *Job) setSignature(sign string) {
	j.Signature = sign
}

//AddExec add exec to job
func (j *Job) AddExec(jobExec IExec) {
	j.Execs = append(j.Execs, jobExec)
}

//GetTask returns task
func (j Job) GetTask() (string, error) {
	decode, err := helpers.Decode64(j.Task)
	if err != nil {
		return "", err
	}
	return string(decode), nil
}

//used to stringify arguments
func toString(x interface{}) (string, error) {
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), nil
	case reflect.String:
		return "\"" + x.(string) + "\"", nil
	case reflect.Slice:
		stringifiy, err := json.Marshal(v.Interface())
		return string(stringifiy), err
	case reflect.Map:
		stringifiy, err := json.Marshal(v.Interface())
		return string(stringifiy), err
	}
	return "", ErrUnableToConvert
}

//returns args as string
func argsStringified(args []interface{}) (string, error) {
	temp := "("
	for i, val := range args {
		stringified, err := toString(val)
		if err != nil {
			return "", err
		}
		if i == len(args)-1 {
			temp += stringified + ""
		} else {
			temp += stringified + ","
		}
	}
	return temp + ")", nil
}

//Execute runs the exec
func (j *Job) Execute(exec IExec, passphrase string) IExec {
	//TODO: kill goroutines running within this function when it exits
	if j.GetPrivate() == true {
		verify, err := j.VerifySignature(exec.GetPub())
		if err != nil {
			exec.SetErr(err)
		}
		if verify == false {
			exec.SetErr(ErrUnverifiedSignature)
			return exec
		}
	}
	glg.Info("Job: Executing job - " + j.GetID())
	start := time.Now()
	cancelClose := make(chan struct{})
	timeoutClose := make(chan struct{})
	execClose := make(chan struct{})
	done := make(chan string)
	execute := make(chan struct{})
	routines := make(map[string]chan struct{})
	routines["cancel"] = cancelClose
	routines["timeout"] = timeoutClose
	routines["exec"] = execClose
	exec.SetStatus(RUNNING)
	exec.SetTimestamp(time.Now().Unix())
	//! watch for cancellation
	go func() {
		select {
		case <-cancelClose:
			return
		case <-exec.GetCancelChan():
			done <- "cancel"
		}
	}()
	//! watch for timeout
	go func() {
		var ttl time.Duration
		if exec.GetTTL() != 0 {
			ttl = exec.GetTTL()
		} else {
			ttl = DefaultMaxTTL
		}
		select {
		case <-timeoutClose:
			return
		case <-time.NewTimer(ttl).C:
			exec.SetStatus(TIMEOUT)
			glg.Warn("Job: Job timeout - " + j.GetID())
			done <- "timeout"
		}
	}()
	//! execute job
	go func() {
		select {
		case <-execClose:
			return
		case <-execute:
			r := exec.GetRetries()
		retry:
			//TODO: support tmp directory for saved files
			//TODO: clean tmp directory every 10 mins
			env := anko_vm.NewEnv()
			anko_core.LoadAllBuiltins(env) //!FIXME: switch to gizo-network/anko
			envs, err := exec.GetEnvsMap(passphrase)
			task, err := j.GetTask()
			var result reflect.Value
			if err == nil {
				env.Define("env", envs)
				if len(exec.GetArgs()) == 0 {
					result, err = env.Execute(task + "\n" + j.GetName() + "()")
				} else {
					args, err := argsStringified(exec.GetArgs())
					if err == nil {
						result, err = env.Execute(task + "\n" + j.GetName() + args)
					}
				}

				if r != 0 && err != nil {
					r--
					time.Sleep(exec.GetBackoff())
					exec.SetStatus(RETRYING)
					exec.IncrRetriesCount()
					glg.Warn("Job: Retrying job - " + j.GetID())
					goto retry
				}
			}
			exec.SetDuration(time.Duration(time.Now().Sub(start).Nanoseconds()))
			exec.SetErr(err)
			exec.SetResult(result.Interface())
			exec.SetHash()
			exec.SetStatus(FINISHED)
			done <- "exec"
		}
	}()
	execute <- struct{}{}
	select {
	case complete := <-done:
		for routine, channel := range routines {
			if routine != complete { // sends close signal to other goroutes
				channel <- struct{}{}
			}
		}
		return exec
	}
}
