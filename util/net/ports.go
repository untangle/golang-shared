package net

import (
	"fmt"
	"strconv"
	"strings"
)

// PortSpecifierString is a string in the form of a port range, or a single port.
type PortSpecifierString string

// Port alias uint16 for readability
type Port uint16

// PortRange is a range of ports, from Start to End inclusive.
type PortRange struct {
	Start Port
	End   Port
}

// Parse returns the parsed specifier as one of:
// -- int : single port.
// -- PortRange -- PortRange, if the specifier was a range <port>-<port>.
// -- error -- if the port specifier was none of these we return an error object.
func (ss PortSpecifierString) Parse() any {
	if strings.Contains(string(ss), "-") {
		parts := strings.Split(string(ss), "-")
		if len(parts) != 2 {
			return fmt.Errorf("invalid port specifier string range, contains too many -: %s",
				ss)
		}
		if start, err := strconv.ParseUint(parts[0], 10, 16); err != nil {
			return fmt.Errorf("invalid port specifier string range, contains bad ports: %s",
				ss)
		} else if end, err := strconv.ParseUint(parts[1], 10, 16); err != nil {
			return fmt.Errorf("invalid port specifier string range, contains bad ports: %s",
				ss)
		} else if start > end {
			return fmt.Errorf("invalid port range, start > end: %s", ss)
		} else {
			return PortRange{Start: Port(start), End: Port(end)}
		}
	} else if port, err := strconv.ParseUint(string(ss), 10, 16); err == nil {
		return Port(port)
	} else {
		return fmt.Errorf("invalid port specifier: %s", ss)
	}
}

// ContainsPort returns true if the port is between the Start and End of r,
// inclusive.
func (r PortRange) Contains(port Port) bool {
	return r.Start <= port && r.End >= port
}
