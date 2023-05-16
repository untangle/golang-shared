package booleval

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExprs(t *testing.T) {
	intComp := IntegerComparable{22}
	stringComp := StringComparable{"toodle"}
	ip := IPComparable{ipaddr: net.ParseIP("192.168.22.22")}
	_, network, err := net.ParseCIDR("192.168.22.2/24")
	require.Nil(t, err)
	ipNet := IPNetComparable{ipnet: *network}
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
		} else if err != nil {
			continue
		}
		assert.Equal(t,
			test.result,
			result,
			"comparison for test %#v failed.", test)
	}
}
