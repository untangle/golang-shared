package discovery

import (
	"context"

	"github.com/untangle/golang-shared/services/logger"
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

type discoveryServer struct {
	disco.UnimplementedDisoverydServer
}

func (s *discoveryServer) ScanNet(ctx context.Context, in *disco.ScanNetRequest) (*disco.RequestResponse, error) {
	logger.Debug("ScanNet called\n")
	callCollectors([]Command{{Command: CmdScanNet, Arguments: in.Net}})
	return &disco.RequestResponse{Result: disco.ResponseCode_OK}, nil
}

func (s *discoveryServer) ScanHost(ctx context.Context, in *disco.ScanHostRequest) (*disco.RequestResponse, error) {
	logger.Debug("ScanHost called\n")
	callCollectors([]Command{{Command: CmdScanHost, Arguments: in.Host}})
	return &disco.RequestResponse{Result: disco.ResponseCode_OK}, nil
}

func (s *discoveryServer) RequestAllEntries(ctx context.Context, in *disco.EmptyParam) (*disco.RequestResponse, error) {
	logger.Debug("RequestAllEntries called\n")
	publishAll()
	return &disco.RequestResponse{Result: disco.ResponseCode_OK}, nil
}
