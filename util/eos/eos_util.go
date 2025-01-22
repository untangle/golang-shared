package eos_util

import "strings"

// TranslateEosInterface returns a kernel interface name from an EOS interface eosInterfaceName
// e.g. Management1/1 to ma1_1
func TranslateEosInterface(eosInterfaceName string) string {
	kernelName := strings.ReplaceAll(eosInterfaceName, "Ethernet", "et")
	kernelName = strings.ReplaceAll(kernelName, "Management", "ma")
	kernelName = strings.ReplaceAll(kernelName, "/", "_")
	return kernelName
}
