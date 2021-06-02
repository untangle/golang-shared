package licensemanager

// State is an enum used for indicating service state
type State int

// Types of commands that we accept.
const (
	StateEnable State = iota
	StateDisable
)
