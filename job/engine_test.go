package job_test

import (
	"testing"

	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestUniqJob(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, err := job.NewJob(`
	func Factorial(n){
	 if(n > 0){
	  result = n * Factorial(n-1)
	  return result
	 }
	 return 1
	}`, "Factorial", false, priv)
	assert.NoError(t, err)
	assert.Len(t, job.UniqJob([]job.Job{*j, *j, *j, *j}), 1)
}

func TestSign(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, err := job.NewJob(`
	func Factorial(n){
	 if(n > 0){
	  result = n * Factorial(n-1)
	  return result
	 }
	 return 1
	}`, "Factorial", false, priv)
	assert.NoError(t, err)
	assert.NoError(t, j.Sign(priv))
	assert.NotNil(t, j.GetSignature())
}

func TestVerifySignature(t *testing.T) {
	_, _pub := crypt.GenKeys()
	priv, pub := crypt.GenKeys()
	j, err := job.NewJob(`
	func Factorial(n){
	 if(n > 0){
	  result = n * Factorial(n-1)
	  return result
	 }
	 return 1
	}`, "Factorial", false, priv)
	assert.NoError(t, err)
	ver, err := j.VerifySignature(pub)
	assert.NoError(t, err)
	assert.True(t, ver)

	ver, err = j.VerifySignature(_pub)
	assert.NoError(t, err)
	assert.False(t, ver)
}

func TestIsEmpty(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, err := job.NewJob(`
	func Factorial(n){
	 if(n > 0){
	  result = n * Factorial(n-1)
	  return result
	 }
	 return 1
	}`, "Factorial", false, priv)

	assert.NoError(t, err)
	empty, err := j.IsEmpty()
	assert.NoError(t, err)
	assert.False(t, empty)

	j2 := job.Job{}
	empty, err = j2.IsEmpty()
	assert.NoError(t, err)
	assert.True(t, empty)
}

func TestNewJob(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, err := job.NewJob(`
	func Factorial(n){
	 if(n > 0){
	  result = n * Factorial(n-1)
	  return result
	 }
	 return 1
	}`, "Factorial", false, priv)
	task, err := j.GetTask()
	assert.NoError(t, err)
	assert.NoError(t, err)
	assert.NotNil(t, j.GetSubmissionTime())
	assert.NotNil(t, j.GetID())
	assert.NotNil(t, j.GetExecs())
	assert.NotNil(t, j.GetName())
	assert.NotNil(t, task)
	assert.NotNil(t, j.GetSignature())
	assert.NotNil(t, j.GetHash())
}

func TestVerify(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, err := job.NewJob(`
	func Factorial(n){
	 if(n > 0){
	  result = n * Factorial(n-1)
	  return result
	 }
	 return 1
	}`, "Factorial", false, priv)
	assert.NoError(t, err)
	ver, err := j.Verify()
	assert.NoError(t, err)
	assert.True(t, ver)

	envs := job.NewEnvVariables(*job.NewEnv("Env", "Anko"), *job.NewEnv("By", "Lobarr"))
	exec, err := job.NewExec([]interface{}{10}, 5, job.NORMAL, 0, 0, 0, 0, "", envs, "test")
	j.AddExec(*exec)
	ver, err = j.Verify()
	assert.NoError(t, err)
	assert.True(t, ver)

	j.Name = "something else"
	ver, err = j.Verify()
	assert.NoError(t, err)
	assert.False(t, ver)
}

func TestGetExec(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, err := job.NewJob(`
	func Factorial(n){
	 if(n > 0){
	  result = n * Factorial(n-1)
	  return result
	 }
	 return 1
	}`, "Factorial", false, priv)
	envs := job.NewEnvVariables(*job.NewEnv("Env", "Anko"), *job.NewEnv("By", "Lobarr"))
	exec, err := job.NewExec([]interface{}{10}, 5, job.NORMAL, 0, 0, 0, 0, "", envs, "test")
	j.AddExec(*exec)

	_exec, err := j.GetExec(exec.GetHash())
	assert.NoError(t, err)
	assert.True(t, assert.ObjectsAreEqualValues(exec, _exec))

	__exec, err := j.GetExec("something else")
	assert.Error(t, err)
	assert.Nil(t, __exec)
}

func TestExecute(t *testing.T) {
	priv, _ := crypt.GenKeys()
	j, err := job.NewJob(`
	func Factorial(n){
	 if(n > 0){
	  result = n * Factorial(n-1)
	  return result
	 }
	 return 1
	}`, "Factorial", false, priv)
	assert.NoError(t, err)
	ver, err := j.Verify()
	assert.NoError(t, err)
	assert.True(t, ver)

	envs := job.NewEnvVariables(*job.NewEnv("Env", "Anko"), *job.NewEnv("By", "Lobarr"))
	exec, err := job.NewExec([]interface{}{10}, 5, job.NORMAL, 0, 0, 0, 0, "", envs, "test")

	_exec := j.Execute(exec, "test")
	assert.Nil(t, _exec.GetErr())
	assert.NotNil(t, _exec.GetResult()) //3628800
	assert.Equal(t, int64(3628800), _exec.GetResult().(int64))
}
