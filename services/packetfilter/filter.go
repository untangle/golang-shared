package packetfilter

import (
	"github.com/untangle/golang-shared/booleval"
	"github.com/untangle/golang-shared/services/settings"
)

// Test that we can do a simple load of the settings file.
func getFilters(settingsFile *settings.SettingsFile) ([]PacketFilter, error) {
	filterSettings := []PacketFilter{}

	if err := settingsFile.UnmarshalSettingsAtPath(&filterSettings, "packetfilters"); err != nil {
		return nil, err
	}

	// Populate the expressions.
	for idx, f := range filterSettings {
		for _, c := range f.Conditions {
			expression := booleval.AtomicExpression{}
			expression.Operator = c.Operator
			if compareValue, err := getCompareValue(c); err != nil {
				return nil, err
			} else {
				expression.CompareValue = compareValue
			}
			expression.ActualValue = c.Value
			filterSettings[idx].AtomicExpressions = append(filterSettings[idx].AtomicExpressions, expression)
		}
	}
	return filterSettings, nil
}

// TODO: Needs to cover all cases.
func getCompareValue(condition PacketCondition) (booleval.Comparable, error) {
	switch condition.Property {
	case "ServerAddress", "ServerAddressV6", "ClientAddress", "ClientAddressV6":
		return booleval.NewIPComparable(condition.Value), nil
	case "ServerPort", "ClientPort":
		if port, err := booleval.NewIntegerComparableFromAny(condition.Value); err != nil {
			return nil, err
		} else {
			return port, nil
		}
	case "ApplicationIDInferred", "ApplicationNameInferred", "ApplicationProtocolInferred", "ApplicationRiskInferred":
		return booleval.NewStringComparable(condition.Value), nil
	case "CertSubjectCN", "CertSubjectDNS", "CertSubjectO", "ClientReverseDNS", "ServerDNSHint":
		return booleval.NewStringComparable(condition.Value), nil
	}
	return nil, nil
}
