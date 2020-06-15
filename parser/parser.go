package parser

import (
	"fmt"
	"strconv"
	"yac/ast"
	"yac/lexer"
	"yac/token"
)

const (
	LOWEST = iota + 1
	SUM
	PRODUCT
	PREFIX
)

var precedences = map[token.Type]int{
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.DIV:      PRODUCT,
	token.ASTERISK: PRODUCT,
}

type (
	prefixParser func() ast.Node
	infixParser  func(node ast.Node) ast.Node
)

type Parser struct {
	lexer         *lexer.Lexer
	errors        []string
	currentToken  token.Token
	peekToken     token.Token
	prefixParsers map[token.Type]prefixParser
	infixParsers  map[token.Type]infixParser
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) missingHandlerError(t token.Type) {
	var err string
	switch t {
	case token.ILLEGAL:
		err = fmt.Sprintf("Illegal token: %s", p.currentToken.Literal)
	default:
		err = fmt.Sprintf("no prefix parse function for %s found", t)
	}
	p.errors = append(p.errors, err)
}

func (p *Parser) registerPrefixParser(t token.Type, fn prefixParser) {
	p.prefixParsers[t] = fn
}

func (p *Parser) registerInfixParser(t token.Type, fn infixParser) {
	p.infixParsers[t] = fn
}

func (p *Parser) Parse() []ast.Node {
	var nodes = make([]ast.Node, 0)
	for !p.curTokenIs(token.EOF) {
		var node = p.parse(LOWEST)
		if node != nil {
			nodes = append(nodes, node)
		}
		p.nextToken()
	}
	return nodes
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parse(precedence int) ast.Node {
	var prefix = p.prefixParsers[p.currentToken.Type]
	if prefix == nil {
		p.missingHandlerError(p.currentToken.Type)
		return nil
	}
	var leftExpression = prefix()
	for !p.peekTokenIs(token.EOF) && p.peekPrecedence() > precedence {
		var infix = p.infixParsers[p.peekToken.Type]
		if infix == nil {
			return leftExpression
		}
		p.nextToken()
		leftExpression = infix(leftExpression)
	}
	return leftExpression
}

func New(l *lexer.Lexer) *Parser {
	var p = &Parser{
		lexer:  l,
		errors: []string{},
	}

	p.prefixParsers = make(map[token.Type]prefixParser)
	p.registerPrefixParser(token.NUMBER, p.parseInteger)
	p.registerPrefixParser(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixParser(token.LPAREN, p.parseGroupedExpression)

	p.infixParsers = make(map[token.Type]infixParser)
	p.registerInfixParser(token.MINUS, p.parseInfixExpression)
	p.registerInfixParser(token.PLUS, p.parseInfixExpression)
	p.registerInfixParser(token.DIV, p.parseInfixExpression)
	p.registerInfixParser(token.ASTERISK, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) parseInteger() ast.Node {
	var node = &ast.Integer{Token: p.currentToken}
	var val, err = strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal))
		return nil
	}
	node.Value = val
	return node
}

func (p *Parser) parsePrefixExpression() ast.Node {
	var expression = &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parse(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Node) ast.Node {
	var expression = &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	var precedence = p.curPrecedence()
	p.nextToken()
	expression.Right = p.parse(precedence)

	return expression
}

func (p *Parser) peekError(t token.Type) {
	var msg = fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) parseGroupedExpression() ast.Node {
	p.nextToken()
	var expression = p.parse(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return expression
}
