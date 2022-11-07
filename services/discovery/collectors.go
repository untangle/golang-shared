package discovery

import (
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

type CollectorName string

const (
	All CollectorName = "all"

	Discovery CollectorName = "discovery"

	Arp CollectorName = "arp" // will be replaced by Neighbour

	Lldp      CollectorName = "lldp"
	Neighbour CollectorName = "neighbour"
	Nmap      CollectorName = "nmap"
)

type CallCollectorsRequest struct {
	Collectors []CollectorName
	Args       []string
}

func toRpcRequest(req CallCollectorsRequest) disco.RPCCallCollectorsRequest {
	strCollectors := make([]string, len(req.Collectors), len(req.Collectors))
	for i, c := range req.Collectors {
		strCollectors[i] = string(c)
	}

	return disco.RPCCallCollectorsRequest{
		Collectors: strCollectors,
		Args:       req.Args,
	}
}

type CallCollectorsResponse struct {
	Result int32
}

func fromRPCResponse(rpcResponse disco.RPCCallCollectorsResponse) CallCollectorsResponse {
	return CallCollectorsResponse{
		int32(rpcResponse.Result),
	}
}
