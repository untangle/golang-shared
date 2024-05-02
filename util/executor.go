package util

import (
	"os/exec"
)

// ExecutorInterface allows us to mock calls to exec.Cmd
type ExecutorInterface interface {
	Run(cmd *exec.Cmd) error
	Output(cmd *exec.Cmd) ([]byte, error)
	CombinedOutput(cmd *exec.Cmd) ([]byte, error)
}

// Executor is the struct used by most of the code: handlerResourceProvider, dynamic_lists, and system
type Executor struct{}

var _ ExecutorInterface = (*Executor)(nil) // type enforcement

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
