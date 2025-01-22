package eos_util

import (
	"testing"
)

func TestRandomizeSlice(t *testing.T) {
	tests := []struct {
		eosName    string
		kernelName string
	}{
		{
			eosName:    "Ethernet1/4",
			kernelName: "et1_4",
		},
		{
			eosName:    "Management1/1",
			kernelName: "ma1_1",
		},
	}

	// setting custom seed wua Seed has been deprecated

	for _, test := range tests {
		t.Run(test.eosName, func(t *testing.T) {
			if kernelName := TranslateEosInterface(test.eosName); kernelName != test.kernelName {
				t.Errorf("Expected %s -- got %s", test.kernelName, kernelName)
			}
		})
	}
}
