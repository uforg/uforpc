package golang

import (
	"github.com/uforg/uforpc/internal/codegen/genkit"
	"github.com/uforg/uforpc/internal/schema"
)

func generatePackage(g *genkit.GenKit, _ schema.Schema, config Config) error {
	g.Inline("// This file has been generated using UFO RPC. DO NOT EDIT.")
	g.Line("// If you edit this file, it will be overwritten the next time it is generated.")
	g.Break()

	g.Linef("// Package %s contains the generated code for the UFO RPC schema", config.PackageName)
	g.Linef("package %s", config.PackageName)
	g.Break()

	g.Line("import (")
	g.Block(func() {
		g.Line(`"encoding/json"`)
		g.Line(`"fmt"`)
		g.Line(`"regexp"`)
		g.Line(`"strings"`)
	})
	g.Line(")")
	g.Break()

	return nil
}
