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

func (state *State) fromString(cmd string) (State, error) {

	switch lowerCmd := strings.ToLower(cmd); lowerCmd {
	case "enable":
		return StateEnable, nil
	case "disable":
		return StateDisable, nil
	default:
		return StateDisable, errors.New("Unable to parse state\n")

	}
}
