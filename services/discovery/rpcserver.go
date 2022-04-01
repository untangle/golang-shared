package discovery

import (
	"github.com/untangle/golang-shared/services/logger"
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

// DiscoveryRPCService is the RPC service for the discovery service
type DiscoveryRPCService int

// ScanNet is a command to scan a network, argument is the networks (CIDR notation)
func (s *DiscoveryRPCService) ScanNet(args *disco.ScanNetRequest, reply *disco.RequestResponse) error {
	logger.Debug("ScanNet called\n")
	callCollectors([]Command{{Command: CmdScanNet, Arguments: args.Net}})
	reply = &disco.RequestResponse{Result: disco.ResponseCode_OK}
	return nil
}

// ScanHost is a command to scan a host, argument is the hostnames
func (s *DiscoveryRPCService) ScanHost(args *disco.ScanHostRequest, reply *disco.RequestResponse) error {
	logger.Debug("ScanHost called\n")
	callCollectors([]Command{{Command: CmdScanHost, Arguments: args.Host}})
	reply = &disco.RequestResponse{Result: disco.ResponseCode_OK}
	return nil
}

// ReqquestAllEntries is a command to request all entries via the ZMQ publisher
func (s *DiscoveryRPCService) RequestAllAntries(args *disco.EmptyParam, reply *disco.RequestResponse) error {
	logger.Debug("RequestAllAntries called\n")
	publishAll()
	reply = &disco.RequestResponse{Result: disco.ResponseCode_OK}
	return nil
}
