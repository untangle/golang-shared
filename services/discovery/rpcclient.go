package discovery

import (
	"net/rpc"

	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

const (
	network = "tcp"
	address = "127.0.0.1:5563"
)

// CallCollectors is a stub for the RPC call
func CallCollectors(args CallCollectorsRequest) (*CallCollectorsResponse, error) {
	logger.Debug("CallCollectors called\n")
	if len(args.Collectors) == 0 {
		logger.Warn("CallCollectors called but no collector specified!")
	}

	client, err := rpc.DialHTTP(network, address)
	if err != nil {
		logger.Err("Failed to connect to discovery service: %s\n", err.Error())
		return nil, err
	}
	defer client.Close()

	rpcRequest := toRpcRequest(args)
	var rpcResponse disco.CallDiscoveryResponse

	if err := client.Call("DiscoveryRPCService.CallDiscovery", &rpcRequest, &rpcResponse); err != nil {
		logger.Err("Failed to call DiscoveryRPCService.CallDiscovery %s\n", err.Error())
		return nil, err
	}

	response := fromRPCResponse(rpcResponse)
	return &response, nil
}
