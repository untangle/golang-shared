package discovery

import (
	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

// CollectorName - type alias used for identifying collector plugins
type CollectorName string

const (
	All CollectorName = "all"

	Discovery CollectorName = "discovery"

	Arp       CollectorName = "arp" // will be replaced by Neighbour
	Lldp      CollectorName = "lldp"
	Neighbour CollectorName = "neighbour"
	Nmap      CollectorName = "nmap"
)

// CallCollectorsRequest - request for the CallCollectors RPC stub, a stronger-typed wrapper for the RPC func's request
type CallCollectorsRequest struct {
	Collectors []CollectorName
	Args       []string
}

// toRpcRequest - converts the wrapper request struct to the one the RPC function expects
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

// CallCollectorsResponse - CallCollectors RPC stub response, wrapper for the RPC func's response
type CallCollectorsResponse struct {
	Result int32
}

// fromRPCResponse - converts the RPC response into a wrapper
func fromRPCResponse(rpcResponse disco.RPCCallCollectorsResponse) CallCollectorsResponse {
	return CallCollectorsResponse{
		int32(rpcResponse.Result),
	}
}
