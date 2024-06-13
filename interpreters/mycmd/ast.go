package mycmd

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

type LetStatement struct {
	Token Token
	Name  *Identifier
	Value Expression
}

type ExecuteStatement struct {
	Token      Token
	ScriptName string
}

type Identifier struct {
	Token Token
	Value string
}

type IntegerLiteral struct {
	Token Token
	Value string
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (es *ExecuteStatement) statementNode()       {}
func (es *ExecuteStatement) TokenLiteral() string { return es.Token.Literal }

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
