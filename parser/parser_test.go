package parser

import (
	"testing"
	"yac/ast"
	"yac/lexer"
	"yac/token"
)

func TestAst(t *testing.T) {
	var node = &ast.InfixExpression{
		Token: token.Token{
			Type:    token.MINUS,
			Literal: "-",
		},
		Operator: "-",
		Left: ast.Integer{
			Token: token.Token{
				Type:    token.NUMBER,
				Literal: "50",
			},
			Value: 50,
		},
		Right: ast.Integer{
			Token: token.Token{
				Type:    token.NUMBER,
				Literal: "10",
			},
			Value: 10,
		},
	}
	var p = New(lexer.New("50-10"))
	var pp = p.Parse()
	if len(pp) != 1 || pp[0].String() != node.String() {
		t.Errorf("Wrong output. expected=%s", node.String())
	}
}

func TestIntegerNode(t *testing.T) {
	var node = ast.Integer{
		Token: token.Token{
			Type:    token.NUMBER,
			Literal: "10",
		},
		Value: 10,
	}
	var p = New(lexer.New("10"))
	var pp = p.Parse()
	if len(pp) != 1 || pp[0].String() != node.String() {
		t.Errorf("Wrong output. expected=%s", node.String())
	}
}

func TestGroupedExpression(t *testing.T) {
	var p = New(lexer.New("132 - 54/32   + -12 * 5"))
	var pp = p.Parse()
	if len(pp) != 1 || "((132-(54/32))+((-12)*5))" != pp[0].String() {
		t.Errorf("Expected val: ((132-(54/32))+((-12)*5))")
	}
}
