package helpers

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/minifyjs/minifyjs/core/internal/version"
)

// Result is what running the CLI produces. The field names mirror
// the exec.ExitError shape so callers that already know Go's exec
// package feel at home.
type Result struct {
	ExitCode int
	Stdout   []byte
	Stderr   []byte
}

// RunOptions configures a single CLI invocation.
type RunOptions struct {
	// Args are the command-line arguments, without the program name.
	Args []string

	// Stdin, if non-nil, is fed to the process's stdin.
	Stdin []byte

	// Dir, if non-empty, is the working directory. Defaults to the
	// current directory.
	Dir string

	// Env, if non-nil, replaces the process environment entirely.
	// Each entry is "KEY=VALUE".
	Env []string
}


func RunWithTimeout(t *testing.T, timeout time.Duration, opts RunOptions) Result {
	t.Helper()
	cmd := exec.Command(CLIPath(t), opts.Args...)
	if opts.Dir != "" {
		cmd.Dir = opts.Dir
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start CLI: %v", err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case <-done:
	case <-time.After(timeout):
		_ = cmd.Process.Kill()
		<-done
	}
	return Result{
		ExitCode: cmd.ProcessState.ExitCode(),
		Stdout:   stdout.Bytes(),
		Stderr:   stderr.Bytes(),
	}
}

// Run executes the compiled CLI with the given options and returns
// the exit code, stdout, and stderr. It fails the test only if the
// process could not be started; a non-zero exit code is a valid
// result and is reported through Result.ExitCode.
func Run(t *testing.T, opts RunOptions) Result {
	t.Helper()

	cmd := exec.Command(CLIPath(t), opts.Args...)
	if opts.Dir != "" {
		cmd.Dir = opts.Dir
	}
	if opts.Env != nil {
		cmd.Env = opts.Env
	} else {
		cmd.Env = os.Environ()
	}
	if opts.Stdin != nil {
		cmd.Stdin = bytes.NewReader(opts.Stdin)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := Result{
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
	}
	if err == nil {
		result.ExitCode = 0
		return result
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		result.ExitCode = exitErr.ExitCode()
		return result
	}
	t.Fatalf("run CLI: %v", err)
	return result
}

// Version returns the version string the current build reports. It
// is used by tests that assert on --version output.
func Version() string {
	return version.Version
}