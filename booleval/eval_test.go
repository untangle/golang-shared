package booleval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInts(t *testing.T) {
	i := IntegerComparable{32}

	tests := []struct {
		value  any
		iserr  bool
		result bool
	}{
		{"32", false, true},
		{32, false, true},
		{32.0, false, true},
		{33, false, false},
		{"33", false, false},
		{"32.0", false, true},
		{"32.01", true, false},
		{"dood", true, false},
	}

	for _, test := range tests {
		wasEqual, err := i.Equal(test.value)
		if test.iserr {
			assert.NotNil(
				t,
				err,
				"this value should result in an error: %v\n",
				test.value)
		} else if err != nil {
			assert.Fail(t, "i.Equal(%v) returned an error: %v and should not have\n",
				err)
		}
		assert.Equal(t, test.result, wasEqual,
			"should get %v for wasEqual of %v and %v(%T)", test.result, i.theInteger, test.value, test.value)
	}
}

func TestSimpleConds(t *testing.T) {
	i := IntegerComparable{22}
	i2 := IntegerComparable{21}

	s := StringComparable{"toodle"}
	tests := []struct {
		cond   Condition
		iserr  bool
		result bool
	}{
		{Condition{"==", []Comparable{i}, "22"}, false, true},
		{Condition{"==", []Comparable{i}, "23"}, false, false},
		{Condition{"==", []Comparable{i, i2}, "21"}, false, true},

		// (i > 21) OR (i2 > 21)
		{Condition{">", []Comparable{i, i2}, "21"}, false, true},

		{Condition{"<", []Comparable{i, i2}, 0}, false, false},
		{Condition{">", []Comparable{i, i2}, 0}, false, true},
		{Condition{"==", []Comparable{s}, "toodle"}, false, true},
	}
	for _, test := range tests {
		result, err := EvalCondition(test.cond)
		if err != nil && !test.iserr {
			assert.Fail(t,
				"Received error from evaluating condition %v: %v\n",
				test.cond, err)
			continue
		} else if err != nil {
			continue
		}
		assert.Equal(t, test.result, result)
	}
}

func TestExprs(t *testing.T) {
	i := IntegerComparable{22}

	s := StringComparable{"toodle"}
	tests := []struct {
		cond   []Condition
		iserr  bool
		result bool
	}{
		// (22 > 2) AND ("toodle" == "toodle")
		{[]Condition{
			Condition{">", []Comparable{i}, "2"},
			Condition{"==", []Comparable{s}, "toodle"}},
			false,
			true},
	}
	for _, test := range tests {
		result, err := EvalConditions(test.cond)
		if err != nil && !test.iserr {
			assert.Fail(t,
				"Received error from evaluating condition %v: %v\n",
				test.cond, err)
			continue
		} else if err != nil {
			continue
		}
		assert.Equal(t, test.result, result)
	}
}
