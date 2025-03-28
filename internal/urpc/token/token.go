package token

type TokenType string

type Token struct {
	Type     TokenType
	Literal  string
	FileName string
	Line     int
	Column   int
}

const (
	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"

	// Identifiers and literals
	IDENT     TokenType = "IDENT"
	STRING    TokenType = "STRING"
	INT       TokenType = "INT"
	FLOAT     TokenType = "FLOAT"
	COMMENT   TokenType = "COMMENT"
	DOCSTRING TokenType = "DOCSTRING"

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

	// Keywords
	VERSION TokenType = "VERSION"
	RULE    TokenType = "RULE"
	TYPE    TokenType = "TYPE"
	PROC    TokenType = "PROC"
	INPUT   TokenType = "INPUT"
	OUTPUT  TokenType = "OUTPUT"
	META    TokenType = "META"
	ERROR   TokenType = "ERROR"
	TRUE    TokenType = "TRUE"
	FALSE   TokenType = "FALSE"
)

// delimiters is a map of delimiters to their corresponding token types.
var delimiters = map[string]TokenType{
	":": COLON,
	",": COMMA,
	"(": LPAREN,
	")": RPAREN,
	"{": LBRACE,
	"}": RBRACE,
	"[": LBRACKET,
	"]": RBRACKET,
	"@": AT,
	"?": QUESTION,
}

// IsDelimiter returns true if the character is a delimiter.
func IsDelimiter(ch byte) bool {
	_, ok := delimiters[string(ch)]
	return ok
}

// GetDelimiterTokenType returns the token type for the given delimiter.
func GetDelimiterTokenType(ch byte) TokenType {
	return delimiters[string(ch)]
}

// keywords is a map of keywords to their corresponding token types.
var keywords = map[string]TokenType{
	"version": VERSION,
	"rule":    RULE,
	"type":    TYPE,
	"proc":    PROC,
	"input":   INPUT,
	"output":  OUTPUT,
	"meta":    META,
	"error":   ERROR,
	"true":    TRUE,
	"false":   FALSE,
}

// IsKeyword returns true if the identifier is a keyword.
func IsKeyword(ident string) bool {
	_, ok := keywords[ident]
	return ok
}

// GetKeywordTokenType returns the token type for the given keyword.
func GetKeywordTokenType(ident string) TokenType {
	return keywords[ident]
}
