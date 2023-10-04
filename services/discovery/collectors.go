package discovery

import (
	"fmt"
	"strings"

	disco "github.com/untangle/golang-shared/structs/protocolbuffers/Discoverd"
)

type CollectorName string // CollectorName - type alias used for identifying collector plugins

const (
	All CollectorName = "all"

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
func toRpcRequest(req CallCollectorsRequest) disco.CallDiscoveryRequest {
	strCollectors := make([]string, len(req.Collectors))
	for i, c := range req.Collectors {
		strCollectors[i] = string(c)
	}

	return disco.CallDiscoveryRequest{
		Collectors: strCollectors,
		Args:       req.Args,
	}
}

// CallCollectorsResponse - CallCollectors RPC stub response, wrapper for the RPC func's response
type CallCollectorsResponse struct {
	Result int32
}

// fromRPCResponse - converts the RPC response into a wrapper
func fromRPCResponse(rpcResponse *disco.CallDiscoveryResponse) CallCollectorsResponse {
	return CallCollectorsResponse{
		int32(rpcResponse.Result),
	}
}

// Normalizes the data in each collector entry
// Returns an error if the data couldn't be normalized or
// if the provided argument isn't a pointer to a collector struct
func NormalizeCollectorEntry(collector interface{}) error {
	switch collectorWithType := collector.(type) {
	case *disco.LLDP:
		collectorWithType.Ip = strings.ToUpper(collectorWithType.Ip)
		collectorWithType.Mac = strings.ToUpper(collectorWithType.Mac)
	case *disco.NEIGH:
		collectorWithType.Ip = strings.ToUpper(collectorWithType.Ip)
		collectorWithType.Mac = strings.ToUpper(collectorWithType.Mac)
	case *disco.NMAP:
		collectorWithType.Ip = strings.ToUpper(collectorWithType.Ip)
		collectorWithType.Mac = strings.ToUpper(collectorWithType.Mac)
	default:
		return fmt.Errorf("provided argument was not a pointer to a collector struct")
	}

	return nil
}

// Wraps an LLDP, NMAP, or NEIGH collector struct in a device entry
// If the collector struct is missing a field needed to initialize the
// collector struct, and error is returned.
func WrapCollectorInDeviceEntry(collector interface{}) (*DeviceEntry, error) {
	var deviceEntry DeviceEntry

	switch collectorWithType := collector.(type) {
	case *disco.LLDP:
		if collectorWithType.Ip == "" {
			return nil, fmt.Errorf("LLDP entry missing IP field")
		}

		deviceEntry.MacAddress = collectorWithType.Mac
		deviceEntry.LastUpdate = collectorWithType.LastUpdate

		deviceEntry.Lldp = make(map[string]*disco.LLDP)
		deviceEntry.Lldp[collectorWithType.Ip] = collectorWithType

	case *disco.NEIGH:
		if collectorWithType.Ip == "" {
			return nil, fmt.Errorf("NEIGH entry missing IP field")
		}

		deviceEntry.MacAddress = collectorWithType.Mac
		deviceEntry.LastUpdate = collectorWithType.LastUpdate

		deviceEntry.Neigh = make(map[string]*disco.NEIGH)
		deviceEntry.Neigh[collectorWithType.Ip] = collectorWithType

	case *disco.NMAP:
		if collectorWithType.Ip == "" {
			return nil, fmt.Errorf(("NMAP entry missing IP field"))
		}

		deviceEntry.MacAddress = collectorWithType.Mac
		deviceEntry.LastUpdate = collectorWithType.LastUpdate

		deviceEntry.Nmap = make(map[string]*disco.NMAP)
		deviceEntry.Nmap[collectorWithType.Ip] = collectorWithType

	default:
		return nil, fmt.Errorf("argument provided to the function was not an LLDP, NMAP, or NEIGH struct")
	}

	return &deviceEntry, nil
}
