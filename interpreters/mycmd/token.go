package mycmd

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1343456

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	PRINT    = "PRINT"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":    FUNCTION,
	"let":   LET,
	"print": PRINT,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
