package booleval

import (
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExprs(t *testing.T) {
	intComp := IntegerComparable{22}
	stringComp := StringComparable{"toodle"}
	ip := NewIPComparable("192.168.22.22")
	ipNet, err := NewIPOrIPNetComparable("192.168.22.2/24")
	require.Nil(t, err)
	location, err := time.LoadLocation("")
	assert.Nil(t, err)
	time1 := time.Date(1999, time.April, int(time.Monday), 0, 0, 0, 0, location)
	timeComparable1999 := TimeComparable{time: time1}
	time2 := time.Date(2000, time.April, int(time.Monday), 0, 0, 0, 0, location)
	timeComparable2000 := TimeComparable{time: time2}
	tests := []struct {
		expr   Expression
		result bool
		iserr  bool
	}{
		// (192.168.22.22 == 192.168.22.22 AND 192.168.22.22 in 192.168.22.2/24 ) OR (1999 == 2000)
		{NewSimpleExpression(
			OrOfAndsMode,
			[][]*AtomicExpression{
				{
					{"==", ip, "192.168.22.22"},
					{"==", ipNet, "192.168.22.22"},
				},
				{
					{"==", timeComparable2000, time1},
				},
			},
		),
			true,
			false,
		},
		{NewSimpleExpression(
			OrOfAndsMode,
			[][]*AtomicExpression{
				{
					{"==", ip, "192.168.22.22"},
					{"==", ipNet, "192.168.22.22"},
				},

				{
					{">=", timeComparable2000, time1},
				},
			},
		),
			true,
			false,
		},

		{NewSimpleExpression(
			OrOfAndsMode,
			[][]*AtomicExpression{
				{
					{"==", ip, "192.168.55.22"},
					{"==", ipNet, "192.168.55.22"},
				},
				{
					{"==", timeComparable2000, time1},
				},
			},
		),
			false,
			false,
		},
		{NewSimpleExpression(
			OrOfAndsMode,
			[][]*AtomicExpression{
				{
					{"!=", ip, "192.168.55.22"},
					{"!=", ipNet, "192.168.55.22"},
				},
				{
					{"<", timeComparable1999, time2},
				},
			},
		),
			true,
			false,
		},

		{NewSimpleExpression(
			OrOfAndsMode,
			[][]*AtomicExpression{
				{
					{"!=", ip, "192.168.55.22"},
					{"!=", ipNet, "192.168.55.22"},
				},
				{
					{">", timeComparable1999, time2},
				},
			},
		),
			true,
			false,
		},
		{NewSimpleExpression(
			AndOfOrsMode,
			[][]*AtomicExpression{
				{
					{"!=", ip, "192.168.55.23"},
					{"!=", ipNet, "192.168.55.23"},
				},
				{
					{"<", timeComparable1999, time2},
				},
			},
		),
			true,
			false,
		},
		{NewSimpleExpression(
			AndOfOrsMode,
			[][]*AtomicExpression{
				{
					{"<=", IntegerComparable{25}, 24},
					{"<=", IntegerComparable{27}, 24},
					{"<=", IntegerComparable{24}, 24},
				},
				{
					{"<=", timeComparable1999, time2},
				},
			},
		),
			true,
			false,
		},
		{NewSimpleExpression(
			AndOfOrsMode,
			[][]*AtomicExpression{
				{
					{"<=", IntegerComparable{25}, 24},
					{"<=", IntegerComparable{27}, 24},
					{"<=", IntegerComparable{24}, 24},
				},
				{
					{"<=", timeComparable1999, time2},
				},
				{
					{"<=", StringComparable{"dood"}, "doodle"},
				},
				{
					{"==", ip, "192.168.22.23"},
					{"==", ipNet, "192.168.22.23"},
				},
			},
		),
			true,
			false,
		},
		{NewSimpleExpression(
			AndOfOrsMode,
			[][]*AtomicExpression{
				{
					{"<=", IntegerComparable{25}, 24},
					{"<=", IntegerComparable{27}, 24},
				},
				{
					{"<=", timeComparable1999, time2},
				},
				{
					{"<=", StringComparable{"dood"}, "doodle"},
				},
				{
					{"==", ip, "192.168.22.23"},
					{"==", ipNet, "192.168.22.23"},
				},
			},
		),
			false,
			false,
		},
		{NewSimpleExpression(
			AndOfOrsMode,
			[][]*AtomicExpression{
				{
					{"<=", IntegerComparable{25}, 24},
					{"<=", IntegerComparable{27}, 24},
					{"<=", IntegerComparable{24}, 24},
				},
				{
					{"<=", timeComparable1999, time2},
				},
				{
					{"<=", StringComparable{"dood"}, "doodle"},
				},
				{
					{"<=", ip, "192.168.22.23"},
					{"<=", ipNet, "192.168.22.23"},
				},
			},
		),
			false,
			true,
		},
		{NewSimpleExpression(
			AndOfOrsMode,
			[][]*AtomicExpression{
				{
					{">", intComp, "2"},
				},
				{
					{"==", stringComp, "toodle"},
				},
			},
		),
			true,
			false,
		},
		{NewSimpleExpression(
			AndOfOrsMode,
			[][]*AtomicExpression{
				{
					{">=", intComp, 88},
				},
				{
					{"==", stringComp, "toodle"},
				},
			},
		),
			false,
			false,
		},
		{NewSimpleExpression(
			33, // invalid mode, test for error.
			[][]*AtomicExpression{
				{
					{">=", intComp, 88},
				},
				{
					{"==", stringComp, "toodle"},
				},
			},
		),
			false,
			true,
		},
		{NewSimpleExpression(
			AndOfOrsMode,
			[][]*AtomicExpression{
				{
					{"oogabooga", intComp, 88},
				},
				{
					{"==", stringComp, "toodle"},
				},
			},
		),
			false,
			true,
		},
		{NewExpressionWithLookupFunc(
			AndOfOrsMode,
			[][]*AtomicExpression{
				{
					{"<=", IntegerComparable{25}, "myInt"},
					{"<=", IntegerComparable{27}, "myInt"},
				},
				{
					{"<=", timeComparable1999, time2},
				},
				{
					{"<=", StringComparable{"dood"}, "myString"},
				},
				{
					{"==", ip, "myIp"},
					{"==", ipNet, "myIp"},
				},
			},
			func(v any) any {
				stringVal := v.(string)
				vars := map[string]any{
					"myIp":     "192.168.22.23",
					"myString": "doodle",
					"myInt":    24,
				}
				return vars[stringVal]
			},
		),
			false,
			false,
		},
	}

	for _, test := range tests {

		result, err := test.expr.Evaluate()
		if err != nil && !test.iserr {
			assert.Fail(t,
				"failure",
				"Received error from evaluating condition %v: %v\n",
				test.expr, err)
			continue
		} else if test.iserr && err == nil {
			assert.Fail(t,
				"failure",
				"Did not receive expected failure %v\n", test.expr)
			continue
		} else if err != nil {
			// okay, we got expected error.
			continue
		} else {
			assert.Equal(t,
				test.result,
				result,
				"comparison for test %#v failed.", test)
		}
	}
}

func BenchmarkEvaller(b *testing.B) {
	ip := NewIPComparable("192.168.22.22")
	ipNet, _ := NewIPOrIPNetComparable("192.168.22.2/24")
	location := time.Local
	time1 := time.Date(1999, time.April, 1, 0, 0, 0, 0, location)
	timeComparable1999 := TimeComparable{time: time1}
	time2 := time.Date(2000, time.April, int(time.Monday), 0, 0, 0, 0, location)
	timeComparable2000 := TimeComparable{time: time2}
	middleTime := time.Date(1999, time.December, 1, 0, 0, 0, 0, location)
	vars := map[string]any{
		"myIp":     netip.AddrFrom4([4]byte{192, 168, 22, 42}),
		"myString": "doodle",
		"time":     middleTime,
		"myInt":    29,
	}
	expr := NewExpressionWithLookupFunc(
		AndOfOrsMode,
		[][]*AtomicExpression{
			{
				{"==", IntegerComparable{25}, "myInt"},
				{"<=", IntegerComparable{27}, "myInt"},
			},
			{
				{"<=", timeComparable1999, "time"},
			},
			{
				{">=", timeComparable2000, "time"},
			},
			{
				{"<=", StringComparable{"dood"}, "myString"},
			},
			{
				{"==", ip, "myIp"},
				{"==", ipNet, "myIp"},
			},
		},
		func(v any) any {
			stringVal := v.(string)

			return vars[stringVal]
		},
	)
	b.ResetTimer()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		if val, _ := expr.Evaluate(); !val {
			b.FailNow()
		}
	}
	b.StopTimer()
}
