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

	// Identifiers
	IDENT TokenType = "IDENT"

	// Literals
	STRING  TokenType = "STRING"
	INT     TokenType = "INT"
	FLOAT   TokenType = "FLOAT"
	BOOLEAN TokenType = "BOOLEAN"

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	COLON     TokenType = ":"
	DOT       TokenType = "."
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	LBRACKET  TokenType = "["
	RBRACKET  TokenType = "]"
	AT        TokenType = "@"
	QUESTION  TokenType = "?"
	NEWLINE   TokenType = "\n"

	// Keywords
	TYPE TokenType = "TYPE"
	PROC TokenType = "PROC"
)
