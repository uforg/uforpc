package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
	"github.com/uforg/uforpc/urpc/internal/urpc/lexer"
)

// Error is an alias for participle.Error
type Error = participle.Error

// Parser is an alias for participle.Parser with ast.Schema as the root node
type Parser = participle.Parser[ast.Schema]

// ParserInstance is a pre-built parser instance for URPC schemas.
var ParserInstance = participle.MustBuild[ast.Schema](
	participle.Lexer(&lexer.ParticipleLexer{}),
	participle.Elide("Newline", "Whitespace"),
	participle.UseLookahead(4),
)
