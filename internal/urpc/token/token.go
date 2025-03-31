package token

type TokenType string

type Token struct {
	// Type is the type of the token.
	Type TokenType
	// Literal is the literal value of the token.
	Literal string
	// FileName is the name of the file where the token was found.
	FileName string
	// LineStart is the line number of the first character of the token.
	LineStart int
	// LineEnd is the line number of the last character of the token.
	LineEnd int
	// ColumnStart is the column number of the first character of the token.
	ColumnStart int
	// ColumnEnd is the column number of the last character of the token.
	ColumnEnd int
}

const (
	// Special tokens
	Eof     TokenType = "Eof"
	Illegal TokenType = "Illegal"

	// Identifiers, comments and docstrings
	Ident     TokenType = "Ident"
	Comment   TokenType = "Comment"
	Docstring TokenType = "Docstring"

	// Identifiers and literals
	StringLiteral TokenType = "StringLiteral"
	IntLiteral    TokenType = "IntLiteral"
	FloatLiteral  TokenType = "FloatLiteral"
	TrueLiteral   TokenType = "TrueLiteral"
	FalseLiteral  TokenType = "FalseLiteral"

	// Operators and delimiters
	Colon    TokenType = "Colon"
	Comma    TokenType = "Comma"
	LParen   TokenType = "LParen"
	RParen   TokenType = "RParen"
	LBrace   TokenType = "LBrace"
	RBrace   TokenType = "RBrace"
	LBracket TokenType = "LBracket"
	RBracket TokenType = "RBracket"
	At       TokenType = "At"
	Question TokenType = "Question"

	// Keywords
	Version TokenType = "Version"
	Rule    TokenType = "Rule"
	Type    TokenType = "Type"
	Extends TokenType = "Extends"
	Proc    TokenType = "Proc"
	Input   TokenType = "Input"
	Output  TokenType = "Output"
	Meta    TokenType = "Meta"
	Error   TokenType = "Error"
	For     TokenType = "For"
	Param   TokenType = "Param"
	String  TokenType = "String"
	Int     TokenType = "Int"
	Float   TokenType = "Float"
	Boolean TokenType = "Boolean"
)

var TokenTypes = []TokenType{
	// Special tokens
	Eof,
	Illegal,

	// Identifiers, comments and docstrings
	Ident,
	Comment,
	Docstring,

	// Literals
	StringLiteral,
	IntLiteral,
	FloatLiteral,
	TrueLiteral,
	FalseLiteral,

	// Operators and delimiters
	Colon,
	Comma,
	LParen,
	RParen,
	LBrace,
	RBrace,
	LBracket,
	RBracket,
	At,
	Question,

	// Keywords
	Version,
	Rule,
	Type,
	Extends,
	Proc,
	Input,
	Output,
	Meta,
	Error,
	For,
	Param,
	String,
	Int,
	Float,
	Boolean,
}

// delimiters is a map of delimiters to their corresponding token types.
var delimiters = map[string]TokenType{
	":": Colon,
	",": Comma,
	"(": LParen,
	")": RParen,
	"{": LBrace,
	"}": RBrace,
	"[": LBracket,
	"]": RBracket,
	"@": At,
	"?": Question,
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
	"version": Version,
	"rule":    Rule,
	"type":    Type,
	"extends": Extends,
	"proc":    Proc,
	"input":   Input,
	"output":  Output,
	"meta":    Meta,
	"error":   Error,
	"for":     For,
	"param":   Param,
	"string":  String,
	"int":     Int,
	"float":   Float,
	"boolean": Boolean,
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
