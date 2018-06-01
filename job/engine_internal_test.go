package job

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToString(t *testing.T) {
	str, err := toString(true)
	assert.NoError(t, err)
	assert.IsType(t, string(""), str)

	str, err = toString(1)
	assert.NoError(t, err)
	assert.IsType(t, string(""), str)

	str, err = toString(1.632)
	assert.NoError(t, err)
	assert.IsType(t, string(""), str)

	str, err = toString(-1)
	assert.NoError(t, err)
	assert.IsType(t, string(""), str)

	str, err = toString("test")
	assert.NoError(t, err)
	assert.IsType(t, string(""), str)
	assert.Equal(t, "\"test\"", str)

	str, err = toString([]string{"test", "test"})
	assert.NoError(t, err)
	assert.IsType(t, string(""), str)

	str, err = toString(map[string]string{"test": "test"})
	assert.NoError(t, err)
	assert.IsType(t, string(""), str)
}

func TestArgsStringified(t *testing.T) {
	args, err := argsStringified([]interface{}{true, 1, 1.632, -1, "test", []string{"test", "test"}, map[string]string{"test": "test"}})
	assert.NoError(t, err)
	assert.NotNil(t, args)
	assert.IsType(t, string(""), args)
}
