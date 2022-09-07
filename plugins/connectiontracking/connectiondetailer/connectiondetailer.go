package connectiondetailer

type ConnectionDetailer interface {
	GetConnectionDetails() *ConnectionDetails
}
