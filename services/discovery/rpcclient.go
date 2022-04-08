package discovery

import (
	"net/rpc"

	"github.com/untangle/golang-shared/services/logger"
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

// RequestNetworkScan is a stub for the RPC call
func RequestNetworkScan(args disco.ScanNetRequest) {
	logger.Info("Requesting network scan\n")
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:5563")
	if err != nil {
		logger.Err("Failed to connect to discovery service: %s\n", err.Error())
		return
	}
	defer client.Close()

	var reply int
	err = client.Call("DiscoveryRPCService.ScanNet", args, &reply)
	if err != nil {
		logger.Err("Failed to call DiscoveryRPCService.ScanNet %s\n", err.Error())
	}
}

// RequestHostScan is a stub for the RPC call
func RequestHostScan(args disco.ScanHostRequest) {
	logger.Info("Requesting host scan\n")
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:5563")
	if err != nil {
		logger.Err("Failed to connect to discovery service: %s\n", err.Error())
		return
	}
	defer client.Close()

	var reply int
	err = client.Call("DiscoveryRPCService.ScanHost", args, &reply)
	if err != nil {
		logger.Err("Failed to call DiscoveryRPCService.ScanHost %s\n", err.Error())
	}
}

// RequestAllEntries is a stub for the RPC call
func RequestAllEntries() {
	logger.Info("Requesting all entries\n")
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:5563")
	if err != nil {
		logger.Err("Failed to connect to discovery service: %s\n", err.Error())
		return
	}
	defer client.Close()

	var reply int
	err = client.Call("DiscoveryRPCService.RequestAllEntries", 0, &reply)
	if err != nil {
		logger.Err("Failed to call DiscoveryRPCService.RequestAllEntries %s\n", err.Error())
	}
}
