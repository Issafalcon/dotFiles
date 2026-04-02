// Package utils provides shared helper functions for the DotFiles TUI application.
//
// This file implements async command execution with streaming output,
// context-based cancellation, and structured result capture.
//
// Key Go concepts used here:
//   - os/exec: Running external commands (https://pkg.go.dev/os/exec)
//   - context.Context: Cancellation and timeout propagation (https://pkg.go.dev/context)
//   - Goroutines: Lightweight concurrent functions (https://go.dev/doc/effective_go#goroutines)
//   - Channels: Typed conduits for communication between goroutines (https://go.dev/doc/effective_go#channels)
//   - select: Multiplexing channel operations (https://go.dev/ref/spec#Select_statements)
//
// See also: https://go.dev/doc/effective_go
package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

// CommandResult holds the complete output of a finished command.
//
// In Go, structs are value types that group related data together.
// Fields starting with an uppercase letter are "exported" (public),
// while lowercase fields are unexported (private to the package).
// See: https://go.dev/doc/effective_go#names
// See: https://go.dev/ref/spec#Struct_types
type CommandResult struct {
	// Stdout contains all lines captured from the command's standard output.
	// []string is a slice — a dynamically-sized, flexible view into an array.
	// See: https://go.dev/doc/effective_go#slices
	Stdout []string

	// Stderr contains all lines captured from the command's standard error.
	Stderr []string

	// ExitCode is the process exit code (0 = success, non-zero = failure).
	// Go uses int as a general-purpose integer type.
	// See: https://go.dev/ref/spec#Numeric_types
	ExitCode int

	// Err holds any Go-level error that occurred (e.g., command not found,
	// context cancelled). This is separate from a non-zero exit code.
	//
	// The error interface is Go's conventional way to represent error conditions.
	// See: https://go.dev/doc/effective_go#errors
	// See: https://pkg.go.dev/errors
	Err error
}

// RunCommand executes a shell command synchronously and returns the collected result.
//
// Parameters:
//   - ctx: A context.Context for cancellation and timeouts. If the context is
//     cancelled or times out, the command process is killed.
//     See: https://pkg.go.dev/context
//   - command: The shell command string to execute (passed to "sh -c").
//
// Returns:
//   - CommandResult: The collected stdout, stderr, exit code, and any error.
//   - error: A Go-level error if the command could not be started at all.
//
// Example usage:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	result, err := utils.RunCommand(ctx, "ls -la")
//
// In Go, functions can return multiple values. This is idiomatic for returning
// both a result and an error. Callers must check the error before using the result.
// See: https://go.dev/doc/effective_go#multiple-returns
func RunCommand(ctx context.Context, command string) (CommandResult, error) {
	// exec.CommandContext creates a command that will be killed if ctx is cancelled.
	// We use "sh -c" to execute the command string through a shell, which allows
	// shell features like pipes (|), redirects (>), and variable expansion ($VAR).
	// See: https://pkg.go.dev/os/exec#CommandContext
	cmd := exec.CommandContext(ctx, "sh", "-c", command)

	// SysProcAttr sets OS-level process attributes. Setpgid=true puts the child
	// process in its own process group, allowing us to kill it and all its children.
	// See: https://pkg.go.dev/syscall#SysProcAttr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// NOTE: We intentionally do NOT set cmd.Stdin = os.Stdin here.
	// Interactive commands (e.g. sudo prompts) are handled via tea.ExecProcess
	// which suspends Bubble Tea and gives the subprocess full terminal control.
	// Setting Stdin here would compete with Bubble Tea for input.

	// StdoutPipe and StderrPipe return io.ReadCloser interfaces connected to the
	// command's stdout/stderr. We read from these pipes to capture output.
	//
	// In Go, interfaces define behavior (method sets) without specifying implementation.
	// io.ReadCloser combines the io.Reader and io.Closer interfaces.
	// See: https://pkg.go.dev/io#ReadCloser
	// See: https://go.dev/doc/effective_go#interfaces
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		// Return a zero-value CommandResult and the error.
		// In Go, every type has a zero value (empty struct for CommandResult).
		// See: https://go.dev/ref/spec#The_zero_value
		return CommandResult{}, fmt.Errorf("creating stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return CommandResult{}, fmt.Errorf("creating stderr pipe: %w", err)
	}

	// Start the command (does not wait for it to finish).
	// See: https://pkg.go.dev/os/exec#Cmd.Start
	if err := cmd.Start(); err != nil {
		return CommandResult{}, fmt.Errorf("starting command: %w", err)
	}

	// We need to read stdout and stderr concurrently because both pipes might
	// fill their OS buffers. If we read one sequentially after the other, the
	// unread pipe could block the child process (deadlock).
	//
	// A sync.WaitGroup waits for a collection of goroutines to finish.
	// It's like a concurrent counter: Add(n) increments, Done() decrements,
	// Wait() blocks until the counter reaches zero.
	// See: https://pkg.go.dev/sync#WaitGroup
	// See: https://go.dev/doc/effective_go#goroutines
	var wg sync.WaitGroup

	// Variables to collect output from goroutines.
	var stdoutLines, stderrLines []string

	// A sync.Mutex provides mutual exclusion — only one goroutine can hold the
	// lock at a time. We use it to protect shared slice access.
	// See: https://pkg.go.dev/sync#Mutex
	var mu sync.Mutex

	// Add 2 to the WaitGroup because we're about to launch 2 goroutines.
	wg.Add(2)

	// Launch a goroutine to read stdout.
	//
	// A goroutine is a lightweight thread managed by the Go runtime. You start one
	// with the "go" keyword followed by a function call. Goroutines are multiplexed
	// onto OS threads and are very cheap to create (~2KB stack).
	// See: https://go.dev/doc/effective_go#goroutines
	// See: https://go.dev/tour/concurrency/1
	go func() {
		// defer schedules a function call to run when the enclosing function returns.
		// This ensures wg.Done() is called even if readLines panics.
		// See: https://go.dev/doc/effective_go#defer
		defer wg.Done()

		lines := readLines(stdoutPipe)
		// Lock the mutex before modifying the shared variable.
		mu.Lock()
		stdoutLines = lines
		// Unlock the mutex so other goroutines can access the variable.
		mu.Unlock()
	}()

	// Launch a goroutine to read stderr.
	go func() {
		defer wg.Done()

		lines := readLines(stderrPipe)
		mu.Lock()
		stderrLines = lines
		mu.Unlock()
	}()

	// Wait for both reader goroutines to finish before calling cmd.Wait().
	// cmd.Wait() closes the pipes, so we must finish reading first.
	wg.Wait()

	// cmd.Wait() waits for the command to exit and releases associated resources.
	// See: https://pkg.go.dev/os/exec#Cmd.Wait
	exitCode := 0
	cmdErr := cmd.Wait()

	if cmdErr != nil {
		// Type assertion: Check if the error is an *exec.ExitError, which contains
		// the process exit code. In Go, type assertions extract the concrete type
		// from an interface value.
		// See: https://go.dev/doc/effective_go#interface_conversions
		// See: https://go.dev/tour/methods/15
		var exitErr *exec.ExitError
		if isExitError(cmdErr, &exitErr) {
			exitCode = exitErr.ExitCode()
		} else {
			// A non-ExitError means something else went wrong (e.g., signal, I/O error).
			return CommandResult{
				Stdout:   stdoutLines,
				Stderr:   stderrLines,
				ExitCode: -1,
				Err:      cmdErr,
			}, cmdErr
		}
	}

	return CommandResult{
		Stdout:   stdoutLines,
		Stderr:   stderrLines,
		ExitCode: exitCode,
		Err:      cmdErr,
	}, nil
}

// isExitError checks if an error is an *exec.ExitError using errors.As-style logic.
//
// We use this helper because errors.As requires import of the errors package and
// this keeps the type assertion pattern visible for learning purposes.
// See: https://pkg.go.dev/errors#As
func isExitError(err error, target **exec.ExitError) bool {
	// A "comma ok" type assertion returns (value, bool) instead of panicking
	// if the assertion fails. This is a safe way to check interface types.
	// See: https://go.dev/ref/spec#Type_assertions
	exitErr, ok := err.(*exec.ExitError)
	if ok {
		*target = exitErr
	}
	return ok
}

// RunCommandStreaming executes a shell command and calls onLine for each line of output
// as it is produced. This enables real-time output display in the TUI.
//
// Parameters:
//   - ctx: Context for cancellation/timeout.
//   - command: Shell command string.
//   - onLine: Callback function invoked for each line of output.
//     The isStderr parameter indicates whether the line came from stderr.
//
// In Go, functions are first-class values — you can pass them as arguments,
// assign them to variables, and return them from other functions.
// See: https://go.dev/doc/effective_go#functions
// See: https://go.dev/tour/moretypes/24 (function closures)
func RunCommandStreaming(ctx context.Context, command string, onLine func(line string, isStderr bool)) error {
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// NOTE: We intentionally do NOT set cmd.Stdin here — see RunCommand above.

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("creating stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("creating stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting command: %w", err)
	}

	// Channels are typed conduits through which you can send and receive values
	// with the channel operator (<-). They synchronize goroutines without explicit
	// locks.
	//
	// make(chan Type) creates an unbuffered channel — sends block until another
	// goroutine receives, and vice versa. This provides synchronization.
	//
	// make(chan Type, capacity) creates a buffered channel — sends only block
	// when the buffer is full.
	//
	// See: https://go.dev/doc/effective_go#channels
	// See: https://go.dev/tour/concurrency/2
	// See: https://go.dev/ref/spec#Channel_types
	done := make(chan error, 2) // buffered with capacity 2 for the two reader goroutines

	// streamLines reads lines from a pipe and calls onLine for each.
	// This is a closure — it captures the onLine and done variables from the
	// enclosing scope.
	// See: https://go.dev/tour/moretypes/25
	streamLines := func(pipe io.ReadCloser, isStderr bool) {
		// bufio.NewScanner creates a scanner that reads input line by line.
		// See: https://pkg.go.dev/bufio#Scanner
		scanner := bufio.NewScanner(pipe)

		// scanner.Scan() advances to the next line and returns true if successful.
		// It returns false at EOF or on error.
		for scanner.Scan() {
			line := scanner.Text()
			onLine(line, isStderr)
		}

		// Send any scanner error (or nil) into the done channel.
		// The <- operator sends a value into a channel.
		done <- scanner.Err()
	}

	// Launch two goroutines to stream stdout and stderr concurrently.
	go streamLines(stdoutPipe, false)
	go streamLines(stderrPipe, true)

	// Wait for both streaming goroutines to finish.
	//
	// The select statement lets a goroutine wait on multiple channel operations.
	// It blocks until one of its cases can proceed, then executes that case.
	// If multiple cases are ready, one is chosen at random.
	// See: https://go.dev/ref/spec#Select_statements
	// See: https://go.dev/tour/concurrency/5
	//
	// Here we don't use select because we need to wait for exactly 2 completions
	// sequentially. But we demonstrate that receiving from a channel (<-done)
	// blocks until a value is available.
	var scanErr error
	for i := 0; i < 2; i++ {
		// Receive from the done channel. This blocks until one of the goroutines sends.
		// The <- operator on the left side of = receives a value from a channel.
		if err := <-done; err != nil {
			scanErr = err
		}
	}

	// Wait for the command to finish after all output has been read.
	cmdErr := cmd.Wait()

	// Return the most relevant error.
	if cmdErr != nil {
		return cmdErr
	}
	return scanErr
}

// readLines reads all lines from an io.Reader and returns them as a string slice.
//
// This is a private (unexported) helper function. In Go, functions starting with
// a lowercase letter are only visible within their package.
// See: https://go.dev/doc/effective_go#names
func readLines(r io.Reader) []string {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		// strings.TrimRight removes trailing whitespace characters.
		// See: https://pkg.go.dev/strings#TrimRight
		lines = append(lines, strings.TrimRight(scanner.Text(), "\r\n"))
	}
	return lines
}
