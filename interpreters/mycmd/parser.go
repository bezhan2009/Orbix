package mycmd

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for p.curToken.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case LET:
		return p.parseLetStatement()
	case IDENT:
		if p.curToken.Literal == "execute" {
			return p.parseExecuteStatement()
		}
	}
	return nil
}

func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression()

	return stmt
}

func (p *Parser) parseExecuteStatement() *ExecuteStatement {
	stmt := &ExecuteStatement{Token: p.curToken}

	if !p.expectPeek(STRING) {
		return nil
	}

	stmt.ScriptName = p.curToken.Literal

	return stmt
}

func (p *Parser) parseExpression() Expression {
	switch p.curToken.Type {
	case INT:
		return &IntegerLiteral{Token: p.curToken, Value: p.curToken.Literal}
	case IDENT:
		return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}
	return nil
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	return false
}
