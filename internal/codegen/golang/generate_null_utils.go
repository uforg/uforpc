package golang

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/uforg/uforpc/internal/codegen/genkit"
	"github.com/uforg/uforpc/internal/schema"
)

//go:embed pieces/nullutils/nullutils.go
var nullUtilsRawPiece string

func generateNullUtils(g *genkit.GenKit, _ schema.Schema, _ Config) error {
	split := strings.Split(nullUtilsRawPiece, "/** START FROM HERE **/")
	if len(split) < 2 {
		return fmt.Errorf("nullutils.go: could not find start marker")
	}

	g.Inline(split[1])
	g.Break()

	return nil
}
