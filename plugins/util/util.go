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
	logger.Debug("Running interrupt\n")
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
		logger.Err("Failure to get %s process: %s\n", executable, err.Error())
		return err
	}
	return process.Signal(syscall.SIGHUP)
}
