package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/uforg/uforpc/internal/urpc/ast"
	"github.com/uforg/uforpc/internal/urpc/lexer"
)

var Parser = participle.MustBuild[ast.Schema](
	participle.Lexer(&lexer.ParticipleLexer{}),
	participle.Elide("Whitespace"),
)
