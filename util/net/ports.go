package net

import (
	"fmt"
	"strconv"
	"strings"
)

// PortSpecifierString is a string in the form of a port range, or a single port.
type PortSpecifierString string

// PortRange is a range of ports, from Start to End inclusive.
type PortRange struct {
	Start int
	End   int
}

// Parse returns the parsed specifier as one of:
// -- int : single port.
// -- PortRange -- PortRange, if the specifier was a range <port>-<port>.
// -- error -- if the port specifier was none of these we return an error object.
func (ps PortSpecifierString) Parse() any {
	var err error
	var startPort, endPort int
	if strings.Contains(string(ps), "-") {
		parts := strings.Split(string(ps), "-")
		if len(parts) != 2 {
			return fmt.Errorf("invalid port specifier string range, contains too many -: %s",
				ps)
		}
		// Convert parts into ints
		if startPort, err = strconv.Atoi(parts[0]); err != nil {
			return fmt.Errorf("invalid port specifier string range, contains bad start port: %s",
				parts[0])
		} else if endPort, err = strconv.Atoi(parts[1]); err != nil {
			return fmt.Errorf("invalid port specifier string range, contains bad end port: %s",
				parts[1])
		}
		if startPort > endPort {
			return fmt.Errorf("invalid port range, start > end: %s", ps)
		} else {
			return PortRange{Start: startPort, End: endPort}
		}
		// Not a range, just a single port
	} else if port, err := strconv.Atoi(string(ps)); err != nil {
		return fmt.Errorf("invalid port specifier: %s", ps)
	} else {
		return port
	}
}

// ContainsPort returns true if the port is between the Start and End of r,
// inclusive.
func (r PortRange) ContainsPort(port int) bool {
	return r.Start <= port && r.End >= port
}
