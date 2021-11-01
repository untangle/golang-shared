package licensemanager

import (
	"github.com/untangle/golang-shared/plugins/util"
	"github.com/untangle/golang-shared/services/logger"
)

// Service struct is used to state/hook of each service
type Service struct {
	State ServiceState `json:"state"`
	Hook  ServiceHook  `json:"hook"`
}

// SetServiceState sets the desired state of an service
// @param save bool - if we should store the service state or not, one off save
// @return error - associated errors
func (s *Service) setServiceState(newAllowedState State) error {
	var err error

	runInterrupt := false
	oldAllowedState := s.State.getAllowedState()
	s.State.setAllowedState(newAllowedState)
	logger.Info("old state: %v\n", oldAllowedState)
	logger.Info("new state: %v\n", s.State.getAllowedState())
	switch newAllowedState {
	case StateEnable:
		// if switching, need to run start
		if oldAllowedState == StateDisable {
			runInterrupt = s.ServiceStart()
		}
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
func (s *Service) ServiceStart() bool {
	if s.Hook.Start == nil {
		logger.Info("No start specified, using sighup\n")
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
func (s *Service) ServiceStop() bool {
	if s.Hook.Stop == nil {
		logger.Info("No stop specified, using sighup\n")
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
