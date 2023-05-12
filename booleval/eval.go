package booleval

import (
	"fmt"
)

type Comparable interface {
	Equal(other any) (bool, error)
	Greater(other any) (bool, error)
}

type Condition struct {
	Operator       string
	PossibleValues []Comparable
	ActualValue    any
}

func EvalConditions(conds []Condition) (bool, error) {
	for _, cond := range conds {
		didPass, err := EvalCondition(cond)
		if !didPass {
			return false, nil
		} else if err != nil {
			return false, err

		}
	}
	return true, nil
}

type boolEvaler func(Condition) (bool, error)

func oneOf(params []boolEvaler, condition Condition) (bool, error) {
	for _, param := range params {
		if val, err := param(condition); err != nil {
			return false, err
		} else if val {
			return true, nil
		}
	}
	return false, nil
}
func noneOf(params []boolEvaler, condition Condition) (bool, error) {
	if wasOneOf, err := oneOf(params, condition); err != nil {
		return false, err
	} else {
		return !wasOneOf, nil
	}
}

func notCond(param boolEvaler, condition Condition) (bool, error) {
	if wasTrue, err := EvalEquals(condition); err != nil {
		return false, err
	} else {
		return !wasTrue, nil
	}
}

func EvalCondition(cond Condition) (bool, error) {
	switch cond.Operator {
	case "==":
		return EvalEquals(cond)
	case "!=":
		return notCond(EvalEquals, cond)
	case "<":
		return noneOf([]boolEvaler{EvalGreater, EvalEquals}, cond)
	case ">":
		return EvalGreater(cond)
	case "<=":
		return notCond(EvalGreater, cond)
	case ">=":
		return oneOf([]boolEvaler{EvalGreater, EvalEquals}, cond)
	}
	return false, fmt.Errorf("booleval: EvalCondition: no such operator %v\n", cond.Operator)
}

func EvalEquals(cond Condition) (returnVal bool, err error) {
	returnVal = false
	for _, value := range cond.PossibleValues {
		if wasTrue, err := value.Equal(cond.ActualValue); err != nil {
			return false, err
		} else if wasTrue {
			return true, nil
		}
	}
	return
}

func EvalGreater(cond Condition) (returnVal bool, err error) {
	returnVal = false
	for _, value := range cond.PossibleValues {
		if wasTrue, err := value.Greater(cond.ActualValue); err != nil {
			return false, err
		} else if wasTrue {
			return true, nil
		}
	}
	return
}
