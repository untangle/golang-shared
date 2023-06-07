package packetfilter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/untangle/golang-shared/booleval"
	"github.com/untangle/golang-shared/services/settings"
)

func TestConstructor(t *testing.T) {
	settingsFile := settings.NewSettingsFile("./testdata/filter_settings.json")
	filters, err := getFilters(settingsFile)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(filters))
	assert.Equal(t, 2, len((filters)[0].Conditions))
}

func TestConditions(t *testing.T) {
	testCondition := booleval.AtomicExpression{
		Operator:     "==",
		CompareValue: booleval.NewIPComparable("8.8.8.8"),
		ActualValue:  "8.8.8.8",
	}

	settingsFile := settings.NewSettingsFile("./testdata/filter_settings.json")
	filters, err := getFilters(settingsFile)
	assert.Nil(t, err)
	assert.EqualValues(t, filters[0].AtomicExpressions[0], testCondition)
}
