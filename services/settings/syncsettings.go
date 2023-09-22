package settings

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// SyncSettings is the struct holding sync-settings information
type SyncSettings struct {
	SettingsFile           string
	DefaultsFile           string
	CurrentFile            string
	OS                     string
	TmpSettingsFile        string
	SyncSettingsExecutable string
	UIDFile                string
}

// NewSyncSettings creates a new settings object
func NewSyncSettings(settingsfile string, defaultsfile string, currentfile string, os string, tmpsettingsfile string, syncsettingsexecutable string, uidfile string) *SyncSettings {
	s := new(SyncSettings)

	s.SettingsFile = settingsfile
	s.DefaultsFile = defaultsfile
	s.CurrentFile = currentfile
	s.OS = os
	s.TmpSettingsFile = tmpsettingsfile
	s.SyncSettingsExecutable = syncsettingsexecutable
	s.UIDFile = uidfile

	return s
}

// CreateDefaults creates the settings defauls.json file
func (s *SyncSettings) CreateDefaults() error {
	// sync the defaults
	cmdArgs := []string{"-o", s.OS, "-c", "-s", "-f", s.TmpSettingsFile}
	err := s.runSyncSettings(cmdArgs)
	if err != nil {
		logger.Warn("Error creating defaults: %s\n", err.Error())
		return err
	}

	// move the defaults. Have to read/write file to avoid docker copy errors
	defaultsBytes, readErr := ioutil.ReadFile(s.TmpSettingsFile)
	if readErr != nil {
		logger.Warn("Failure copying defaults over: %s\n", readErr.Error())
		return readErr
	}

	writeErr := ioutil.WriteFile(s.DefaultsFile, defaultsBytes, 0660)
	if writeErr != nil {
		logger.Warn("Failure copying defaults over: %s\n", writeErr.Error())
		return writeErr
	}

	removeErr := os.Remove(s.TmpSettingsFile)
	if removeErr != nil {
		logger.Warn("Could not remove default tmp file: %s. Continueing\n", removeErr.Error())
	}

	return nil
}

// NormalSync runs sync settings with OS and filename specified
func (s *SyncSettings) NormalSync() error {
	cmdArgs := []string{"-o", s.OS, "-f", s.SettingsFile}
	err := s.runSyncSettings(cmdArgs)
	if err != nil {
		logger.Warn("Error running sync-settings: %s\n", err.Error())
		return err
	}
	return nil
}

// SimulateSync will run sync-settings with simulation flag on the given filePath
// This will not write any files out or restart any services
// but will get the return result as if the file was run properly
func (s *SyncSettings) SimulateSync(filePath string) error {
	cmdArgs := []string{"-o", s.OS, "-s", "-f", filePath}
	err := s.runSyncSettings(cmdArgs)
	if err != nil {
		logger.Warn("Error running sync-settings with simulate flag : %s\n", err.Error())
		return err
	}

	return nil
}

// FirstSyncSettingsRun will create the settings file if it doesn't exist, or rerun sync-settings for good measure
func (s *SyncSettings) FirstSyncSettingsRun() error {
	cmdArgs := []string{"-o", s.OS, "-n"}

	// check if settings.json exists, if not create it
	info, checkErr := os.Stat(s.SettingsFile)
	if os.IsNotExist(checkErr) {
		cmdArgs = append(cmdArgs, "-c")
	} else if info.IsDir() {
		logger.Warn("File is a directory, that's wrong\n")
		return errors.New("Settings file is a directory")
	} else if checkErr != nil {
		logger.Warn("Something went wrong creating settings file: %s\n", checkErr.Error())
		return checkErr
	}

	err := s.runSyncSettings(cmdArgs)
	if err != nil {
		logger.Warn("Error running sync-settings: %s\n", err.Error())
		return err
	}
	return nil
}

// runSyncSettings runs sync settings with given cmd args
func (s *SyncSettings) runSyncSettings(cmdArgs []string) error {
	cmd := exec.Command(s.SyncSettingsExecutable, cmdArgs...)
	outbytes, err := cmd.CombinedOutput()
	output := string(outbytes)

	if err != nil {
		// If just a non-zero exit code, just use standard language
		// Otherwise use the real error message
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() != 0 {
					logger.Warn("Failed to run sync-settings: %v\n", err.Error())
					return errors.New("Failed to save settings")
				}
			}
		}

		logger.Err("Failed to run sync-settings: %v\n", err.Error())
		return err
	}

	if outputErr := s.logSyncSettingsOutput(output, nil); outputErr != nil {
		return outputErr
	}

	return nil
}

// logSyncSettingsOutput logs the output from a sync-settings run
func (s *SyncSettings) logSyncSettingsOutput(output string, err error) error {
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		if logger.IsDebugEnabled() {
			logger.Debug("sync-settings: %v\n", scanner.Text())
		}
	}
	if err != nil {
		logger.Warn("sync-settings return an error: %v\n", err.Error())
		return err
	}
	return nil
}
