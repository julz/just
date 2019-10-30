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
func Run(cmd *exec.Cmd, options ...CmdOption) {
	for _, option := range options {
		option(cmd)
	}

	Check(cmd.Run())
}

// RunSh just runs the given command in a shell.
// If the command fails, the `FailHandler` is called with a nice error.
func RunSh(cmd string, options ...CmdOption) {
	Run(Sh(cmd), options...)
}

// GetStdout runs the given `exec.Cmd` and returns the resulting stdout.
// If the command fails, the `FailHandler` is called with a nice error.
func GetStdout(cmd *exec.Cmd, options ...CmdOption) []byte {
	var stdout, stderr bytes.Buffer
	options = append(options, Out(&stdout), Err(&stderr))

	for _, option := range options {
		option(cmd)
	}

	if err := cmd.Run(); err != nil {
		CheckP("output", fmt.Errorf("%s; stderr: %s", err, string(stderr.Bytes())))
	}

	return stdout.Bytes()
}

// DecodeJSON decodes JSON from `r` into `result` and calls the FailHandler on
// errors
func DecodeJSON(r io.Reader, result interface{}) {
	CheckP("decode json", json.NewDecoder(r).Decode(result))
}

// DecodeJSONOutput runs the given `exec.Cmd` and decodes the stdout as json into `result`.
// If either running the command or decoding the json errors, the FailHandler is called
// with a nicely wrapped error.
func DecodeJSONOutput(cmd *exec.Cmd, result interface{}, options ...CmdOption) {
	DecodeJSON(bytes.NewReader(GetStdout(cmd, options...)), result)
}

// DecodeOutput runs the given `cmd` in a shell and decodes the stdout as json
// into `result`.  If either running the command or decoding the json errors,
// the FailHandler is called with a nicely wrapped error.
func DecodeJSONOutputSh(cmd string, result interface{}, options ...CmdOption) {
	DecodeJSON(bytes.NewReader(GetStdout(Sh(cmd), options...)), result)
}

// Sh returns an *exec.Cmd that runs the given `cmd` argument in a shell
func Sh(cmd string) *exec.Cmd {
	return exec.Command("sh", "-c", cmd)
}

type CmdOption func(*exec.Cmd)

// Out sets the Stdout of the command the the given writer
func Out(w io.Writer) func(*exec.Cmd) {
	return func(c *exec.Cmd) {
		if c.Stdout != nil {
			w = io.MultiWriter(w, c.Stdout)
		}

		c.Stdout = w
	}
}

// Out sets the Stderr of the command the the given writer
func Err(w io.Writer) func(*exec.Cmd) {
	return func(c *exec.Cmd) {
		if c.Stderr != nil {
			w = io.MultiWriter(w, c.Stderr)
		}

		c.Stderr = w
	}
}
