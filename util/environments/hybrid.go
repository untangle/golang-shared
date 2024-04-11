package environments

import (
	"fmt"
	"os"
)

const (
	HybridConfigPath = "/mnt/flash/mfw-settings/hybrid"
)

// Checks if running in the hybrid environment(packetd running in EOS with other daemons in an openWRT BST)
func IsHybrid() (bool, error) {
	_, err := os.Stat(HybridConfigPath)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}

	return false, fmt.Errorf("could not determine if file denoting a hybrid exists: %w", err)
}
