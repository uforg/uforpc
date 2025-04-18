package golang

import (
	_ "embed"
	"fmt"

	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/util/strutil"
)

//go:embed pieces/required_validator.go
var requiredValidatorRawPiece string

func generateValidator(_ schema.Schema, config Config) (string, error) {
	piece := strutil.GetStrAfter(requiredValidatorRawPiece, "/** START FROM HERE **/")
	if piece == "" {
		return "", fmt.Errorf("required_validator.go: could not find start delimiter")
	}
	return piece, nil
}
