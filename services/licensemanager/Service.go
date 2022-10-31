package licensemanager

import (
	"github.com/untangle/golang-shared/plugins/util"
	"github.com/untangle/golang-shared/services/logger"
)

// Service struct is used to store state/hook of each service
type Service struct {
	Name  string       `json:"name"`
	State ServiceState `json:"state"`
	Hook  ServiceHook  `json:"hook"`
}

// setServiceState sets the desired state of an service
// @param State newAllowedState - new allowed state
// @return error - associated errors
func (s *Service) setServiceState(newAllowedState State) error {
	var err error

	runInterrupt := false
	oldAllowedState := s.State.getAllowedState()
	s.State.setAllowedState(newAllowedState)

	logger.Debug("old state of %s: %v\n", s.Name, oldAllowedState)
	logger.Debug("new state of %s: %v\n", s.Name, s.State.getAllowedState())

	switch newAllowedState {
	case StateEnable:
		// always run start
		runInterrupt = s.ServiceStart()
	case StateDisable:
		// always run stop
		runInterrupt = s.ServiceStop()
	}

	// api called the sighup, so enable/disable service
	if runInterrupt {
		err := util.RunSighup(config.Executable)
		if err != nil {
			return err
		}
	}

	return err
}

// ServiceStart starts the service, either via sighup or normal start
// @return bool - if sighup should be run or not
func (s *Service) ServiceStart() bool {
	// if no Start() in hook, run sighup and write out file
	if s.Hook.Start == nil {
		logger.Debug("No start specified, using sighup\n")
		err := s.State.writeOutServiceToEnableOrDisable()
		if err != nil {
			logger.Warn("Failure to write out service start: %s\n", err.Error())
			return false
		}
		return true
	}
	s.Hook.Start()
	return false
}

// ServiceStop stops the service, either via sighup or normal stop
// @return bool on if sighup should be run or not
func (s *Service) ServiceStop() bool {
	// if no Stop() in hook, run sighup and write out file
	if s.Hook.Stop == nil {
		logger.Debug("No stop specified, using sighup\n")
		err := s.State.writeOutServiceToEnableOrDisable()
		if err != nil {
			logger.Warn("Failure to write out service stop: %s\n", err.Error())
			return false
		}
		return true
	}
	s.Hook.Stop()
	return false
}
