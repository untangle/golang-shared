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
	// op -- one of eq/gt, which operation to try, eq -> Equal(), gt -> Greater()
	op OpType

	// value -- the 'actual' value.
	value any

	// result -- result of comparison.
	result bool

	// iserr -- should we expect an error from the method call to Equal/Greater?
	iserr bool
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

/*
Below are table-driven tests for all the Comparable objects in this package.
We generally:

1. Construct some comparable

2. Give a table of potential comparisons (Equal/Greater) and
'actual' values that we may want to compare against the Comparable.

3. Pass the comparable from (1) and the table to testDriver which conducts the tests.
*/
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

func TestTimeOfDay(t *testing.T) {
	tod := TimeOfDayComparable{
		timeSinceDayStart: 12*time.Hour + 30*time.Minute,
	}

	tests := []valueCondTest{
		{eq, "12:30PM", true, false},
		{eq, "12:30:00pm", true, false},
		{eq, "12:30", true, false},
		{eq, "bloohblahblah", false, true},
		{eq, 12*time.Hour + 30*time.Minute, true, false},
		{gt, "13:00", false, false},
		{gt, 12*time.Hour + 29*time.Minute, true, false},
		{gt, 12*time.Hour + 31*time.Minute, false, false},
	}
	testDriver(t, tod, tests)

	tod, err := NewTimeOfDayFromTimeString("9:00AM")
	assert.Nil(t, err)
	tests = []valueCondTest{
		{eq, time.Date(1999, time.April, 5, 9, 0, 1, 0, time.Local), true, false},
		{eq, time.Date(1999, time.April, 5, 9, 1, 1, 0, time.Local), false, false},
		{gt, time.Date(1999, time.April, 5, 8, 1, 1, 0, time.Local), false, false},
	}
	testDriver(t, tod, tests)
}

func TestDayOfWeek(t *testing.T) {
	weekday := DayOfWeekComparable{
		dayOfWeek: time.Monday,
	}

	tests := []valueCondTest{
		{eq, "monday", true, false},
		{eq, time.Monday, true, false},
		{eq, time.Date(1999, time.April, 5, 1, 0, 0, 0, time.Local), true, false},
		{eq, time.Date(1999, time.April, 5, 1, 0, 0, 0, time.Local).Unix(), true, false},
		{eq, "bloohblahblah", false, true},
		{gt, "sunday", true, false},
		{gt, "tuesday", false, false},
	}
	testDriver(t, weekday, tests)

	weekday, err := NewDayOfWeekFromString("thursday")
	assert.Nil(t, err)
	tests = []valueCondTest{
		{eq, "monday", false, false},
		{eq, time.Monday, false, false},
		{eq, "thursday", true, false},
	}

	testDriver(t, weekday, tests)

	_, err = NewDayOfWeekFromString("fakeday")
	assert.NotNil(t, err)

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
		{gt, "dood.dood", false, true},
		{gt, "dood:dood", false, true},
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

	_, err = NewIPOrIPNetComparable("192.168.123/24")
	assert.NotNil(t, err)

	ip = NewIPComparable("ABCD:EF01:2345:6789:ABCD:EF01:2345:6789")

	tests = []valueCondTest{
		{eq, "123.123.1.1", false, false},
		{eq, "123.123.1.2", false, false},
		{eq, 6, false, true},
		{eq, "dood", false, false},
		{eq, "ABCD:EF01:2345:6789:ABCD:EF01:2345:6789", true, false},
		{eq, "2001:0db8::2", false, false},
		{eq, "2001:fe99::0", false, false},
		{eq, "2001:DB8:0:0:8:800:200C:417A", false, false},
		{eq, "FF01:0:0:0:0:0:0:101", false, false},
		{eq, "::13.1.68.3", false, false},
		{eq, "::FFFF:129.144.52.38", false, false},
		{eq, "::FFFF:216.151.130.153", false, false},
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
		{eq, 6, false, true},
		{eq, "dood", false, false},
		{eq, "ABCD:EF01:2345:6789:ABCD:EF01:2345:6789", false, false},
		{eq, "2001:0db8::2", false, false},
		{eq, "2001:fe99::0", false, false},
		{eq, "2001:DB8:0:0:8:800:200C:417A", false, false},
		{eq, "FF01:0:0:0:0:0:0:101", false, false},
		{eq, "::13.1.68.3", false, false},
		{eq, "::FFFF:129.144.52.38", false, false},
		{eq, "::FFFF:216.151.130.153", false, false},
	}
	testDriver(t, ipnet, tests)

	_, val, _ = net.ParseCIDR("ABCD:EF01:2345:6789:ABCD:EF01:2345:6789/24")
	ipnet = IPNetComparable{ipnet: *val}

	tests = []valueCondTest{
		{eq, "132.1.23.1", false, false},
		{eq, "132.1.24.1", false, false},
		{eq, "132.22.1.1", false, false},
		{eq, 6, false, true},
		{eq, "dood", false, false},
		{eq, "ABCD:EF01:2345:6789:ABCD:EF01:2345:6789", true, false},
		{eq, "2001:0db8::2", false, false},
		{eq, "2001:fe99::0", false, false},
		{eq, "2001:DB8:0:0:8:800:200C:417A", false, false},
		{eq, "FF01:0:0:0:0:0:0:101", false, false},
		{eq, "::13.1.68.3", false, false},
		{eq, "::FFFF:129.144.52.38", false, false},
		{eq, "::FFFF:216.151.130.153", false, false},
	}
	testDriver(t, ipnet, tests)
}

func TestTimeComparable(t *testing.T) {
	time1 := time.Date(1999, time.April, 1, 0, 0, 0, 0, time.UTC)
	timeComparable := TimeComparable{time: time1}
	tests := []valueCondTest{
		{eq, time1.Unix(), true, false},
		{eq, time1, true, false},
		{eq, "foop", false, true},
		{gt, "foop", false, true},
		{eq, "01 Apr 99 00:00 UTC", true, false},
		{gt, "01 Apr 98 00:00 UTC", true, false},
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

func TestStringArrays(t *testing.T) {
	comp := NewStringArrayComparable([]string{"a_one", "a_two"})
	tests := []valueCondTest{
		{eq, "a_one", true, false},
		{eq, "a_two", true, false},
		{eq, "a", false, false},
		{eq, &valueCondTest{}, false, false},
		{eq, 1, false, false}}
	testDriver(t, comp, tests)
}
