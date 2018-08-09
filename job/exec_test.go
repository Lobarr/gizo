package job_test

import (
	"testing"

	"github.com/gizo-network/gizo/job"
	"github.com/stretchr/testify/assert"
)

func TestNewExec(t *testing.T) {
	envs := job.NewEnvVariables(*job.NewEnv("Env", "Anko"), *job.NewEnv("By", "Lobarr"))
	_, err := job.NewExec([]interface{}{10}, 5, job.NORMAL, 0, 0, 0, 0, "", envs, "test")
	assert.NoError(t, err)

	_, err = job.NewExec([]interface{}{10}, 10, job.NORMAL, 0, 0, 0, 0, "", envs, "test")
	assert.Error(t, err)
}

func TestGetEnvs(t *testing.T) {
	envs := job.NewEnvVariables(*job.NewEnv("Env", "Anko"), *job.NewEnv("By", "Lobarr"))
	exec, err := job.NewExec([]interface{}{10}, 5, job.NORMAL, 0, 0, 0, 0, "", envs, "test")
	assert.NoError(t, err)
	_envs, err := exec.GetEnvs("test")
	assert.NoError(t, err)
	assert.True(t, assert.ObjectsAreEqual(envs, _envs))

	_, err = exec.GetEnvs("something else")
	assert.Error(t, err)
}

func TestUniqExec(t *testing.T) {
	envs := job.NewEnvVariables(*job.NewEnv("Env", "Anko"), *job.NewEnv("By", "Lobarr"))
	exec, err := job.NewExec([]interface{}{10}, 5, job.NORMAL, 0, 0, 0, 0, "", envs, "test")
	assert.NoError(t, err)
	assert.Len(t, job.UniqExec([]job.Exec{*exec, *exec, *exec, *exec}), 1)
}
