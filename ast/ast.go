package ast

import (
	"bytes"
	"yac/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Node
}

func (p *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}

func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Node
	Right    Node
}

func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(i.Operator)
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}

type Integer struct {
	Token token.Token
	Value int64
}

func (i Integer) TokenLiteral() string {
	return i.Token.Literal
}

func (i Integer) String() string {
	return i.Token.Literal
}
