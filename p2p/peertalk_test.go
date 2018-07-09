package p2p_test

import (
	"encoding/hex"
	"testing"

	"github.com/gizo-network/gizo/helpers"

	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job"

	"github.com/gizo-network/gizo/p2p"
	"github.com/stretchr/testify/assert"
)

func TestHelloMessage(t *testing.T) {
	hm, err := p2p.HelloMessage([]byte("test"))
	assert.NoError(t, err)
	assert.NotNil(t, hm)
}

func TestInvalidMessage(t *testing.T) {
	im, err := p2p.InvalidMessage()
	assert.NoError(t, err)
	assert.NotNil(t, im)
}

func TestConnFullMessage(t *testing.T) {
	cf, err := p2p.ConnFullMessage()
	assert.NoError(t, err)
	assert.NotNil(t, cf)
}

func TestJobMessage(t *testing.T) {
	priv, pub := crypt.GenKeys()
	privBytes, err := hex.DecodeString(priv)
	assert.NoError(t, err)
	j, err := job.NewJob(`
	func Factorial(n){
	 if(n > 0){
	  result = n * Factorial(n-1)
	  return result
	 }
	 return 1
	}`, "Factorial", false, priv)
	assert.NoError(t, err)
	assert.NotNil(t, pub)
	res := make(chan qItem.Item)
	envs := job.NewEnvVariables(*job.NewEnv("Env", "Anko"), *job.NewEnv("By", "Lobarr"))
	exec, err := job.NewExec([]interface{}{2}, 5, job.NORMAL, 0, 0, 0, 0, pub, envs, "test")
	item := qItem.NewItem(*j, exec, res, exec.GetCancelChan())
	itemBytes, err := helpers.Serialize(item)
	assert.NoError(t, err)
	jmBytes, err := p2p.JobMessage(itemBytes, privBytes)
	assert.NoError(t, err)
	var jm p2p.PeerMessage
	err = helpers.Deserialize(jmBytes, &jm)
	assert.NoError(t, err)
	verify, err := jm.VerifySignature(pub)
	assert.NoError(t, err)
	assert.True(t, verify)
}
