package util

import (
	"os/exec"

	"github.com/stretchr/testify/assert"
)

// ExecutorInterface allows us to mock calls to exec.Cmd
type ExecutorInterface interface {
	Run(cmd *exec.Cmd) error
	Output(cmd *exec.Cmd) ([]byte, error)
	CombinedOutput(cmd *exec.Cmd) ([]byte, error)
}

// Executor is the struct used by most of the code: handlerResourceProvider, dynamic_lists, and system
type Executor struct{}

// Run will run the passed cmd command in the terminal
func (*Executor) Run(cmd *exec.Cmd) error {
	return cmd.Run()
}

// Output runs the cmd command in the terminal and returns its standard output
func (*Executor) Output(cmd *exec.Cmd) ([]byte, error) {
	return cmd.Output()
}

// CombinedOutput runs the cmd command in the terminal and returns its combined standard output and standard error
func (*Executor) CombinedOutput(cmd *exec.Cmd) ([]byte, error) {
	return cmd.CombinedOutput()
}

// testExecutorNoErr is a struct used for tests, where Run, Output, and CombinedOutput all return with no errors
type TestExecutorNoErr struct{}

func (*TestExecutorNoErr) Run(cmd *exec.Cmd) error {
	return nil
}
func (*TestExecutorNoErr) Output(cmd *exec.Cmd) ([]byte, error) {
	return []byte{}, nil
}
func (*TestExecutorNoErr) CombinedOutput(cmd *exec.Cmd) ([]byte, error) {
	return []byte{}, nil
}

// testExecutorRunErr is a struct used for tests where Run returns an error, but Output and CombinedOutput do not
type TestExecutorRunErr struct{}

func (*TestExecutorRunErr) Run(cmd *exec.Cmd) error {
	return assert.AnError
}
func (*TestExecutorRunErr) Output(cmd *exec.Cmd) ([]byte, error) {
	return []byte{}, nil
}
func (*TestExecutorRunErr) CombinedOutput(cmd *exec.Cmd) ([]byte, error) {
	return []byte{}, nil
}

// testExecutorCombinedOutputErr is a struct used for tetss where CombinedOutput returns an error, but Run and Output do not
type TestExecutorCombinedOutputErr struct{}

func (*TestExecutorCombinedOutputErr) Run(cmd *exec.Cmd) error {
	return nil
}
func (*TestExecutorCombinedOutputErr) Output(cmd *exec.Cmd) ([]byte, error) {
	return []byte{}, nil
}
func (*TestExecutorCombinedOutputErr) CombinedOutput(cmd *exec.Cmd) ([]byte, error) {
	return []byte{}, assert.AnError
}
