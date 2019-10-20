package just_test

import (
	"os/exec"
	"testing"

	"github.com/julz/just"
	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
)

func TestOutput(t *testing.T) {
	// success case
	result := just.GetStdout(exec.Command("echo", "hello"))
	assert.Equal(t, "hello\n", string(result))

	// failure case
	var gotError error
	just.FailHandler = func(err error) {
		gotError = err
	}

	just.GetStdout(exec.Command("bash", "-c", "exit 3"))
	assert.Check(t, cmp.ErrorContains(gotError, "exit status 3"))
}

func TestDecodeJSONOutput(t *testing.T) {
	// success case
	result := make(map[string]interface{})
	just.DecodeJSONOutput(exec.Command("echo", `{"test": 3}`), &result)
	assert.Equal(t, 3.0, result["test"])

	// failure case
	var gotError error
	just.FailHandler = func(err error) {
		gotError = err
	}

	result = make(map[string]interface{})
	just.DecodeJSONOutput(exec.Command("echo", `invalid json`), &result)
	assert.Check(t, cmp.ErrorContains(gotError, "decode json: invalid character"))
}
