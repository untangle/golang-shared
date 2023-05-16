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
		{eq, "dood", false, false},
	}
	testDriver(t, ip, tests)
}

func TestIPNets(t *testing.T) {
	_, val, _ := net.ParseCIDR("132.1.23.0/24")
	ipnet := IPNetComparable{ipnet: *val}
	tests := []valueCondTest{
		{eq, "132.1.23.1", true, false},
		{eq, "132.1.24.1", false, false},
		{eq, "132.22.1.1", false, false},
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
		{gt, time1.Add(time.Hour), false, false},
		{gt, time1.Add(-time.Hour), true, false},
		{gt, time1.Unix() - 1, true, false},
		{gt, time1.Unix() + 1, false, false}}
	testDriver(t, timeComparable, tests)
}
