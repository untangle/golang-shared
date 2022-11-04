package discovery

import (
	"net/rpc"

	"github.com/untangle/golang-shared/services/logger"
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

const (
	network = "tcp"
	address = "127.0.0.1:5563"
)

// RequestCallCollectors is a stub for the RPC call
func RequestCallCollectors(args disco.ScanRequest) {
	logger.Info("RequestCallCollectors called\n")
	if len(args.Collectors) == 0 {
		logger.Warn("RequestHostScan called but no collector specified!")
	}

	client, err := rpc.DialHTTP(network, address)
	if err != nil {
		logger.Err("Failed to connect to discovery service: %s\n", err.Error())
		return
	}
	defer client.Close()

	var reply disco.ScanResponse
	err = client.Call("DiscoveryRPCService.CallCollectors", args, &reply)
	if err != nil {
		logger.Err("Failed to call DiscoveryRPCService.CallCollectors %s\n", err.Error())
	}
}
