package util

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/untangle/golang-shared/services/logger"
)

// PluginStartup is placeholder for starting util
func PluginStartup() {

}

// PluginShutdown is placeholder for shutting down util
func PluginShutdown() {

}

// RunSighup will take the given executable and run sighup on it
// @param executable - executable to run sighup on
// @return any error from running
func RunSighup(executable string) error {
	return SendSignal(executable, syscall.SIGHUP)
}

// RunSigusr1 will take the given executable and run Sigusr1 on it
// @param executable - executable to run Sigusr1 on
// @return any error from running
func RunSigusr1(executable string) error {
	return SendSignal(executable, syscall.SIGUSR1)
}

// SendSignal will use a pid and send a signal to that pid
// @param excutable string - the binary process name
// @param signal syscall.Signal - the signal type to send
func SendSignal(executable string, signal syscall.Signal) error {

	logger.Debug("Sending %s to %s\n", signal, executable)

	pidStr, err := exec.Command("pgrep", executable).CombinedOutput()
	if err != nil {
		logger.Err("Failure to get %s pid: %s\n", executable, err.Error())
		return err
	}
	logger.Debug("Pid: %s\n", pidStr)

	pid, err := strconv.Atoi(strings.TrimSpace(string(pidStr)))
	if err != nil {
		logger.Err("Failure converting pid for %s: %s\n", executable, err.Error())
		return err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		logger.Err("Failure to get %d process: %s\n", pid, err.Error())
		return err
	}
	return process.Signal(signal)
}
