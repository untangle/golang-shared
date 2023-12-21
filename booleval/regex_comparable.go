package booleval

import (
	"fmt"
	"regexp"
)

// data type to support a regex pattern match
type RegexComparable struct {
	GreaterNotApplicable
	pattern string
	regex   *regexp.Regexp
}

// Create a new regex pattern match comparable
func NewRegexComparable(apattern string) (*RegexComparable, error) {
	if aregex, err := regexp.Compile(apattern); err != nil {
		return nil, err
	} else {
		regexComp := RegexComparable{
			pattern: apattern,
			regex:   aregex,
		}
		return &regexComp, nil
	}
}

// Test a string against a RegexComparable
func (regexComp *RegexComparable) Equal(other any) (bool, error) {
	switch avalue := other.(type) {
	case []byte:
		return regexComp.regex.Match(avalue), nil
	case string:
		return regexComp.regex.Match([]byte(avalue)), nil
	default:
		return false, fmt.Errorf("could not interpret %v for matching against pattern: %s", avalue, regexComp.pattern)
	}
}
