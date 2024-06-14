package mycmd

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[string]int{
	"==": EQUALS,
	"!=": EQUALS,
	"<":  LESSGREATER,
	">":  LESSGREATER,
	"+":  SUM,
	"-":  SUM,
	"*":  PRODUCT,
	"/":  PRODUCT,
	"(":  CALL,
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseFunctionLiteral() Expression {
	lit := &FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*Identifier {
	identifiers := []*Identifier{}

	if p.peekTokenIs(RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end TokenType) []Expression {
	list := []Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}


package mycmd

import (
	"fmt"
)

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:            l,
		errors:       []string{},
		prefixParseFns: make(map[TokenType]prefixParseFn),
		infixParseFns:  make(map[TokenType]infixParseFn),
	}

	p.registerPrefix(IDENT, p.parseIdentifier)
	p.registerPrefix(INT, p.parseIntegerLiteral)
	p.registerPrefix(BANG, p.parsePrefixExpression)
	p.registerPrefix(MINUS, p.parsePrefixExpression)
	p.registerPrefix(TRUE, p.parseBoolean)
	p.registerPrefix(FALSE, p.parseBoolean)
	p.registerPrefix(LPAREN, p.parseGroupedExpression)
	p.registerPrefix(IF, p.parseIfExpression)
	p.registerPrefix(FUNCTION, p.parseFunctionLiteral)

	p.registerInfix(PLUS, p.parseInfixExpression)
	p.registerInfix(MINUS, p.parseInfixExpression)
	p.registerInfix(SLASH, p.parseInfixExpression)
	p.registerInfix(ASTERISK, p.parseInfixExpression)
	p.registerInfix(EQ, p.parseInfixExpression)
	p.registerInfix(NOT_EQ, p.parseInfixExpression)
	p.registerInfix(LT, p.parseInfixExpression)
	p.registerInfix(GT, p.parseInfixExpression)
	p.registerInfix(LPAREN, p.parseCallExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) parseFunctionLiteral() Expression {
	lit := &FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*Identifier {
	identifiers := []*Identifier{}

	if p.peekTokenIs(RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return identifiers
}


package mycmd

import (
	"fmt"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[string]int{
	"==": EQUALS,
	"!=": EQUALS,
	"<":  LESSGREATER,
	">":  LESSGREATER,
	"+":  SUM,
	"-":  SUM,
	"*":  PRODUCT,
	"/":  PRODUCT,
	"(":  CALL,
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseFunctionLiteral() Expression {
	lit := &FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*Identifier {
	identifiers := []*Identifier{}

	if p.peekTokenIs(RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end TokenType) []Expression {
	list := []Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}
