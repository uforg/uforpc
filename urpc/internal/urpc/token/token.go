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
	Ident        TokenType = "Ident"
	Comment      TokenType = "Comment"      // Single line comment with //
	CommentBlock TokenType = "CommentBlock" // Multiline comment with /* */
	Docstring    TokenType = "Docstring"

	// Literals
	StringLiteral TokenType = "StringLiteral"
	IntLiteral    TokenType = "IntLiteral"
	FloatLiteral  TokenType = "FloatLiteral"
	TrueLiteral   TokenType = "TrueLiteral"
	FalseLiteral  TokenType = "FalseLiteral"

	// Operators and delimiters
	Newline    TokenType = "Newline"
	Whitespace TokenType = "Whitespace"
	Colon      TokenType = "Colon"
	Comma      TokenType = "Comma"
	LParen     TokenType = "LParen"
	RParen     TokenType = "RParen"
	LBrace     TokenType = "LBrace"
	RBrace     TokenType = "RBrace"
	LBracket   TokenType = "LBracket"
	RBracket   TokenType = "RBracket"
	At         TokenType = "At"
	Question   TokenType = "Question"

	// Keywords
	Version    TokenType = "Version"
	Import     TokenType = "Import"
	Deprecated TokenType = "Deprecated"
	Rule       TokenType = "Rule"
	Type       TokenType = "Type"
	Extends    TokenType = "Extends"
	Proc       TokenType = "Proc"
	Input      TokenType = "Input"
	Output     TokenType = "Output"
	Meta       TokenType = "Meta"
	Error      TokenType = "Error"
	For        TokenType = "For"
	Param      TokenType = "Param"
	String     TokenType = "String"
	Int        TokenType = "Int"
	Float      TokenType = "Float"
	Bool       TokenType = "Bool"
	Datetime   TokenType = "Datetime"
)

var TokenTypes = []TokenType{
	// Special tokens
	Eof,
	Illegal,

	// Identifiers, comments and docstrings
	Ident,
	Comment,
	CommentBlock,
	Docstring,

	// Literals
	StringLiteral,
	IntLiteral,
	FloatLiteral,
	TrueLiteral,
	FalseLiteral,

	// Operators and delimiters
	Newline,
	Whitespace,
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
	Import,
	Deprecated,
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
	Bool,
	Datetime,
}

// delimiters is a map of delimiters to their corresponding token types.
var delimiters = map[rune]TokenType{
	'\n': Newline,
	':':  Colon,
	',':  Comma,
	'(':  LParen,
	')':  RParen,
	'{':  LBrace,
	'}':  RBrace,
	'[':  LBracket,
	']':  RBracket,
	'@':  At,
	'?':  Question,
}

// IsDelimiter returns true if the character is a delimiter.
func IsDelimiter(ch rune) bool {
	_, ok := delimiters[ch]
	return ok
}

// GetDelimiterTokenType returns the token type for the given delimiter.
func GetDelimiterTokenType(ch rune) TokenType {
	return delimiters[ch]
}

// keywords is a map of keywords to their corresponding token types.
var keywords = map[string]TokenType{
	"version":    Version,
	"import":     Import,
	"deprecated": Deprecated,
	"rule":       Rule,
	"type":       Type,
	"extends":    Extends,
	"proc":       Proc,
	"input":      Input,
	"output":     Output,
	"meta":       Meta,
	"error":      Error,
	"for":        For,
	"param":      Param,
	"string":     String,
	"int":        Int,
	"float":      Float,
	"bool":       Bool,
	"datetime":   Datetime,
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
