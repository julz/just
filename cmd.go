package just

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
)

var FailHandler = func(err error) {
	log.Fatal(err)
}

func Check(err error) {
	if err != nil {
		FailHandler(err)
	}
}

func CheckP(prefix string, err error) {
	if err != nil {
		FailHandler(fmt.Errorf("%s: %s", prefix, err))
	}
}

// Run just runs the given `exec.Cmd`.
// If the command fails, the `FailHandler` is called with a nice error.
func Run(cmd *exec.Cmd) {
	Check(cmd.Run())
}

// GetStdout runs the given `exec.Cmd` and returns the resulting stdout.
// If the command fails, the `FailHandler` is called with a nice error.
func GetStdout(cmd *exec.Cmd) (out []byte) {
	out, err := cmd.Output()
	if err != nil {
		switch eerr := err.(type) {
		case *exec.ExitError:
			CheckP("output", fmt.Errorf("%s; stderr: %s", err, string(eerr.Stderr)))
		default:
			CheckP("output", err)
		}
	}

	return out
}

// DecodeJSON decodes JSON from `r` into `result` and calls the FailHandler on
// errors
func DecodeJSON(r io.Reader, result interface{}) {
	CheckP("decode json", json.NewDecoder(r).Decode(result))
}

// DecodeOutput runs the given `exec.Cmd` and decodes the stdout as json into `result`.
// If either running the command or decoding the json errors, the FailHandler is called
// with a nicely wrapped error.
func DecodeJSONOutput(cmd *exec.Cmd, result interface{}) {
	DecodeJSON(bytes.NewReader(GetStdout(cmd)), result)
}
