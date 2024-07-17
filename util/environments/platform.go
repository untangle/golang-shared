package environments

import "os"

const eosReleaseFile = "/etc/Eos-release"

// Checks if running in EOS
func IsEOS() bool {
	_, err := os.Stat(eosReleaseFile)
	return (err == nil)
}
