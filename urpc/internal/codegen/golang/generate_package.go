package golang

import (
	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/schema"
)

func generatePackage(_ schema.Schema, config Config) (string, error) {
	g := genkit.NewGenKit().WithTabs()

	g.Line("// This file has been generated using UFO RPC. DO NOT EDIT.")
	g.Line("// If you edit this file, it will be overwritten the next time it is generated.")
	g.Line("//nolint")
	g.Break()

	g.Linef("// Package %s contains the generated code for the UFO RPC schema", config.PackageName)
	g.Linef("package %s", config.PackageName)
	g.Break()

	g.Line("import (")
	g.Block(func() {
		g.Line(`"encoding/json"`)
		g.Line(`"fmt"`)
		g.Line(`"io"`)
		g.Line(`"time"`)
		g.Line(`"slices"`)
	})
	g.Line(")")
	g.Break()

	return g.String(), nil
}
