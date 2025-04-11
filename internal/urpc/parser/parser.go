package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/uforg/uforpc/internal/urpc/ast"
	"github.com/uforg/uforpc/internal/urpc/lexer"
)

type Parser = participle.Parser[ast.Schema]

var ParserInstance = participle.MustBuild[ast.Schema](
	participle.Lexer(&lexer.ParticipleLexer{}),
	participle.Elide("Newline", "Whitespace"),
)
