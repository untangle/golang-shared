package licensemanager

import "github.com/untangle/golang-shared/services/logger"

// ServiceHook struct is used to indicate start/stop/enabled hooks for services
type ServiceHook struct {
	Start   func()
	Stop    func()
	Enabled func() bool
}

func (s *ServiceHook) ServiceStart() bool {
	if s.Start == nil {
		logger.Info("No start specified, using sighup\n")
		return true
	}
	s.Start()
	return false
}

func (s *ServiceHook) ServiceStop() bool {
	if s.Stop == nil {
		logger.Info("No stop specified, using sighup\n")
		return true
	}
	s.Stop()
	return false
}
