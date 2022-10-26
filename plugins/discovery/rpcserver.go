package discovery

import (
	"github.com/untangle/golang-shared/services/logger"
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

// DiscoveryRPCService is the RPC service for the discovery service
type DiscoveryRPCService disco.NmapRequest

// ScanNet is a command to scan a network, argument is the networks (CIDR notation)
func (s *DiscoveryRPCService) ScanNet(args *disco.NmapRequest, reply *disco.NmapResponse) error {
	logger.Debug("ScanNet called\n")
	NewDiscovery().callCollectors([]Command{{Command: CmdScanNet, Arguments: args.Net}})
	reply = &disco.NmapResponse{Result: disco.ResponseCode_OK}
	return nil
}

// ScanHost is a command to scan a host, argument is the hostnames
func (s *DiscoveryRPCService) ScanHost(args *disco.NmapRequest, reply *disco.NmapResponse) error {
	logger.Debug("ScanHost called\n")
	NewDiscovery().callCollectors([]Command{{Command: CmdScanHost, Arguments: args.Host}})
	reply = &disco.NmapResponse{Result: disco.ResponseCode_OK}
	return nil
}
