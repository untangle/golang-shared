package mocks

import (
	"os/exec"

	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/util"
)

// MockExecutorNoErr is a struct used for tests, where Run, Output, and CombinedOutput all return with no errors
type MockExecutorNoErr struct{}

// type enforcement
var _ util.ExecutorInterface = (*MockExecutorNoErr)(nil)

func (m *MockExecutorNoErr) Run(cmd *exec.Cmd) error {
	return nil
}

func (*MockExecutorNoErr) Output(cmd *exec.Cmd) ([]byte, error) {
	return []byte{}, nil
}

func (*MockExecutorNoErr) CombinedOutput(cmd *exec.Cmd) ([]byte, error) {
	return []byte{}, nil
}

// MockExecutorRunErr is a struct used for tests where Run returns an error, but Output and CombinedOutput do not
type MockExecutorRunErr struct{ MockExecutorNoErr }

// type enforcement
var _ util.ExecutorInterface = (*MockExecutorRunErr)(nil)

func (*MockExecutorRunErr) Run(cmd *exec.Cmd) error {
	return assert.AnError
}

// MockExecutorCombinedOutputErr is a struct used for tetss where CombinedOutput returns an error, but Run and Output do not
type MockExecutorCombinedOutputErr struct{ MockExecutorNoErr }

// type enforcement
var _ util.ExecutorInterface = (*MockExecutorCombinedOutputErr)(nil)

func (*MockExecutorCombinedOutputErr) CombinedOutput(cmd *exec.Cmd) ([]byte, error) {
	return []byte{}, assert.AnError
}
