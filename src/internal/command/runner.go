package command

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

type Result struct {
	Stdout string
	Stderr string
	Code   int
}

type Runner interface {
	Run(ctx context.Context, name string, args ...string) Result
	LookPath(name string) bool
}

type LocalRunner struct {
	Timeout time.Duration
}

func (r LocalRunner) Run(ctx context.Context, name string, args ...string) Result {
	timeout := r.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	res := Result{Stdout: stdout.String(), Stderr: stderr.String()}
	if err == nil {
		return res
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		res.Code = exit.ExitCode()
	} else {
		res.Code = 127
		if res.Stderr == "" {
			res.Stderr = err.Error()
		}
	}
	return res
}

func (r LocalRunner) LookPath(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func FirstLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
