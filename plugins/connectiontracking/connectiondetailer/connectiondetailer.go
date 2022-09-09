package connectiondetailer

// Interface for fetching system connection details
type ConnectionDetailer interface {
	GetConnectionList() ([]*ConnectionInfo, error)
	FetchSystemConnections() error
}
