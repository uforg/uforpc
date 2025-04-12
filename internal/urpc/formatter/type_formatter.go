package formatter

import (
	"strings"

	"github.com/uforg/uforpc/internal/genkit"
	"github.com/uforg/uforpc/internal/urpc/ast"
)

type typeFormatter struct {
	g        *genkit.GenKit
	typeDecl *ast.TypeDecl
}

func newTypeFormatter(typeDecl *ast.TypeDecl) *typeFormatter {
	if typeDecl == nil {
		typeDecl = &ast.TypeDecl{}
	}

	return &typeFormatter{
		g:        genkit.NewGenKit().WithSpaces(2),
		typeDecl: typeDecl,
	}
}

// format formats the entire typeDecl, handling spacing and EOL comments.
//
// Returns the formatted genkit.GenKit.
func (f *typeFormatter) format() *genkit.GenKit {
	if f.typeDecl.Docstring != nil {
		f.g.Linef(`"""%s"""`, f.typeDecl.Docstring.Value)
	}
	f.g.Inlinef(`type %s `, f.typeDecl.Name)

	if len(f.typeDecl.Extends) > 0 {
		joinedExtends := strings.Join(f.typeDecl.Extends, ", ")
		f.g.Inlinef(`extends %s `, joinedExtends)
	}

	fieldsFormatter := newFieldsFormatter(f.typeDecl, f.typeDecl.Children)
	f.g.Line(strings.TrimSpace(fieldsFormatter.format().String()))

	return f.g
}
