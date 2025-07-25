package formatter

import (
	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

type typeFormatter struct {
	g        *genkit.GenKit
	typeDecl *ast.TypeDecl
}

func newTypeFormatter(g *genkit.GenKit, typeDecl *ast.TypeDecl) *typeFormatter {
	if typeDecl == nil {
		typeDecl = &ast.TypeDecl{}
	}

	return &typeFormatter{
		g:        g,
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

	if f.typeDecl.Deprecated != nil {
		if f.typeDecl.Deprecated.Message == nil {
			f.g.Inline("deprecated ")
		}
		if f.typeDecl.Deprecated.Message != nil {
			f.g.Linef("deprecated(\"%s\")", strutil.EscapeQuotes(*f.typeDecl.Deprecated.Message))
		}
	}

	// Force strict pascal case
	f.g.Inlinef(`type %s `, strutil.ToPascalCase(f.typeDecl.Name))

	fieldsFormatter := newFieldsFormatter(f.g, f.typeDecl, f.typeDecl.Children)
	fieldsFormatter.format()

	return f.g
}
