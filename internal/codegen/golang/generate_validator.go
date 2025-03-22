package golang

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/uforg/uforpc/internal/codegen/genkit"
	"github.com/uforg/uforpc/internal/schema"
)

//go:embed pieces/validator/validator.go
var validatorRawPiece string

func generateValidator(g *genkit.GenKit, _ schema.Schema, config Config) error {
	if config.OmitClientRequestValidation && config.OmitServerRequestValidation {
		return nil
	}

	split := strings.Split(validatorRawPiece, "/** START FROM HERE **/")
	if len(split) < 2 {
		return fmt.Errorf("validator.go: could not find start marker")
	}

	g.Inline(split[1])
	g.Break()

	return nil
}
