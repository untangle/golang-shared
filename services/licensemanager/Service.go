package licensemanager

import (
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
// @param string executable - executable to use when enabling/disabling the service
// @return error - associated errors
func (s *Service) setServiceState(newAllowedState State) error {
	var err error

	oldAllowedState := s.State.getAllowedState()
	s.State.setAllowedState(newAllowedState)
	logger.Info("********* old state of %s: %v\n", s.Name, oldAllowedState)
	logger.Info("********* new state of %s: %v\n", s.Name, s.State.getAllowedState())

	return err
}

// ServiceStart starts the service.
func (s *Service) ServiceStart() bool {
	// if no Start() in hook, run sighup and write out file
	if s.Hook.Start == nil {
		return true
	}
	s.Hook.Start()
	return false
}

// ServiceStop stops the service.
func (s *Service) ServiceStop() bool {
	// if no Stop() in hook, run sighup and write out file
	if s.Hook.Stop == nil {
		return true
	}
	s.Hook.Stop()
	return false
}
