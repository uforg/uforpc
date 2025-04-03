package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/uforg/uforpc/internal/urpc/ast"
	"github.com/uforg/uforpc/internal/urpc/lexer"
)

var Parser = participle.MustBuild[ast.URPCSchema](
	participle.Lexer(&lexer.ParticipleLexer{}),
)
