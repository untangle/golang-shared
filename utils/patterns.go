package utils

import "regexp"

const (
	// IPv4Regex  matches IPV4 addresses.
	IPv4Regex = `(\d+\.\d+\.\d+\.\d+)`

	// MacRegex is a regex matching MAC addresses.
	MacRegex = `((?:[0-9A-Fa-f][0-9A-Fa-f]:){5,5}[0-9A-Fa-f][0-9A-Fa-f])`

	// HexRegex is a generic 0x hex number regex.
	HexRegex = `(0x[a-fA-F\d]+)`

	// MaskRegex is a genric ARP mask regex.
	MaskRegex = `(\*)`

	// DeviceRegex is a network interface device regex (like eth0 and so on).
	DeviceRegex = `([a-zA-Z]+[a-zA-Z0-9]+)`
)

// IsMacAddress returns true if the addr string is a MAC.
func IsMacAddress(addr string) bool {
	fullRegex := `^` + MacRegex + `$`
	didMatch, err := regexp.MatchString(fullRegex, addr)
	return didMatch && err == nil
}
