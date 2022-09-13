package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettings(t *testing.T) {
	type settingsObject struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}
	s := NewSettingsFile("./testdata/settings.json")
	value := settingsObject{}
	err := s.UnmarshalSettingsAtPath(&value, "a", "b")
	assert.Nil(t, err)
	assert.Equal(
		t,
		value,
		settingsObject{
			Foo: "hello",
			Bar: 1})
}
