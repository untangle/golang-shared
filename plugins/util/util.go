package util

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"

	loggerModel "github.com/untangle/golang-shared/logger"
	//protobuf "github.com/untangle/golang-shared/structs/protocolbuffers/ActiveSessions"
	//grpc "google.golang.org/grpc"
	//"google.golang.org/grpc/credentials/insecure"
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

	// This should normally work on OpenWRT
	pidStr, err := exec.Command("pidof", executable).CombinedOutput()
	if err != nil {
		logger.Debug("Failure to get %s pid: %s\n", executable, err.Error())
		return SendSignalViaGRPC(executable, signal, err)
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

// Copied from /src/EfwSfeModules/ModuleConstants.h
const PACKETD_CONFIG_UPDATE = "configUpdate"

func SendSignalViaGRPC(executable string, signal syscall.Signal, err error) error {
	/*
		// There won't be a packetd PID on EOS
		if strings.Contains(executable, "packetd") {
			message := *protobuf.PacketdConfigArg{signal: int(signal)}
			// If we didn't find packetd then assume we will find local Sfe at the default bess port
			conn, newErr := grpc.Dial("localhost:10514", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if newErr != nil {
				logger.Warn("Could not Dial: %v\n", newErr)
				return newErr
			} else {
				defer conn.Close()
				conn.Connect()

				if newErr = conn.Invoke(context.Background(), PACKETD_CONFIG_UPDATE, message, grpc.WithDefaultCallOptions()); newErr != nil {
					logger.Warn("Could not call %s: %v\n", PACKETD_CONFIG_UPDATE, newErr)
					return newErr
				}
				logger.Debug("Sent signal %d to %s\n", int(signal), PACKETD_CONFIG_UPDATE)
				return nil
			}
		}
	*/
	return err
}
