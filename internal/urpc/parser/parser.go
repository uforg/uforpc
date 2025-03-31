package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/uforg/uforpc/internal/urpc/ast"
)

var Parser = participle.MustBuild[ast.URPCSchema](
	participle.Lexer(&AdaptedLexer{}),
	participle.Elide("Comment"),
)
