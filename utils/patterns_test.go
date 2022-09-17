package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMacPattern(t *testing.T) {
	assert.True(t, IsMacAddress("00:11:22:33:44:55"))
	assert.True(t, IsMacAddress("00:aa:22:33:44:55"))
	assert.True(t, IsMacAddress("00:FF:22:33:44:55"))
	assert.True(t, IsMacAddress("00:FF:22:33:44:00"))
	assert.False(t, IsMacAddress("00:FF:z2:33:44:00"))
	assert.False(t, IsMacAddress("00:FF:22:33:"))
	assert.False(t, IsMacAddress("00:FF:22:33:999"))
}

func TestIpv4Pattern(t *testing.T) {
	assert.True(t, IsIpv4Address("192.168.0.1"))
	assert.True(t, IsIpv4Address("1.1.0.1"))
	assert.True(t, IsIpv4Address("0.0.0.0"))
	assert.True(t, IsIpv4Address("192.168.0.1"))
	assert.False(t, IsIpv4Address("00:FF:z2:33:44:00"))
	assert.False(t, IsIpv4Address("192.168.1"))
	assert.False(t, IsIpv4Address("192.16a.0.1"))
}
