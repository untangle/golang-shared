package booleval

import (
	"fmt"
	"net"
	"net/netip"
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

	weekday, err = NewDayOfWeekFromString("fakeday")
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
	time1 := time.Date(1999, time.April, 1, 0, 0, 0, 0, time.Local)
	timeComparable := TimeComparable{time: time1}
	tests := []valueCondTest{
		{eq, time1.Unix(), true, false},
		{eq, time1, true, false},
		{eq, "foop", false, true},
		{gt, "foop", false, true},
		{eq, "01 Apr 99 00:00 MST", true, false},
		{gt, "01 Apr 98 00:00 MST", true, false},
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

func testBenchmarkOldDriver(b *testing.B, comparable Comparable, tests []valueCondTest) {
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
				b,
				err,
				"this value should result in an error: %v\n",
				test.value)
		} else if err != nil {
			assert.Fail(b, "i.Equal(%v) returned an error: %v and should not have\n",
				err)
		}
		assert.Equal(b, test.result, result,
			"should get %v for wasEqual of %v and %v(%T)", test.result, comparable, test.value, test.value)
	}
}

func BenchmarkOldIP(b *testing.B) {
	//If there is any setup then uncomment this
	//b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.Run("test", func(b *testing.B) {
			//Copied from TestIPs above
			//It would be nice to use the code in place but it's not
			//clear how to make that work with testing.B
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
			testBenchmarkOldDriver(b, ip, tests)

			ip = NewIPComparable("123.123.1.1")
			tests = []valueCondTest{
				{eq, "123.123.1.1", true, false},
				{eq, "123.123.1.2", false, false},
				{eq, 6, false, true},
				{eq, "dood", false, false}}
			testBenchmarkOldDriver(b, ip, tests)

			ipComparable, err := NewIPOrIPNetComparable("192.168.123.1")
			assert.Nil(b, err)
			result, err := ipComparable.Equal(net.IPv4(192, 168, 123, 1))
			assert.Nil(b, err)
			assert.True(b, result)

			ipComparable, err = NewIPOrIPNetComparable("192.168.123.1/24")
			assert.Nil(b, err)
			result, err = ipComparable.Equal(net.IPv4(192, 168, 123, 4))
			assert.Nil(b, err)
			assert.True(b, result)

			result, err = ipComparable.Equal(net.IPv4(192, 168, 123, 6))
			assert.Nil(b, err)
			assert.True(b, result)

			_, err = NewIPOrIPNetComparable("192.168.123/24")
			assert.NotNil(b, err)

			//Copied from TestIPNets above
			//It would be nice to use the code in place but it's not
			//clear how to make that work with testing.B
			_, val, _ := net.ParseCIDR("132.1.23.0/24")
			ipnet := IPNetComparable{ipnet: *val}
			tests = []valueCondTest{
				{eq, "132.1.23.1", true, false},
				{eq, "132.1.24.1", false, false},
				{eq, "132.22.1.1", false, false},
				{eq, 6, false, true},
				{eq, "dood", false, false}}

			testBenchmarkOldDriver(b, ipnet, tests)
		})
	}
}

func testBenchmarkNetIPDriver(b *testing.B, comparable any, tests []valueCondTest) {
	var result bool
	var err error
	var isComparableAddr bool = true
	switch comparable.(type) {
	case netip.Prefix:
		isComparableAddr = false
	}
	testAddr := netip.Addr{}
	testPrefix := netip.Prefix{}
	for _, test := range tests {
		switch test.value.(type) {
		case string:
			testAddr, err = netip.ParseAddr(test.value.(string))
			if err != nil {
				// Assume that it is a CIDR
				testPrefix, err = netip.ParsePrefix(test.value.(string))
				if err != nil {
					test.iserr = true
					result = false
				} else {
					switch test.op {
					case eq:
						if isComparableAddr {
							result = testPrefix.Contains(comparable.(netip.Addr))
						} else {
							result = testPrefix == comparable
						}
					case gt:
						if isComparableAddr {
							result = !testPrefix.Contains(comparable.(netip.Addr))
						} else {
							result = !(testPrefix == comparable)
						}
					}
				}
			} else {
				switch test.op {
				case eq:
					if isComparableAddr {
						result = comparable == testAddr
					} else {
						result = comparable.(netip.Prefix).Contains(testAddr)
					}
				case gt:
					if isComparableAddr {
						result = comparable != testAddr && !comparable.(netip.Addr).Less(testAddr)
					} else {
						result = !comparable.(netip.Prefix).Contains(testAddr)
					}
				}
			}
		default:
			err = fmt.Errorf("Problem translating test value: %v\n", test.value)
			test.iserr = true
		}
		if test.iserr {
			assert.NotNil(
				b,
				err,
				"this value should result in an error: %v\n",
				test.value)
		} else if err != nil {
			assert.Fail(b, "i.Equal(%v) returned an error: %v and should not have\n",
				err)
		}
		if test.result != result {
			fmt.Printf("Debug here\n")
		}
		assert.Equal(b, test.result, result,
			"should get %v for wasEqual of %v and %v(%T)", test.result, comparable, test.value, test.value)
	}
}

func BenchmarkNetIP(b *testing.B) {
	//If there is any setup then uncomment this
	//b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.Run("test", func(b *testing.B) {
			//Copied from TestIPs above
			//It would be nice to use the code in place but it's not
			//clear how to make that work with testing.B
			ip, _ := netip.ParseAddr("23.23.1.1")
			tests := []valueCondTest{
				{eq, "23.23.1.1", true, false},
				{eq, "23.23.1.1", true, false},
				{eq, "1.1.1.1", false, false},
				{eq, "fe80::1", false, false},
				{eq, "23.23.0.0/16", true, false},
				{eq, "23.23.0/16", false, true},
				{eq, "dood", false, false},
			}
			testBenchmarkNetIPDriver(b, ip, tests)

			ip, _ = netip.ParseAddr("123.123.1.1")
			tests = []valueCondTest{
				{eq, "123.123.1.1", true, false},
				{eq, "123.123.1.2", false, false},
				{eq, 6, false, true},
				{eq, "dood", false, false}}
			testBenchmarkNetIPDriver(b, ip, tests)

			ipComparable, err := NewIPOrIPNetComparable("192.168.123.1")
			assert.Nil(b, err)
			result, err := ipComparable.Equal(net.IPv4(192, 168, 123, 1))
			assert.Nil(b, err)
			assert.True(b, result)

			ipComparable, err = NewIPOrIPNetComparable("192.168.123.1/24")
			assert.Nil(b, err)
			result, err = ipComparable.Equal(net.IPv4(192, 168, 123, 4))
			assert.Nil(b, err)
			assert.True(b, result)

			result, err = ipComparable.Equal(net.IPv4(192, 168, 123, 6))
			assert.Nil(b, err)
			assert.True(b, result)

			_, err = NewIPOrIPNetComparable("192.168.123/24")
			assert.NotNil(b, err)

			//Copied from TestIPNets above
			//It would be nice to use the code in place but it's not
			//clear how to make that work with testing.B
			val, _ := netip.ParsePrefix("132.1.23.0/24")
			tests = []valueCondTest{
				{eq, "132.1.23.1", true, false},
				{eq, "132.1.24.1", false, false},
				{eq, "132.22.1.1", false, false},
				{eq, 6, false, true},
				{eq, "dood", false, false}}

			testBenchmarkNetIPDriver(b, val, tests)
		})
	}
}

func BenchmarkBothIP(b *testing.B) {
	fmt.Printf("BenchnarkOldIP semantics:\n")
	BenchmarkOldIP(b)
	fmt.Printf("BenchnarkNetIP semantics:\n")
	BenchmarkNetIP(b)
}
