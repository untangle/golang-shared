package booleval

import (
	"fmt"
)

// Comparable -- a Comparable is a simple interface to allow the
// evaluator to make comparisons between objects. Comparables usually
// are used on the left hand side of an expression.
type Comparable interface {
	Equal(other any) (bool, error)
	Greater(other any) (bool, error)
}

// GreaterNotApplicable is a struct for embedding in objects where
// they should not be ordered, like IP addresses. It will return an
// error for any application of Greater().
type GreaterNotApplicable struct{}

// Greater returns an error if it is used, for embedding in other
// Comparables that are not ordered.
func (GreaterNotApplicable) Greater(other any) (bool, error) {
	return false, fmt.Errorf("this type does not support ordering")
}

// AtomicExpression is a single comparsion -- an operator, the value
// to compare against (the lefthand side), and the actual value.
type AtomicExpression struct {
	Operator     string
	CompareValue Comparable
	ActualValue  any
}

const (
	// AndOfOrsMode is an evaluator mode -- and we AND each
	// condition, and the possible values are ORed
	AndOfOrsMode = iota

	// OrOfAndsMode is an evaluator mode -- we OR each condition
	// and the possible values are ANDed.
	OrOfAndsMode
)

// EvaluatorMode is either AndOfOrsMode or OrOfAndsMode, and is the
// logical connectives/format of the expression to be evaluated.
type EvaluatorMode int

// Expression -- this represents a simple boolean expression which is
// either and and of ors:
//
// (e.g. ((x OR y)AND (a OR b OR c OR d) AND z ...))
//
// or and or of ands:
// ((x AND y) OR (b AND c AND d) ...)
//
// and can be used as a 'backend' for evaluating generic boolean
// expressions.
//
// It contains a LookupFunc, which is used to 'look up' each
// ActualValue in an AtomicExpression during evaluation, to allow for
// variable-like strings to be used for ActualValue.
type Expression struct {
	// ExpressionConnective is the 'mode' -- and of ors, or or of ands.
	ExpressionConnective EvaluatorMode

	// AtomicExpressions is the nested list of expressions,
	// like:
	// [[x OR y] AND [a OR b OR c]]
	// [x OR y] is the first list in the outer list.
	// The inner lists  are called 'clauses'
	AtomicExpressions [][]*AtomicExpression

	// LookupFunc is used to look up a replacement for any
	// ActualValue in an AtomicExpression during evaluation.
	LookupFunc func(any) any
}

// NewSimpleExpression returns an expression that uses the identity
// function for lookups.
func NewSimpleExpression(
	connective EvaluatorMode,
	exprs [][]*AtomicExpression) Expression {
	return Expression{
		ExpressionConnective: connective,
		AtomicExpressions:    exprs,
		LookupFunc:           func(v any) any { return v },
	}
}

// ExpressionCopyWithLookupFunc creates a copy of this expression (not
// a deep copy of the actual expression), using the given lookup
// function.
func ExpressionCopyWithLookupFunc(
	expr Expression,
	lookupFunc func(any) any) Expression {
	return Expression{
		ExpressionConnective: expr.ExpressionConnective,
		AtomicExpressions:    expr.AtomicExpressions,
		LookupFunc:           lookupFunc,
	}
}

// NewExpressionWithLookupFunc creates a new expression with the given
// connective, list of expressions, and lookup function.
func NewExpressionWithLookupFunc(
	connective EvaluatorMode,
	exprs [][]*AtomicExpression,
	lookupFunc func(any) any) Expression {
	return Expression{
		ExpressionConnective: connective,
		AtomicExpressions:    exprs,
		LookupFunc:           lookupFunc}
}

type boolEvaler func(any) (bool, error)

func anyOf[P any, F func(P) (bool, error)](eval F, params []P) (bool, error) {
	for _, item := range params {
		if val, err := eval(item); err != nil {
			return false, err
		} else if val {
			return true, nil
		}
	}
	return false, nil
}

func allOf[P any, F func(P) (bool, error)](eval F, params []P) (bool, error) {
	for _, item := range params {
		if val, err := eval(item); err != nil {
			return false, err
		} else if !val {
			return false, nil
		}
	}
	return true, nil

}

// Evaluate() returns the value of the expression, or an error if
// something was malformed.
func (expr Expression) Evaluate() (bool, error) {
	return expr.evalExpressionClauses(expr.AtomicExpressions)
}

func noneOf[P any, F func(P) (bool, error)](eval F, params []P) (bool, error) {
	if wasOneOf, err := anyOf(eval, params); err != nil {
		return false, err
	} else {
		return !wasOneOf, nil
	}
}

func notOfResult(result bool, err error) (bool, error) {
	if err != nil {
		return false, err
	} else {
		return !result, nil
	}

}
func (e Expression) evalExpressionClauses(expr [][]*AtomicExpression) (bool, error) {
	switch e.ExpressionConnective {
	case AndOfOrsMode:
		return allOf(e.evalClause, expr)
	case OrOfAndsMode:
		return anyOf(e.evalClause, expr)
	}
	return false, fmt.Errorf("booleval: unknown mode passed to evaluator: %v", e.ExpressionConnective)
}

func (e Expression) evalClause(clause []*AtomicExpression) (bool, error) {
	switch e.ExpressionConnective {
	case AndOfOrsMode:
		return anyOf(e.evalAtomicExpression, clause)
	case OrOfAndsMode:
		return allOf(e.evalAtomicExpression, clause)
	}
	return false, fmt.Errorf("booleval: unknown mode passed to evaluator: %v", e.ExpressionConnective)
}

func (e Expression) evalAtomicExpression(cond *AtomicExpression) (bool, error) {
	switch cond.Operator {
	case "==":
		return cond.CompareValue.Equal(e.LookupFunc(cond.ActualValue))
	case "!=":
		return notOfResult(cond.CompareValue.Equal(e.LookupFunc(cond.ActualValue)))
	case "<":
		return noneOf(
			func(evaluator boolEvaler) (bool, error) {
				return evaluator(e.LookupFunc(cond.ActualValue))
			},
			[]boolEvaler{cond.CompareValue.Equal, cond.CompareValue.Greater},
		)

	case ">":
		return cond.CompareValue.Greater(e.LookupFunc(cond.ActualValue))
	case "<=":
		return notOfResult(cond.CompareValue.Greater(e.LookupFunc(cond.ActualValue)))
	case ">=":
		return anyOf(
			func(evaluator boolEvaler) (bool, error) {
				return evaluator(e.LookupFunc(cond.ActualValue))
			},
			[]boolEvaler{cond.CompareValue.Greater, cond.CompareValue.Greater})
	}
	return false, fmt.Errorf("booleval: EvalCondition: no such operator %v", cond.Operator)
}
