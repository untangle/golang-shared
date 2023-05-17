package booleval

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type OpType string

const (
	eq = "=="
	gt = ">"
)

type valueCondTest struct {
	op     OpType
	value  any
	result bool
	iserr  bool
}

func testDriver(t *testing.T, comparable Comparable, tests []valueCondTest) {
	var result bool
	var err error
	for _, test := range tests {
		switch test.op {
		case eq:
			result, err = comparable.Equal(test.value)
		case gt:
			result, err = comparable.Greater(test.value)
		}

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
		assert.Equal(t, test.result, result,
			"should get %v for wasEqual of %v and %v(%T)", test.result, comparable, test.value, test.value)
	}
}

func TestInts(t *testing.T) {
	i := IntegerComparable{32}

	tests := []valueCondTest{
		{eq, "32", true, false},
		{eq, 32, true, false},
		{eq, 32.0, true, false},
		{eq, 33, false, false},
		{eq, uint8(16), false, false},
		{eq, int8(32), true, false},
		{gt, uint16(16), true, false},
		{eq, float32(32.0), true, false},
		{gt, "dood", false, true},
		{eq, "33", false, false},
		{eq, "32.0", true, false},
		{eq, "32.01", false, true},
		{eq, "dood", false, true},
	}
	testDriver(t, i, tests)

	var myInt uint8 = 3
	i2 := NewIntegerComparableFromIntType(myInt)
	assert.EqualValues(t, i2.theInteger, 3)

	i3, err := NewIntegerComparableFromAny("33.0")
	assert.Nil(t, err)
	assert.EqualValues(t, i3.theInteger, 33)
}

func TestIPs(t *testing.T) {
	ip := IPComparable{ipaddr: net.ParseIP("23.23.1.1")}
	tests := []valueCondTest{
		{eq, "23.23.1.1", true, false},
		{eq, net.IPv4(23, 23, 1, 1), true, false},
		{eq, "1.1.1.1", false, false},
		{eq, "fe80::1", false, false},
		{eq, "23.23.0.0/16", true, false},
		{eq, "23.23.0/16", false, true},
		{eq, "dood", false, false},
	}
	testDriver(t, ip, tests)

	ip = NewIPComparable("123.123.1.1")
	tests = []valueCondTest{
		{eq, "123.123.1.1", true, false},
		{eq, "123.123.1.2", false, false},
		{eq, 6, false, true},
		{eq, "dood", false, false}}
	testDriver(t, ip, tests)

	ipComparable, err := NewIPOrIPNetComparable("192.168.123.1")
	assert.Nil(t, err)
	result, err := ipComparable.Equal(net.IPv4(192, 168, 123, 1))
	assert.Nil(t, err)
	assert.True(t, result)

	ipComparable, err = NewIPOrIPNetComparable("192.168.123.1/24")
	assert.Nil(t, err)
	result, err = ipComparable.Equal(net.IPv4(192, 168, 123, 4))
	assert.Nil(t, err)
	assert.True(t, result)

	result, err = ipComparable.Equal(net.IPv4(192, 168, 123, 6))
	assert.Nil(t, err)
	assert.True(t, result)

	ipComparable, err = NewIPOrIPNetComparable("192.168.123/24")
	assert.NotNil(t, err)
}

func TestIPNets(t *testing.T) {
	_, val, _ := net.ParseCIDR("132.1.23.0/24")
	ipnet := IPNetComparable{ipnet: *val}
	tests := []valueCondTest{
		{eq, "132.1.23.1", true, false},
		{eq, "132.1.24.1", false, false},
		{eq, "132.22.1.1", false, false},
		{eq, 6, false, true},
		{eq, "dood", false, false}}

	testDriver(t, ipnet, tests)

}

func TestTimeComparable(t *testing.T) {
	location, err := time.LoadLocation("")
	assert.Nil(t, err)
	time1 := time.Date(1999, time.April, int(time.Monday), 0, 0, 0, 0, location)
	timeComparable := TimeComparable{time: time1}
	tests := []valueCondTest{
		{eq, time1.Unix(), true, false},
		{eq, time1, true, false},
		{eq, "foop", false, true},
		{gt, "foop", false, true},
		{gt, time1.Add(time.Hour), false, false},
		{gt, time1.Add(-time.Hour), true, false},
		{gt, time1.Unix() - 1, true, false},
		{gt, time1.Unix() + 1, false, false}}
	testDriver(t, timeComparable, tests)
}

func TestStrings(t *testing.T) {
	comp := NewStringComparable("a_one")
	tests := []valueCondTest{
		{eq, "a_one", true, false},
		{gt, "a", true, false},
		{eq, &valueCondTest{}, false, false},
		{eq, 1, false, false}}
	testDriver(t, comp, tests)

	comp = NewStringComparable("192.168.123.1")
	tests = []valueCondTest{
		{eq, net.IPv4(192, 168, 123, 1), true, false},
		{gt, "a", false, false},
		{gt, 1, false, true},
		{eq, 1, false, false}}
	testDriver(t, comp, tests)
}
