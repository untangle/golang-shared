package booleval

import (
	"fmt"
)

type Comparable interface {
	Equal(other any) (bool, error)
	Greater(other any) (bool, error)
}

type GreaterNotApplicable struct{}

func (GreaterNotApplicable) Greater(other any) (bool, error) {
	return false, fmt.Errorf("this type does not support ordering")
}

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

type EvaluatorMode int

type Expression struct {
	ExpressionConnective EvaluatorMode
	AtomicExpressions    [][]*AtomicExpression
	LookupFunc           func(any) any
}

func NewSimpleExpression(
	connective EvaluatorMode,
	exprs [][]*AtomicExpression) Expression {
	return Expression{
		ExpressionConnective: connective,
		AtomicExpressions:    exprs,
		LookupFunc:           func(v any) any { return v },
	}
}

func ExpressionCopyWithLookupFunc(
	expr Expression,
	lookupFunc func(any) any) Expression {
	return Expression{
		ExpressionConnective: expr.ExpressionConnective,
		AtomicExpressions:    expr.AtomicExpressions,
		LookupFunc:           lookupFunc,
	}
}

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

func (expr Expression) Evaluate() (bool, error) {
	return expr.EvalExpressionClauses(expr.AtomicExpressions)
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
func (e Expression) EvalExpressionClauses(expr [][]*AtomicExpression) (bool, error) {
	switch e.ExpressionConnective {
	case AndOfOrsMode:
		return allOf(e.EvalClause, expr)
	case OrOfAndsMode:
		return anyOf(e.EvalClause, expr)
	}
	return false, fmt.Errorf("booleval: unknown mode passed to evaluator: %v", e.ExpressionConnective)
}

func (e Expression) EvalClause(clause []*AtomicExpression) (bool, error) {
	switch e.ExpressionConnective {
	case AndOfOrsMode:
		return anyOf(e.EvalAtomicExpression, clause)
	case OrOfAndsMode:
		return allOf(e.EvalAtomicExpression, clause)
	}
	return false, fmt.Errorf("booleval: unknown mode passed to evaluator: %v", e.ExpressionConnective)
}

func (e Expression) EvalAtomicExpression(cond *AtomicExpression) (bool, error) {
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
