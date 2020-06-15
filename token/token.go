package token

type Type string

const (
	ILLEGAL  Type = "ILLEGAL"
	EOF      Type = "EOF"
	NUMBER   Type = "NUMBER"
	PLUS     Type = "+"
	MINUS    Type = "-"
	ASTERISK Type = "*"
	DIV      Type = "/"
	LPAREN   Type = "("
	RPAREN   Type = ")"
)

type Token struct {
	Type    Type
	Literal string
}
