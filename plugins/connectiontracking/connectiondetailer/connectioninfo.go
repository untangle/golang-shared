package connectiondetailer

// Structs containing information on each connection
type ConnectionInfo struct {
	Original    *Connection
	Reply       *Connection
	Independent *Independent
}

type Connection struct {
	Layer3 *Layer3
	Layer4 *Layer4
}

type Independent struct {
	Timeout int32
	Mark    int32
	Use     int32
	Id      int64
}

type Layer3 struct {
	Protonum  int32
	Protoname string
	Src       string
	Dst       string
}

type Layer4 struct {
	Protonum  int32
	Protoname string
	SPort     int32
	DPort     int32
}
