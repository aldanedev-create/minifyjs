package main

import "fmt"

// cliError is a structured error for the CLI's own validation (as
// opposed to errors coming out of the engine, which are already
// well-typed). It exists so future versions can attach exit codes
// and error IDs to CLI-level failures without re-threading the type
// through every function.
type cliError struct {
	Code    string // e.g. "cli.missing-value"
	Message string
}

func (e cliError) Error() string {
	if e.Code == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func newCLIError(code, format string, args ...any) error {
	return cliError{Code: code, Message: fmt.Sprintf(format, args...)}
}