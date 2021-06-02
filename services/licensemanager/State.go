package licensemanager

import (
	"errors"
	"strings"
)

// State is an enum used for indicating service state
type State int

// Types of commands that we accept.
const (
	StateEnable State = iota
	StateDisable
)

func (state *State) fromString(cmd string) error {

	switch lowerCmd := strings.ToLower(cmd); lowerCmd {
	case "enable":
		*state = StateEnable
	case "disable":
		*state = StateDisable
	default:
		return errors.New("Unable to parse state\n")
	}

	return nil
}
