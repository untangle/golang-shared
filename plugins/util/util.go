package util

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"

	loggerModel "github.com/untangle/golang-shared/logger"
	"github.com/untangle/golang-shared/util/environments"
)

var logger loggerModel.LoggerLevels
var once sync.Once

// Startup is placeholder for starting util
func Startup(loggerInstance loggerModel.LoggerLevels) {
	once.Do(func() {
		logger = loggerInstance
	})
}

// Shutdown is placeholder for shutting down util
func Shutdown() {

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
// @param executable string - the binary process name
// @param signal syscall.Signal - the signal type to send
func SendSignal(executable string, signal syscall.Signal) error {
	logger.Debug("Sending %s to %s\n", signal, executable)

	if environments.IsEOS() &&
		strings.Contains(executable, "packetd") {
		return SendSignalViaSysdb(executable, signal)
	}
	// This should normally work on OpenWRT
	pidStr, err := exec.Command("pidof", executable).CombinedOutput()
	if err != nil {
		logger.Debug("Failure to get %s pid: %s\n", executable, err.Error())
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

// Check if the executable is packetd on EOS and if so, try to send it the
// specified signal using Sysdb.
// TODO This only supports packetd signals at this point.
// @param executable - executable to run against
// @param signal syscall.Signal - the signal type to send
// @return any error from running
func SendSignalViaSysdb(executable string,
	signal syscall.Signal) error {

	// This is only relevant for packetd
	arg := ""
	switch signal {
	case syscall.SIGHUP:
		arg = "--sighup"
	case syscall.SIGUSR1:
		arg = "--sigusr1"
	default:
		return fmt.Errorf("unknown signal %v", signal)
	}
	logger.Debug("/usr/bin/updateSysdbSignal %s\n", signal)
	// This is the same script used by sync-settings
	// We are not expecting a failure
	if err := exec.Command("/usr/bin/updateSysdbSignal", arg).Run(); err != nil {
		logger.Err("Failed to run updateSysdbSignal %s with error: %+v\n", arg, err)
		return err
	}
	return nil
}
