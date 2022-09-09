package connectiondetailer

import "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"

// Interface for fetching system connection details
type ConnectionDetailer interface {
	GetConnectionList() ([]*Discoverd.ConnectionTracking, error)
	FetchSystemConnections() error
}
