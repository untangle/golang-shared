package connectiondetailer

import "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"

// Interface for fetching system connection details
type ConnectionDetailer interface {
	GetDeviceToConnections() (map[string][]*Discoverd.ConnectionTracking, error)
	FetchSystemConnections() error
}
