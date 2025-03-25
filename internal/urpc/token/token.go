package token

type TokenType string

type Token struct {
	Type     TokenType
	Literal  string
	FileName string
	Line     int
	Column   int
}

func NewToken(typ TokenType, lit byte, fileName string, line int, column int) Token {
	return Token{
		Type:     typ,
		Literal:  string(lit),
		FileName: fileName,
		Line:     line,
		Column:   column,
	}
}

const (
	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"

	// Identifiers and literals
	IDENT  TokenType = "IDENT"
	STRING TokenType = "STRING"
	INT    TokenType = "INT"
	FLOAT  TokenType = "FLOAT"
	TRUE   TokenType = "TRUE"
	FALSE  TokenType = "FALSE"

	// Operators and delimiters
	COLON    TokenType = ":"
	COMMA    TokenType = ","
	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"
	AT       TokenType = "@"
	QUESTION TokenType = "?"
	NEWLINE  TokenType = "\n"

	// Keywords
	VERSION TokenType = "VERSION"
	TYPE    TokenType = "TYPE"
	PROC    TokenType = "PROC"
	INPUT   TokenType = "INPUT"
	OUTPUT  TokenType = "OUTPUT"
	META    TokenType = "META"
	ERROR   TokenType = "ERROR"
)
