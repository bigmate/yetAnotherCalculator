package evaluator

import (
	"fmt"
	"yac/ast"
)

type Type string

const (
	ERROR   Type = "ERROR"
	INTEGER Type = "INTEGER"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type {
	return INTEGER
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Error struct {
	Value string
}

func (e *Error) Type() Type {
	return ERROR
}

func (e *Error) Inspect() string {
	return e.Value
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Value: fmt.Sprintf(format, a...)}
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR
	}
	return false
}

func Eval(node ast.Node) Object {
	switch node := node.(type) {
	case *ast.Integer:
		return &Integer{Value: node.Value}

	case *ast.PrefixExpression:
		var right = Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		var left = Eval(node.Left)
		if isError(left) {
			return left
		}

		var right = Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	}
	return nil
}

func evalPrefixExpression(operator string, right Object) Object {
	if operator != "-" {
		return newError("Unknown operator: %s", operator)
	}
	if right.Type() != INTEGER {
		return newError("Unknown operand: %s", right.Type())
	}
	var integer = right.(*Integer)
	integer.Value *= -1
	return integer
}

func evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case left.Type() != right.Type():
		return newError("Type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	default:
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right Object) Object {
	var l = left.(*Integer).Value
	var r = right.(*Integer).Value
	switch operator {
	case "+":
		return &Integer{Value: l + r}
	case "-":
		return &Integer{Value: l - r}
	case "*":
		return &Integer{Value: l * r}
	case "/":
		return &Integer{Value: l / r}
	default:
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}
